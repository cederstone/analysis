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

### nakedreturn

If a function has named return values Go let's you omit their names when
calling `return` in that function. However, doing so increases the burden of
comprehension on the reader of the code as the code is less explicit and harder
to reason about in reverse. The benefits to so-called "naked returns" rarely
outweigh the cost. This pass flags naked return statements as errors.
