// Copyright (c) 2014 by Christoph Hack <christoph@tux21b.org>
// All rights reserved. Distributed under the Simplified BSD License.

package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
)

type Bruno struct {
	globals map[string]interface{}
}

func NewBruno() *Bruno {
	b := &Bruno{}
	b.reset()
	return b
}

func (b *Bruno) reset() {
	b.globals = map[string]interface{}{
		"quit": func() {
			fmt.Println("Bye.")
			os.Exit(0)
		},
		"reset": func() {
			b.reset()
		},
		"p": func(expr Expr) (Expr, error) {
			return NewPolynomial(expr)
		},
		"multicoeff": func(p *Polynomial, vars, exp Expr) (Expr, error) {
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
		"multicoeff2": func(p, term *Polynomial) (Expr, error) {
			if len(term.items) != 1 {
				return nil, fmt.Errorf("invalid term")
			}
			varlist := term.vars
			explist := make([]Num, len(term.vars))
			for i := range explist {
				explist[i] = Num{&term.items[0].T[i]}
			}
			return p.MultiCoeff(varlist, explist), nil
		},
		"support": func(p *Polynomial, vars []string) (Expr, error) {
			s := p.Support(vars)
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
		"lexorder": func(p *Polynomial) Expr {
			SortMonomial(p.items, LexTermOrder)
			return p
		},
		"totalorder": func(p *Polynomial) Expr {
			SortMonomial(p.items, TotalTermOrder)
			return p
		},
		"lpp": func(p *Polynomial) Expr {
			return p.LPP()
		},
		"lc": func(p *Polynomial) Expr {
			return p.LC()
		},
		"lm": func(p *Polynomial) Expr {
			return p.LM()
		},
		"higher": func(p *Polynomial, term Expr) (Expr, error) {
			t, err := convertTerm(p, term)
			if err != nil {
				return nil, err
			}
			return p.Higher(t), nil
		},
		"lower": func(p *Polynomial, term Expr) (Expr, error) {
			t, err := convertTerm(p, term)
			if err != nil {
				return nil, err
			}
			return p.Lower(t), nil
		},
		"between": func(p *Polynomial, term1, term2 Expr) (Expr, error) {
			t1, err := convertTerm(p, term1)
			if err != nil {
				return nil, err
			}
			t2, err := convertTerm(p, term2)
			if err != nil {
				return nil, err
			}
			return p.Between(t1, t2), nil
		},
		"remainder": func(p *Polynomial) Expr {
			return p.Remainder()
		},
	}
}

func (b *Bruno) executeCall(call Call) (Expr, error) {
	fn, ok := b.globals[string(call.Ident)]
	if !ok {
		return nil, fmt.Errorf("undefined %q", call.Ident)
	}
	v := reflect.ValueOf(fn)
	fnT := v.Type()

	if fnT.NumIn() != len(call.Args) {
		return nil, fmt.Errorf("invalid number of args. expected %d, got %d.\n",
			fnT.NumIn(), len(call.Args))
	}

	args := make([]reflect.Value, len(call.Args))
	for i := 0; i < len(args); i++ {
		gotV := reflect.ValueOf(call.Args[i])
		gotT := gotV.Type()
		wantT := fnT.In(i)
		switch {
		case gotT.AssignableTo(wantT):
			args[i] = gotV
		case wantT == reflect.TypeOf(&Polynomial{}):
			p, err := NewPolynomial(call.Args[i])
			if err != nil {
				return nil, fmt.Errorf("invalid parameter %d: %v", i+1, err)
			}
			args[i] = reflect.ValueOf(p)
		case wantT == reflect.TypeOf([]string{}):
			vars, err := convertVars(call.Args[i])
			if err != nil {
				return nil, fmt.Errorf("invalid parameter %d: %v", i+1, err)
			}
			args[i] = reflect.ValueOf(vars)
		default:
			return nil, fmt.Errorf("invalid parameter %d.", i+1)
		}
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

func (b *Bruno) ExecExpr(expr Expr) (Expr, error) {
	switch x := expr.(type) {
	case Ident:
		if e, ok := b.globals[string(x)].(Expr); ok {
			return e, nil
		}
	case Call:
		args, err := b.ExecExpr(x.Args)
		if err != nil {
			return nil, err
		}
		call := Call{x.Ident, args.(List)}
		return b.executeCall(call)
	case Assign:
		v, err := b.ExecExpr(x.Expr)
		if err != nil {
			return nil, err
		}
		b.globals[string(x.Ident)] = v
		return Assign{x.Ident, v}, nil
	case Add:
		a, err := b.ExecExpr(x.A)
		if err != nil {
			return nil, err
		}
		b, err := b.ExecExpr(x.B)
		if err != nil {
			return nil, err
		}
		return Add{a, b}, nil
	case Sub:
		a, err := b.ExecExpr(x.A)
		if err != nil {
			return nil, err
		}
		b, err := b.ExecExpr(x.B)
		if err != nil {
			return nil, err
		}
		return Sub{a, b}, nil
	case Mul:
		a, err := b.ExecExpr(x.A)
		if err != nil {
			return nil, err
		}
		b, err := b.ExecExpr(x.B)
		if err != nil {
			return nil, err
		}
		return Mul{a, b}, nil
	case Div:
		a, err := b.ExecExpr(x.A)
		if err != nil {
			return nil, err
		}
		b, err := b.ExecExpr(x.B)
		if err != nil {
			return nil, err
		}
		return Div{a, b}, nil
	case Pow:
		a, err := b.ExecExpr(x.A)
		if err != nil {
			return nil, err
		}
		b, err := b.ExecExpr(x.B)
		if err != nil {
			return nil, err
		}
		return Pow{a, b}, nil
	case List:
		result := make(List, len(x))
		for i := range x {
			v, err := b.ExecExpr(x[i])
			if err != nil {
				return nil, err
			}
			result[i] = v
		}
		return result, nil
	}
	return expr, nil
}

func (b *Bruno) Exec(input string) (Expr, error) {
	expr, err := Parse(input)
	if err != nil {
		return nil, err
	}
	return b.ExecExpr(expr)
}

func convertTerm(p *Polynomial, expr Expr) (Term, error) {
	q := &Polynomial{vars: p.vars, order: p.order}
	if err := q.convert(expr); err != nil {
		return nil, err
	}
	if len(q.items) != 1 {
		return nil, fmt.Errorf("invalid term")
	}
	return q.items[0].T, nil
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

func main() {
	fmt.Println("Bruno 0.1 (2014-03-22) -- \"Übungszettel 1\"")
	fmt.Println("Copyright (c) 2014 by Christoph Hack <christoph@tux21b.org>")
	fmt.Println("Type 'quit()' to quit Bruno.")
	fmt.Println()

	bruno := NewBruno()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		result, err := bruno.Exec(scanner.Text())
		if err != nil {
			fmt.Println("error:", err)
		} else if result != nil {
			fmt.Println(result)
		}
	}
}
