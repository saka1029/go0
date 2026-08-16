package main

import (
	"fmt"
	"testing"
)

func evalTest(t *testing.T, c *Context, expected Evaluable, e Evaluable) {
	fmt.Println("TestEval:", print(e))
	actual := eval(e, c)
	if actual != expected {
		t.Errorf("eval %s -> %s not %s", print(e), print(actual), print(expected))
	}
}

func quote(e Evaluable) Evaluable {
	return list(sym("quote"), e)
}

func TestEval(t *testing.T) {
	c := context()
	evalTest(t, c, "abc", "abc")
	evalTest(t, c, 123, 123)
	evalTest(t, c, true, true)
	evalTest(t, c, sym("a"), quote(sym("a")))
	evalTest(t, c, list(1, 2, 3, 4), quote(list(1, 2, 3, 4)))
	evalTest(t, c, 1, list(sym("car"), quote(list(1, 2, 3, 4))))
	evalTest(t, c, list(2, 3, 4), list(sym("cdr"), quote(list(1, 2, 3, 4))))
	evalTest(t, c, list(sym("a"), 1, 2, 3, 4), list(sym("cons"), quote(sym("a")), quote(list(1, 2, 3, 4))))
	evalTest(t, c, cons(sym("a"), sym("b")), list(sym("cons"), quote(sym("a")), quote(sym("b"))))
}
func TestArithmetic(t *testing.T) {
	c := context()
	evalTest(t, c, 0, list(sym("+")))
	evalTest(t, c, 1, list(sym("+"), 1))
	evalTest(t, c, 3, list(sym("+"), 1, 2))
	evalTest(t, c, 10, list(sym("+"), 1, 2, 3, 4))
	evalTest(t, c, 15, list(sym("+"), 1, 2, list(sym("+"), 3, 4), 5))
	evalTest(t, c, 0, list(sym("-")))
	evalTest(t, c, -1, list(sym("-"), 1))
	evalTest(t, c, 1, list(sym("-"), 3, 2))
	evalTest(t, c, -8, list(sym("-"), 1, 2, 3, 4))
	evalTest(t, c, 1, list(sym("*")))
	evalTest(t, c, 3, list(sym("*"), 3))
	evalTest(t, c, 6, list(sym("*"), 2, 3))
	evalTest(t, c, 24, list(sym("*"), 1, 2, 3, 4))
	evalTest(t, c, 70, list(sym("*"), 1, 2, list(sym("+"), 3, 4), 5))
	evalTest(t, c, 1, list(sym("/")))
	evalTest(t, c, 0, list(sym("/"), 10))
	evalTest(t, c, 1, list(sym("/"), 3, 2))
	evalTest(t, c, 13, list(sym("/"), 1001, 7, 11))
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

func stringReaderTest(t *testing.T, s string, expected int) {
	actual := stringLength(s[0])
	fmt.Printf("%s length=%d\n", s, actual)
	if expected != actual {
		t.Errorf("TestString %s -> %d not %d", s, actual, expected)
	}
}
func stringReaderTest2(t *testing.T, sr *StringReader, expected string) {
	actual := sr.Get()
	fmt.Printf("%s expected=%s\n", actual, expected)
	if expected != actual {
		t.Errorf("TestString %s not %s", actual, expected)
	}
}

func TestStringReader(t *testing.T) {
	stringReaderTest(t, "a", 1)
	stringReaderTest(t, "ä", 2)
	stringReaderTest(t, "あ", 3)
	stringReaderTest(t, "𩸽", 4)
	sr := NewStringReader("aäあ𩸽")
	stringReaderTest2(t, sr, "a")
	stringReaderTest2(t, sr, "ä")
	stringReaderTest2(t, sr, "あ")
	stringReaderTest2(t, sr, "𩸽")
	stringReaderTest2(t, sr, "")
}
