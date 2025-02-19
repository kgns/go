// run -gcflags=-G=3

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"
)

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~complex64 | ~complex128
}

// numericAbs matches a struct containing a numeric type that has an Abs method.
type numericAbs[T Numeric] interface {
	~struct{ Value T }
	Abs() T
}

// AbsDifference computes the absolute value of the difference of
// a and b, where the absolute value is determined by the Abs method.
func absDifference[T Numeric, U numericAbs[T]](a, b U) T {
	d := a.Value - b.Value
	dt := U{Value: d}
	return dt.Abs()
}

// orderedNumeric matches numeric types that support the < operator.
type orderedNumeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Complex matches the two complex types, which do not have a < operator.
type Complex interface {
	~complex64 | ~complex128
}

// orderedAbs is a helper type that defines an Abs method for
// a struct containing an ordered numeric type.
type orderedAbs[T orderedNumeric] struct {
	Value T
}

func (a orderedAbs[T]) Abs() T {
	if a.Value < 0 {
		return -a.Value
	}
	return a.Value
}

// complexAbs is a helper type that defines an Abs method for
// a struct containing a complex type.
type complexAbs[T Complex] struct {
	Value T
}

func (a complexAbs[T]) Abs() T {
	r := float64(real(a.Value))
	i := float64(imag(a.Value))
	d := math.Sqrt(r*r + i*i)
	return T(complex(d, 0))
}

// OrderedAbsDifference returns the absolute value of the difference
// between a and b, where a and b are of an ordered type.
func OrderedAbsDifference[T orderedNumeric](a, b T) T {
	return absDifference(orderedAbs[T]{a}, orderedAbs[T]{b})
}

// ComplexAbsDifference returns the absolute value of the difference
// between a and b, where a and b are of a complex type.
func ComplexAbsDifference[T Complex](a, b T) T {
	return absDifference(complexAbs[T]{a}, complexAbs[T]{b})
}

func main() {
	if got, want := OrderedAbsDifference(1.0, -2.0), 3.0; got != want {
		panic(fmt.Sprintf("got = %v, want = %v", got, want))
	}
	if got, want := OrderedAbsDifference(-1.0, 2.0), 3.0; got != want {
		panic(fmt.Sprintf("got = %v, want = %v", got, want))
	}
	if got, want := OrderedAbsDifference(-20, 15), 35; got != want {
		panic(fmt.Sprintf("got = %v, want = %v", got, want))
	}

	if got, want := ComplexAbsDifference(5.0+2.0i, 2.0-2.0i), 5+0i; got != want {
		panic(fmt.Sprintf("got = %v, want = %v", got, want))
	}
	if got, want := ComplexAbsDifference(2.0-2.0i, 5.0+2.0i), 5+0i; got != want {
		panic(fmt.Sprintf("got = %v, want = %v", got, want))
	}
}
