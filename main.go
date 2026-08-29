package main

import (
	"log"
	"os"
	"slices"
	"time"

	"golang.org/x/term"
)

const (
	BOARD_WIDTH  = 20
	BOARD_HEIGHT = 40
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
	block := PickRandomBlock()

	// read keys interferes in the getPos() function...
	keys := make(chan byte, 1)
	go readKeys(keys)

	startTime := time.Now()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case key, ok := <-keys:
			if !ok {
				log.Println("key channel closed, exiting")
				return
			}

			switch key {
			case 'a':
				terminal.PrintAt(1, 3, "left ")
			case 'd':
				terminal.PrintAt(1, 3, "right")
			case 'q':
			case 'Q':
				log.Println("exit key pressed, exiting")
				return
			case '\x03': // Ctrl+C
				log.Println("ctrl+c pressed, exiting")
				return
			case '\x1b': // ESC
				log.Println("escape key pressed, exiting")
				return
			}

		case <-ticker.C:
			// This continues running even when no key is pressed.
			terminal.Timer(time.Since(startTime).Truncate(time.Second).String())
			terminal.StatusBar("updated updated updated updated updated updated updated updated updated updated updated updated updated")
			if optionRotateOnly {
				board.PlaceBlockAtCenter(block)
				board.RemoveBlockAtCenter(block)
				block.Transpose()
				board.PlaceBlockAtCenter(block)
				continue
			}

			if board.HasBlock() {
				board.MoveBlock()
			} else {
				board.PlaceBlockAtBottom(block)
				board.BlockDirection = BLOCK_DIRECTION_UP
			}
		}
	}

}
