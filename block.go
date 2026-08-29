package main

import (
	"fmt"
	"math/rand"
)

const (
	BLOCK_DIRECTION_UP   = -1
	BLOCK_DIRECTION_DOWN = 1
)

// THE BLOCK
// hold a single Block data
// with position and form
type Block struct {
	x         int
	y         int
	Direction int
	Shape     [][]string
}

func EmptyBlock() Block {
	return Block{
		x:     0,
		y:     0,
		Shape: [][]string{},
	}
}

func (b *Block) IsEmpty() bool {
	return len(b.Shape) == 0
}

func (b *Block) Print() {

	for _, line := range b.Shape {
		for _, c := range line {
			fmt.Printf("%s ", c)
		}
		fmt.Println()
	}
}

func (b *Block) Width() int {
	if len(b.Shape) == 0 {
		return 0
	}
	return len(b.Shape[0])
}

func (b *Block) Height() int {
	return len(b.Shape)
}

func (b *Block) Transpose() {

	// rows turn columns
	cols := len(b.Shape)
	// columns to rows
	rows := len(b.Shape[0])

	// initialize new Block
	t := make([][]string, rows)
	for i := range t {
		t[i] = make([]string, cols)
	}

	// transpose
	for i, row := range b.Shape {
		for j, s := range row {
			// for proper visual effect
			// this is backwards
			t[len(row)-j-1][i] = s
		}
	}

	// update block with transposed one
	b.Shape = t
}

// Initializes and array of blocks
func InitBlocks() []Block {

	b := []Block{
		Block{
			x: 0,
			y: 0,
			Shape: [][]string{
				[]string{"*", "*", "*", " ", " "},
				[]string{" ", " ", "*", "*", "*"},
			},
		},
		Block{
			x: 0,
			y: 0,
			Shape: [][]string{
				[]string{"#", "#", "#"},
				[]string{"#", "#", "#"},
				[]string{"#", "#", "#"},
			},
		},
		Block{
			x: 0,
			y: 0,
			Shape: [][]string{
				[]string{"*", "*", "*", "*", "*", "*"},
			},
		},
		Block{
			x: 0,
			y: 0,
			Shape: [][]string{
				[]string{"*", "*", "*", "*", "*"},
				[]string{" ", " ", "*", " ", " "},
				[]string{" ", " ", "*", " ", " "},
			},
		},
	}

	return b
}

func PickRandomBlock() Block {
	blocks := InitBlocks()
	return blocks[rand.Intn(len(blocks))]
}

// Debug the Block Rendering and Transposition
func CheckBlocksRendering(availableBlocks []Block) {
	for _, b := range availableBlocks {
		b.Print()
		fmt.Printf("\n\n")
		for f := 0; f < 3; f++ {
			b.Transpose()
			b.Print()
			fmt.Printf("\n\n")
		}
	}
}
