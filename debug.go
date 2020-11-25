package main

import "fmt"

func DisassembleChunk(c *Chunk, name string) {
	fmt.Printf("== %s ==\n", name)
	for offset := 0; offset < c.Count(); {
		offset = disassembleInstruction(c, offset)
	}
}

func disassembleInstruction(chunk *Chunk, offset int) int {
	fmt.Printf("%04d ", offset)
	if offset > 0 && chunk.lines[offset] == chunk.lines[offset-1] {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%4d ", chunk.lines[offset])
	}
	instruction := chunk.Code[offset]
	switch OpCode(instruction) {
	case OP_RETURN:
		return simpleInstruction("OP_RETURN", offset)
	case OP_CONSTANT:
		return constantInstruction("OP_CONSTANT", chunk, offset)
	default:
		fmt.Printf("Unknown opcode %d\n", instruction)
		return offset + 1
	}
}

func constantInstruction(name string, chunk *Chunk, offset int) int {
	constant := chunk.Code[offset+1]
	fmt.Printf("%-16s %4d '", name, constant)
	chunk.constants.Values[constant].Print()
	fmt.Printf("'\n")
	return offset + 2
}

func simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}
