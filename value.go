package main

import "fmt"

type Value float64

func (v Value) Print() {
	fmt.Printf("%g", v)
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
