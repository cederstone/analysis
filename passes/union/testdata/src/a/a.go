package a

type Foo interface { // want Foo:"Union"
	tag()
}

type Member1 struct{}

func (*Member1) tag() {}

type Member2 struct{}

func (*Member2) tag() {}

func main() {
	var a Foo

	switch a.(type) { // want "non-total type switch over union: "
	case *Member1:
		return
	}

	switch b := a.(type) { // want "non-total type switch over union: "
	case *Member1:
		_ = b
		return
	}

	switch a.(type) { // is total
	case *Member1, *Member2:
		return
	}

	switch a.(type) { // want "non-total type switch over union: "
	case *Member1:
		return
	default:
	}

}
