package main

import (
	"fmt"
)

type Context struct {
	globals map[Symbol]Evaluable
}

func context() *Context {
	return &Context{map[Symbol]Evaluable{}}
}

type Evaluable interface {
	eval(context *Context) Evaluable
}

type Symbol string
func (symbol Symbol)eval(context *Context) Evaluable {
	return context.globals[symbol]
}

type Integer int
func (value Integer)eval(context *Context) Evaluable {
	return value
}

func main() {
	var context *Context = context()
	context.globals[Symbol("symbol")] = Symbol("value")
	fmt.Println(Symbol("symbol").eval(context))
	fmt.Println(Integer(123).eval(context))
}
