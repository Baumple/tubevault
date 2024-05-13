package cli

import (
	"fmt"

	"github.com/baumple/watchvault/data"
	tea "github.com/charmbracelet/bubbletea"
)

type msgSearchedPlaylists struct {
	playlists []data.Playlist
}

type msgSearchedResult struct {
	playlist data.Playlist
}

type searchModel struct {
	foundPlaylists []data.Playlist
	cursor         int

	searchFocused bool
	text         string

	yt *data.YouTubeApi
	dr data.DataRetriever
}

func (s searchModel) Init() tea.Cmd {
	return nil
}

func (s searchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+p", "up":
			if s.cursor > 0 {
				s.cursor--
			}
		case "ctrl+n", "down":
			if s.cursor < len(s.foundPlaylists)-1 {
				s.cursor++
			}
		case "backspace":
			if len(s.text) > 0 {
				s.text = s.text[:len(s.text)-1]
			}
		case "esc":
                        return nil, nil

		case "enter":
                        return s, func() tea.Msg {
                                playlists := s.yt.GetYoutubePlaylistsBySearch(s.text)
                                return msgSearchedPlaylists{playlists}
                        }
                case "tab":
                        return nil, func() tea.Msg {
                                if len(s.foundPlaylists) <= 0 {
                                        return nil
                                }
                                selectedPlaylist := s.foundPlaylists[s.cursor]
                                videos := s.yt.GetAllPlaylistVideos(selectedPlaylist.Id)
                                selectedPlaylist.Videos = videos
                                return msgSearchedResult{selectedPlaylist}
                        }
		default:
			msg := msg.String()
			if len(msg) == 1 {
				s.text += msg
			}
		}
	case msgSearchedPlaylists:
		playlists := msg.playlists
		s.foundPlaylists = playlists
	}
	return s, nil
}

func (s searchModel) View() string {
	text := "Search playlist\n\n"
	text += "Enter a name: " + s.text + CURSOR + "\n"
	for idx, playlist := range s.foundPlaylists {
		cursor := " "
		if idx == s.cursor {
			cursor = ">"
		}
		text += fmt.Sprintf(" %s Title: %s Description: %s->\n", cursor, playlist.Title, playlist.Description)
	}

        text += "\n\nKeymaps:\n"
        text += "  * <enter> -> Search keyword\n"
        text += "  * <tab>   -> Add playlist at cursor\n"

	return text
}
