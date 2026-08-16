package main

import (
	"fmt"
	"slices"
	"strings"
)

type Context struct {
	globals map[Symbol]Evaluable
}

func evlis(args Evaluable, c *Context) Evaluable {
	result := []Evaluable{}
	for {
		if cons, ok := args.(Cons); ok {
			result = append(result, eval(cons.car, c))
			args = cons.cdr
		} else {
			return listSlice(result)
		}
	}
}

func intArithmetic(args Evaluable, unit int, f func(a, b int) int, c *Context) int {
	args = evlis(args, c)
	for {
		if cons, ok := args.(Cons); ok {
			unit = f(unit, cons.car.(int))
			args = cons.cdr
		} else {
			return unit
		}
	}
}
func intArithmeticMinus(args Evaluable, unit int, f func(a, b int) int, c *Context) int {
	args = evlis(args, c)
	var count, prev int
	for {
		if cons, ok := args.(Cons); ok {
			value := cons.car.(int)
			if count == 1 {
				unit = prev
			}
			unit = f(unit, value)
			count++
			prev = value
			args = cons.cdr
		} else {
			return unit
		}
	}
}

func context() *Context {
	context := &Context{map[Symbol]Evaluable{}}
	context.globals[Symbol("quote")] = func(args Evaluable, c *Context) Evaluable {
		return args.(Cons).car
	}
	context.globals[Symbol("car")] = func(args Evaluable, c *Context) Evaluable {
		return evlis(args, c).(Cons).car.(Cons).car
	}
	context.globals[Symbol("cdr")] = func(args Evaluable, c *Context) Evaluable {
		return evlis(args, c).(Cons).car.(Cons).cdr
	}
	context.globals[Symbol("cons")] = func(args Evaluable, c *Context) Evaluable {
		evaled := evlis(args, c)
		return cons(evaled.(Cons).car, evaled.(Cons).cdr.(Cons).car)
	}
	context.globals[Symbol("+")] = func(args Evaluable, c *Context) Evaluable {
		return intArithmetic(args, 0, func(a, b int) int { return a + b }, c)
	}
	context.globals[Symbol("-")] = func(args Evaluable, c *Context) Evaluable {
		return intArithmeticMinus(args, 0, func(a, b int) int { return a - b }, c)
	}
	context.globals[Symbol("*")] = func(args Evaluable, c *Context) Evaluable {
		return intArithmetic(args, 1, func(a, b int) int { return a * b }, c)
	}
	context.globals[Symbol("/")] = func(args Evaluable, c *Context) Evaluable {
		return intArithmeticMinus(args, 1, func(a, b int) int { return a / b }, c)
	}
	return context
}

type Evaluable interface {
}

type Symbol string

func sym(name string) Symbol {
	return Symbol(name)
}

type Cons struct {
	car, cdr Evaluable
}

func cons(car, cdr Evaluable) Cons {
	return Cons{car, cdr}
}

func list(elements ...Evaluable) Evaluable {
	var result Evaluable = nil
	for _, e := range slices.Backward(elements) {
		result = Cons{e, result}
	}
	return result
}

func listSlice(elements []Evaluable) Evaluable {
	var result Evaluable = nil
	for _, e := range slices.Backward(elements) {
		result = Cons{e, result}
	}
	return result
}

func eval(e Evaluable, c *Context) Evaluable {
	switch v := e.(type) {
	case int, string, bool:
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
		if cddr, ok := c.cdr.(Cons); ok && cddr.cdr == nil {
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
	case int, bool:
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
	fmt.Println(print(cons(Symbol("quote"), Symbol("b"))))
	fmt.Println(print(list(Symbol("quote"), Symbol("b"), 123)))
	// fmt.Println(print(eval(cons(Symbol("quote"), Symbol("b")), context)))
	// fmt.Println(print(eval(2.3, context)))
}
