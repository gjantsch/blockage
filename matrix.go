package main

import "fmt"

// THE MATRIX
// the matrix is the main board where the blocks falls
// and the game happens
type Matrix struct {
	width          int
	height         int
	content        [][]string
	terminal       *Terminal
	block          Block
	blockX         int
	blockY         int
	BlockDirection int
}

func NewMatrix(width int, height int, term *Terminal) Matrix {
	m := Matrix{width: width, height: height, terminal: term}
	m.content = make([][]string, height)
	for i := range m.content {
		m.content[i] = make([]string, width)
		for j := range m.content[i] {
			m.content[i][j] = " "
		}
	}
	return m
}

func (m *Matrix) HasBlock() bool {
	return !m.block.IsEmpty()
}
func (m *Matrix) Render() {
	m.terminal.DrawBoard(m.width, m.height)
}

func (m *Matrix) GetCenter(block Block) (x, y int) {
	// calculate the center position
	centerX := (m.width - block.Width()) / 2
	centerY := (m.height - block.Height()) / 2
	return centerX, centerY
}

func (m *Matrix) RemoveBlockAtCenter(block Block) {
	centerX, centerY := m.GetCenter(block)

	// remove the block at the center position
	for i, row := range block.Shape {
		for j := range row {
			m.UpdateContent(centerX+j, centerY+i, " ")
		}
	}
}

func (m *Matrix) RemoveBlock() {
	for i, row := range m.block.Shape {
		for j := range row {
			m.UpdateContent(m.blockX+j, m.blockY+i, " ")
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

func (m *Matrix) MoveBlock() {
	m.RemoveBlock()
	m.blockY += m.BlockDirection

	m.PutBlock()
	m.terminal.StatusBar(fmt.Sprintf("%d %d", m.blockX, m.blockY))
}

func (m *Matrix) PlaceBlockAtBottom(block Block) {
	bottomX := (m.width - block.Width()) / 2
	bottomY := m.height - block.Height()
	m.PutBlockAt(block, bottomX, bottomY)

	m.blockX = bottomX
	m.blockY = bottomY
	m.block = block
}

func (m *Matrix) PlaceBlockAtTop(block Block) {
	topX := (m.width - block.Width()) / 2
	topY := 0
	m.PutBlockAt(block, topX, topY)
	m.blockX = topX
	m.blockY = topY
	m.block = block
}

func (m *Matrix) PlaceBlockAtCenter(block Block) {
	centerX, centerY := m.GetCenter(block)

	// place the block at the center position
	for i, row := range block.Shape {
		for j, c := range row {
			m.UpdateContent(centerX+j, centerY+i, c)
		}
	}

	m.blockX = centerX
	m.blockY = centerY
	m.block = block
}

func (m *Matrix) PutBlockAt(block Block, x, y int) {
	// place the block at position
	for i, row := range block.Shape {
		for j, c := range row {
			m.UpdateContent(x+j, y+i, c)
		}
	}
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
			if c == " " {
				c = "."
			}
			fmt.Printf("%s", c)
		}
		fmt.Printf("\r\n")
	}
}
