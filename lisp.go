package main

import (
	"fmt"
	"strconv"
)

type Context struct {
	globals map[Symbol]Evaluable
}

func context() *Context {
	return &Context{map[Symbol]Evaluable{}}
}

type Evaluable interface {
	eval(context *Context) Evaluable
	apply(args Evaluable, context *Context) Evaluable
	toString() string
}

type Symbol string
func (this Symbol)eval(context *Context) Evaluable {
	return context.globals[this]
}
func (this Symbol)apply(args Evaluable, context *Context) Evaluable {
	panic("can't apply")
}
func (this Symbol)toString() string {
	return string(this)
}

type Integer int
func (this Integer)eval(context *Context) Evaluable {
	return this
}
func (this Integer)apply(args Evaluable, context *Context) Evaluable {
	panic("can't apply");
}
func (this Integer)toString() string {
	return strconv.Itoa(int(this))
}

type Cons struct {
	car, cdr Evaluable
}
func (this Cons)eval(context *Context) Evaluable {
	return this.car.eval(context).apply(this.cdr, context)
}
func (this Cons)apply(args Evaluable, context *Context) Evaluable {
	panic("can't apply");
}
func (this Cons)toString() string {
	return fmt.Sprintf("(%s . %s)", this.car.toString(), this.cdr.toString())
}

func main() {
	var context *Context = context()
	context.globals[Symbol("symbol")] = Symbol("value")
	fmt.Println(Symbol("symbol").eval(context))
	fmt.Println(Integer(123).eval(context))
	fmt.Println(Cons{Integer(345), Symbol("abc")}.toString())
}

