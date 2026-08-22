package main

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func assertEquals[T any](t *testing.T, name string, expected T, actual T) {
	fmt.Println(name)
	if !reflect.DeepEqual(actual, expected) {
		t.Error("  -> not equal expected=", expected, "actual=", actual)
	}
}

func quote(e Evaluable) Evaluable {
	return list(sym("quote"), e)
}

func TestEval(t *testing.T) {
	e := globalEnv()
	assertEquals(t, "abc", "abc", eval("abc", e))
	assertEquals(t, "123", 123, eval(123, e))
	assertEquals(t, "true", true, eval(true, e))
	assertEquals(t, "'a", sym("a"), eval(quote(sym("a")), e).(Symbol))
	assertEquals(t, "'(1 2 3 4)", list(1, 2, 3, 4), eval(quote(list(1, 2, 3, 4)), e))
	assertEquals(t, "(car '(1 2 3 4))", 1, eval(list(sym("car"), quote(list(1, 2, 3, 4))), e))
	assertEquals(t, "(cdr '(1 2 3 4))", list(2, 3, 4), eval(list(sym("cdr"), quote(list(1, 2, 3, 4))), e))
	assertEquals(t, "(cons 'a '(1 2 3 4))", list(sym("a"), 1, 2, 3, 4), eval(list(sym("cons"), quote(sym("a")), quote(list(1, 2, 3, 4))), e))
	assertEquals(t, "(cons 'a 'b)", cons(sym("a"), sym("b")), eval(list(sym("cons"), quote(sym("a")), quote(sym("b"))), e).(Cons))
}

func TestArithmetic(t *testing.T) {
	e := globalEnv()
	assertEquals(t, "(+)", 0, eval(list(sym("+")), e))
	assertEquals(t, "(+ 1)", 1, eval(list(sym("+"), 1), e))
	assertEquals(t, "(+ 1 2)", 3, eval(list(sym("+"), 1, 2), e))
	assertEquals(t, "(+ 1 2 3 4)", 10, eval(list(sym("+"), 1, 2, 3, 4), e))
	assertEquals(t, "(+ 1 2 (+ 3 4) 5)", 15, eval(list(sym("+"), 1, 2, list(sym("+"), 3, 4), 5), e))
	assertEquals(t, "(-)", 0, eval(list(sym("-")), e))
	assertEquals(t, "(- 1)", -1, eval(list(sym("-"), 1), e))
	assertEquals(t, "(- 3 2)", 1, eval(list(sym("-"), 3, 2), e))
	assertEquals(t, "(- 1 2 3 4)", -8, eval(list(sym("-"), 1, 2, 3, 4), e))
	assertEquals(t, "(*)", 1, eval(list(sym("*")), e))
	assertEquals(t, "(* 3)", 3, eval(list(sym("*"), 3), e))
	assertEquals(t, "(* 2 3)", 6, eval(list(sym("*"), 2, 3), e))
	assertEquals(t, "(* 1 2 3 4)", 24, eval(list(sym("*"), 1, 2, 3, 4), e))
	assertEquals(t, "(* 1 2 (+ 3 4) 5)", 70, eval(list(sym("*"), 1, 2, list(sym("+"), 3, 4), 5), e))
	assertEquals(t, "(/)", 1, eval(list(sym("/")), e))
	assertEquals(t, "(/ 10)", 0, eval(list(sym("/"), 10), e))
	assertEquals(t, "(/ 3 2)", 1, eval(list(sym("/"), 3, 2), e))
	assertEquals(t, "(/ 1001 7 11)", 13, eval(list(sym("/"), 1001, 7, 11), e))
}

func TestPrint(t *testing.T) {
	assertEquals(t, "abc", "abc", print(sym("abc")))
	assertEquals(t, "123", "123", print(123))
	assertEquals(t, "true", "true", print(true))
	assertEquals(t, "()", "()", print(list()))
	assertEquals(t, "(1 a)", "(1 a)", print(list(1, sym("a"))))
	assertEquals(t, "'(1 a)", "'(1 a)", print(quote(list(1, sym("a")))))
	assertEquals(t, "(1 . a)", "(1 . a)", print(cons(1, sym("a"))))
	assertEquals(t, "(quote . a)", "(quote . a)", print(cons(sym("quote"), sym("a"))))
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
	assertEquals(t, "Sizeof(2)", 8, unsafe.Sizeof(2))
	assertEquals(t, "Sizeof(\"abcd\")", 16, unsafe.Sizeof("abcd"))
	assertEquals(t, "Sizeof(Cons0Interface{})", 0, unsafe.Sizeof(Cons0Interface{}))
	assertEquals(t, "Sizeof(Cons2Interface{})", 32, unsafe.Sizeof(Cons2Interface{}))
	assertEquals(t, "Sizeof(Cons2Interface{2, 3})", 32, unsafe.Sizeof(Cons2Interface{2, 3}))
	assertEquals(t, "Sizeof(&Cons2Interface{2, 3})", 8, unsafe.Sizeof(&Cons2Interface{2, 3}))
	assertEquals(t, "Sizeof(Cons2Interface{\"a\", \"b\"})", 32, unsafe.Sizeof(Cons2Interface{"a", "b"}))
	assertEquals(t, "Sizeof(Cons2Interface{\"a\", Cons2Interface{\"b\", \"c\"}})", 32, unsafe.Sizeof(Cons2Interface{"a", Cons2Interface{"b", "c"}}))
	assertEquals(t, "Sizeof(Cons3Interface{})", 48, unsafe.Sizeof(Cons3Interface{}))
	assertEquals(t, "Sizeof(Cons3Interface{2, 3, 4})", 48, unsafe.Sizeof(Cons3Interface{2, 3, 4}))
}

func TestDefine(t *testing.T) {
	e := globalEnv()
	assertEquals(t, "(define a (+ 1 2))", nil, eval(list(sym("define"), sym("a"), list(sym("+"), 1, 2)), e))
	assertEquals(t, "a", 3, eval(sym("a"), e))

}

func TestIsSymbolRest(t *testing.T) {
	assertEquals(t, "a", true, isSymbolRest('a'))
	assertEquals(t, "space", false, isSymbolRest(' '))
	assertEquals(t, ".", true, isSymbolRest('.'))
}

func TestLambda(t *testing.T) {
	e := globalEnv()
	lambda := (list(sym("lambda"), list(sym("a"), sym("b")), list(sym("+"), sym("a"), sym("b"))))
	body := list(lambda, 1, 2)
	assertEquals(t, "((lambda (a b) (+ a b)) 1 2)", 3, eval(body, e))
	assertEquals(t, "(define add (lambda (a b) (+ a b)) 1 2)", nil, eval(list(sym("define"), sym("add"), lambda), e))
	assertEquals(t, "(define a 100)", nil, eval(list(sym("define"), sym("a"), 100), e))
	assertEquals(t, "(add (+ a 1) 4)", 105, eval(list(sym("add"), list(sym("+"), sym("a"), 1), 4), e))
}

func TestReadSymbol(t *testing.T) {
	r := NewReader("  abc  𩸽  dot.dot")
	assertEquals(t, "abc", sym("abc"), r.read().(Symbol))
	assertEquals(t, "𩸽", sym("𩸽"), r.read().(Symbol))
	assertEquals(t, "dot.dot", sym("dot.dot"), r.read().(Symbol))
	assertEquals(t, "EOF", EOF, r.read().(rune))
}

func TestReadInt(t *testing.T) {
	r := NewReader("  123  -456 +78 ")
	assertEquals(t, "123", 123, r.read().(int))
	assertEquals(t, "-456", -456, r.read().(int))
	assertEquals(t, "+78", 78, r.read().(int))
	assertEquals(t, "EOF", EOF, r.read().(rune))
}

func TestReadList(t *testing.T) {
	r := NewReader("  ( abc 123 )  ")
	assertEquals(t, "(abc 123)", list(sym("abc"), 123), r.read())
	assertEquals(t, "EOF", EOF, r.read().(rune))
}

func TestReadQuote(t *testing.T) {
	r := NewReader("'(abc 123)")
	assertEquals(t, "123", list(sym("quote"), list(sym("abc"), 123)), r.read())
	assertEquals(t, "EOF", EOF, r.read().(rune))
}

func TestReadQuote2(t *testing.T) {
	r := NewReader("('abc)")
	assertEquals(t, "('abc)", list(list(sym("quote"), sym("abc"))), r.read())
	assertEquals(t, "EOF", EOF, r.read().(rune))
}

func TestReadQuote3(t *testing.T) {
	r := NewReader("'(a . b)")
	assertEquals(t, "'(a . b))", list(sym("quote"), cons(sym("a"), sym("b"))), r.read())
	assertEquals(t, "EOF", EOF, r.read().(rune))
}

func rep(source string, e *Env) string {
	return print(eval(NewReader(source).read(), e))
}

func TestReadEvalPrint(t *testing.T) {
	e := globalEnv()
	assertEquals(t, "'a", "a", rep("'a", e))
	assertEquals(t, "123", "123", rep("123", e))
	assertEquals(t, "'(a . b)", "(a . b)", rep("'(a . b)", e))
	assertEquals(t, "(cons 1 2)", "(1 . 2)", rep("(cons 1 2)", e))
	assertEquals(t, "(cons 'a 'b)", "(a . b)", rep("(cons 'a 'b)", e))
	assertEquals(t, "(car '(a b))", "a", rep("(car '(a b))", e))
	assertEquals(t, "(cdr '(a b))", "(b)", rep("(cdr '(a b))", e))
	assertEquals(t, "((lambda (a b) (+ a b)) 1 2)", "3", rep("((lambda (a b) (+ a b)) 1 2))", e))
	assertEquals(t, "'(define add (lambda (a b) (+ a b)))", "()", rep("(define add (lambda (a b) (+ a b)))", e))
	assertEquals(t, "(add '7 '8)", "15", rep("(add '7 '8)", e))
}
