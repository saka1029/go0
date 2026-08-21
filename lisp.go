package main

import (
	"fmt"
	"slices"
	"strings"
)

type Env struct {
	globals  map[Symbol]Evaluable
	previous *Env
}

func (this *Env) Get(key Symbol) Evaluable {
	for cur := this; cur != nil; cur = cur.previous {
		if value, ok := cur.globals[key]; ok {
			return value
		}
	}
	panic(fmt.Sprint("Get: cannot get the value of ", key))
}

func (this *Env) Set(key Symbol, value Evaluable) {
	for cur := this; cur != nil; cur = cur.previous {
		if _, ok := cur.globals[key]; ok {
			cur.globals[key] = value
			return
		}
	}
	panic(fmt.Sprint("Set: cannot set the value of ", key))
}

func (this *Env) Define(key Symbol, value Evaluable) {
	this.globals[key] = value
}

func evlis(args Evaluable, c *Env) Evaluable {
	result := []Evaluable{}
	for {
		if cons, ok := args.(Cons); ok {
			result = append(result, eval(cons.car, c))
			args = cons.cdr
		} else {
			return list(result...)
		}
	}
}

func intArithmetic(args Evaluable, start int, f func(a, b int) int, c *Env) int {
	args = evlis(args, c)
	var count, prev int
	for {
		if cons, ok := args.(Cons); ok {
			value := cons.car.(int)
			if count == 1 {
				start = prev
			}
			start = f(start, value)
			count++
			prev = value
			args = cons.cdr
		} else {
			return start
		}
	}
}

func pairlis(parms Evaluable, args Evaluable, env *Env) {
	for {
		if parm, ok := parms.(Cons); ok {
			env.Define(parm.car.(Symbol), args.(Cons).car)
			parms = parm.cdr
			args = args.(Cons).cdr
		} else {
			break
		}
	}
	if parms != nil {
		env.Define(parms.(Symbol), args)
	}
}

func progn(body Evaluable, env *Env) Evaluable {
	if body.(Cons).cdr == nil {
		return eval(body.(Cons).car, env)
	}
	eval(body.(Cons).car, env)
	return progn(body.(Cons).cdr, env)
}

func env(prev ...*Env) *Env {
	switch len(prev) {
	case 0:
		return &Env{map[Symbol]Evaluable{}, nil}
	case 1:
		return &Env{map[Symbol]Evaluable{}, prev[0]}
	default:
		panic("env() : too many arguments")
	}
}

func globalEnv() *Env {
	e := env()
	e.Define(Symbol("quote"), func(args Evaluable, c *Env) Evaluable {
		return args.(Cons).car
	})
	e.Define(Symbol("car"), func(args Evaluable, c *Env) Evaluable {
		return evlis(args, c).(Cons).car.(Cons).car
	})
	e.Define(Symbol("cdr"), func(args Evaluable, c *Env) Evaluable {
		return evlis(args, c).(Cons).car.(Cons).cdr
	})
	e.Define(Symbol("cons"), func(args Evaluable, c *Env) Evaluable {
		evaled := evlis(args, c)
		return cons(evaled.(Cons).car, evaled.(Cons).cdr.(Cons).car)
	})
	e.Define(Symbol("+"), func(args Evaluable, c *Env) Evaluable {
		return intArithmetic(args, 0, func(a, b int) int { return a + b }, c)
	})
	e.Define(Symbol("-"), func(args Evaluable, c *Env) Evaluable {
		return intArithmetic(args, 0, func(a, b int) int { return a - b }, c)
	})
	e.Define(Symbol("*"), func(args Evaluable, c *Env) Evaluable {
		return intArithmetic(args, 1, func(a, b int) int { return a * b }, c)
	})
	e.Define(Symbol("/"), func(args Evaluable, c *Env) Evaluable {
		return intArithmetic(args, 1, func(a, b int) int { return a / b }, c)
	})
	e.Define(Symbol("lambda"), func(args Evaluable, c *Env) Evaluable {
		parms := args.(Cons).car
		body := args.(Cons).cdr
		return func(args Evaluable, e *Env) Evaluable {
			n := env(c)
			pairlis(parms, evlis(args, c), n)
			return progn(body, n)
		}
	})
	e.Define(Symbol("define"), func(args Evaluable, c *Env) Evaluable {
		c.Define(args.(Cons).car.(Symbol), eval(args.(Cons).cdr.(Cons).car, c))
		return nil
	})
	return e
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
		result = cons(e, result)
	}
	return result
}

func eval(e Evaluable, c *Env) Evaluable {
	switch v := e.(type) {
	case int, int8, int16, int32, int64, string, bool:
		return v
	case Symbol:
		return c.Get(v)
	case Cons:
		return eval(v.car, c).(func(args Evaluable, c *Env) Evaluable)(v.cdr, c)
	default:
		panic(fmt.Sprint("eval: unknown type ", v))
	}
}

func printCons(c Cons) string {
	var sb strings.Builder
	if cdr, ok := c.cdr.(Cons); ok && c.car == Symbol("quote") && cdr.cdr == nil {
		sb.WriteString("'")
		sb.WriteString(print(cdr.car))
		return sb.String()
	}
	sb.WriteString("(")
	sb.WriteString(print(c.car))
	e := c.cdr
	for {
		if cons, ok := e.(Cons); ok {
			sb.WriteString(" ")
			sb.WriteString(print(cons.car))
			e = cons.cdr
		} else {
			break
		}
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
	case int, int8, int16, int32, int64, bool:
		return fmt.Sprint(v)
	case string:
		return v
	case Symbol:
		return string(v)
	case Cons:
		return printCons(v)
	default:
		if e == nil {
			return "()"
		}
		panic(fmt.Sprint("print: unknown type ", v))
	}
}

// func eq(left Evaluable, right Evaluable) bool {
// 	switch v := left.(type) {
// 	case Symbol, int, bool, string:
// 		return left == right
// 	case Cons:
// 		if c, ok := right.(Cons); ok {
// 			fmt.Println("v=", &v, "c=", &c)
// 			return &v == &c
// 		} else {
// 			return false
// 		}
// 	default:
// 		fmt.Println("left=", &left, "right=", &right, "v=", &v)
// 		return &left == &right
// 	}
// }

func main() {
	var context *Env = globalEnv()
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
