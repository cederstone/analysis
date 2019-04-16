package test

func test1() (foo int) {
	return // want "return values not explicitly specified"
}

func test2() (foo int) {
	if false {
		return 1 // not naked
	}
	if false {
		return // want "return values not explicitly specified"
	}
	return // want "return values not explicitly specified"
}

func test3() func() int {
	return func() (n int) {
		return // want "return values not explicitly specified"
	}
}

func test4() func() int {
	return func() (n int) {
		func() {
			return // not a naked return
		}()
		return // want "return values not explicitly specified"
	}
}
