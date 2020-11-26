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
	switch instruction {
	case OP_RETURN:
		return simpleInstruction("OP_RETURN", offset)
	case OP_CONSTANT:
		return constantInstruction("OP_CONSTANT", chunk, offset)
	case OP_NIL:
		return simpleInstruction("OP_NIL", offset)
	case OP_TRUE:
		return simpleInstruction("OP_TRUE", offset)
	case OP_FALSE:
		return simpleInstruction("OP_FALSE", offset)
	case OP_EQUAL:
		return simpleInstruction("OP_EQUAL", offset)
	case OP_GREATER:
		return simpleInstruction("OP_GREATER", offset)
	case OP_LESS:
		return simpleInstruction("OP_LESS", offset)
	case OP_ADD:
		return simpleInstruction("OP_ADD", offset)
	case OP_SUBTRACT:
		return simpleInstruction("OP_SUBTRACT", offset)
	case OP_MULTIPLY:
		return simpleInstruction("OP_MULTIPLY", offset)
	case OP_DIVIDE:
		return simpleInstruction("OP_DIVIDE", offset)
	case OP_NOT:
		return simpleInstruction("OP_NOT", offset)
	case OP_NEGATE:
		return simpleInstruction("OP_NEGATE", offset)
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
