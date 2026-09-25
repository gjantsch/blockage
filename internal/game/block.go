package game

// THE BLOCK
// hold a single Block data
// with position and form
type Block struct {
	X     int
	Y     int
	Shape [][]bool
}

func (b *Block) IsEmpty() bool {
	return len(b.Shape) == 0
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

// Rotates a block doing a matrix
// transposition
func (b *Block) Rotate() {

	// rows turn columns
	cols := len(b.Shape)
	// columns to rows
	rows := len(b.Shape[0])

	// initialize new Block
	t := make([][]bool, rows)
	for i := range t {
		t[i] = make([]bool, cols)
	}

	// transpose
	for i, row := range b.Shape {
		for j, s := range row {
			// for proper visual effect the
			// transposition is backwards
			t[len(row)-j-1][i] = s
		}
	}

	// update block with transposed one
	b.Shape = t
}

// Initializes and array of blocks
func InitBlocks() []Block {

	b := []Block{
		{
			X: 0,
			Y: 0,
			Shape: [][]bool{
				{true},
			},
		},
		{
			X: 0,
			Y: 0,
			Shape: [][]bool{
				{true, true, false},
				{true, true, true},
			},
		},
		{
			X: 0,
			Y: 0,
			Shape: [][]bool{
				{true, true},
				{true, true},
				{true, true},
			},
		},
		{
			X: 0,
			Y: 0,
			Shape: [][]bool{
				{true, true, true},
			},
		},
		{
			X: 0,
			Y: 0,
			Shape: [][]bool{
				{true, true, true},
				{false, true, false},
			},
		},
		{
			X: 0,
			Y: 0,
			Shape: [][]bool{
				{true, true, true},
				{true, false, false},
			},
		},
		{
			X: 0,
			Y: 0,
			Shape: [][]bool{
				{true, true, true},
				{false, false, true},
			},
		},
		{
			X: 0,
			Y: 0,
			Shape: [][]bool{
				{true, true, true},
				{true, false, true},
			},
		},
		{
			X: 0,
			Y: 0,
			Shape: [][]bool{
				{true, true, true},
				{false, true, true},
				{false, false, true},
			},
		},
		{
			X: 0,
			Y: 0,
			Shape: [][]bool{
				{true, true},
				{false, true},
			},
		},
		{
			X: 0,
			Y: 0,
			Shape: [][]bool{
				{true, true},
				{true, true},
			},
		},
		{
			X: 0,
			Y: 0,
			Shape: [][]bool{
				{true},
				{true},
			},
		},
	}

	return b
}

func (b Block) Clone() Block {
	clone := b
	clone.Shape = make([][]bool, len(b.Shape))

	for i := range b.Shape {
		clone.Shape[i] = make([]bool, len(b.Shape[i]))
		copy(clone.Shape[i], b.Shape[i])
	}

	return clone
}
