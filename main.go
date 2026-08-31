package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gjantsch/fun02/internal/game"
	"github.com/gjantsch/fun02/internal/render"
)

const (
	// board dimensions
	BOARD_WIDTH  = 8
	BOARD_HEIGHT = 32

	// used key codes
	KEY_UP     = "\x1b[A"
	KEY_DOWN   = "\x1b[B"
	KEY_RIGHT  = "\x1b[C"
	KEY_LEFT   = "\x1b[D"
	KEY_CTRL_C = "\x03"
	KEY_ESC    = "\x1b"
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
		log.Fatalf("error initializing terminal: %v", err)
		return
	}
	defer terminal.Close()

	// The matrix render must run before readKeys starts,
	// or GetPos and readKeys will race on stdin
	board := game.NewMatrix(BOARD_WIDTH, BOARD_HEIGHT, &terminal)
	err = board.Render()
	if err != nil {
		log.Fatalf("error rendering board: %v", err)
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
				log.Println("key channel closed, exiting")
				return
			}

			switch key {
			case KEY_LEFT:
				board.MoveBlockX(game.Left)
			case KEY_RIGHT:
				board.MoveBlockX(game.Right)
			case KEY_UP:
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
			case KEY_CTRL_C:
				terminal.Print(" \r\nctrl+c pressed, exiting\r\n")
				return
			case KEY_ESC:
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
