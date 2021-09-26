// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains simple golden tests for various examples.
// Besides validating the results when the implementation changes,
// it provides a way to look at the generated code without having
// to execute the print statements in one's head.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Golden represents a test case.
type Golden struct {
	name   string
	input  string // input; the package clause is provided when running the test.
	output string // expected output.
}

var golden = []Golden{
	{name: "day", input: day_in, output: day_out},
	{name: "offset", input: offset_in, output: offset_out},
	{name: "gap", input: gap_in, output: gap_out},
	{name: "num", input: num_in, output: num_out},
	{name: "unum", input: unum_in, output: unum_out},
	{name: "unumpos", input: unumpos_in, output: unumpos_out},
	{name: "prime", input: prime_in, output: prime_out},
	{name: "prefix", input: prefix_in, output: prefix_out},
	{name: "tokens", input: tokens_in, output: tokens_out},
}

// Each example starts with "type XXX [u]int", with a single space separating them.

// Simple test: enumeration of type int starting at 0.
const day_in = `type Day int
const (
	Monday Day = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)
`

const day_out = `
// Enum returns a list of values declared for a type.
func (Day) Enum() []interface{} {
	return []interface{}{
		Monday,
		Tuesday,
		Wednesday,
		Thursday,
		Friday,
		Saturday,
		Sunday,
	}
}
`

// Enumeration with an offset.
// Also includes a duplicate.
const offset_in = `type Number int
const (
	_ Number = iota
	One
	Two
	Three
	AnotherOne = One  // Duplicate; note that AnotherOne doesn't appear below.
)
`

const offset_out = `
// Enum returns a list of values declared for a type.
func (Number) Enum() []interface{} {
	return []interface{}{
		One,
		Two,
		Three,
	}
}
`

// Gaps and an offset.
const gap_in = `type Gap int
const (
	Two Gap = 2
	Three Gap = 3
	Five Gap = 5
	Six Gap = 6
	Seven Gap = 7
	Eight Gap = 8
	Nine Gap = 9
	Eleven Gap = 11
)
`

const gap_out = `
// Enum returns a list of values declared for a type.
func (Gap) Enum() []interface{} {
	return []interface{}{
		Two,
		Three,
		Five,
		Six,
		Seven,
		Eight,
		Nine,
		Eleven,
	}
}
`

// Signed integers spanning zero.
const num_in = `type Num int
const (
	m_2 Num = -2 + iota
	m_1
	m0
	m1
	m2
)
`

const num_out = `
// Enum returns a list of values declared for a type.
func (Num) Enum() []interface{} {
	return []interface{}{
		m_2,
		m_1,
		m0,
		m1,
		m2,
	}
}
`

// Unsigned integers spanning zero.
const unum_in = `type Unum uint
const (
	m_2 Unum = iota + 253
	m_1
)

const (
	m0 Unum = iota
	m1
	m2
)
`

const unum_out = `
// Enum returns a list of values declared for a type.
func (Unum) Enum() []interface{} {
	return []interface{}{
		m_2,
		m_1,
		m0,
		m1,
		m2,
	}
}
`

// Unsigned positive integers.
const unumpos_in = `type Unumpos uint
const (
	m253 Unumpos = iota + 253
	m254
)

const (
	m1 Unumpos = iota + 1
	m2
	m3
)
`

const unumpos_out = `
// Enum returns a list of values declared for a type.
func (Unumpos) Enum() []interface{} {
	return []interface{}{
		m253,
		m254,
		m1,
		m2,
		m3,
	}
}
`

// Enough gaps to trigger a map implementation of the method.
// Also includes a duplicate to test that it doesn't cause problems.
const prime_in = `type Prime int
const (
	p2 Prime = 2
	p3 Prime = 3
	p5 Prime = 5
	p7 Prime = 7
	p77 Prime = 7 // Duplicate; note that p77 doesn't appear below.
	p11 Prime = 11
	p13 Prime = 13
	p17 Prime = 17
	p19 Prime = 19
	p23 Prime = 23
	p29 Prime = 29
	p37 Prime = 31
	p41 Prime = 41
	p43 Prime = 43
)
`

const prime_out = `
// Enum returns a list of values declared for a type.
func (Prime) Enum() []interface{} {
	return []interface{}{
		p2,
		p3,
		p5,
		p7,
		p77,
		p11,
		p13,
		p17,
		p19,
		p23,
		p29,
		p37,
		p41,
		p43,
	}
}
`

const prefix_in = `type Type int
const (
	TypeInt Type = iota
	TypeString
	TypeFloat
	TypeRune
	TypeByte
	TypeStruct
	TypeSlice
)
`

const prefix_out = `
// Enum returns a list of values declared for a type.
func (Type) Enum() []interface{} {
	return []interface{}{
		TypeInt,
		TypeString,
		TypeFloat,
		TypeRune,
		TypeByte,
		TypeStruct,
		TypeSlice,
	}
}
`

const tokens_in = `type Token int
const (
	And Token = iota // &
	Or               // |
	Add              // +
	Sub              // -
	Ident
	Period // .

	// not to be used
	SingleBefore
	// not to be used
	BeforeAndInline // inline
	InlineGeneral /* inline general */
)
`

const tokens_out = `
// Enum returns a list of values declared for a type.
func (Token) Enum() []interface{} {
	return []interface{}{
		And,
		Or,
		Add,
		Sub,
		Ident,
		Period,
		SingleBefore,
		BeforeAndInline,
		InlineGeneral,
	}
}
`

func TestGolden(t *testing.T) {
	dir, err := ioutil.TempDir("", "enumer")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir)

	for _, test := range golden {
		g := Generator{}
		input := "package test\n" + test.input
		file := test.name + ".go"
		absFile := filepath.Join(dir, file)
		err := ioutil.WriteFile(absFile, []byte(input), 0o600)
		if err != nil {
			t.Error(err)
		}

		g.parsePackage([]string{absFile}, nil)
		// Extract the name and type of the constant from the first line.
		tokens := strings.SplitN(test.input, " ", 3)
		if len(tokens) != 3 {
			t.Fatalf("%s: need type declaration on first line", test.name)
		}
		g.generate(tokens[1])
		got := string(g.format())
		if got != test.output {
			println(got)
			// return
			t.Errorf("%s: got(%d)\n====\n%q====\nexpected(%d)\n====%q", test.name, len(got), got, len(test.output), test.output)
		}
	}
}
