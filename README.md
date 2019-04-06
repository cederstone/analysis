# analysis

Static analysis of Go code through golang.org/x/tools/go/analysis

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
