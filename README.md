# analysis

Static analysis of Go code through golang.org/x/tools/go/analysis

## How to run

```bash
$ go get github.com/cederstone/analysis/...
$ enum github.com/...
$ keyedlit github.com/...
```

## Passes

### keyedlit

The `keyedlit` pass flags places where variables are declared using keyed
composite literals without specifying `Timeout` or `KeepAlive`-related fields.

For example, the following code will be reported as bad.

```go
package main

import "net/http"

func main() {
	c := &http.Client{
		Transport: http.DefaultTransport,
	}
}
```

The `*http.Client` type has a `Timeout` field which, if left unset, defaults to
infinity. A timeout of infinity is a poor default and has been the cause of
several production issues.

Timeouts and KeepAlives should be carefully thought through.

Additionally, the keyedlit pass can be run in `strict` mode by setting the
`strict` flag to `true`. In this mode all fields must be specified when
creating a keyed composite literal. This means that all fields must be
considered when inintializing a new variable using keyed composite literal
syntax.

The `strict` flag guards against fields where the zero-value is a poor
default. This is especially useful when updating dependencies, where
behaviourally significant fields have been added to existing structs.

```go
// BAD
func main() {
	c := &http.Client{
		Transport: http.DefaultTransport,
	}
}

// GOOD
func main() {
	c := &http.Client{
		Transport: http.DefaultTransport,
		Timeout: 0,
		KeepAlive: 0,
	}
}
```

### enum

Go vaguely supports enums through the following `const`/`iota` pattern:

```go
package main

type MyEnum int

const (
	MyEnum1 MyEnum = iota
	MyEnum2
	MyEnum3
)

func main() {
	val := getEnum(...)
	switch val {
	case MyEnum1:
		// do something for MyEnum1
	case MyEnum2, MyEnum3:
		// do something for MyEnum2 or MyEnum3
	}
}
```

However, adding a new `MyEnum4` member to the enum requires one to find all
switch statements in your codebase and ensure that they cover all cases. Even
more, if a dependency is updated one needs to carefully determine whether the
dependency has expanded its list of enum members to ensure one's own switch
statements are still total.

Missing switch statements when expanding an enum is a common cause of subtle
bugs. This pass helps avoid them.

The `enum` pass considers a type an `enum` if

* its base type is `int`
* its members are defined in a const block
* its const members are all defined using the iota pattern

```go
type MyEnum int

const (
	MyEnum1 MyEnum = iota
	MyEnum2
	MyEnum3
)

// BAD
func bad() {
	val := getEnum(...)
	switch val {
	case MyEnum1:
		return
	default:
	}
}

// GOOD
func good() {
	val := getEnum(...)
	switch val {
	case MyEnum1, MyEnum2:
		return
	case MyEnum3:
	}
}
```

### nakedreturn

If a function has named return values Go let's you omit their names when
calling `return` in that function. However, doing so increases the burden of
comprehension on the reader of the code as the code is less explicit and harder
to reason about in reverse. The benefits to so-called "naked returns" rarely
outweigh the cost. This pass flags naked return statements as errors.

```go
// BAD
func foo() (n int) {
	n = 3
	return
}

// GOOD
func foo() (n int) {
	n = 3
	return n
}
```

### union

Go supports unions by having an exported interface contain an unexprted 'tag'
method. By adding compile-time type assertions (i.e., `var _ Interface =
new(Obj)` this allows one to ensure that all members of a union satisfy the
union's interface. This strategy is used in the `go/ast` package among other
places, where all types that implement the `ast.Expr` interface must implement
the `exprNode()` method:
https://golang.org/src/go/ast/ast.go?s=1432:1473#L31. This is a useful trick
for implementing closed unions.

The `union` pass checks that whenever there is a type switch on a variable of
the union interface type, all values of that type are explicitly handled in
`case`-statements.

The `union` pass treats any exported interface that includes an unexported
method that has no parameters and returns no values as a tagged union.

The pass checks imported packages and is aware of type aliases.

```go
type Letter interface {
    String() string
	isLetter()
}

type A struct{}

func (*A) String() string { return "a" }
func (*A) isLetter() {}

type B struct{}

func (*B) String() string { return "a" }
func (*B) isLetter() {}

// BAD
func bad() {
	letter := getLetter()
	switch letter.(type) {
	case *A:
		fmt.Println("Yay we have an A!")
	default:
	}
}

// GOOD
func good() {
	letter := getLetter()
	switch letter.(type) {
	case *A:
		fmt.Println("Yay we have an A!")
	case *B:
		fmt.Println("Yay we have a B!")
	default:
	}
}
```
