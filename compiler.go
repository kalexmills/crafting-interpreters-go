package main

import (
	"fmt"
	"os"
	"strconv"
)

func compile(source string, chunk *Chunk) bool {
	initScanner(source)
	compilingChunk = chunk
	advanceParser()
	expression()
	consume(TOKEN_EOF, "Expect end of expression.")
	endCompiler()
	return !parser.HadError
}

type Parser struct {
	Current   Token
	Previous  Token
	HadError  bool
	PanicMode bool
}

type Precedence uint8

const (
	PREC_NONE Precedence = iota
	PREC_ASSIGNMENT
	PREC_OR
	PREC_AND
	PREC_EQUALITY
	PREC_COMPARISON
	PREC_TERM
	PREC_FACTOR
	PREC_UNARY
	PREC_CALL
	PREC_PRIMARY
)

type ParseRule struct {
	Prefix     ParseFn
	Infix      ParseFn
	Precedence Precedence
}

type ParseFn = func()

var parser Parser

func consume(tokenType TokenType, msg string) {
	if parser.Current.Type == tokenType {
		advanceParser()
		return
	}
	errorAtCurrent(msg)
}

func advanceParser() {
	parser.Previous = parser.Current
	for {
		parser.Current = scanToken()
		if parser.Current.Type != TOKEN_ERROR {
			break
		}
		errorAtCurrent(*parser.Current.Source)
	}
}

func errorAtCurrent(msg string) {
	errorAt(&parser.Current, msg)
}

func errorRpt(msg string) {
	errorAt(&parser.Previous, msg)
}

func errorAt(token *Token, msg string) {
	if parser.PanicMode {
		return
	}
	parser.PanicMode = true
	fmt.Fprintf(os.Stderr, "[line %d] Error", token.Line)
	if token.Type == TOKEN_EOF {
		fmt.Fprintf(os.Stderr, " at end")
	} else if token.Type == TOKEN_ERROR {
		// nothing
	} else {
		fmt.Fprintf(os.Stderr, " at '%s'", (*token.Source)[token.Start:token.Start+token.Length])
	}
	fmt.Fprintf(os.Stderr, ": %s\n", msg)
	parser.HadError = true
}

func expression() {
	parsePrecedence(PREC_ASSIGNMENT)
}

func parsePrecedence(precedence Precedence) {
	advanceParser()
	prefixRule := rules[parser.Previous.Type].Prefix
	if prefixRule == nil {
		errorRpt("expect expression.")
		return
	}
	prefixRule()

	for precedence <= rules[parser.Current.Type].Precedence {
		advanceParser()
		infixRule := rules[parser.Previous.Type].Infix
		infixRule()
	}
}

func emitConstant(value Value) {
	emitBytes(OP_CONSTANT, makeConstant(value))
}

func makeConstant(value Value) byte {
	constant := currentChunk().AddConstant(value)
	if constant > 256 {
		errorRpt("too many constants in one chunk.")
		return 0
	}
	return byte(constant)
}

var compilingChunk *Chunk

func currentChunk() *Chunk {
	return compilingChunk
}

func emitByte(b byte) {
	currentChunk().Write(b, parser.Previous.Line)
}

func emitReturn() {
	emitByte(OP_RETURN)
}

func emitBytes(b1, b2 byte) {
	emitByte(b1)
	emitByte(b2)
}

func endCompiler() {
	emitReturn()
	if DEBUG_PRINT_CODE {
		if !parser.HadError {
			DisassembleChunk(currentChunk(), "code")
		}
	}
}

func compileBinary() {
	operatorType := parser.Previous.Type
	rule := rules[operatorType]
	parsePrecedence(Precedence(rule.Precedence + 1))

	switch operatorType {
	case TOKEN_PLUS:
		emitByte(OP_ADD)
	case TOKEN_MINUS:
		emitByte(OP_SUBTRACT)
	case TOKEN_STAR:
		emitByte(OP_MULTIPLY)
	case TOKEN_SLASH:
		emitByte(OP_DIVIDE)
	}
}

func compileGrouping() {
	expression()
	consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
}

func compileNumber() {
	value, _ := strconv.ParseFloat((*parser.Previous.Source)[parser.Previous.Start:parser.Previous.Start+parser.Previous.Length], 64)
	emitConstant(NumberVal(value))
}

func compileUnary() {
	operatorType := parser.Previous.Type

	parsePrecedence(PREC_UNARY)

	switch operatorType {
	case TOKEN_MINUS:
		emitByte(OP_NEGATE)
	default:
		return
	}
}

var rules map[TokenType]ParseRule

func init() {
	rules = map[TokenType]ParseRule{
		TOKEN_LEFT_PAREN:    {compileGrouping, nil, PREC_NONE},
		TOKEN_RIGHT_PAREN:   {nil, nil, PREC_NONE},
		TOKEN_LEFT_BRACE:    {nil, nil, PREC_NONE},
		TOKEN_RIGHT_BRACE:   {nil, nil, PREC_NONE},
		TOKEN_COMMA:         {nil, nil, PREC_NONE},
		TOKEN_DOT:           {nil, nil, PREC_NONE},
		TOKEN_MINUS:         {compileUnary, compileBinary, PREC_TERM},
		TOKEN_PLUS:          {nil, compileBinary, PREC_TERM},
		TOKEN_SEMICOLON:     {nil, nil, PREC_NONE},
		TOKEN_SLASH:         {nil, compileBinary, PREC_FACTOR},
		TOKEN_STAR:          {nil, compileBinary, PREC_FACTOR},
		TOKEN_BANG:          {nil, nil, PREC_NONE},
		TOKEN_BANG_EQUAL:    {nil, nil, PREC_NONE},
		TOKEN_EQUAL:         {nil, nil, PREC_NONE},
		TOKEN_EQUAL_EQUAL:   {nil, nil, PREC_NONE},
		TOKEN_GREATER:       {nil, nil, PREC_NONE},
		TOKEN_GREATER_EQUAL: {nil, nil, PREC_NONE},
		TOKEN_LESS:          {nil, nil, PREC_NONE},
		TOKEN_LESS_EQUAL:    {nil, nil, PREC_NONE},
		TOKEN_IDENTIFIER:    {nil, nil, PREC_NONE},
		TOKEN_STRING:        {nil, nil, PREC_NONE},
		TOKEN_NUMBER:        {compileNumber, nil, PREC_NONE},
		TOKEN_AND:           {nil, nil, PREC_NONE},
		TOKEN_CLASS:         {nil, nil, PREC_NONE},
		TOKEN_ELSE:          {nil, nil, PREC_NONE},
		TOKEN_FALSE:         {nil, nil, PREC_NONE},
		TOKEN_FOR:           {nil, nil, PREC_NONE},
		TOKEN_FUN:           {nil, nil, PREC_NONE},
		TOKEN_IF:            {nil, nil, PREC_NONE},
		TOKEN_NIL:           {nil, nil, PREC_NONE},
		TOKEN_OR:            {nil, nil, PREC_NONE},
		TOKEN_PRINT:         {nil, nil, PREC_NONE},
		TOKEN_RETURN:        {nil, nil, PREC_NONE},
		TOKEN_SUPER:         {nil, nil, PREC_NONE},
		TOKEN_THIS:          {nil, nil, PREC_NONE},
		TOKEN_TRUE:          {nil, nil, PREC_NONE},
		TOKEN_VAR:           {nil, nil, PREC_NONE},
		TOKEN_WHILE:         {nil, nil, PREC_NONE},
		TOKEN_ERROR:         {nil, nil, PREC_NONE},
		TOKEN_EOF:           {nil, nil, PREC_NONE},
	}
}
