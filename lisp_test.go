package main

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func assertEquals[T any](t *testing.T, method string, name string, expected T, actual T) {
	fmt.Println(method, name)
	if !reflect.DeepEqual(actual, expected) {
		t.Error("  -> not equal expected=", expected, "actual=", actual)
	}
}

func quote(e Evaluable) Evaluable {
	return list(sym("quote"), e)
}

func TestEval(t *testing.T) {
	e := env()
	assertEquals(t, "TestEval", "abc", "abc", eval("abc", e))
	assertEquals(t, "TestEval", "123", 123, eval(123, e))
	assertEquals(t, "TestEval", "true", true, eval(true, e))
	assertEquals(t, "TestEval", "'a", sym("a"), eval(quote(sym("a")), e).(Symbol))
	assertEquals(t, "TestEval", "'(1 2 3 4)", list(1, 2, 3, 4), eval(quote(list(1, 2, 3, 4)), e))
	assertEquals(t, "TestEval", "(car '(1 2 3 4))", 1, eval(list(sym("car"), quote(list(1, 2, 3, 4))), e))
	assertEquals(t, "TestEval", "(cdr '(1 2 3 4))", list(2, 3, 4), eval(list(sym("cdr"), quote(list(1, 2, 3, 4))), e))
	assertEquals(t, "TestEval", "(cons 'a '(1 2 3 4))", list(sym("a"), 1, 2, 3, 4), eval(list(sym("cons"), quote(sym("a")), quote(list(1, 2, 3, 4))), e))
	assertEquals(t, "TestEval", "(cons 'a 'b)", cons(sym("a"), sym("b")), eval(list(sym("cons"), quote(sym("a")), quote(sym("b"))), e).(Cons))
}

func TestArithmetic(t *testing.T) {
	e := env()
	assertEquals(t, "TestArithmetic", "(+)", 0, eval(list(sym("+")), e))
	assertEquals(t, "TestArithmetic", "(+ 1)", 1, eval(list(sym("+"), 1), e))
	assertEquals(t, "TestArithmetic", "(+ 1 2)", 3, eval(list(sym("+"), 1, 2), e))
	assertEquals(t, "TestArithmetic", "(+ 1 2 3 4)", 10, eval(list(sym("+"), 1, 2, 3, 4), e))
	assertEquals(t, "TestArithmetic", "(+ 1 2 (+ 3 4) 5)", 15, eval(list(sym("+"), 1, 2, list(sym("+"), 3, 4), 5), e))
	assertEquals(t, "TestArithmetic", "(-)", 0, eval(list(sym("-")), e))
	assertEquals(t, "TestArithmetic", "(- 1)", -1, eval(list(sym("-"), 1), e))
	assertEquals(t, "TestArithmetic", "(- 3 2)", 1, eval(list(sym("-"), 3, 2), e))
	assertEquals(t, "TestArithmetic", "(- 1 2 3 4)", -8, eval(list(sym("-"), 1, 2, 3, 4), e))
	assertEquals(t, "TestArithmetic", "(*)", 1, eval(list(sym("*")), e))
	assertEquals(t, "TestArithmetic", "(* 3)", 3, eval(list(sym("*"), 3), e))
	assertEquals(t, "TestArithmetic", "(* 2 3)", 6, eval(list(sym("*"), 2, 3), e))
	assertEquals(t, "TestArithmetic", "(* 1 2 3 4)", 24, eval(list(sym("*"), 1, 2, 3, 4), e))
	assertEquals(t, "TestArithmetic", "(* 1 2 (+ 3 4) 5)", 70, eval(list(sym("*"), 1, 2, list(sym("+"), 3, 4), 5), e))
	assertEquals(t, "TestArithmetic", "(/)", 1, eval(list(sym("/")), e))
	assertEquals(t, "TestArithmetic", "(/ 10)", 0, eval(list(sym("/"), 10), e))
	assertEquals(t, "TestArithmetic", "(/ 3 2)", 1, eval(list(sym("/"), 3, 2), e))
	assertEquals(t, "TestArithmetic", "(/ 1001 7 11)", 13, eval(list(sym("/"), 1001, 7, 11), e))
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

func TestSize(t *testing.T) {
	assertEquals(t, "TestSize", "Sizeof(2)", 8, unsafe.Sizeof(2))
	assertEquals(t, "TestSize", "Sizeof(\"abcd\")", 16, unsafe.Sizeof("abcd"))
	assertEquals(t, "TestSize", "Sizeof(Cons0Interface{})", 0, unsafe.Sizeof(Cons0Interface{}))
	assertEquals(t, "TestSize", "Sizeof(Cons2Interface{2, 3})", 32, unsafe.Sizeof(Cons2Interface{2, 3}))
	assertEquals(t, "TestSize", "Sizeof(&Cons2Interface{2, 3})", 8, unsafe.Sizeof(&Cons2Interface{2, 3}))
	assertEquals(t, "TestSize", "Sizeof(Cons2Interface{\"a\", \"b\"})", 32, unsafe.Sizeof(Cons2Interface{"a", "b"}))
	assertEquals(t, "TestSize", "Sizeof(Cons2Interface{\"a\", Cons2Interface{\"b\", \"c\"}})", 32, unsafe.Sizeof(Cons2Interface{"a", Cons2Interface{"b", "c"}}))
	assertEquals(t, "TestSize", "Sizeof(Cons3Interface{2, 3, 4})", 48, unsafe.Sizeof(Cons3Interface{2, 3, 4}))
}
