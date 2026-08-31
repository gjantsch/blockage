package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gjantsch/fun02/internal/game"
	"github.com/gjantsch/fun02/internal/render"
)

const (
	// board dimensions
	BoardWidth  = 8
	BoardHeight = 32

	// used key codes
	KeyUp    = "\x1b[A"
	KeyDown  = "\x1b[B"
	KeyRight = "\x1b[C"
	KeyLeft  = "\x1b[D"
	KeyCtrlC = "\x03"
	KeyEsc   = "\x1b"
)

// Read keys from stdin
func readKeys(keys chan<- string) {
	buffer := make([]byte, 3)

	for {
		n, err := os.Stdin.Read(buffer)
		if err != nil {
			close(keys)
			return
		}

		keys <- string(buffer[:n])
	}
}

// render game over message
func gameOver(terminal *render.Terminal) {
	terminal.Print("\r\n!!!!!!!!!!!!!!!!!!!")
	terminal.Print("\r\n!!!              !!")
	terminal.Print("\r\n!!!  GAME OVER   !!")
	terminal.Print("\r\n!!!              !!")
	terminal.Print("\r\n!!!!!!!!!!!!!!!!!!!")
	terminal.Print("\r\n")
}

func main() {

	terminal, err := render.NewTerminal()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initializing terminal: %v\n", err)
		return
	}
	defer terminal.Close()

	// The matrix render must run before readKeys starts,
	// or GetPos and readKeys will race on stdin
	board := game.NewMatrix(BoardWidth, BoardHeight, &terminal)
	err = board.Render()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error rendering board: %v\n", err)
		return
	}
	board.PickRandomBlock()
	board.PlaceBlockAtBottom()

	keys := make(chan string, 1)
	go readKeys(keys)

	startTime := time.Now()

	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case key, ok := <-keys:
			if !ok {
				fmt.Println("key channel closed, exiting")
				return
			}

			switch key {
			case KeyLeft:
				board.MoveBlockX(game.Left)
			case KeyRight:
				board.MoveBlockX(game.Right)
			case KeyUp:
				for board.MoveBlock() == game.Moved {
					// no-op
				}
			case "r", "R":
				board.RotateBlock()
			case "s", "S":
				board.ChangeBlockShape()
			case "e", "E":
				board.BlockShapeToDot()
			case "q", "Q":
				terminal.Print("\r\nexit key pressed, exiting\r\n")
				return
			case KeyCtrlC:
				terminal.Print(" \r\nctrl+c pressed, exiting\r\n")
				return
			case KeyEsc:
				terminal.Print("\r\nescape key pressed, exiting\r\n")
				return
			}

		case <-ticker.C:
			// update the timer
			terminal.Timer(time.Since(startTime).Truncate(time.Second).String())

			if !board.HasBlock() {
				break
			}

			outcome := board.MoveBlock()
			switch outcome {
			case game.GameOver:
				gameOver(&terminal)
				return

			case game.Landed:
				// check for full rows
				rowsCompleted := board.CheckForFullRows()
				terminal.StatusBar(fmt.Sprintf("%d rows completed", rowsCompleted))

				// place a new block at the bottom
				board.PickRandomBlock()
				board.PlaceBlockAtBottom()
			}

		}
	}

}
