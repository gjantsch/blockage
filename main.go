package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	// board dimensions
	BOARD_WIDTH  = 8
	BOARD_HEIGHT = 32

	// block components
	BL_NIL = " "
	BL_FIL = "*"

	// used key codes
	KEY_UP     = "\x1b[A"
	KEY_DOWN   = "\x1b[B"
	KEY_RIGHT  = "\x1b[C"
	KEY_LEFT   = "\x1b[D"
	KEY_CTRL_C = "\x03"
	KEY_ESC    = "\x1b"

	// directions
	DIR_UP    = -1
	DIR_DOWN  = 1
	DIR_LEFT  = -1
	DIR_RIGHT = 1
	KEY_Q     = "q"
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

func gameOver(terminal Terminal) {
	terminal.Print("\r\n!!!!!!!!!!!!!!!!!!!")
	terminal.Print("\r\n!!!              !!")
	terminal.Print("\r\n!!!  GAME OVER   !!")
	terminal.Print("\r\n!!!              !!")
	terminal.Print("\r\n!!!!!!!!!!!!!!!!!!!")
	terminal.Print("\r\n")
}

func main() {

	err, terminal := NewTerminal()
	defer terminal.Close()
	if err != nil {
		log.Fatalf("error initializing terminal: %v", err)
		return
	}

	board := NewMatrix(BOARD_WIDTH, BOARD_HEIGHT, &terminal)
	err = board.Render()
	if err != nil {
		log.Fatalf("error rendering board: %v", err)
		return
	}
	board.PickRandomBlock()
	board.PlaceBlockAtBottom()
	board.BlockDirection = DIR_UP

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
				board.MoveBlockX(DIR_LEFT)
			case KEY_RIGHT:
				board.MoveBlockX(DIR_RIGHT)
			case KEY_UP:
				for !board.BlockHitTop() && !board.Collided {
					board.MoveBlock()
				}
			case "r", "R":
				board.RotateBlock()
			case "s", "S":
				board.ChangeBlockShape()
			case "e", "E":
				board.BlockShapeToDot()
			case "q", "Q":
				fmt.Printf("\r\nexit key pressed, exiting\r\n")
				return
			case KEY_CTRL_C:
				fmt.Printf(" \r\nctrl+c pressed, exiting\r\n")
				return
			case KEY_ESC:
				fmt.Printf("\r\nescape key pressed, exiting\r\n")
				return
			}

		case <-ticker.C:
			// update the timer
			terminal.Timer(time.Since(startTime).Truncate(time.Second).String())

			if !board.HasBlock() {
				break
			}

			board.MoveBlock()
			if board.BlockHitTop() || board.Collided {
				if board.MovesCount == 1 && board.Collided {
					gameOver(terminal)
					return
				}
				board.CheckForFullRows()
				board.PickRandomBlock()
				board.PlaceBlockAtBottom()
			}

		}
	}

}
