package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

const (
	// ANSI Device Status Report sequence
	DSR_QUERY_POSITION = "\x1b[6n"
	DSR_CLEAR_SCREEN   = "\x1b[2J\x1b[H"
	DSR_GO_XY          = "\x1b[%d;%dH"
)

// THE TERMINAL
type Terminal struct {
	fd             int
	width          int
	height         int
	boardWidth     int
	boardHeight    int
	boardTimerRow  int
	boardStatusRow int
	boardStartRow  int
	boardStartCol  int
	clockPtr       int
}

func NewTerminal() (error, Terminal) {
	t := Terminal{}
	t.fd = int(os.Stdout.Fd())

	if !term.IsTerminal(t.fd) {
		return fmt.Errorf("must run on a terminal environment"), Terminal{}
	}

	w, h, err := term.GetSize(t.fd)
	if err != nil {
		return fmt.Errorf("failed to get terminal size: %v", err), Terminal{}
	}
	t.width = w
	t.height = h

	return nil, t
}

func (t *Terminal) NewLine() {
	fmt.Printf("\r\n")
}

func (t *Terminal) GetPos() (row, col int, err error) {
	// request cursor position
	_, err = os.Stdout.Write([]byte(DSR_QUERY_POSITION))
	if err != nil {
		return 0, 0, err
	}

	// read the response from Stdin in the format: \x1b[row;colR
	buf := make([]byte, 32)
	n, err := os.Stdin.Read(buf)
	if err != nil {
		return 0, 0, err
	}

	// parse the bytes
	resp := string(buf[:n])
	if !strings.HasPrefix(resp, "\x1b[") && !strings.HasSuffix(resp, "R") {
		return 0, 0, fmt.Errorf("terminal: unexpected response format: %s [%v]", resp, buf[:n])
	}

	trimmed := strings.TrimSuffix(resp, "R")
	trimmed = strings.TrimPrefix(trimmed, "\x1b[")
	parts := strings.Split(trimmed, ";")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("terminal: malformed cursor coordinates")
	}

	row, _ = strconv.Atoi(parts[0])
	col, _ = strconv.Atoi(parts[1])

	return row, col, nil
}

func (t *Terminal) ClearScreen() {
	fmt.Printf(DSR_CLEAR_SCREEN)
}

func (t *Terminal) ClearLine(row int) {
	t.PrintAt(row, 0, strings.Repeat(" ", t.width))
	t.GoBottom()
}

func (t *Terminal) PrintAt(row, col int, c string) {
	fmt.Fprintf(os.Stdout, "\x1b[%d;%dH", row, col)
	fmt.Print(c)
}

// Usually, the bottom left corner of the screen is a good
// place to rest the cursor
func (t *Terminal) GoBottom() {
	fmt.Fprintf(os.Stdout, "\x1b[%d;%dH", t.height, t.width)
}

// Turn off the cursor
func (t *Terminal) SetCursorOff() {
	fmt.Print("\x1b[?25l")
}

// Tur on the cursor
func (t *Terminal) SetCursorOn() {
	fmt.Print("\x1b[?25h")
}

// Draw the board on the screen and start the basic calculations
// to find the initial coordinates to render objects
func (t *Terminal) DrawBoard(width, height int) {

	// this is the timer row
	fmt.Printf("[%s]\r\n", strings.Repeat("  ", width))

	// first frame row +---...---+
	fmt.Printf("+%s+\r\n", strings.Repeat("--", width))
	for i := 1; i <= height; i++ {
		fmt.Printf("|%s|\r\n", strings.Repeat("  ", width))
	}
	// last frame row +---...---+
	fmt.Printf("+%s+\r\n", strings.Repeat("--", width))
	// the status bar
	fmt.Printf("...")

	// get the position of the status bar row
	row, _, _ := t.GetPos()

	t.boardStatusRow = row
	t.boardStartRow = t.boardStatusRow - height - 1
	t.boardTimerRow = t.boardStartRow - 2
	t.boardStartCol = 2
	t.boardHeight = height
	t.boardWidth = width

	t.Timer("")
	t.StatusBar("")
}

func (t *Terminal) Timer(timeString string) {
	chars := []string{"|", "/", "-", "\\"}
	content := fmt.Sprintf("[%d] %s %s", t.boardStatusRow, chars[t.clockPtr], timeString)
	t.PrintAt(t.boardTimerRow, 2, content)
	t.clockPtr = (t.clockPtr + 1) % len(chars)
}

func (t *Terminal) StatusBar(c string) {
	content := fmt.Sprintf("start row/col: %d, %d | timer row: %d | %s", t.boardStartRow, t.boardStartCol, t.boardTimerRow, c)
	t.PrintAt(t.boardStatusRow, 0, content)
}

// PrintAtBoard receives x and y coordinates RELATIVE TO THE BOARD
// and not the screen, so it does the proper calculation about where
// is the right place to put the string
func (t *Terminal) PrintAtBoard(x, y int, c string) {
	// adjust row and col to board coordinates
	boardRow := t.boardStartRow + y + 1
	boardCol := t.boardStartCol + (x * 2) + 1

	fmt.Fprintf(os.Stdout, "\x1b[%d;%dH", boardRow, boardCol)
	fmt.Print(c)
}

// Dump some info
func (t *Terminal) Info() {
	fmt.Printf("screen size: %dw x %dh", t.width, t.height)
}
