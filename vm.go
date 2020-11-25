package main

import "fmt"

const DEBUG_TRACE_EXECUTION = true // N.B. this does not use conditional compilation; it's handled at runtime.

var vm VM

const STACK_MAX = 256

type VM struct {
	// N.B. uses slice indices instead 'real C-pointers', to avoid the unsafe package.
	chunk    *Chunk
	ip       int
	stack    [STACK_MAX]Value
	stackTop int
}

type InterpretResult byte

const (
	INTERPRET_OK InterpretResult = iota
	INTERPRET_COMPILE_ERROR
	INTERPRET_RUNTIME_ERROR
)

func Interpret(chunk *Chunk) InterpretResult {
	vm.chunk = chunk
	vm.ip = 0
	return run()
}

func run() InterpretResult {
	for {
		if DEBUG_TRACE_EXECUTION {
			fmt.Printf("          ")
			for i := 0; i < vm.stackTop; i++ {
				fmt.Printf("[ ")
				vm.stack[i].Print()
				fmt.Printf(" ]")
			}
			println()
			disassembleInstruction(vm.chunk, vm.ip)
		}
		instruction := vm.readByte()
		switch instruction {
		case OP_CONSTANT:
			constant := vm.readConstant()
			vm.push(constant)
		case OP_RETURN:
			vm.pop().Print()
			println()
			return INTERPRET_OK
		}
	}
}

func (v *VM) readByte() byte {
	result := v.chunk.Code[v.ip]
	v.ip++
	return result
}

func (v *VM) readConstant() Value {
	result := v.chunk.constants.Values[v.readByte()]
	return result
}

func (v *VM) resetStack() {
	vm.stackTop = 0
}

func (v *VM) push(value Value) {
	v.stack[v.stackTop] = value
	v.stackTop++
}

func (v *VM) pop() Value {
	v.stackTop--
	return v.stack[v.stackTop]
}
