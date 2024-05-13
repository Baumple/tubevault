package cli

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/baumple/watchvault/data"
	"github.com/baumple/watchvault/utility"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	SECONDS_DAY = 86400
)

type playlistModel struct {
	width  int
	height int

	itemsPerPage int

	cursor   int
	playlist *data.Playlist

	visualMode  bool
	visualStart int

	dr data.DataRetriever

	currentModel tea.Model
}

func (p playlistModel) Init() tea.Cmd {
	return nil
}

func (p playlistModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if p.currentModel != nil {
		model, cmd := p.currentModel.Update(msg)
		p.currentModel = model
		return p, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.height = msg.Height
		p.width = msg.Width
		p.itemsPerPage = p.height / 3

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if p.playlist.Length() > 0 {
				p.currentModel = newVideoModel(&p.playlist.Videos[p.cursor], p.width, p.height)
			}

		case "q", "ctrl+c":
			return p, tea.Quit

		case "esc":
			if p.visualMode {
				p.visualMode = false
			} else {
				return nil, nil
			}

		case "down", "j":
			if p.cursor < len(p.playlist.Videos)-1 {
				p.cursor++
			}

		case "up", "k":
			if p.cursor > 0 {
				p.cursor--
			}
		case "v":
			p.visualMode = !p.visualMode

		case "ctrl+u":
			p.cursor = max(p.cursor-15, 0)

		case "ctrl+d":
			p.cursor = min(p.cursor+15, len(p.playlist.Videos)-1)

		case "G":
			p.cursor = p.playlist.Length() - 1
		case "g":
			p.cursor = 0

		case " ":
			if p.visualMode {
				p.visualMode = false
			}
			selection := p.getSelectionIndices()
			p.visualStart = p.cursor

			for i := selection.start; i < selection.end; i++ {
				video := &p.playlist.Videos[i]
				video.Watched = !video.Watched
			}

			return p, func() tea.Msg {
				for i := selection.start; i < selection.end; i++ {
					err := p.dr.UpdateVideoWatched(p.playlist.Id, p.playlist.Videos[i].Id, p.playlist.Videos[i].Watched)
					if err != nil {
						log.Fatal(err)
					}
				}
				return nil
			}
		}
	}

	if !p.visualMode {
		p.visualStart = p.cursor
	}
	return p, nil
}

func (p playlistModel) View() string {
	if p.currentModel != nil {
		return p.currentModel.View()
	}

	text := makeTobBarTitle("Playlist", p.width)
	text += makeLine(" "+p.playlist.Title, p.width)
	text += makeSeparatorTitle("Published at", p.width)
	text += makeLine(" "+p.playlist.PublishedAt.String(), p.width)
	text += makeSeparatorTitle("Description", p.width)

	if len(p.playlist.Description) <= 0 {
		text += makeLine(" ...", p.width)
	}

	lines := utility.SplitEveryN(p.playlist.Description, p.width)

	for _, line := range lines {
		text += makeLine(" "+line, p.width)
	}

	text += makeSeparator(p.width)

	nVideos := len(p.playlist.Videos)
	selection := p.getSelectionIndices()

	currentTime := time.Now().Unix()
	windowIndices := GetWindow(p.cursor, p.itemsPerPage)

	// ratio between current window and total elements
	startProgress := float32(windowIndices.Start) / float32(p.playlist.Length())

	// take that ratio and multiply it with height of item list
	// -> relative start index
	barStart := windowIndices.Start + int(float32(p.height/2)*startProgress)

	// same for progressbar
	endProgress := float32(windowIndices.End) / float32(p.playlist.Length())
	barEnd := windowIndices.Start + int(float32(p.height/2)*endProgress)

	for i := windowIndices.Start; i < windowIndices.End; i++ {
		if i < 0 {
			continue
		}

		leftBar := VERTICAL_BAR

		if i >= barStart && i <= barEnd {
			leftBar = "\033[34m" + ALT_VERTICAL_BAR + "\033[0m"
		}

		if i >= nVideos {
			text += makeLineBorders("", leftBar, VERTICAL_BAR, p.width)
			continue
		}

		video := &p.playlist.Videos[i]

		cursor := " "
		if p.cursor == i {
			cursor = ">"
		}

		watched := "[ ]"
		if video.Watched {
			watched = "[X]"
		}

		modifier := ""
		if p.visualMode && i >= selection.start && i < selection.end {
			modifier = "\033[;5m"
		}

		// if the video is newer than three days mark it as "NEW"
		newText := "     "
		if currentTime-video.PublishedAt.Unix() < SECONDS_DAY*3 {
			newText = ">NEW<"
		}

                text += leftBar
		text += fmt.Sprintf(
			"%s %s %s %s %s\033[0m",
			modifier,
			cursor,
			watched,
			newText,
			video.Title,
		)

		whiteSpaceLeft := p.width -
			len(video.Title) -
			1 - // cursor
			3 - // watched
			2 // borders

		text += strings.Repeat(" ", whiteSpaceLeft)
		text += VERTICAL_BAR
		text += "\n"

	}

	text += makeBottomBar(p.width)

	text += makeTobBarTitle("Keymaps", p.width)
	text += makeLine("  * <esc>  -> return", p.width)
	text += makeLine("  * <space> -> toggle watched", p.width)
	text += makeLine("  * <v>     -> visual mode", p.width)
	text += makeBottomBar(p.width)

	return text
}

func NewPlaylistModel(dr data.DataRetriever, width int, height int, playlist *data.Playlist) playlistModel {
	playlist.Sort()

	return playlistModel{
		width:        width,
		height:       height,
		itemsPerPage: height / 2,
		cursor:       0,
		playlist:     playlist,
		visualMode:   false,
		visualStart:  0,
		dr:           dr,
	}
}

// TODO: REmove duplicate im too lazy rn
func (p *playlistModel) getSelectionIndices() struct {
	start int
	end   int
} {
	selectionStart := p.visualStart
	selectionEnd := p.cursor

	if selectionStart > selectionEnd {
		temp := selectionEnd
		selectionEnd = selectionStart
		selectionStart = temp
	}

	return struct {
		start int
		end   int
	}{
		selectionStart,
		selectionEnd + 1,
	}
}

func GetWindow(cursor int, width int) struct {
	Start int
	End   int
} {
	start := 0
	end := 0
	if cursor >= width/2 {
		start = cursor - width/2
		end = cursor + width/2
	} else {
		end += width - 1
	}

	return struct {
		Start int
		End   int
	}{start, end}
}
