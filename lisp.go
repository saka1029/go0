package main

import (
	"fmt"
	// "strconv"
)

type Context struct {
	globals map[string]Evaluable
}

func context() *Context {
	return &Context{map[string]Evaluable{}}
}

type Evaluable interface {
// 	// eval(context *Context) Evaluable
// 	// apply(args Evaluable, context *Context) Evaluable
// 	// toString() string
}

// type Symbol string
// func (this Symbol)eval(context *Context) Evaluable {
// 	return context.globals[this]
// }
// func (this Symbol)apply(args Evaluable, context *Context) Evaluable {
// 	panic("can't apply")
// }
// func (this Symbol)toString() string {
// 	return string(this)
// }

// type Integer int
// func (this Integer)eval(context *Context) Evaluable {
// 	return this
// }
// func (this Integer)apply(args Evaluable, context *Context) Evaluable {
// 	panic("can't apply");
// }
// func (this Integer)toString() string {
// 	return strconv.Itoa(int(this))
// }

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
// func (this Cons)eval(context *Context) Evaluable {
// 	return this.car.eval(context).apply(this.cdr, context)
// }
// func (this Cons)apply(args Evaluable, context *Context) Evaluable {
// 	panic("can't apply");
// }
// func (this Cons)toString() string {
// 	return fmt.Sprintf("(%s . %s)", this.car.toString(), this.cdr.toString())
// }

func check(v interface{}) {
	switch x := v.(type) {
	case int:
		fmt.Println("int ", x)
	// case Integer:
	// 	fmt.Println("Integer %d\n", x)
	case string:
		fmt.Println("string ", x)
	default:
		fmt.Println("unknown ", x)
	}
}

func eval(e Evaluable, c *Context) Evaluable {
	switch v := e.(type) {
	case int:
		return v
	case string:
		return c.globals[v]
	case Cons:
		return apply(eval(v.car, c), v.cdr, c)
	default:
		panic(fmt.Sprintf("unknown type %T", v))
	}
}

func apply(e Evaluable, args Evaluable, c *Context) Evaluable {
	return nil
}

func main() {
	var context *Context = context()
	fmt.Println(context)
	context.globals["symbol"] = "value"
	fmt.Println(eval("symbol", context))
	fmt.Println(eval(123, context))
	fmt.Println(list(123, "a", 3, 4))
	// fmt.Println(Symbol("symbol").eval(context))
	// fmt.Println(Integer(123).eval(context))
	// fmt.Println(Cons{Integer(345), Symbol("abc")}.toString())
	// check(123)
	// check("abc")
	// check(Cons{2,3})
	// check(Integer(123))
	// check(Symbol("sym"))
}

