package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/baumple/watchvault/data"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	ALT_VERTICAL_BAR = "┃"
	VERTICAL_BAR     = "│"
	HORIZONTAL_BAR   = "─"
	VERT_CROSS_RIGHT = "├"
	VERT_CROSS_LEFT  = "┤"
	CENTER           = "┼"
	CORNER_TL        = "┌"
	CORNER_TR        = "┐"
	CORNER_BL        = "└"
	CORNER_BR        = "┘"
	ARROW_UD         = "⇳"
	ARROW_DWN        = "⇩"
	ARROW_UP         = "⇧"
)

func makeTobBar(width int) string {
	return makeLineBorders(strings.Repeat(HORIZONTAL_BAR, width-3), CORNER_TL, CORNER_TR, width)
}

func makeTobBarTitle(title string, width int) string {
	text := title +
		strings.Repeat(HORIZONTAL_BAR, width-len(title)-3) // fill whitespace with horizontal bars
	return makeLineBorders(text, CORNER_TL, CORNER_TR, width)
}

func makeSeparator(width int) string {
	return makeLineBorders(strings.Repeat(HORIZONTAL_BAR, width-3), VERT_CROSS_RIGHT, VERT_CROSS_LEFT, width)
}

func makeSeparatorTitle(title string, width int) string {
	text := title +
		strings.Repeat(HORIZONTAL_BAR, width-len(title)-3) // fill whitespace with horizontal bars

	return makeLine("", width) + // Create an empty line above so there is more padding
		makeLineBorders(text, VERT_CROSS_RIGHT, VERT_CROSS_LEFT, width)
}

func makeBottomBar(width int) string {
	return makeLineBorders(strings.Repeat(HORIZONTAL_BAR, width-3), CORNER_BL, CORNER_BR, width)
}

func makeLine(text string, width int) string {
	res := VERTICAL_BAR
	res += text

	for i := 0; i < width-len(text)-3; i++ {
		res += " "
	}

	res += VERTICAL_BAR
	return res + "\n"
}

func makeLineBorders(text string, borderLeft string, borderRight string, width int) string {
	res := borderLeft
	res += text

	for i := 0; i < width-len(text)-3; i++ {
		res += " "
	}

	res += borderRight
	return res + "\n"
}

func getDR() data.DataRetriever {
	return &data.JsonRetriever{}
}

// StartCLI starts the command line interface
func StartCLI() {
	fmt.Print("\033[2J")
	fmt.Print("\033[2;1H")

	yt := data.NewYouTubeApi()
	mainModel := initialModel()
	mainModel.dr = getDR()
	mainModel.yt = yt

	p := tea.NewProgram(mainModel)
	p.SetWindowTitle("watchvault")
	if _, err := p.Run(); err != nil {
		fmt.Printf("Could not start program: %v\n", err)
		os.Exit(1)
	}
}
