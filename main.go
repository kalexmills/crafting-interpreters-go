package main

func main() {
	chunk := Chunk{}
	constant := chunk.AddConstant(1.2)
	chunk.Write(OP_CONSTANT, 123)
	chunk.Write(byte(constant), 123)

	constant = chunk.AddConstant(3.4)
	chunk.Write(OP_CONSTANT, 123)
	chunk.Write(byte(constant), 123)

	chunk.Write(OP_ADD, 123)

	constant = chunk.AddConstant(5.6)
	chunk.Write(OP_CONSTANT, 123)
	chunk.Write(byte(constant), 123)

	chunk.Write(OP_DIVIDE, 123)
	chunk.Write(OP_NEGATE, 123)
	chunk.Write(OP_RETURN, 123)
	DisassembleChunk(&chunk, "test chunk")
	Interpret(&chunk)
}
