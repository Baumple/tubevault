package data

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const PLAYLIST_DIR = "playlists"

type JsonRetriever struct {
	saveDir     string
	playlistDir string
}

// getSaveDirPath returns a path to the save dir
func (jr *JsonRetriever) getSaveDirPath() (string, error) {
	if jr.saveDir != "" {
		return jr.saveDir, nil
	} else {
		return GetSaveDirPath()
	}
}

// getPlaylistDir returns a path to where playlists are stored.
// Usually playlists are stored in HOME_DIR/.tubevault/playlists
// But also "cashes" the value
func (jr *JsonRetriever) getPlaylistDir() (string, error) {
	if jr.playlistDir != "" {
		return jr.playlistDir, nil
	} else {
		return GetPlaylistDir()
	}

}

// GetPlaylists reads, parses and returns a slice of stored playlists
func (jr *JsonRetriever) GetPlaylists() ([]Playlist, error) {
	playlistDir, err := jr.getPlaylistDir()

	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(playlistDir)
	if err != nil {
		return nil, err
	}

	playlists := []Playlist{}
	for _, entry := range entries {
		playlistPath := filepath.Join(playlistDir, entry.Name())
		file, err := os.Open(playlistPath)
		if err != nil {
			return nil, err
		}

		playlist := Playlist{}
		err = json.NewDecoder(file).Decode(&playlist)
		if err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

// DeletePlaylist deletes the playlist with the given id
func (jr *JsonRetriever) DeletePlaylist(id string) error {
	playlistDir, err := jr.getPlaylistDir()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(playlistDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.Name() == id+".json" {
			err := os.Remove(filepath.Join(playlistDir, entry.Name()))
			if err != nil {
				return err
			}

		}
	}

	return err
}

// SavePlaylist writes creates a playlistDir and stores the playlist as json
func (jr *JsonRetriever) SavePlaylist(playlist *Playlist) error {
	saveDir, err := jr.getPlaylistDir()
	if err != nil {
		return err
	}

	playlistPath := filepath.Join(saveDir, playlist.Id+".json")
	playlistFile, err := os.OpenFile(playlistPath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	err = json.NewEncoder(playlistFile).Encode(playlist)

	return err
}

// UpdateVideoWatched implements DataRetriever.
func (jr *JsonRetriever) UpdateVideoWatched(playlistId string, videoId string, watched bool) error {
	playlistDir, err := jr.getPlaylistDir()
	if err != nil {
		return err
	}

	playlistPath := filepath.Join(playlistDir, playlistId+".json")
	playlistFile, err := os.OpenFile(playlistPath, os.O_RDWR|os.O_CREATE, 0777)
	defer playlistFile.Close()

	if err != nil {
		return err
	}

	playlist := Playlist{}
	err = json.NewDecoder(playlistFile).Decode(&playlist)
	if err != nil {
		return err
	}

	for idx, video := range playlist.Videos {
		if video.Id == videoId {
			playlist.Videos[idx].Watched = watched
		}
	}

	err = jr.SavePlaylist(&playlist)

	return err
}

func (jr *JsonRetriever) Close() {}
