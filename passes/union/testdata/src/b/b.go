package b

import "a"

type Bar = a.Foo

type Bar1 = *a.Member1
type Bar2 = *a.Member2

func main() {
	var f a.Foo
	switch f.(type) { // want "non-total type switch over union: "
	case *a.Member1:
		return
	}

	switch f.(type) { // total using default clause
	case *a.Member1:
		return
	default:
	}

	switch b := f.(type) { // want "non-total type switch over union: "
	case *a.Member1:
		_ = b
		return
	}

	switch f.(type) { // is total
	case *a.Member1, *a.Member2:
		return
	}

	switch f.(type) { // want "non-total type switch over union: "
	case *a.Member1:
		return
	}

	var g Bar
	switch g.(type) { // One type alias used
	case *a.Member1, Bar2:
		return
	}

	switch g.(type) { // want "non-total type switch over union: "
	case Bar2:
		return
	}

	switch g.(type) { // total using type aliases
	case Bar1, Bar2:
		return
	}

	switch g.(type) { // total using type aliases and default clause
	case Bar1:
		return
	default:
		return
	}

	switch g.(type) { // total using type aliases split cases
	case Bar1:
		return
	case Bar2:
		return
	}

	// not a union!
	x := map[string]interface{}{"a": 1}
	switch x["a"].(type) {
	case int:
		return
	}
}
