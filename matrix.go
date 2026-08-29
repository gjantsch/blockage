package main

import (
	"fmt"
	"math/rand"
)

// THE MATRIX
// the matrix is the main board where the blocks falls
// and the game happens
type Matrix struct {
	width           int
	height          int
	content         [][]string
	terminal        *Terminal
	block           Block
	blockX          int
	blockY          int
	BlockDirection  int
	blocks          []Block
	currentBlockPtr int
}

func NewMatrix(width int, height int, term *Terminal) Matrix {
	m := Matrix{width: width, height: height, terminal: term}
	m.content = make([][]string, height)
	for i := range m.content {
		m.content[i] = make([]string, width)
		for j := range m.content[i] {
			m.content[i][j] = BLOCK_EMPTY
		}
	}
	m.blocks = InitBlocks()

	return m
}

func (m *Matrix) PickRandomBlock() {
	m.currentBlockPtr = rand.Intn(len(m.blocks))
	b := m.blocks[m.currentBlockPtr]
	m.block = b.Clone()
}

// Rotate the block in The Matrix
// the block nows how to rotate (transpose)
// but on The Matrix it must be able to:
//   - proper erase it self
//   - transpose
//   - redraw it self keeping the center
//     to have smooth rotation look and feel
func (m *Matrix) RotateBlock() {

	oldWidth := m.block.Width()
	oldHeight := m.block.Height()

	m.RemoveBlock()
	m.block.Transpose()

	newWidth := m.block.Width()
	newHeight := m.block.Height()

	// adjust block position to keep the center
	m.blockX += (oldWidth - newWidth) / 2
	m.blockY += (oldHeight - newHeight) / 2

	m.PutBlock()
}

func (m *Matrix) ChangeBlockShape() {
	m.RemoveBlock()
	m.currentBlockPtr = (m.currentBlockPtr + 1) % len(m.blocks)
	b := m.blocks[m.currentBlockPtr]

	b.x = m.blockX
	b.y = m.blockY

	m.block = b.Clone()
	m.PutBlock()
}

func (m *Matrix) HasBlock() bool {
	return !m.block.IsEmpty()
}

func (m *Matrix) BlockHitTop() bool {
	return m.blockY == 0
}

func (m *Matrix) Render() {
	m.terminal.DrawBoard(m.width, m.height)
}

func (m *Matrix) GetBlockCenter() (x, y int) {
	// calculate the center position
	centerX := (m.width - m.block.Width()) / 2
	centerY := (m.height - m.block.Height()) / 2
	return centerX, centerY
}

func (m *Matrix) RemoveBlockAtCenter() {
	centerX, centerY := m.GetBlockCenter()

	// remove the block at the center position
	for i, row := range m.block.Shape {
		for j := range row {
			m.UpdateContent(centerX+j, centerY+i, BLOCK_EMPTY)
		}
	}
}

func (m *Matrix) RemoveBlock() {
	for i, row := range m.block.Shape {
		for j := range row {
			m.UpdateContent(m.blockX+j, m.blockY+i, BLOCK_EMPTY)
		}
	}
}

func (m *Matrix) PutBlock() {
	for i, row := range m.block.Shape {
		for j := range row {
			m.UpdateContent(m.blockX+j, m.blockY+i, m.block.Shape[i][j])
		}
	}
}

func (m *Matrix) MoveBlockX(direction int) {
	m.RemoveBlock()
	if direction == BLOCK_DIRECTION_LEFT && m.blockX > 0 {
		m.blockX--
	}

	if direction == BLOCK_DIRECTION_RIGHT && m.blockX < m.width-m.block.Width() {
		m.blockX++
	}

	m.PutBlock()
}

func (m *Matrix) MoveBlockY(direction int) {
	m.RemoveBlock()
	if direction == BLOCK_DIRECTION_UP && m.blockY > 0 {
		m.blockY--
	}

	if direction == BLOCK_DIRECTION_DOWN && m.blockY < m.height-m.block.Height() {
		m.blockY++
	}

	m.PutBlock()
}

func (m *Matrix) MoveBlock() {
	m.MoveBlockY(m.BlockDirection)
	m.terminal.StatusBar(fmt.Sprintf("%d %d", m.blockX, m.blockY))
}

func (m *Matrix) PlaceBlockAtBottom() {
	bottomX := ((m.width - m.block.Width()) / 2)
	bottomY := m.height - m.block.Height()
	m.PutBlockAt(bottomX, bottomY)
}

func (m *Matrix) PlaceBlockAtTop() {
	topX := (m.width - m.block.Width()) / 2
	topY := 0
	m.PutBlockAt(topX, topY)
}

func (m *Matrix) PlaceBlockAtCenter() {
	centerX, centerY := m.GetBlockCenter()
	m.PutBlockAt(centerX, centerY)
}

func (m *Matrix) PutBlockAt(x, y int) {
	// place the block at position
	for i, row := range m.block.Shape {
		for j, c := range row {
			m.UpdateContent(x+j, y+i, c)
		}
	}
	m.blockX = x
	m.blockY = y
}

// update content matrix and terminal
func (m *Matrix) UpdateContent(x, y int, c string) {
	m.content[y][x] = c
	m.terminal.PrintAtBoard(x, y, c)
}

func (m *Matrix) Dump() {
	fmt.Printf("\r\n")
	for i := range m.content {
		for j := range m.content[i] {
			c := m.content[i][j]
			if c == BLOCK_EMPTY {
				c = "."
			}
			fmt.Printf("%s", c)
		}
		fmt.Printf("\r\n")
	}
}
