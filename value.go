package main

import "fmt"

type ValueType uint8

const (
	VAL_BOOL ValueType = iota
	VAL_NIL
	VAL_NUMBER
	VAL_OBJ
)

type Value interface { // N.B. go doesn't have sum-types, so we use an interface; this is more verbose
	Type() ValueType
	AsBoolean() bool
	AsNumber() float64
	AsObj() Obj
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

func isObj(v Value) bool {
	return v.Type() == VAL_OBJ
}

type NilVal struct{}

func (nv NilVal) Type() ValueType {
	return VAL_NIL
}

func (nv NilVal) AsBoolean() bool {
	panic("nil value is not a boolean!") // N.B. panicking is one choice... returning the zero value is another...
}

func (nv NilVal) AsNumber() float64 {
	panic("nil value is not a number!")
}

func (nv NilVal) AsObj() Obj {
	panic("nil value is not an object!")
}

func (nv NilVal) Print() {
	fmt.Printf("nil")
}

type BoolVal bool

func (bv BoolVal) Type() ValueType {
	return VAL_BOOL
}

func (bv BoolVal) AsBoolean() bool {
	return bool(bv)
}

func (bv BoolVal) AsNumber() float64 {
	panic("bool value is not a number!")
}

func (bv BoolVal) AsObj() Obj {
	panic("bool value is not an object!")
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

func (nv NumberVal) AsObj() Obj {
	panic("number value is not an object!")
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
