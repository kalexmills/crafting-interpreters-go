package main

import "fmt"

type ValueType uint8

const (
	VAL_BOOL ValueType = iota
	VAL_NIL
	VAL_NUMBER
)

type Value interface { // N.B. go doesn't have sum-types, so we use an interface; this is more verbose
	Type() ValueType
	AsBoolean() bool
	AsNumber() float64
	Print()
}

func isNumber(v Value) bool {
	return v.Type() == VAL_NUMBER
}
func isBool(v Value) bool {
	return v.Type() == VAL_BOOL
}
func isNil(v Value) bool {
	return v.Type() == VAL_NIL
}

type BoolVal bool

func (bv BoolVal) Type() ValueType {
	return VAL_BOOL
}

func (bv BoolVal) AsBoolean() bool {
	return bool(bv)
}

func (bv BoolVal) AsNumber() float64 {
	panic("bool value is not a number!") // N.B. panicking is one choice... returning the zero value is another...
}

func (bv BoolVal) Print() {
	fmt.Printf("%t", bool(bv))
}

type NumberVal float64

func (nv NumberVal) Type() ValueType {
	return VAL_NUMBER
}

func (nv NumberVal) AsBoolean() bool {
	panic("number value is not a boolean!")
}

func (nv NumberVal) AsNumber() float64 {
	return float64(nv)
}

func (nv NumberVal) Print() {
	fmt.Printf("%g", float64(nv))
}

type ValueArray struct {
	Values []Value
}

func (a ValueArray) Count() int {
	return len(a.Values)
}

func (a *ValueArray) WriteValue(v Value) error {
	a.Values = append(a.Values, v)
	return nil
}
