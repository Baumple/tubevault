package cli

import (
	"github.com/baumple/watchvault/data"
	tea "github.com/charmbracelet/bubbletea"
)

type videoModel struct {
	width  int
	height int
	video  *data.Video

	pageIndex int
}

func (v videoModel) Init() tea.Cmd {
	return nil
}

func (v videoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return nil, tea.Quit
		case "esc", "q":
			return nil, nil
		}
	}
	return v, nil
}

func (v videoModel) View() string {
        text := ""
        text += v.video.Description
	return text
}

func newVideoModel(video *data.Video, width int, height int) videoModel {
        return videoModel {
                width: width,
                height: height,
                video: video,
        }
}
