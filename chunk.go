package main

const (
	OP_RETURN byte = iota
	OP_CONSTANT
)

// Chunk is a chunk
type Chunk struct {
	// N.B. the dynamic array implementation in go handles all the features mentioned in the book.
	Code      []byte
	lines     []int
	constants ValueArray
}

func (c Chunk) Count() int {
	return len(c.Code)
}

func (c *Chunk) AddConstant(v Value) int {
	c.constants.WriteValue(v)
	return c.constants.Count() - 1
}

func (c *Chunk) Write(b byte, line int) {
	c.Code = append(c.Code, b)
	c.lines = append(c.lines, line)
}
