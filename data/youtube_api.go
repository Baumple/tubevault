package data

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YouTubeApi struct {
	youtubeService *youtube.Service
}

type config struct {
	ApiKey string
}

func getApiKey() string {
	file, err := GetSaveFile()
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	c := config{}

	err = json.NewDecoder(file).Decode(&c)

	if c.ApiKey == "" {
		fmt.Print("No api key provided.\nPlease enter your api key: ")
		fmt.Scan(&c.ApiKey)
		err = json.NewEncoder(file).Encode(&c)
		if err != nil {
			log.Fatalf("Could not store api key: %s", err.Error())
		}

	} else if err != nil {
		log.Fatal(err)
	}

	return c.ApiKey
}

func NewYouTubeApi() YouTubeApi {
	// TODO: Prompt user to enter id and store
	// it in a config file

	apiKey := getApiKey()
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Could not initiate youtube api: %v\n", err.Error())
	}
	return YouTubeApi{
		youtubeService: youtubeService,
	}
}

func (yt *YouTubeApi) GetYoutubePlaylistsById(id string) []Playlist {
	playlistsResp, err := yt.youtubeService.Playlists.List([]string{"snippet", "id"}).Id(id).Do()
	if err != nil {
		log.Fatal(err)
	}

	playlists := []Playlist{}
	for _, playlistResp := range playlistsResp.Items {
		time, err := time.Parse(time.RFC3339, playlistResp.Snippet.PublishedAt)
		if err != nil {
			log.Fatal(err)
		}

		var playlist Playlist
		playlist.Id = playlistResp.Id
		playlist.Title = playlistResp.Snippet.Title
		playlist.PublishedAt = time
		playlist.Description = playlistResp.Snippet.Description
		playlists = append(playlists, playlist)
	}

	return playlists
}

func (yt *YouTubeApi) GetYoutubePlaylistsBySearch(search string) []Playlist {
	playlistsResp, err := yt.youtubeService.Search.List([]string{"id",
		"snippet"}).
		Q(search).
		Type("playlist").
		Do()

	if err != nil {
		log.Fatal(err)
	}

	playlists := []Playlist{}
	for _, playlistResp := range playlistsResp.Items {
		publishedAt, err := time.Parse(time.RFC3339, playlistResp.Snippet.PublishedAt)
		if err != nil {
			log.Fatal(err)
		}

		var playlist Playlist
		playlist.Title = playlistResp.Snippet.Title
		playlist.Id = playlistResp.Id.PlaylistId
		playlist.Description = playlistResp.Snippet.Description
		playlist.PublishedAt = publishedAt

		playlists = append(playlists, playlist)
	}

	return playlists
}

func (yt *YouTubeApi) GetAllPlaylistVideos(id string) []Video {
	videos := []Video{}

	nextPageToken := ""

	for {
		videosResp, err := yt.
			youtubeService.
			PlaylistItems.
			List([]string{"id", "snippet"}).
			PlaylistId(id).
			MaxResults(100).
			PageToken(nextPageToken).
			Do()

		if err != nil {
			log.Fatal(err)
		}

		for _, videoResp := range videosResp.Items {
			publishedAt, err := time.Parse(time.RFC3339, videoResp.Snippet.PublishedAt)
			if err != nil {
				log.Fatal("Could not parse field PublishedAt: " + err.Error())
			}
			videos = append(videos, Video{
				Id:          videoResp.Id,
				Title:       videoResp.Snippet.Title,
				Description: videoResp.Snippet.Description,
				PublishedAt: publishedAt,
				PlaylistId:  videoResp.Snippet.PlaylistId,
				Watched:     false,
			})
		}

		nextPageToken = videosResp.NextPageToken
		if nextPageToken == "" {
			break
		}
	}

	return videos
}

type Playlist struct {
	Id          string
	Title       string
	Description string
	PublishedAt time.Time
	Videos      []Video
	Updated     bool `json:"-"`
}

func (p *Playlist) String() string {
	return fmt.Sprintf("Playlist: { Id: %s, Title: %s, Description: %s, "+
		"Publish: %s }", p.Id, p.Title, p.Description, p.PublishedAt)
}

func (p *Playlist) Length() int {
	return len(p.Videos)
}

func (p *Playlist) Sort() {
	quicksort(0, p.Length()-1, p.Videos)

}

// quicksort sorts a given video slice
func quicksort(lowerBounds int, upperBounds int, videos []Video) {
	if lowerBounds >= upperBounds || lowerBounds < 0 {
		return
	}

	p := partition(lowerBounds, upperBounds, videos)

	quicksort(lowerBounds, p-1, videos)
	quicksort(p+1, upperBounds, videos)
}

// 1 3 2 4 pivot = 4, pivotI = 0
// -> 4 3 2 1 -> sort => pivot = 1 pivotI = 0
// -> 4<->1 because 4(the pivotI) > 1 (the pivot)

func partition(lowerBounds int, upperBounds int, videos []Video) int {
	pivotElement := videos[upperBounds]
	pivotPos := lowerBounds

	for i := lowerBounds; i < upperBounds; i++ {
		comp := videos[i].PublishedAt.Compare(pivotElement.PublishedAt)
		if comp <= 0 {
			swap(i, pivotPos, videos)
			pivotPos++
		}
	}

	// this is required to avoid an endless loop when
	// the right most item is the biggest item -> the pivot wont change and
	// nothing will move
	swap(pivotPos, upperBounds, videos)
	return pivotPos
}

func swap[T any](a int, b int, nums []T) {
	temp := nums[a]
	nums[a] = nums[b]
	nums[b] = temp
}

// FetchUpdate checks if the playlist has new videos
// it sets the flag Video.New to true
func (p *Playlist) FetchUpdate(yt *YouTubeApi) {
	videos := yt.GetAllPlaylistVideos(p.Id) // Get the newer playlist information

	// go through every video and check whether it is already in the list
	for idx, video := range videos {
		isNew := false
		for idx := range p.Videos {
                        knownVideo := &p.Videos[idx]
			isNew = isNew || video.Id == knownVideo.Id
                        if knownVideo.Id == video.Id {
                                knownVideo.Title = video.Title
                                knownVideo.Description = video.Description
                        }
		}
		if !isNew { // if it is not, append it
			p.Updated = true
			p.Videos = append(p.Videos, videos[idx])
		}
	}
}

type Video struct {
	Id          string
	Title       string
	Description string
	PublishedAt time.Time
	PlaylistId  string
	Watched     bool
}

func (v *Video) String() string {
	return fmt.Sprintf("Playlist: { Id: %s, Title: %s, Description: %s, "+
		"Publish: %s }", v.Id, v.Title, v.Description, v.PublishedAt)
}
