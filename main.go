package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"time"

	"golang.org/x/term"
)

const (
	BOARD_WIDTH  = 20
	BOARD_HEIGHT = 40
	BLOCK_EMPTY  = "."
	BLOCK_FILLED = "*"
)

func readKeys(keys chan<- byte) {
	buffer := make([]byte, 1)

	for {
		_, err := os.Stdin.Read(buffer)
		if err != nil {
			close(keys)
			return
		}

		keys <- buffer[0]
	}
}

func main() {

	optionRender := len(os.Args) > 1 && slices.Contains(os.Args[1:], "-r")
	optionRotateOnly := len(os.Args) > 1 && slices.Contains(os.Args[1:], "-o")

	err, terminal := NewTerminal()
	if err != nil {
		log.Fatalf("error initializing terminal: %v", err)
		return
	}

	board := NewMatrix(BOARD_WIDTH, BOARD_HEIGHT, &terminal)

	if optionRender {
		CheckBlocksRendering(InitBlocks())
		board.Render()
		terminal.Info()
		return
	}

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
	board.BlockDirection = BLOCK_DIRECTION_UP

	// read keys interferes in the getPos() function...
	keys := make(chan byte, 1)
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
			case 'r':
				board.RotateBlock()
			case 's':
				board.ChangeBlockShape()
			case 'q':
			case 'Q':
				fmt.Printf("\r\nexit key pressed, exiting\r\n")
				return
			case '\x03': // Ctrl+C
				fmt.Printf(" \r\nctrl+c pressed, exiting\r\n")
				return
			case '\x1b': // ESC
				fmt.Printf("\r\nescape key pressed, exiting\r\n")
				return
			}

		case <-ticker.C:
			// This continues running even when no key is pressed.
			terminal.Timer(time.Since(startTime).Truncate(time.Second).String())
			if optionRotateOnly {
				board.PlaceBlockAtCenter()
				board.RemoveBlockAtCenter()
				board.RotateBlock()
				board.PlaceBlockAtCenter()
				continue
			}

			if board.HasBlock() {
				board.MoveBlock()
			}
		}
	}

}
