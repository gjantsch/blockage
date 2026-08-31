package frame

import (
	"fmt"
	"strings"
)

func DrawFrame(width, height int) string {

	var frame string

	// this is the timer row
	frame = fmt.Sprintf("[%s]\r\n", strings.Repeat("  ", width))

	// first frame row +---...---+
	frame += fmt.Sprintf("+%s+\r\n", strings.Repeat("--", width))
	for i := 1; i <= height; i++ {
		frame += fmt.Sprintf("|%s|\r\n", strings.Repeat("  ", width))
	}
	// last frame row +---...---+
	frame += fmt.Sprintf("+%s+\r\n", strings.Repeat("--", width))
	// the status bar
	frame += "..."

	return frame
}
