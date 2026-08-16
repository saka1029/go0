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

func stringLength(b byte) int {
	switch {
	case b&0b10000000 == 0b00000000:
		return 1
	case b&0b11100000 == 0b11000000:
		return 2
	case b&0b11110000 == 0b11100000:
		return 3
	case b&0b11111000 == 0b11110000:
		return 4
	default:
		return -999
	}
}

type StringReader struct {
	s     string
	index int
}

func newStringReader(s string) *StringReader {
	return &StringReader{s, 0}
}

func (this *StringReader) get() string {
	if this.index >= len(this.s) {
		return ""
	}
	length := stringLength(this.s[this.index])
	result := this.s[this.index : this.index+length]
	this.index += length
	return result
}

func stringTest(t *testing.T, s string, expected int) {
	actual := stringLength(s[0])
	fmt.Printf("%s length=%d\n", s, actual)
	if expected != actual {
		t.Errorf("TestString %s -> %d not %d", s, actual, expected)
	}
}
func stringTest2(t *testing.T, sr *StringReader, expected string) {
	actual := sr.get()
	fmt.Printf("%s expected=%s\n", actual, expected)
	if expected != actual {
		t.Errorf("TestString %s not %s", actual, expected)
	}
}

func TestString(t *testing.T) {
	stringTest(t, "a", 1)
	stringTest(t, "ä", 2)
	stringTest(t, "あ", 3)
	stringTest(t, "𩸽", 4)
	sr := newStringReader("aäあ𩸽")
	stringTest2(t, sr, "a")
	stringTest2(t, sr, "ä")
	stringTest2(t, sr, "あ")
	stringTest2(t, sr, "𩸽")
	stringTest2(t, sr, "")
}
