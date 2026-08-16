package main

import (
	"fmt"
	"strings"
)

type Context struct {
	globals map[Symbol]Evaluable
}

func context() *Context {
	return &Context{map[Symbol]Evaluable{}}
}

type Evaluable interface {
}

type Symbol string

type Cons struct {
	car, cdr Evaluable
}
func list(elements ...Evaluable) Evaluable {
	var result Evaluable = nil
	for i := len(elements) - 1; i >= 0; i-- {
		result = Cons{elements[i], result}
	}
	return result
}

func eval(e Evaluable, c *Context) Evaluable {
	switch v := e.(type) {
	case int, string:
		return v
	case Symbol:
		return c.globals[v]
	case Cons:
		return apply(eval(v.car, c), v.cdr, c)
	default:
		panic(fmt.Sprint("unknown type", v))
	}
}

func apply(e Evaluable, args Evaluable, c *Context) Evaluable {
	return nil
}

func printCons(c Cons) string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(print(c.car))
	e := c.cdr;
	for true {
		cons, ok := e.(Cons)
		if !ok {
			break;
		}
		sb.WriteString(" " + print(cons.car))
		e = cons.cdr
	}
	if (e != nil) {
		sb.WriteString(" . ")
		sb.WriteString(print(e))
	}
	sb.WriteString(")")
	return sb.String()
}

func print(e Evaluable) string {
	switch v := e.(type) {
	case int:
		return fmt.Sprint(v)
	case string:
		return v
	case Cons:
		return printCons(v)
	default:
		panic(fmt.Sprint("unknown type", v))
	}
}

func main() {
	var context *Context = context()
	fmt.Println(context)
	context.globals["symbol"] = "value"
	fmt.Println(print(eval("string", context)))
	fmt.Println(print(eval(Symbol("symbol"), context)))
	fmt.Println(print(eval(123, context)))
	fmt.Println(print(list(123, "a", 3, 4)))
	fmt.Println(print(Cons{123, "a"}))
}

