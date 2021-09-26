# enumer

[![Build Status](https://github.com/vearutop/enumer/workflows/test-unit/badge.svg)](https://github.com/vearutop/enumer/actions?query=branch%3Amaster+workflow%3Atest-unit)
[![Coverage Status](https://codecov.io/gh/vearutop/enumer/branch/master/graph/badge.svg)](https://codecov.io/gh/vearutop/enumer)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/vearutop/enumer)
[![Time Tracker](https://wakatime.com/badge/github/vearutop/enumer.svg)](https://wakatime.com/badge/github/vearutop/enumer)
![Code lines](https://sloc.xyz/github/vearutop/enumer/?category=code)
![Comments](https://sloc.xyz/github/vearutop/enumer/?category=comments)

A Go tool to generate enumerations of type constants.

## Fork

This is a fork of [stringer](https://github.com/golang/tools/tree/master/cmd/stringer)
changed to generate a map of names and values instead of `String() string`.

## Installation

```
go install github.com/vearutop/enumer@latest
```

## Usage

Having `day.go` with this contents:

```go
//go:generate enumer -type=Day

type Day int
const (
	Monday Day = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)
```

Alternatively, in go1.17 and newer, you can set up generation with
```
//go:generate go run github.com/vearutop/enumer@latest -type=Day
```

After running `go generate` you will get `day_enum.go` with 

```go
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
```

that implements [`jsonschema.Enum`](https://pkg.go.dev/github.com/swaggest/jsonschema-go#Enum).