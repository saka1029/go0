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
	panic(fmt.Sprint("Get: cannot get the value of", key))
}

func (this *Env) Set(key Symbol, value Evaluable) {
	for cur := this; cur != nil; cur = cur.previous {
		if _, ok := cur.globals[key]; ok {
			cur.globals[key] = value
			return
		}
	}
	panic(fmt.Sprint("Set: cannot set the value of", key))
}

func (this *Env) Define(key Symbol, value Evaluable) {
	this.globals[key] = value
}

func evlis(args Evaluable, c *Env) Evaluable {
	result := []Evaluable{}
	for {
		if cons, ok := args.(Cons); ok {
			result = append(result, eval(cons.Car(), c))
			args = cons.Cdr()
		} else {
			return list(result...)
		}
	}
}

func intArithmetic(args Evaluable, unit int, f func(a, b int) int, c *Env) int {
	args = evlis(args, c)
	var count, prev int
	for {
		if cons, ok := args.(Cons); ok {
			value := cons.Car().(int)
			if count == 1 {
				unit = prev
			}
			unit = f(unit, value)
			count++
			prev = value
			args = cons.Cdr()
		} else {
			return unit
		}
	}
}

func env(previous ...*Env) *Env {
	var env *Env
	switch len(previous) {
	case 0:
		env = &Env{map[Symbol]Evaluable{}, nil}
	case 1:
		env = &Env{map[Symbol]Evaluable{}, previous[0]}
	default:
		panic("env() Too many args")
	}
	env.Define(Symbol("quote"), func(args Evaluable, c *Env) Evaluable {
		return args.(Cons).Car()
	})
	env.Define(Symbol("car"), func(args Evaluable, c *Env) Evaluable {
		return evlis(args, c).(Cons).Car().(Cons).Car()
	})
	env.Define(Symbol("cdr"), func(args Evaluable, c *Env) Evaluable {
		return evlis(args, c).(Cons).Car().(Cons).Cdr()
	})
	env.Define(Symbol("cons"), func(args Evaluable, c *Env) Evaluable {
		evaled := evlis(args, c)
		return cons(evaled.(Cons).Car(), evaled.(Cons).Cdr().(Cons).Car())
	})
	env.Define(Symbol("+"), func(args Evaluable, c *Env) Evaluable {
		return intArithmetic(args, 0, func(a, b int) int { return a + b }, c)
	})
	env.Define(Symbol("-"), func(args Evaluable, c *Env) Evaluable {
		return intArithmetic(args, 0, func(a, b int) int { return a - b }, c)
	})
	env.Define(Symbol("*"), func(args Evaluable, c *Env) Evaluable {
		return intArithmetic(args, 1, func(a, b int) int { return a * b }, c)
	})
	env.Define(Symbol("/"), func(args Evaluable, c *Env) Evaluable {
		return intArithmetic(args, 1, func(a, b int) int { return a / b }, c)
	})
	return env
}

type Evaluable interface {
}

func Equal(left, right Evaluable) bool {
	if lcons, ok := left.(Cons); ok {
		return lcons.Equal(right)
	} else {
		return left == right
	}
}

type Symbol string

func sym(name string) Symbol {
	return Symbol(name)
}

type Cons struct {
	car, cdr *Evaluable
}

func cons(car, cdr Evaluable) Cons {
	return Cons{&car, &cdr}
}

func (this Cons) Car() Evaluable {
	return *this.car
}

func (this Cons) Cdr() Evaluable {
	return *this.cdr
}

func (this Cons) Equal(right Evaluable) bool {
	if rcons, ok := right.(Cons); ok {
		return Equal(this.Car(), rcons.Car()) && Equal(this.Cdr(), rcons.Cdr())
	} else {
		return false
	}
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
	case int, string, bool:
		return v
	case Symbol:
		return c.globals[v]
	case Cons:
		return eval(v.Car(), c).(func(args Evaluable, c *Env) Evaluable)(v.Cdr(), c)
	default:
		panic(fmt.Sprint("eval: unknown type ", v))
	}
}

func printCons(c Cons) string {
	var sb strings.Builder
	if cdr, ok := c.Cdr().(Cons); ok && c.Car() == Symbol("quote") && cdr.Cdr() == nil {
		sb.WriteString("'")
		sb.WriteString(print(cdr.Car()))
		return sb.String()
	}
	sb.WriteString("(")
	sb.WriteString(print(c.Car()))
	e := c.Cdr()
	for {
		if cons, ok := e.(Cons); ok {
			sb.WriteString(" ")
			sb.WriteString(print(cons.Car()))
			e = cons.Cdr()
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
	case int, bool:
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

func main() {
	var context *Env = env()
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
