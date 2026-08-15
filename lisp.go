package main

import (
	"fmt"
)

type Symbol string

type Context struct {
	globals map[Symbol]Evaluable
}

func context() *Context {
	return &Context{map[Symbol]Evaluable{}}
}

type Evaluable interface {
	eval(context *Context) Evaluable
}

func (symbol Symbol)eval(context *Context) Evaluable {
	return context.globals[symbol]
}

func main() {
	var context *Context = context()
	symbol := Symbol("symbol")
	context.globals[symbol] = Symbol("value")
	fmt.Println(symbol.eval(context))
}
