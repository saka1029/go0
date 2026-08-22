package main

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"unicode"
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

func arithmetic[T int | int8 | int16 | int32 | int64 | float32 | float64](
	args Evaluable, start T, f func(a, b T) T) Evaluable {
	var count, prev T
	for {
		if cons, ok := args.(Cons); ok {
			value := cons.car.(T)
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
			env.Define(parm.car.(Symbol), car(args))
			parms = parm.cdr
			args = cdr(args)
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
		return car(args)
	})
	e.Define(Symbol("car"), func(args Evaluable, c *Env) Evaluable {
		return car(car(evlis(args, c)))
	})
	e.Define(Symbol("cdr"), func(args Evaluable, c *Env) Evaluable {
		return cdr(car(evlis(args, c)))
	})
	e.Define(Symbol("cons"), func(args Evaluable, c *Env) Evaluable {
		evaled := evlis(args, c)
		return cons(car(evaled), car(cdr(evaled)))
	})
	e.Define(Symbol("+"), func(args Evaluable, c *Env) Evaluable {
		return arithmetic(evlis(args, c), 0, func(a, b int) int { return a + b })
	})
	e.Define(Symbol("-"), func(args Evaluable, c *Env) Evaluable {
		return arithmetic(evlis(args, c), 0, func(a, b int) int { return a - b })
	})
	e.Define(Symbol("*"), func(args Evaluable, c *Env) Evaluable {
		return arithmetic(evlis(args, c), 1, func(a, b int) int { return a * b })
	})
	e.Define(Symbol("/"), func(args Evaluable, c *Env) Evaluable {
		return arithmetic(evlis(args, c), 1, func(a, b int) int { return a / b })
	})
	e.Define(Symbol("lambda"), func(args Evaluable, c *Env) Evaluable {
		parms := car(args)
		body := cdr(args)
		return func(args Evaluable, e *Env) Evaluable {
			n := env(c)
			pairlis(parms, evlis(args, c), n)
			return progn(body, n)
		}
	})
	e.Define(Symbol("define"), func(args Evaluable, c *Env) Evaluable {
		c.Define(car(args).(Symbol), eval(car(cdr(args)), c))
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

func car(this Evaluable) Evaluable {
	return this.(Cons).car
}

func cdr(this Evaluable) Evaluable {
	return this.(Cons).cdr
}

func listDot(elements ...Evaluable) Evaluable {
	i := len(elements) - 1
	result := elements[i]
	for i--; i >= 0; i-- {
		result = cons(elements[i], result)
	}
	return result
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

type Reader struct {
	runeReader *strings.Reader
	ch         rune
	buffer     []rune
}

func NewReader(source string) *Reader {
	rr := strings.NewReader(source)
	reader := &Reader{rr, 0, []rune{}}
	reader.get()
	return reader
}

func (this *Reader) append(ch rune) {
	this.buffer = append(this.buffer, ch)
}

func (this *Reader) get() rune {
	if this.ch == EOF {
		return this.ch
	}
	ch, _, e := this.runeReader.ReadRune()
	if e == io.EOF {
		this.ch = EOF
	} else if e != nil {
		panic("get() : " + e.Error())
	} else {
		this.ch = ch
	}
	// EOFの場合もthis.bufferにappendする
	this.append(this.ch)
	return this.ch
}

func (this *Reader) getClear() rune {
	this.buffer = nil
	return this.get()
}

func (this *Reader) spaces() {
	last := this.ch
	for unicode.IsSpace(this.ch) {
		last = this.get()
	}
	this.buffer = nil
	this.append(last)
}

const EOF = rune(-1)

func (this *Reader) readList() Evaluable {
	this.getClear()
	slice := []Evaluable{}
	for {
		this.spaces()
		switch this.ch {
		case ')':
			this.getClear()
			return list(slice...)
		case EOF:
			panic("readList: unexpected EOF")
		case '.':
			this.getClear()
			last := this.read()
			this.spaces()
			if this.ch != ')' {
				panic("readList: ')' expected")
			}
			this.getClear()
			slice = append(slice, last)
			return listDot(slice)
		default:
			slice = append(slice, this.read())
		}
	}
}

func (this *Reader) readNumber() int {
	for unicode.IsDigit(this.ch) {
		this.get()
	}
	if result, err := strconv.Atoi(string(this.buffer[0 : len(this.buffer)-1])); err == nil {
		return result
	} else {
		panic("readNumber: " + err.Error())
	}
}

func isSymbolFirst(ch rune) bool {
	switch ch {
	case -1, '(', ')', '.':
		return false
	default:
		return !unicode.IsSpace(ch) && !unicode.IsDigit(ch)
	}
}

func isSymbolRest(ch rune) bool {
	return isSymbolFirst(ch) || unicode.IsDigit(ch) || ch == '.'
}

func (this *Reader) readSymbol() Symbol {
	for isSymbolRest(this.ch) {
		this.get()
	}
	result := string(this.buffer[0 : len(this.buffer)-1])
	return sym(result)
}

func (this *Reader) read() Evaluable {
	this.spaces()
	switch this.ch {
	case -1:
		return EOF
	case '(':
		return this.readList()
	case '\'':
		this.getClear()
		return list(sym("quote"), this.read())
	case '+', '-':
		if unicode.IsDigit(this.get()) {
			return this.readNumber()
		} else {
			return this.readSymbol()
		}
	case '.':
		panic("read: unexpected '.'")
	default:
		if isSymbolFirst(this.ch) {
			return this.readSymbol()
		} else if unicode.IsDigit(this.ch) {
			return this.readNumber()
		} else {
			panic("read: unexpected char '" + string(this.ch) + "'")
		}
	}

}

func main() {
}
