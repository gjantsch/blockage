package game

import (
	"math/rand"
)

type XDirection int

const (
	Left  XDirection = -1
	Right XDirection = 1
)

type YDirection int

const (
	Up   YDirection = -1
	Down YDirection = 1
)

type MoveOutcome int

const (
	Moved MoveOutcome = iota
	Landed
	GameOver
)

const (
	// board cell fill state; kept as named symbols (not just " "/"*") so they
	// can be swapped for more visible debug chars without touching logic
	BL_NIL = " "
	BL_FIL = "*"
)

// THE MATRIX
// the matrix is the main board where the blocks falls
// and the game happens
type Matrix struct {
	width           int
	height          int
	content         [][]string
	terminal        Renderer
	block           Block
	blockX          int
	blockY          int
	blocks          []Block
	currentBlockPtr int
	collided        bool
	movesCount      int
	rowsCompleted   int
}

func NewMatrix(width int, height int, term Renderer) Matrix {
	m := Matrix{width: width, height: height, terminal: term}
	m.content = make([][]string, height)
	for i := range m.content {
		m.content[i] = make([]string, width)
		for j := range m.content[i] {
			m.content[i][j] = BL_NIL
		}
	}
	m.blocks = InitBlocks()

	return m
}

func (m *Matrix) PickRandomBlock() {
	m.collided = false
	m.movesCount = 0
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

	// TODO: check if the rotation is possible
	if !m.BlockCanRotate() {
		return
	}

	m.RemoveBlock()
	m.block.Rotate()

	newWidth := m.block.Width()
	newHeight := m.block.Height()

	// adjust block position to keep the center
	m.blockX += (oldWidth - newWidth) / 2
	m.blockY += (oldHeight - newHeight) / 2

	m.PutBlock()
}

// when block rotates the matrix will be transposed
// and when we put the block on the board the position
// is based on its height minus the board height
func (m *Matrix) BlockCanRotate() bool {

	newWidth := m.block.Height()
	newHeight := m.block.Width()

	// so, if we rotate, X axis must not flood
	if m.blockX < 0 || m.blockX+newWidth > m.width {
		return false
	}

	// and Y axis must not flood as well
	if m.blockY < 0 || m.blockY+newHeight > m.height {
		return false
	}

	return true
}

func (m *Matrix) ChangeBlockShape() {
	m.RemoveBlock()

	m.currentBlockPtr = (m.currentBlockPtr + 1) % len(m.blocks)
	b := m.blocks[m.currentBlockPtr]
	b.X = m.blockX
	b.Y = m.blockY
	m.block = b.Clone()

	m.PutBlock()
}

func (m *Matrix) BlockShapeToDot() {
	m.RemoveBlock()

	b := m.blocks[0]
	b.X = m.blockX
	b.Y = m.blockY
	m.block = b.Clone()

	m.PutBlock()
}

func (m *Matrix) HasBlock() bool {
	return !m.block.IsEmpty()
}

func (m *Matrix) BlockHitTop() bool {
	return m.blockY == 0
}

func (m *Matrix) Render() error {
	return m.terminal.DrawBoard(m.width, m.height)
}

func (m *Matrix) RemoveBlock() {
	for i, row := range m.block.Shape {
		for j := range row {
			if m.block.Shape[i][j] == BL_FIL {
				m.UpdateContent(m.blockX+j, m.blockY+i, BL_NIL)
			}
		}
	}
}

func (m *Matrix) PutBlock() {
	for i, row := range m.block.Shape {
		for j := range row {
			if m.block.Shape[i][j] == BL_FIL {
				m.UpdateContent(m.blockX+j, m.blockY+i, m.block.Shape[i][j])
			}
		}
	}
}

// --- COLISION DETECTION
func (m *Matrix) WillCollide(x, y int) bool {
	for i, row := range m.block.Shape {
		for j, c := range row {
			if c == BL_FIL {
				newX := x + j
				newY := y + i
				if m.content[newY][newX] == BL_FIL {
					return true
				}
			}
		}
	}
	return false
}

// --- MOVE BLOCK FUNCTIONS
// --- IS LEGAL MOVE or DO WE HIT ANOTHER BLOCK?
//
// In a simple and generic way we can move it to the next place and check that
// no BLOCK_FILLED on the shape hits a BLOCK_FILLED on the matrix on the next
// position.
// Key decision:
// - MoveBlockX/Y will do the actual move.
// - NextX/Y will return the next coordinate
func (m *Matrix) NextX(direction XDirection) int {
	nextX := m.blockX
	if direction == Left && m.blockX > 0 {
		nextX--
	}

	if direction == Right && m.blockX < m.width-m.block.Width() {
		nextX++
	}

	if m.WillCollide(nextX, m.blockY) {
		m.collided = true
		return m.blockX
	}

	return nextX
}

func (m *Matrix) NextY(direction YDirection) int {
	nextY := m.blockY

	if direction == Up && m.blockY > 0 {
		nextY--
	}

	if direction == Down && m.blockY < m.height-m.block.Height() {
		nextY++
	}

	if m.WillCollide(m.blockX, nextY) {
		m.collided = true
		return m.blockY
	}

	return nextY
}

func (m *Matrix) CheckForFullRows() int {
	removed := true
	for removed {
		removed = false
		for y := 0; y < m.height; y++ {
			full := true
			for x := 0; x < m.width && full; x++ {
				full = full && m.content[y][x] != BL_NIL
			}
			if full {
				removed = true
				for ny := y; ny < m.height-1; ny++ {
					for nx := 0; nx < m.width; nx++ {
						uy := ny + 1
						m.UpdateContent(nx, ny, m.content[uy][nx])
						m.UpdateContent(nx, uy, BL_NIL)
					}
				}
				m.rowsCompleted++
			}
		}
	}
	return m.rowsCompleted
}

func (m *Matrix) MoveBlockX(direction XDirection) {
	m.RemoveBlock()
	m.blockX = m.NextX(direction)
	m.PutBlock()
}

func (m *Matrix) MoveBlockY(direction YDirection) {
	m.RemoveBlock()
	m.blockY = m.NextY(direction)
	m.PutBlock()
}

func (m *Matrix) MoveBlock() MoveOutcome {
	m.movesCount++
	m.MoveBlockY(Up)

	if m.movesCount == 1 && m.collided {
		return GameOver
	}

	if m.BlockHitTop() || m.collided {
		return Landed
	}

	return Moved
}

func (m *Matrix) PlaceBlockAtBottom() {
	bottomX := ((m.width - m.block.Width()) / 2)
	bottomY := m.height - m.block.Height()
	m.PutBlockAt(bottomX, bottomY)
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
	if x < 0 || x >= m.width || y < 0 || y >= m.height {
		return
	}
	m.content[y][x] = c
	m.terminal.PrintAtBoard(x, y, c)
}
