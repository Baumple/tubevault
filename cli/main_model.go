package cli

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/baumple/watchvault/data"
	tea "github.com/charmbracelet/bubbletea"
)

const CURSOR = "█"

type msgListInsert struct {
	playlist data.Playlist
}

type msgListUpdated struct {
	playlists []data.Playlist
}

type mainModel struct {
	width  int
	height int

	currentModel tea.Model

	cursor           int
	trackedPlaylists []data.Playlist

	dr data.DataRetriever
	yt data.YouTubeApi
}

func initialModel() mainModel {
	return mainModel{}
}

func (s mainModel) Init() tea.Cmd {
	return func() tea.Msg {
		playlists, err := s.dr.GetPlaylists()
		for idx := range playlists {
                        playlists[idx].FetchUpdate(&s.yt)
                        if playlists[idx].Updated {
                                s.dr.SavePlaylist(&playlists[idx])
                        }
		}
		if err != nil {
			log.Fatal(err)
		}

		return msgListUpdated{playlists}
	}
}

func (s mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if s.currentModel != nil {
		model, cmd := s.currentModel.Update(msg)
		s.currentModel = model
		return s, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height

	case tea.KeyMsg:
		return s.HandleInput(msg.String())

	case msgListUpdated:
		s.trackedPlaylists = msg.playlists
		s.cursor = 0

	case msgSearchedResult:
		playlist := msg.playlist

		s.addPlaylist(playlist)

		s.currentModel = nil
	}
	return s, nil
}

func (s *mainModel) HandleInput(key string) (tea.Model, tea.Cmd) {
	switch key {

	case "up", "k":
		if s.cursor > 0 {
			s.cursor--
		}
	case "down", "j":
		if s.cursor < len(s.trackedPlaylists)-1 {
			s.cursor++
		}

	case tea.KeyF5.String():
		if len(s.trackedPlaylists) <= 0 {
			break
		}
		return s, func() tea.Msg {
			err := s.dr.DeletePlaylist(s.trackedPlaylists[s.cursor].Id)
			if err != nil {
				log.Fatal(err)
			}
			playlists, err := s.dr.GetPlaylists()
			if err != nil {
				log.Fatal(err)
			}

			return msgListUpdated{playlists}
		}

	case "s":
		searchModel := searchModel{
			foundPlaylists: []data.Playlist{},
			cursor:         0,
			text:           "",
			yt:             &s.yt,
			dr:             s.dr,
			searchFocused:  true,
		}
		s.currentModel = searchModel

	case "enter":
		if len(s.trackedPlaylists) <= 0 {
			break
		}
		return s, func() tea.Msg {
			_, err := exec.Command("firefox",
				"https://youtube.com/playlist?list="+s.trackedPlaylists[s.cursor].
					Id).
				Output()
			if err != nil {
				return err
			}
			return nil
		}

	case " ":
		if len(s.trackedPlaylists) <= 0 {
			break
		}
		playlistModel := NewPlaylistModel(s.dr, s.width, s.height, &s.trackedPlaylists[s.cursor])
		s.currentModel = playlistModel
		return s, nil

	case "q", "esc":
		return s, tea.Quit
	}
	return s, nil
}

func (s mainModel) View() string {
	// the width being 0 means we have not yet received
	// the size of the terminal
	if s.width == 0 {
		return ""
	}

	if s.currentModel != nil {
		return s.currentModel.View()
	}
	text := "Tracked playlists:\n\n"

	maxLenTitle := 0
	for _, playlist := range s.trackedPlaylists {
		maxLenTitle = max(len(playlist.Title), maxLenTitle)
	}

	text += "Name:" + strings.Repeat(" ", maxLenTitle) + "Description:\n"
	for i := 0; i < s.width; i++ {
		borderChar := "─"
		if i == maxLenTitle+7 {
			borderChar = "┬"
		}
		text += borderChar
	}
	text += "\n"

	for i, playlist := range s.trackedPlaylists {
		cursor := " "
		if s.cursor == i {
			cursor = ">"
		}

                updatedText := " "
                if playlist.Updated {
                        updatedText = "*"
                }
		descriptionPadding := maxLenTitle - len(playlist.Title) + 1
		text += fmt.Sprintf(
			"%s %s %s %s │ %s",
			cursor,
			playlist.Title,
			strings.Repeat(" ", descriptionPadding),
                        updatedText,
			playlist.Description,
		) + "\n"
	}

	for i := 0; i < s.width; i++ {
		if i == maxLenTitle+7 {
			text += "┴"
		} else {
			text += "─"
		}
	}
	text += "\n"

	text += makeTobBarTitle("Keymaps", s.width)
	text += makeLine(" * <q>     -> quit", s.width)
	text += makeLine(" * <F5>    -> remove playlist", s.width)
	text += makeLine(" * <s>     -> search playlist", s.width)
	text += makeLine(" * <space> -> view playlist", s.width)
	text += makeBottomBar(s.width)

	return text
}

// isTracked returns whether the given playlist (id) is already in the list
func (s *mainModel) isTracked(playlistId string) bool {
	for _, tracked := range s.trackedPlaylists {
		if tracked.Id == playlistId {
			return true
		}
	}
	return false
}

// addPlaylist checks if the given playlist is not in the list
// and then adds it. It will also update cursor position.
func (s *mainModel) addPlaylist(playlist data.Playlist) {
	if !s.isTracked(playlist.Id) {
		err := s.dr.SavePlaylist(&playlist)
		if err != nil {
			log.Fatal(err)
		}
		s.trackedPlaylists = append(s.trackedPlaylists, playlist)
		s.cursor = len(s.trackedPlaylists) - 1
	}
}
