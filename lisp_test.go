package main

import (
	"fmt"
	"testing"
)

func evalTest(t *testing.T, c *Context, expected string, e Evaluable) {
	fmt.Println("TestEval:", print(e))
	actual := eval(e, c)
	actualStr := print(eval(e, c))
	if actualStr != expected {
		t.Errorf("eval %s -> %s not %s", print(e), print(actual), expected)
	}
}

func TestEval(t *testing.T) {
	c := context()
	evalTest(t, c, "(1 2 3 4)", list(Symbol("quote"), list(1, 2, 3, 4)))
	evalTest(t, c, "1", list(Symbol("car"), list(Symbol("quote"), list(1, 2, 3, 4))))
}
