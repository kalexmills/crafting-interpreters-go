package main

func main() {
	chunk := Chunk{}
	constant := chunk.AddConstant(1.2)
	chunk.Write(OP_CONSTANT, 123)
	chunk.Write(byte(constant), 123)
	chunk.Write(OP_RETURN, 123)
	DisassembleChunk(&chunk, "test chunk")
	Interpret(&chunk)
}
