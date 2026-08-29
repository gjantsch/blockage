package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/term"
)

const (
	// board dimensions
	BOARD_WIDTH  = 9
	BOARD_HEIGHT = 20

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

func main() {

	err, terminal := NewTerminal()
	if err != nil {
		log.Fatalf("error initializing terminal: %v", err)
		return
	}

	board := NewMatrix(BOARD_WIDTH, BOARD_HEIGHT, &terminal)

	// the cursor is noisy, so turn off
	terminal.SetCursorOff()
	defer terminal.SetCursorOn()

	// we must go raw mode
	previousState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), previousState)

	board.Render()
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
			case "r":
				board.RotateBlock()
			case "s":
				board.ChangeBlockShape()
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
					fmt.Printf("\r\n!!!!!!!!!!!!!!!!!!!")
					fmt.Printf("\r\n!!!              !!")
					fmt.Printf("\r\n!!!  GAME OVER   !!")
					fmt.Printf("\r\n!!!              !!")
					fmt.Printf("\r\n!!!!!!!!!!!!!!!!!!!")
					fmt.Printf("\r\n")
					return
				}
				board.CheckForFullRows()
				board.PickRandomBlock()
				board.PlaceBlockAtBottom()
			}

		}
	}

}
