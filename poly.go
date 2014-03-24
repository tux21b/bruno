// Copyright (c) 2014 by Christoph Hack <christoph@tux21b.org>
// All rights reserved. Distributed under the Simplified BSD License.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"sort"
)

type Polynomial struct {
	vars  []string
	order TermOrder
	items []Monomial
}

func NewPolynomial(expr Expr) (*Polynomial, error) {
	if p, ok := expr.(*Polynomial); ok {
		return p, nil
	}
	p := &Polynomial{order: LexTermOrder}
	p.vars = collectVars(expr)
	if err := p.convert(expr); err != nil {
		return nil, err
	}
	p.normalize()
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
	m := Monomial{*big.NewRat(1, 1), make([]big.Rat, len(p.vars))}
	if err := p.convertMonomial(expr, &m); err != nil {
		return err
	}
	p.items = append(p.items, m)
	p.normalize()
	return nil
}

func (p *Polynomial) convertMonomial(expr Expr, m *Monomial) error {
	switch x := expr.(type) {
	case Num:
		m.C.Mul(&m.C, x.Rat)
		return nil
	case Mul:
		if err := p.convertMonomial(x.A, m); err != nil {
			return err
		}
		if err := p.convertMonomial(x.B, m); err != nil {
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
			m.T[idx].Add(&m.T[idx], exp.Rat)
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
			m.T[idx].Add(&m.T[idx], ratOne)
			return nil
		}
	}
	return errors.New("invalid polynomial")
}

var ratZero = big.NewRat(0, 1)
var ratOne = big.NewRat(1, 1)

func (p *Polynomial) String() string {
	buf := &bytes.Buffer{}
	for i, t := range p.items {
		if i > 0 {
			buf.WriteString(" + ")
		}
		buf.WriteString(t.C.RatString())
		for j := range p.vars {
			if t.T[j].Sign() != 0 {
				buf.WriteByte('*')
				buf.WriteString(p.vars[j])
				if t.T[j].Cmp(ratOne) != 0 {
					buf.WriteByte('^')
					buf.WriteString(t.T[j].RatString())
				}
			}
		}
	}
	if len(p.items) == 0 {
		buf.WriteString("0")
	}
	return buf.String()
}

func (p *Polynomial) MultiCoeff(vars []string, exp []Num) *Polynomial {
	rval := &Polynomial{vars: p.vars, order: p.order}
	idx := rval.indexVars(vars)
	for _, term := range p.items {
		valid := true
		for i := range idx {
			if idx[i] < 0 || term.T[idx[i]].Cmp(exp[i].Rat) != 0 {
				valid = false
				break
			}
		}
		if valid {
			m := Monomial{term.C, make(Term, len(rval.vars))}
			for i := range m.T {
				m.T[i].Set(&term.T[i])
			}
			for i := range idx {
				m.T[idx[i]].SetInt64(0)
			}
			rval.items = append(rval.items, m)
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
	s := make([][]Num, len(p.items))
	for i := 0; i < len(p.items); i++ {
		s[i] = make([]Num, len(vars))
		for j := 0; j < len(s[i]); j++ {
			s[i][j].Rat = new(big.Rat)
			if idx[j] >= 0 {
				s[i][j].Rat.Set(&p.items[i].T[idx[j]])
			} else {
				s[i][j].Rat.SetInt64(0)
			}
		}
	}
	return s
}

func (p *Polynomial) LPP() *Polynomial {
	rval := &Polynomial{vars: p.vars, order: p.order}
	if len(p.items) > 0 {
		rval.items = append(rval.items, Monomial{*ratOne, p.items[0].T})
	}
	return rval
}

func (p *Polynomial) LC() Num {
	if len(p.items) > 0 {
		return Num{&p.items[0].C}
	}
	return Num{big.NewRat(0, 1)}
}

func (p *Polynomial) LM() *Polynomial {
	rval := &Polynomial{vars: p.vars, order: p.order}
	if len(p.items) > 0 {
		rval.items = p.items[:1]
	}
	return rval
}

func (p *Polynomial) Higher(t Term) *Polynomial {
	n := sort.Search(len(p.items), func(i int) bool {
		return !p.order(t, p.items[i].T)
	})
	return &Polynomial{vars: p.vars, order: p.order, items: p.items[:n]}
}

func (p *Polynomial) Lower(t Term) *Polynomial {
	n := sort.Search(len(p.items), func(i int) bool {
		return p.order(p.items[i].T, t)
	})
	return &Polynomial{vars: p.vars, order: p.order, items: p.items[n:]}
}

func (p *Polynomial) Between(t1, t2 Term) *Polynomial {
	return p.Lower(t2).Higher(t1)
}

func (p *Polynomial) Remainder() *Polynomial {
	rval := &Polynomial{vars: p.vars, order: p.order}
	if len(p.items) > 0 {
		rval.items = p.items[1:]
	}
	return rval
}

func (p *Polynomial) ReduceTerm(f *Polynomial, t Term) (*Polynomial, error) {
	if len(p.vars) != len(t) {
		return nil, fmt.Errorf("invalid term")
	}
	if len(f.items) == 0 {
		return nil, fmt.Errorf("invalid polynomial f")
	}
	idx := p.findTerm(t)
	if idx < 0 {
		return nil, fmt.Errorf("invalid term (not in support)")
	}
	u := Monomial{T: make(Term, len(p.vars))}
	u.C.Quo(&p.items[idx].C, &f.items[0].C)
	u.C.Neg(&u.C)
	for j := 0; j < len(p.vars); j++ {
		u.T[j].Sub(&p.items[idx].T[j], &f.items[0].T[j])
	}
	h := &Polynomial{vars: p.vars, order: p.order}
	for i := 0; i < len(f.items); i++ {
		m := Monomial{T: make(Term, len(p.vars))}
		m.C.Mul(&u.C, &f.items[i].C)
		for j := 0; j < len(p.vars); j++ {
			m.T[j].Add(&u.T[j], &f.items[i].T[j])
		}
		h.items = append(h.items, m)
	}
	for i := 0; i < len(p.items); i++ {
		pos := h.findTerm(p.items[i].T)
		if pos < 0 {
			h.items = append(h.items, p.items[i])
		} else {
			h.items[pos].C.Add(&h.items[pos].C, &p.items[i].C)
		}
	}
	if !h.valid() {
		return nil, fmt.Errorf("invalid reduction %v", h)
	}
	h.normalize()
	return h, nil
}

func (p *Polynomial) Reduce(f *Polynomial) *Polynomial {
	for i := 0; i < len(p.items); i++ {
		if h, err := p.ReduceTerm(f, p.items[i].T); err == nil {
			return h
		}
	}
	return p
}

func (p *Polynomial) ReduceAny(fns []*Polynomial) *Polynomial {
	for i := range fns {
		h := p.Reduce(fns[i])
		if !p.Equal(h) {
			fmt.Printf("reduced by %v to %v\n", fns[i], h)
			return h
		}
	}
	return p
}

func (p *Polynomial) Equal(q *Polynomial) bool {
	if p == q {
		return true
	}
	if len(p.vars) != len(q.vars) || len(p.items) != len(q.items) {
		return false
	}
	for i := 0; i < len(p.items); i++ {
		if p.items[i].C.Cmp(&q.items[i].C) != 0 {
			return false
		}
		for j := 0; j < len(p.vars); j++ {
			if p.items[i].T[j].Cmp(&q.items[i].T[j]) != 0 {
				return false
			}
		}
	}
	return true
}

func (p *Polynomial) normalize() {
	for i := 0; i < len(p.items); i++ {
		if p.items[i].C.Sign() == 0 {
			n := len(p.items) - 1
			p.items[i], p.items[n] = p.items[n], p.items[i]
			p.items = p.items[:n]
			i--
		}
	}
	SortMonomial(p.items, p.order)
}

func (p *Polynomial) valid() bool {
	for i := 0; i < len(p.items); i++ {
		for j := 0; j < len(p.vars); j++ {
			if p.items[i].T[j].Sign() < 0 {
				return false
			}
		}
	}
	return true
}

func (p *Polynomial) findTerm(t Term) int {
	i := sort.Search(len(p.items), func(i int) bool {
		return !p.order(t, p.items[i].T)
	})
	if i < 0 || i >= len(p.items) {
		return -1
	}
	for j := 0; j < len(p.vars); j++ {
		if p.items[i].T[j].Cmp(&t[j]) != 0 {
			return -1
		}
	}
	return i
}

type Monomial struct {
	C big.Rat
	T Term
}

type Term []big.Rat

type TermOrder func(a, b Term) bool

func LexTermOrder(a, b Term) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		x := a[i].Cmp(&b[i])
		if x < 0 {
			return true
		} else if x > 0 {
			return false
		}
	}
	return false
}

func LexTermOrderRev(a, b Term) bool {
	if len(a) != len(b) {
		return false
	}
	for i := len(a) - 1; i >= 0; i-- {
		x := a[i].Cmp(&b[i])
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
	for i := 0; i < len(a); i++ {
		sumA.Add(sumA, &a[i])
		sumB.Add(sumB, &b[i])
	}
	x := sumA.Cmp(sumB)
	if x < 0 {
		return true
	} else if x > 0 {
		return false
	}
	return LexTermOrder(a, b)
}

type monomialSorter struct {
	items []Monomial
	order TermOrder
}

func (s monomialSorter) Less(i, j int) bool {
	return s.order(s.items[j].T, s.items[i].T)
}

func (s monomialSorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

func (s monomialSorter) Len() int {
	return len(s.items)
}

func SortMonomial(items []Monomial, order TermOrder) {
	sort.Sort(monomialSorter{items, order})
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
