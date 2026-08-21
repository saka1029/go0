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
	m := "TestEval"
	e := globalEnv()
	assertEquals(t, m, "abc", "abc", eval("abc", e))
	assertEquals(t, m, "123", 123, eval(123, e))
	assertEquals(t, m, "true", true, eval(true, e))
	assertEquals(t, m, "'a", sym("a"), eval(quote(sym("a")), e).(Symbol))
	assertEquals(t, m, "'(1 2 3 4)", list(1, 2, 3, 4), eval(quote(list(1, 2, 3, 4)), e))
	assertEquals(t, m, "(car '(1 2 3 4))", 1, eval(list(sym("car"), quote(list(1, 2, 3, 4))), e))
	assertEquals(t, m, "(cdr '(1 2 3 4))", list(2, 3, 4), eval(list(sym("cdr"), quote(list(1, 2, 3, 4))), e))
	assertEquals(t, m, "(cons 'a '(1 2 3 4))", list(sym("a"), 1, 2, 3, 4), eval(list(sym("cons"), quote(sym("a")), quote(list(1, 2, 3, 4))), e))
	assertEquals(t, m, "(cons 'a 'b)", cons(sym("a"), sym("b")), eval(list(sym("cons"), quote(sym("a")), quote(sym("b"))), e).(Cons))
}

func TestArithmetic(t *testing.T) {
	m := "TestArithmetic"
	e := globalEnv()
	assertEquals(t, m, "(+)", 0, eval(list(sym("+")), e))
	assertEquals(t, m, "(+ 1)", 1, eval(list(sym("+"), 1), e))
	assertEquals(t, m, "(+ 1 2)", 3, eval(list(sym("+"), 1, 2), e))
	assertEquals(t, m, "(+ 1 2 3 4)", 10, eval(list(sym("+"), 1, 2, 3, 4), e))
	assertEquals(t, m, "(+ 1 2 (+ 3 4) 5)", 15, eval(list(sym("+"), 1, 2, list(sym("+"), 3, 4), 5), e))
	assertEquals(t, m, "(-)", 0, eval(list(sym("-")), e))
	assertEquals(t, m, "(- 1)", -1, eval(list(sym("-"), 1), e))
	assertEquals(t, m, "(- 3 2)", 1, eval(list(sym("-"), 3, 2), e))
	assertEquals(t, m, "(- 1 2 3 4)", -8, eval(list(sym("-"), 1, 2, 3, 4), e))
	assertEquals(t, m, "(*)", 1, eval(list(sym("*")), e))
	assertEquals(t, m, "(* 3)", 3, eval(list(sym("*"), 3), e))
	assertEquals(t, m, "(* 2 3)", 6, eval(list(sym("*"), 2, 3), e))
	assertEquals(t, m, "(* 1 2 3 4)", 24, eval(list(sym("*"), 1, 2, 3, 4), e))
	assertEquals(t, m, "(* 1 2 (+ 3 4) 5)", 70, eval(list(sym("*"), 1, 2, list(sym("+"), 3, 4), 5), e))
	assertEquals(t, m, "(/)", 1, eval(list(sym("/")), e))
	assertEquals(t, m, "(/ 10)", 0, eval(list(sym("/"), 10), e))
	assertEquals(t, m, "(/ 3 2)", 1, eval(list(sym("/"), 3, 2), e))
	assertEquals(t, m, "(/ 1001 7 11)", 13, eval(list(sym("/"), 1001, 7, 11), e))
}

func TestPrint(t *testing.T) {
	m := "TestPrint"
	assertEquals(t, m, "abc", "abc", print(sym("abc")))
	assertEquals(t, m, "123", "123", print(123))
	assertEquals(t, m, "true", "true", print(true))
	assertEquals(t, m, "()", "()", print(list()))
	assertEquals(t, m, "(1 a)", "(1 a)", print(list(1, sym("a"))))
	assertEquals(t, m, "'(1 a)", "'(1 a)", print(quote(list(1, sym("a")))))
	assertEquals(t, m, "(1 . a)", "(1 . a)", print(cons(1, sym("a"))))
	assertEquals(t, m, "(quote . a)", "(quote . a)", print(cons(sym("quote"), sym("a"))))
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
	m := "TestSize"
	assertEquals(t, m, "Sizeof(2)", 8, unsafe.Sizeof(2))
	assertEquals(t, m, "Sizeof(\"abcd\")", 16, unsafe.Sizeof("abcd"))
	assertEquals(t, m, "Sizeof(Cons0Interface{})", 0, unsafe.Sizeof(Cons0Interface{}))
	assertEquals(t, m, "Sizeof(Cons2Interface{})", 32, unsafe.Sizeof(Cons2Interface{}))
	assertEquals(t, m, "Sizeof(Cons2Interface{2, 3})", 32, unsafe.Sizeof(Cons2Interface{2, 3}))
	assertEquals(t, m, "Sizeof(&Cons2Interface{2, 3})", 8, unsafe.Sizeof(&Cons2Interface{2, 3}))
	assertEquals(t, m, "Sizeof(Cons2Interface{\"a\", \"b\"})", 32, unsafe.Sizeof(Cons2Interface{"a", "b"}))
	assertEquals(t, m, "Sizeof(Cons2Interface{\"a\", Cons2Interface{\"b\", \"c\"}})", 32, unsafe.Sizeof(Cons2Interface{"a", Cons2Interface{"b", "c"}}))
	assertEquals(t, m, "Sizeof(Cons3Interface{})", 48, unsafe.Sizeof(Cons3Interface{}))
	assertEquals(t, m, "Sizeof(Cons3Interface{2, 3, 4})", 48, unsafe.Sizeof(Cons3Interface{2, 3, 4}))
}

func TestDefine(t *testing.T) {
	m := "TestDefine"
	e := globalEnv()
	assertEquals(t, m, "(define a (+ 1 2))", nil, eval(list(sym("define"), sym("a"), list(sym("+"), 1, 2)), e))
	assertEquals(t, m, "a", 3, eval(sym("a"), e))

}

func TestLambda(t *testing.T) {
	m := "TestLambda"
	e := globalEnv()
	lambda := (list(sym("lambda"), list(sym("a"), sym("b")), list(sym("+"), sym("a"), sym("b"))))
	body := list(lambda, 1, 2)
	assertEquals(t, m, "((lambda (a b) (+ a b)) 1 2)", 3, eval(body, e))
	assertEquals(t, m, "(define add (lambda (a b) (+ a b)) 1 2)", nil, eval(list(sym("define"), sym("add"), lambda), e))
	assertEquals(t, m, "(add 3 4)", 7, eval(list(sym("add"), 3, 4), e))
}
