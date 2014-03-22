// Copyright (c) 2014 by Christoph Hack <christoph@tux21b.org>
// All rights reserved. Distributed under the Simplified BSD License.

package main

import (
	"bytes"
	"errors"
	"math/big"
	"sort"
)

type Polynomial struct {
	vars  []string
	terms []Term
}

func NewPolynomial(expr Expr) (*Polynomial, error) {
	if p, ok := expr.(*Polynomial); ok {
		return p, nil
	}
	p := &Polynomial{}
	p.vars = collectVars(expr)
	if err := p.convert(expr); err != nil {
		return nil, err
	}
	SortTerms(p.terms, LexTermOrder)
	return p, nil
}

func (p *Polynomial) convert(expr Expr) error {
	if add, ok := expr.(Add); ok {
		if err := p.convert(add.A); err != nil {
			return err
		}
		if err := p.convert(add.B); err != nil {
			return err
		}
		return nil

	}
	term := Term{c: Num{big.NewRat(1, 1)}, s: make([]Num, len(p.vars))}
	for i := 0; i < len(p.vars); i++ {
		term.s[i].Rat = new(big.Rat)
	}
	if err := p.convertTerm(expr, &term); err != nil {
		return err
	}
	p.terms = append(p.terms, term)
	return nil
}

func (p *Polynomial) convertTerm(expr Expr, t *Term) error {
	switch x := expr.(type) {
	case Num:
		t.c.Rat.Mul(t.c.Rat, x.Rat)
		return nil
	case Mul:
		if err := p.convertTerm(x.A, t); err != nil {
			return err
		}
		if err := p.convertTerm(x.B, t); err != nil {
			return err
		}
		return nil
	case Pow:
		ident, ok1 := x.A.(Ident)
		exp, ok2 := x.B.(Num)
		idx := -1
		for i := 0; i < len(p.vars); i++ {
			if p.vars[i] == string(ident) {
				idx = i
				break
			}
		}
		if ok1 && ok2 && idx >= 0 {
			t.s[idx].Rat.Add(t.s[idx].Rat, exp.Rat)
			return nil
		}
	case Ident:
		idx := -1
		for i := 0; i < len(p.vars); i++ {
			if p.vars[i] == string(x) {
				idx = i
				break
			}
		}
		if idx >= 0 {
			t.s[idx].Rat.Add(t.s[idx].Rat, ratOne)
			return nil
		}
	}
	return errors.New("invalid polynomial")
}

var ratOne = big.NewRat(1, 1)

func (p *Polynomial) String() string {
	buf := &bytes.Buffer{}
	for i, t := range p.terms {
		if i > 0 {
			buf.WriteString(" + ")
		}
		buf.WriteString(t.c.String())
		for j := range p.vars {
			if t.s[j].Sign() != 0 {
				buf.WriteByte('*')
				buf.WriteString(p.vars[j])
				if t.s[j].Rat.Cmp(ratOne) != 0 {
					buf.WriteByte('^')
					buf.WriteString(t.s[j].String())
				}
			}
		}
	}
	if len(p.terms) == 0 {
		buf.WriteString("0")
	}
	return buf.String()
}

func (p *Polynomial) MultiCoeff(vars []string, exp []Num) *Polynomial {
	rval := &Polynomial{vars: p.vars}
	idx := rval.indexVars(vars)
	for _, term := range p.terms {
		valid := true
		for i := range idx {
			if idx[i] < 0 || term.s[idx[i]].Cmp(exp[i].Rat) != 0 {
				valid = false
				break
			}
		}
		if valid {
			s := make([]Num, len(rval.vars))
			for i := range s {
				s[i] = Num{new(big.Rat)}
				s[i].Rat.Set(term.s[i].Rat)
			}
			for i := range idx {
				s[idx[i]].Rat.SetInt64(0)
			}
			rval.terms = append(rval.terms, Term{term.c, s})
		}
	}
	return rval
}

func (p *Polynomial) indexVars(vars []string) []int {
	idx := make([]int, len(vars))
	for i := range idx {
		idx[i] = -1
	}
	for i := range vars {
		for j := range p.vars {
			if vars[i] == p.vars[j] {
				idx[i] = j
				break
			}
		}
	}
	return idx
}

func (p *Polynomial) Support(vars []string) [][]Num {
	idx := p.indexVars(vars)
	s := make([][]Num, len(p.terms))
	for i := 0; i < len(p.terms); i++ {
		s[i] = make([]Num, len(vars))
		for j := 0; j < len(s[i]); j++ {
			s[i][j].Rat = new(big.Rat)
			if idx[j] >= 0 {
				s[i][j].Rat.Set(p.terms[i].s[idx[j]].Rat)
			} else {
				s[i][j].Rat.SetInt64(0)
			}
		}
	}
	return s
}

type Term struct {
	c Num
	s []Num
}

type TermOrder func(a, b Term) bool

func LexTermOrder(a, b Term) bool {
	for i := 0; i < len(a.s); i++ {
		x := a.s[i].Rat.Cmp(b.s[i].Rat)
		if x < 0 {
			return true
		} else if x > 0 {
			return false
		}
	}
	return false
}

func TotalTermOrder(a, b Term) bool {
	sumA := big.NewRat(0, 1)
	sumB := big.NewRat(0, 1)
	for i := 0; i < len(a.s); i++ {
		sumA.Add(sumA, a.s[i].Rat)
		sumB.Add(sumB, b.s[i].Rat)
	}
	x := sumA.Cmp(sumB)
	if x < 0 {
		return true
	} else if x > 0 {
		return false
	}
	return LexTermOrder(a, b)
}

type termSorter struct {
	terms []Term
	order TermOrder
}

func (s termSorter) Less(i, j int) bool {
	return s.order(s.terms[i], s.terms[j])
}

func (s termSorter) Swap(i, j int) {
	s.terms[i], s.terms[j] = s.terms[j], s.terms[i]
}

func (s termSorter) Len() int {
	return len(s.terms)
}

func SortTerms(terms []Term, order TermOrder) {
	sort.Sort(termSorter{terms, order})
}

func collectVars(expr Expr) []string {
	vars := make(map[Ident]struct{})
	collectVars2(expr, vars)
	rval := make([]string, 0, len(vars))
	for v := range vars {
		rval = append(rval, string(v))
	}
	sort.Strings(rval)
	return rval
}

func collectVars2(expr Expr, vars map[Ident]struct{}) {
	switch x := expr.(type) {
	case Num:
	case Ident:
		vars[x] = struct{}{}
	case Add:
		collectVars2(x.A, vars)
		collectVars2(x.B, vars)
	case Sub:
		collectVars2(x.A, vars)
		collectVars2(x.B, vars)
	case Mul:
		collectVars2(x.A, vars)
		collectVars2(x.B, vars)
	case Div:
		collectVars2(x.A, vars)
		collectVars2(x.B, vars)
	case Pow:
		collectVars2(x.A, vars)
		collectVars2(x.B, vars)
	case Call:
		collectVars2(x.Args, vars)
	case List:
		for i := range x {
			collectVars2(x[i], vars)
		}
	}
}
