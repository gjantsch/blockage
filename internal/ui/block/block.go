package block

const (
	BL_NIL = " "
	BL_FIL = "*"
)

// THE BLOCK
// hold a single Block data
// with position and form
type Block struct {
	X         int
	Y         int
	Direction int
	Shape     [][]string
}

func EmptyBlock() Block {
	return Block{
		X:     0,
		Y:     0,
		Shape: [][]string{},
	}
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
	t := make([][]string, rows)
	for i := range t {
		t[i] = make([]string, cols)
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
		Block{
			X: 0,
			Y: 0,
			Shape: [][]string{
				[]string{BL_FIL},
			},
		},
		Block{
			X: 0,
			Y: 0,
			Shape: [][]string{
				[]string{BL_FIL, BL_FIL, BL_NIL},
				[]string{BL_FIL, BL_FIL, BL_FIL},
			},
		},
		Block{
			X: 0,
			Y: 0,
			Shape: [][]string{
				[]string{BL_FIL, BL_FIL},
				[]string{BL_FIL, BL_FIL},
				[]string{BL_FIL, BL_FIL},
			},
		},
		Block{
			X: 0,
			Y: 0,
			Shape: [][]string{
				[]string{BL_FIL, BL_FIL, BL_FIL},
			},
		},
		Block{
			X: 0,
			Y: 0,
			Shape: [][]string{
				[]string{BL_FIL, BL_FIL, BL_FIL},
				[]string{BL_NIL, BL_FIL, BL_NIL},
			},
		},
		Block{
			X: 0,
			Y: 0,
			Shape: [][]string{
				[]string{BL_FIL, BL_FIL, BL_FIL},
				[]string{BL_NIL, BL_FIL, BL_NIL},
			},
		},
		Block{
			X: 0,
			Y: 0,
			Shape: [][]string{
				[]string{BL_FIL, BL_FIL, BL_FIL},
				[]string{BL_NIL, BL_FIL, BL_FIL},
				[]string{BL_NIL, BL_NIL, BL_FIL},
			},
		},
		Block{
			X: 0,
			Y: 0,
			Shape: [][]string{
				[]string{BL_FIL, BL_FIL, BL_FIL},
				[]string{BL_NIL, BL_NIL, BL_FIL},
			},
		},
		Block{
			X: 0,
			Y: 0,
			Shape: [][]string{
				[]string{BL_FIL, BL_FIL},
				[]string{BL_NIL, BL_FIL},
			},
		},
		Block{
			X: 0,
			Y: 0,
			Shape: [][]string{
				[]string{BL_FIL, BL_FIL},
				[]string{BL_FIL, BL_FIL},
			},
		},
		Block{
			X: 0,
			Y: 0,
			Shape: [][]string{
				[]string{BL_FIL},
				[]string{BL_FIL},
			},
		},
		Block{
			X: 0,
			Y: 0,
			Shape: [][]string{
				[]string{BL_FIL, BL_FIL, BL_FIL},
				[]string{BL_FIL, BL_NIL, BL_FIL},
			},
		}, Block{
			X: 0,
			Y: 0,
			Shape: [][]string{
				[]string{BL_FIL, BL_FIL, BL_FIL},
				[]string{BL_NIL, BL_NIL, BL_FIL},
			},
		},
	}

	return b
}

func (b Block) Clone() Block {
	clone := b
	clone.Shape = make([][]string, len(b.Shape))

	for i := range b.Shape {
		clone.Shape[i] = make([]string, len(b.Shape[i]))
		copy(clone.Shape[i], b.Shape[i])
	}

	return clone
}
