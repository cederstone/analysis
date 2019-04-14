package main

import(
	"time"
)

type Foo1 struct {
	Timeout time.Duration
	Bar struct{}
}

type Foo2 struct {
	FooTimeout time.Duration
	Bar struct{}
}

type Foo3 struct {
	KeepAlive time.Duration
	Bar struct{}
}

type Foo4 struct {
	FooKeepAlive time.Duration
	Bar struct{}
}

type Foo5 struct {
	Timeout time.Duration
	KeepAlive time.Duration
	Bar struct{}
}

func main() {
	a := Foo1{ // want "unspecified field Timeout of Foo1"
		Bar: struct{}{},
	}
	a = Foo1{} // not specifying any fields
	_ = a
	b := Foo2{ // want "unspecified field FooTimeout of Foo2"
		Bar: struct{}{},
	}
	b = Foo2{} // not specifying any fields
	_ = b
	c := Foo3{ // want "unspecified field KeepAlive of Foo3"
		Bar: struct{}{},
	}
	c = Foo3{} // not specifying any fields
	_ = c
	d := Foo4{ // want "unspecified field FooKeepAlive of Foo4"
		Bar: struct{}{},
	}
	d = Foo4{} // not specifying any fields
	_ = d
	e := Foo5{ // want "unspecified field Timeout of Foo5" "unspecified field KeepAlive of Foo5"
		Bar: struct{}{},
	}
	e = Foo5{} // not specifying any fields
	_ = e
}
