// Copyright (c) 2014 by Christoph Hack <christoph@tux21b.org>
// All rights reserved. Distributed under the Simplified BSD License.

package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
)

var globals = map[string]interface{}{
	"echo": func(x interface{}) {
		fmt.Println("foo", x)
	},
	"quit": func() {
		fmt.Println("Bye.")
		os.Exit(0)
	},
	"p": func(expr Expr) {
		p, err := NewPolynomial(expr)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Println(p)
	},
	"multicoeff": func(expr, vars, exp Expr) (Expr, error) {
		p, err := NewPolynomial(expr)
		if err != nil {
			return nil, err
		}
		var varlist []string
		if v, ok := vars.(List); ok {
			varlist = make([]string, len(v))
			for i := 0; i < len(v); i++ {
				if x, ok := v[i].(Ident); ok {
					varlist[i] = string(x)
				} else {
					return nil, fmt.Errorf("invalid vars list")
				}
			}
		} else {
			return nil, fmt.Errorf("invalid vars list")
		}
		var explist []Num
		if v, ok := exp.(List); ok {
			explist = make([]Num, len(v))
			for i := 0; i < len(v); i++ {
				if x, ok := v[i].(Num); ok {
					explist[i] = x
				} else {
					return nil, fmt.Errorf("invalid exp list")
				}
			}
		} else {
			return nil, fmt.Errorf("invalid exp list")
		}
		return p.MultiCoeff(varlist, explist), nil
	},
}

func executeCall(call Call) error {
	fn, ok := globals[string(call.Ident)]
	if !ok {
		return fmt.Errorf("undefined %q", call.Ident)
	}
	v := reflect.ValueOf(fn)
	args := make([]reflect.Value, len(call.Args))
	for i := 0; i < len(args); i++ {
		args[i] = reflect.ValueOf(call.Args[i])
	}
	t := v.Type()
	if t.NumIn() != len(args) {
		return fmt.Errorf("invalid number of args. expected %d, got %d.\n",
			t.NumIn(), len(args))
	}
	result := v.Call(args)
	var (
		val Expr
		err error
	)
	if len(result) >= 1 && !result[0].IsNil() {
		x, ok := result[0].Interface().(Expr)
		if !ok {
			return fmt.Errorf("invalid result")
		}
		val = x
	}
	if len(result) >= 2 && !result[1].IsNil() {
		x, ok := result[1].Interface().(error)
		if !ok {
			return fmt.Errorf("invalid result")
		}
		err = x
	}
	if len(result) >= 3 {
		return fmt.Errorf("invalid result")
	}
	if err == nil && val != nil {
		fmt.Println(val)
	}
	return err
}

func main() {
	fmt.Println("Bruno 0.1 (2014-03-22) -- \"Übungszettel 1\"")
	fmt.Println("Copyright (c) 2014 by Christoph Hack <christoph@tux21b.org>")
	fmt.Println("Type 'quit()' to quit Bruno.")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		expr, err := Parse(scanner.Text())
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		if call, ok := expr.(Call); ok {
			if err := executeCall(call); err != nil {
				fmt.Println("error:", err)
			}
		}
	}
}
