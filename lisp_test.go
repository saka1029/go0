package main

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func evalTest(t *testing.T, c *Env, expected Evaluable, e Evaluable) {
	fmt.Println("TestEval:", print(e))
	actual := eval(e, c)
	// if actual != expected {
	// if reflect.DeepEqual(actual, expected) {
	if actual != expected {
		t.Errorf("eval %s -> %s not %s", print(e), print(actual), print(expected))
	}
}

func quote(e Evaluable) Evaluable {
	return list(sym("quote"), e)
}

func TestEval(t *testing.T) {
	e := env()
	evalTest(t, e, "abc", "abc")
	evalTest(t, e, 123, 123)
	evalTest(t, e, true, true)
	evalTest(t, e, sym("a"), quote(sym("a")))
	evalTest(t, e, list(1, 2, 3, 4), quote(list(1, 2, 3, 4)))
	evalTest(t, e, 1, list(sym("car"), quote(list(1, 2, 3, 4))))
	evalTest(t, e, list(2, 3, 4), list(sym("cdr"), quote(list(1, 2, 3, 4))))
	evalTest(t, e, list(sym("a"), 1, 2, 3, 4), list(sym("cons"), quote(sym("a")), quote(list(1, 2, 3, 4))))
	evalTest(t, e, cons(sym("a"), sym("b")), list(sym("cons"), quote(sym("a")), quote(sym("b"))))
}
func TestArithmetic(t *testing.T) {
	e := env()
	evalTest(t, e, 0, list(sym("+")))
	evalTest(t, e, 1, list(sym("+"), 1))
	evalTest(t, e, 3, list(sym("+"), 1, 2))
	evalTest(t, e, 10, list(sym("+"), 1, 2, 3, 4))
	evalTest(t, e, 15, list(sym("+"), 1, 2, list(sym("+"), 3, 4), 5))
	evalTest(t, e, 0, list(sym("-")))
	evalTest(t, e, -1, list(sym("-"), 1))
	evalTest(t, e, 1, list(sym("-"), 3, 2))
	evalTest(t, e, -8, list(sym("-"), 1, 2, 3, 4))
	evalTest(t, e, 1, list(sym("*")))
	evalTest(t, e, 3, list(sym("*"), 3))
	evalTest(t, e, 6, list(sym("*"), 2, 3))
	evalTest(t, e, 24, list(sym("*"), 1, 2, 3, 4))
	evalTest(t, e, 70, list(sym("*"), 1, 2, list(sym("+"), 3, 4), 5))
	evalTest(t, e, 1, list(sym("/")))
	evalTest(t, e, 0, list(sym("/"), 10))
	evalTest(t, e, 1, list(sym("/"), 3, 2))
	evalTest(t, e, 13, list(sym("/"), 1001, 7, 11))
}

func printTest(t *testing.T, expected string, e Evaluable) {
	actual := print(e)
	fmt.Println("TestPrint:", actual)
	if actual != expected {
		t.Errorf("print %s not %s", actual, expected)
	}
}

func TestPrint(t *testing.T) {
	printTest(t, "abc", sym("abc"))
	printTest(t, "123", 123)
	printTest(t, "true", true)
	printTest(t, "()", list())
	printTest(t, "(1 a)", list(1, sym("a")))
	printTest(t, "'(1 a)", quote(list(1, sym("a"))))
	printTest(t, "(1 . a)", cons(1, sym("a")))
	printTest(t, "(quote . a)", cons(sym("quote"), sym("a")))
}

type Cons0Interface struct {
}

type Cons2Interface struct {
	car, cdr Evaluable
}

type Cons3Interface struct {
	car, cdr, other Evaluable
}

func assertEquals[T any](t *testing.T, name string, expected T, actual T) {
	fmt.Println("test:", name)
	if !reflect.DeepEqual(actual, expected) {
		t.Error("  -> not equal expected=", expected, "actual=", actual)
	}
}

func TestSize(t *testing.T) {
	assertEquals(t, "Sizeof(2)", 8, unsafe.Sizeof(2))
	assertEquals(t, "Sizeof(\"abcd\")", 16, unsafe.Sizeof("abcd"))
	assertEquals(t, "Sizeof(Cons0Interface{})", 0, unsafe.Sizeof(Cons0Interface{}))
	assertEquals(t, "Sizeof(Cons2Interface{2, 3})", 32, unsafe.Sizeof(Cons2Interface{2, 3}))
	assertEquals(t, "Sizeof(&Cons2Interface{2, 3})", 8, unsafe.Sizeof(&Cons2Interface{2, 3}))
	assertEquals(t, "Sizeof(Cons2Interface{\"a\", \"b\"})", 32, unsafe.Sizeof(Cons2Interface{"a", "b"}))
	assertEquals(t, "Sizeof(Cons2Interface{\"a\", Cons2Interface{\"b\", \"c\"}})", 32, unsafe.Sizeof(Cons2Interface{"a", Cons2Interface{"b", "c"}}))
	assertEquals(t, "Sizeof(Cons3Interface{2, 3, 4})", 48, unsafe.Sizeof(Cons3Interface{2, 3, 4}))
}
