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
	"quit": func() {
		fmt.Println("Bye.")
		os.Exit(0)
	},
	"p": func(expr Expr) (Expr, error) {
		return NewPolynomial(expr)
	},
	"multicoeff": func(expr, vars, exp Expr) (Expr, error) {
		p, err := NewPolynomial(expr)
		if err != nil {
			return nil, err
		}
		varlist, err := convertVars(vars)
		if err != nil {
			return nil, err
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
	"support": func(expr, vars Expr) (Expr, error) {
		p, err := NewPolynomial(expr)
		if err != nil {
			return nil, err
		}
		varlist, err := convertVars(vars)
		if err != nil {
			return nil, err
		}
		s := p.Support(varlist)
		result := make(List, len(s))
		for i := range result {
			lst := make(List, len(s[i]))
			for j := range s[i] {
				lst[j] = s[i][j]
			}
			result[i] = lst
		}
		return result, nil
	},
	"lexorder": func(expr Expr) (Expr, error) {
		p, err := NewPolynomial(expr)
		if err != nil {
			return nil, err
		}
		SortTerms(p.terms, LexTermOrder)
		return p, nil
	},
	"totalorder": func(expr Expr) (Expr, error) {
		p, err := NewPolynomial(expr)
		if err != nil {
			return nil, err
		}
		SortTerms(p.terms, TotalTermOrder)
		return p, nil
	},
	"lpp": func(expr Expr) (Expr, error) {
		p, err := NewPolynomial(expr)
		if err != nil {
			return nil, err
		}
		return p.LPP(), nil
	},
	"lc": func(expr Expr) (Expr, error) {
		p, err := NewPolynomial(expr)
		if err != nil {
			return nil, err
		}
		return p.LC(), nil
	},
	"lm": func(expr Expr) (Expr, error) {
		p, err := NewPolynomial(expr)
		if err != nil {
			return nil, err
		}
		return p.LM(), nil
	},
}

func executeCall(call Call) (Expr, error) {
	fn, ok := globals[string(call.Ident)]
	if !ok {
		return nil, fmt.Errorf("undefined %q", call.Ident)
	}
	v := reflect.ValueOf(fn)
	args := make([]reflect.Value, len(call.Args))
	for i := 0; i < len(args); i++ {
		args[i] = reflect.ValueOf(call.Args[i])
	}
	t := v.Type()
	if t.NumIn() != len(args) {
		return nil, fmt.Errorf("invalid number of args. expected %d, got %d.\n",
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
			return nil, fmt.Errorf("invalid result")
		}
		val = x
	}
	if len(result) >= 2 && !result[1].IsNil() {
		x, ok := result[1].Interface().(error)
		if !ok {
			return nil, fmt.Errorf("invalid result")
		}
		err = x
	}
	if len(result) >= 3 {
		return nil, fmt.Errorf("invalid result")
	}
	return val, err
}

func convertVars(expr Expr) ([]string, error) {
	var list []string
	if v, ok := expr.(List); ok {
		list = make([]string, len(v))
		for i := 0; i < len(v); i++ {
			if x, ok := v[i].(Ident); ok {
				list[i] = string(x)
			} else {
				return nil, fmt.Errorf("invalid vars list")
			}
		}
	} else {
		return nil, fmt.Errorf("invalid vars list")
	}
	return list, nil
}

func execute(expr Expr) (Expr, error) {
	switch x := expr.(type) {
	case Ident:
		if e, ok := globals[string(x)].(Expr); ok {
			return e, nil
		}
	case Call:
		args, err := execute(x.Args)
		if err != nil {
			return nil, err
		}
		call := Call{x.Ident, args.(List)}
		return executeCall(call)
	case Assign:
		v, err := execute(x.Expr)
		if err != nil {
			return nil, err
		}
		globals[string(x.Ident)] = v
		return Assign{x.Ident, v}, nil
	case Add:
		a, err := execute(x.A)
		if err != nil {
			return nil, err
		}
		b, err := execute(x.B)
		if err != nil {
			return nil, err
		}
		return Add{a, b}, nil
	case Sub:
		a, err := execute(x.A)
		if err != nil {
			return nil, err
		}
		b, err := execute(x.B)
		if err != nil {
			return nil, err
		}
		return Sub{a, b}, nil
	case Mul:
		a, err := execute(x.A)
		if err != nil {
			return nil, err
		}
		b, err := execute(x.B)
		if err != nil {
			return nil, err
		}
		return Mul{a, b}, nil
	case Div:
		a, err := execute(x.A)
		if err != nil {
			return nil, err
		}
		b, err := execute(x.B)
		if err != nil {
			return nil, err
		}
		return Div{a, b}, nil
	case Pow:
		a, err := execute(x.A)
		if err != nil {
			return nil, err
		}
		b, err := execute(x.B)
		if err != nil {
			return nil, err
		}
		return Pow{a, b}, nil
	case List:
		result := make(List, len(x))
		for i := range x {
			v, err := execute(x[i])
			if err != nil {
				return nil, err
			}
			result[i] = v
		}
		return result, nil
	}
	return expr, nil
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
		if result, err := execute(expr); err != nil {
			fmt.Println("error:", err)
		} else if result != nil {
			fmt.Println(result)
		}
	}
}
