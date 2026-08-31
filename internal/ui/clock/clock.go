package clock

type Clock struct {
	ptr   int
	chars []string
}

func NewClock() *Clock {
	return &Clock{
		ptr:   0,
		chars: []string{"|", "/", "-", "\\"},
	}
}

func (c *Clock) Next() string {
	char := c.chars[c.ptr]
	c.ptr = (c.ptr + 1) % len(c.chars)
	return char
}
