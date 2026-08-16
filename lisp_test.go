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

func TestEval(t *testing.T) {
	c := context()
	evalTest(t, c, list(1, 2, 3, 4), list(Symbol("quote"), list(1, 2, 3, 4)))
	evalTest(t, c, 1, list(Symbol("car"), list(Symbol("quote"), list(1, 2, 3, 4))))
	evalTest(t, c, list(2, 3, 4), list(Symbol("cdr"), list(Symbol("quote"), list(1, 2, 3, 4))))
}
