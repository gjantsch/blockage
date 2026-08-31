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
	ANSI_CLEAR_SCREEN  = "\x1b[2J\x1b[H"
	ANSI_GOTO_XY       = "\x1b[%d;%dH"
	ANSI_CURSOR_OFF    = "\x1b[?25l"
	ANSI_CURSOR_ON     = "\x1b[?25h"
)

// THE TERMINAL
type Terminal struct {
	previousState *term.State
	fd            int
	width         int
	height        int
	boardWidth    int
	boardHeight   int
	boardStartRow int
	boardStartCol int
	timerRow      int
	statusRow     int
	clock         Clock
	frame         Frame
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

	// we must go raw mode
	previousState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to set terminal to raw mode: %v", err), Terminal{}
	}
	t.previousState = previousState

	t.SetCursorOff()

	t.clock = *NewClock()
	t.frame = NewFrame(t.width, t.height)

	return nil, t
}

func (t *Terminal) Close() {
	t.SetCursorOn()
	term.Restore(t.fd, t.previousState)
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

func (t *Terminal) PrintAt(row, col int, c string) {
	fmt.Fprintf(os.Stdout, "\x1b[%d;%dH", row, col)
	fmt.Print(c)
}

func (t *Terminal) ClearScreen() {
	fmt.Printf(ANSI_CLEAR_SCREEN)
}

// Turn off the cursor
func (t *Terminal) SetCursorOff() {
	fmt.Print(ANSI_CURSOR_OFF)
}

// Tur on the cursor
func (t *Terminal) SetCursorOn() {
	fmt.Print(ANSI_CURSOR_ON)
}

// Draw the board on the screen and start the basic calculations
// to find the initial coordinates to render objects
func (t *Terminal) DrawBoard(width, height int) error {
	t.frame.Draw()
	err := t.computeBoardLayout(width, height)
	if err != nil {
		return fmt.Errorf("failed to compute board layout: %v", err)
	}
	t.Timer("")
	t.StatusBar("")
	return nil
}

func (t *Terminal) computeBoardLayout(width int, height int) error {
	// get the position of the status bar row
	row, _, err := t.GetPos()
	if err != nil {
		return fmt.Errorf("failed to get terminal position: %v", err)
	}

	t.statusRow = row
	t.boardStartRow = t.statusRow - height - 1
	t.timerRow = t.boardStartRow - 2
	t.boardStartCol = 2
	t.boardHeight = height
	t.boardWidth = width

	return nil
}

func (t *Terminal) Timer(timeString string) {
	t.PrintAt(t.timerRow, 2, fmt.Sprintf("%s %s", t.clock.Next(), timeString))
}

func (t *Terminal) StatusBar(content string) {
	t.PrintAt(t.statusRow, 0, content)
}

// PrintAtBoard receives x and y coordinates RELATIVE TO THE BOARD
// and not the screen, so it does the proper calculation about where
// is the right place to put the string
func (t *Terminal) PrintAtBoard(x, y int, c string) {
	// adjust row and col to board coordinates
	boardRow := t.boardStartRow + y
	boardCol := t.boardStartCol + (x * 2) + 1

	fmt.Fprintf(os.Stdout, ANSI_GOTO_XY, boardRow, boardCol)
	fmt.Print(c)
}
