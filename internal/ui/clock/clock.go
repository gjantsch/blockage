package clock

type Clock struct {
	ptr int
}

func NewClock() *Clock {
	return &Clock{
		ptr: 0,
	}
}

func (c *Clock) Next() string {
	spinner := [4]string{"|", "/", "-", "\\"}
	char := spinner[c.ptr]
	c.ptr = (c.ptr + 1) % 4
	return char
}
