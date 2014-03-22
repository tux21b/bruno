// Copyright (c) 2014 by Christoph Hack <christoph@tux21b.org>
// All rights reserved. Distributed under the Simplified BSD License.

package main

import (
	"fmt"
	"math/big"
)

type Expr interface {
	String() string
}

type Num struct {
	*big.Rat
}

func (n Num) String() string {
	return n.Rat.RatString()
}

type Ident string

func (i Ident) String() string {
	return string(i)
}

type Add struct {
	A, B Expr
}

func (a Add) String() string {
	return fmt.Sprintf("(%v + %v)", a.A, a.B)
}

type Sub struct {
	A, B Expr
}

func (s Sub) String() string {
	return fmt.Sprintf("(%v - %v)", s.A, s.B)
}

type Mul struct {
	A, B Expr
}

func (m Mul) String() string {
	return fmt.Sprintf("(%v * %v)", m.A, m.B)
}

type Div struct {
	A, B Expr
}

func (d Div) String() string {
	return fmt.Sprintf("(%v / %v)", d.A, d.B)
}

type Pow struct {
	A, B Expr
}

func (p Pow) String() string {
	return fmt.Sprintf("(%v ^ %v)", p.A, p.B)
}

type List []Expr

func (l List) String() string {
	return fmt.Sprintf("%v", []Expr(l))
}

type Call struct {
	Ident Ident
	Args  List
}

func (c Call) String() string {
	return fmt.Sprintf("%v(%v)", c.Ident, c.Args)
}

type Assign struct {
	Ident Ident
	Expr  Expr
}

func (a Assign) String() string {
	return fmt.Sprintf("%s = %v", a.Ident, a.Expr)
}
