package main

import (
	"fmt"
	"strings"
)

type Context struct {
	globals map[Symbol]Evaluable
}

func context() *Context {
	context := &Context{map[Symbol]Evaluable{}}
	context.globals[Symbol("quote")] = func(args Evaluable, c *Context) Evaluable {
		return args.(Cons).car
	}
	return context
}

type Evaluable interface {
}

type Symbol string

type Cons struct {
	car, cdr Evaluable
}

func cons(car, cdr Evaluable) Cons {
	return Cons{car, cdr}
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
		return eval(v.car, c).(func(args Evaluable, c *Context) Evaluable)(v.cdr, c)
	default:
		panic(fmt.Sprint("eval: unknown type ", v))
	}
}

func printCons(c Cons) string {
	var sb strings.Builder
	if c.car == Symbol("quote") {
		cddr, ok := c.cdr.(Cons)
		if ok && cddr.cdr == nil {
			sb.WriteString("'")
			sb.WriteString(print(cddr.car))
			return sb.String()
		}
	}
	sb.WriteString("(")
	sb.WriteString(print(c.car))
	e := c.cdr
	for true {
		cons, ok := e.(Cons)
		if !ok {
			break
		}
		sb.WriteString(" ")
		sb.WriteString(print(cons.car))
		e = cons.cdr
	}
	if e != nil {
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
	case Symbol:
		return string(v)
	case Cons:
		return printCons(v)
	default:
		panic(fmt.Sprint("print: unknown type ", v))
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
	fmt.Println(print(cons(123, "a")))
	fmt.Println(print(list(Symbol("quote"), list("a"))))
	fmt.Println(print(eval(list(Symbol("quote"), list("a")), context)))
	fmt.Println(print(eval(list(Symbol("quote"), Symbol("b")), context)))
	fmt.Println(print(cons(Symbol("quote"), Symbol("b"))))
	// fmt.Println(print(eval(cons(Symbol("quote"), Symbol("b")), context)))
	// fmt.Println(print(eval(2.3, context)))
}
