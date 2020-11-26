package main

import (
	"fmt"
	"os"
)

const DEBUG_TRACE_EXECUTION = true // N.B. this does not use conditional compilation; it's handled at runtime.
const DEBUG_PRINT_CODE = true

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

func interpret(source string) InterpretResult {
	var chunk Chunk

	if !compile(source, &chunk) {
		return INTERPRET_COMPILE_ERROR
	}

	vm.chunk = &chunk
	vm.ip = 0
	result := run()

	return result
}

func Interpret(chunk *Chunk) InterpretResult {
	vm.chunk = chunk
	vm.ip = 0
	return run()
}

var ADD = func(a, b float64) Value { return NumberVal(a + b) }
var SUB = func(a, b float64) Value { return NumberVal(a - b) }
var MUL = func(a, b float64) Value { return NumberVal(a * b) }
var DIV = func(a, b float64) Value { return NumberVal(a / b) }
var GT = func(a, b float64) Value { return BoolVal(a > b) }
var LT = func(a, b float64) Value { return BoolVal(a < b) }

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
		case OP_NEGATE:
			if !isNumber(vm.peek(0)) {
				runtimeError("Operand must be a number.")
				return INTERPRET_RUNTIME_ERROR
			}
			vm.push(NumberVal(-vm.pop().AsNumber()))
		case OP_EQUAL:
			b, a := vm.pop(), vm.pop()
			vm.push(BoolVal(valuesEqual(a, b)))
		case OP_GREATER:
			if !vm.binaryOp(GT) {
				return INTERPRET_COMPILE_ERROR
			}
		case OP_LESS:
			if !vm.binaryOp(LT) {
				return INTERPRET_COMPILE_ERROR
			}
		case OP_ADD:
			if !vm.binaryOp(ADD) {
				return INTERPRET_RUNTIME_ERROR
			}
		case OP_SUBTRACT:
			if !vm.binaryOp(SUB) {
				return INTERPRET_RUNTIME_ERROR
			}
		case OP_MULTIPLY:
			if !vm.binaryOp(MUL) {
				return INTERPRET_RUNTIME_ERROR
			}
		case OP_DIVIDE:
			if !vm.binaryOp(DIV) {
				return INTERPRET_RUNTIME_ERROR
			}
		case OP_NOT:
			vm.push(BoolVal(isFalsey(vm.pop())))
		case OP_NIL:
			vm.push(NilVal(struct{}{}))
		case OP_TRUE:
			vm.push(BoolVal(true))
		case OP_FALSE:
			vm.push(BoolVal(false))
		case OP_RETURN:
			vm.pop().Print()
			println()
			return INTERPRET_OK
		}
	}
}

func valuesEqual(a, b Value) bool {
	if a.Type() != b.Type() {
		return false
	}
	switch a.Type() {
	case VAL_BOOL:
		return a.AsBoolean() == b.AsBoolean()
	case VAL_NIL:
		return true
	case VAL_NUMBER:
		return a.AsNumber() == b.AsNumber()
	}
	return false
}

func isFalsey(value Value) bool {
	return isNil(value) || (isBool(value) && !value.AsBoolean())
}

func (v *VM) binaryOp(op func(float64, float64) Value) bool {
	if !isNumber(v.peek(0)) || !isNumber(v.peek(1)) {
		runtimeError("Operands must be numbers.")
		return false
	}
	b, a := v.pop().AsNumber(), v.pop().AsNumber()
	v.push(op(a, b))
	return true
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

func (v VM) peek(distance int) Value {
	return vm.stack[vm.stackTop-distance-1]
}

func (v *VM) push(value Value) {
	v.stack[v.stackTop] = value
	v.stackTop++
}

func (v *VM) pop() Value {
	v.stackTop--
	return v.stack[v.stackTop]
}

func runtimeError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintln(os.Stderr)
	instruction := vm.ip - 1
	line := vm.chunk.lines[instruction]
	fmt.Fprintf(os.Stderr, "[line %d] in script\n", line)
}
