package game

// Renderer is the drawing surface Matrix needs; satisfied implicitly by render.Terminal.
type Renderer interface {
	DrawBoard(width, height int) error
	PrintAtBoard(x, y int, c string)
	StatusBar(content string)
}
