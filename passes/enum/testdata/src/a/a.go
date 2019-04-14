package main

type Foo int

const (
	Foo1 Foo = iota
	Foo2
)

type Bar uint

const (
	Bar1 Bar = iota
	Bar2
)

type FooA int
type FooB int
type FooC int
type FooD int

const (
	_ FooA = iota
	FooA1
	FooA2

	Bar3 Bar = iota

	FooB1 FooB = iota
	FooB2

	FooC1 FooC = 1

	FooD1 FooD = iota
	FooD2
)

func main() {
	a := Foo1
	switch a { // want "non-total switch over enum"
	case Foo1:
	}
	b := Bar1
	switch b { // not an int
	case Bar1:
	}
	c := FooA1
	switch c { // want "non-total switch over enum"
	case FooA1:
	}
	d := FooA1
	switch c { // anonynmous member is ignored
	case FooA1, FooA2:
	}
	e := FooB1
	switch e { // want "non-total switch over enum"
	case FooB1:
	}
	f := FooC1
	switch f { // not using iota
	case FooC1:
	}
	g := FooD1
	switch g { // fully specified
	case FooD1:
	case FooD2:
	}
}
