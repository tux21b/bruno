//line grammar.y:1

// Copyright (c) 2014 by Christoph Hack <christoph@tux21b.org>
// All rights reserved. Distributed under the Simplified BSD License.

package main

import __yyfmt__ "fmt"

//line grammar.y:6
import (
	"fmt"
	"math/big"
	"strings"
	"unicode"
	"unicode/utf8"
)

//line grammar.y:19
type yySymType struct {
	yys int
	val Expr
}

const NUM = 57346
const ID = 57347
const NEG = 57348

var yyToknames = []string{
	"NUM",
	"ID",
	" =",
	" +",
	" -",
	" *",
	" /",
	"NEG",
	" ^",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line grammar.y:60

type Lexer struct {
	input  string
	pos    int
	err    error
	result Expr
}

func (l *Lexer) Lex(lval *yySymType) int {
	for l.pos < len(l.input) {
		r, n := utf8.DecodeRuneInString(l.input[l.pos:])
		switch {
		case strings.ContainsRune("+-*/^(),[]=", r):
			l.pos += n
			return int(r)
		case unicode.IsSpace(r):
			l.pos += n
		case unicode.IsDigit(r) || r == '.':
			i := l.pos
			dot := false
			for {
				r, n = utf8.DecodeRuneInString(l.input[i:])
				if r == '.' && !dot {
					dot = true
				} else if !unicode.IsDigit(r) {
					break
				}
				i += n
			}
			v := new(big.Rat)
			v.SetString(l.input[l.pos:i])
			lval.val = Num{v}
			l.pos = i
			return NUM
		case unicode.IsLetter(r):
			i := l.pos + n
			for {
				r, n = utf8.DecodeRuneInString(l.input[i:])
				if !unicode.IsDigit(r) && !unicode.IsLetter(r) {
					break
				}
				i += n
			}
			lval.val = Ident(l.input[l.pos:i])
			l.pos = i
			return ID
		default:
			return int(r)
		}
	}
	return 0
}

func (l *Lexer) Error(s string) {
	l.err = SyntaxError(s)
}

type SyntaxError string

func (s SyntaxError) Error() string {
	return fmt.Sprintf("syntax error: %v", string(s))
}

func Parse(input string) (Expr, error) {
	l := &Lexer{input: input}
	yyParse(l)
	if l.err != nil {
		return nil, l.err
	}
	return l.result, nil
}

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 20
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 47

var yyAct = []int{

	20, 3, 31, 9, 10, 11, 12, 16, 13, 21,
	22, 23, 24, 25, 26, 27, 5, 17, 30, 18,
	8, 9, 10, 11, 12, 6, 13, 7, 29, 5,
	4, 32, 33, 8, 14, 28, 11, 12, 6, 13,
	7, 15, 15, 13, 19, 2, 1,
}
var yyPact = []int{

	25, -1000, -1000, -4, 28, -1000, 12, 12, 12, 12,
	12, 12, 12, 12, 12, 12, 14, 29, 2, -15,
	-4, 31, 27, 27, 31, 31, 31, -4, 17, -1000,
	-1000, 12, -1000, -4,
}
var yyPgo = []int{

	0, 46, 45, 0, 19, 44,
}
var yyR1 = []int{

	0, 1, 1, 2, 2, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 4, 4, 5, 5,
}
var yyR2 = []int{

	0, 0, 1, 1, 3, 1, 1, 4, 3, 3,
	3, 3, 3, 3, 2, 3, 0, 1, 1, 3,
}
var yyChk = []int{

	-1000, -1, -2, -3, 5, 4, 13, 15, 8, 7,
	8, 9, 10, 12, 6, 13, -3, 5, -4, -5,
	-3, -3, -3, -3, -3, -3, -3, -3, -4, 14,
	16, 17, 14, -3,
}
var yyDef = []int{

	1, -2, 2, 3, 6, 5, 0, 16, 0, 0,
	0, 0, 0, 0, 0, 16, 0, 6, 0, 17,
	18, 14, 10, 11, 12, 13, 15, 4, 0, 8,
	9, 0, 7, 19,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	13, 14, 9, 7, 17, 8, 3, 10, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 6, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 15, 3, 16, 12,
}
var yyTok2 = []int{

	2, 3, 4, 5, 11,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(yyToknames) {
		if yyToknames[c-4] != "" {
			return yyToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(c), uint(char))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		//line grammar.y:34
		{
			yylex.(*Lexer).result = nil
		}
	case 2:
		//line grammar.y:35
		{
			yylex.(*Lexer).result = yyS[yypt-0].val
		}
	case 3:
		//line grammar.y:39
		{
			yyVAL.val = yyS[yypt-0].val
		}
	case 4:
		//line grammar.y:40
		{
			yyVAL.val = Assign{yyS[yypt-2].val.(Ident), yyS[yypt-0].val}
		}
	case 5:
		//line grammar.y:42
		{
			yyVAL.val = yyS[yypt-0].val
		}
	case 6:
		//line grammar.y:43
		{
			yyVAL.val = yyS[yypt-0].val
		}
	case 7:
		//line grammar.y:44
		{
			yyVAL.val = Call{yyS[yypt-3].val.(Ident), yyS[yypt-1].val.(List)}
		}
	case 8:
		//line grammar.y:45
		{
			yyVAL.val = yyS[yypt-1].val
		}
	case 9:
		//line grammar.y:46
		{
			yyVAL.val = yyS[yypt-1].val
		}
	case 10:
		//line grammar.y:47
		{
			yyVAL.val = Add{yyS[yypt-2].val, yyS[yypt-0].val}
		}
	case 11:
		//line grammar.y:48
		{
			yyVAL.val = Sub{yyS[yypt-2].val, yyS[yypt-0].val}
		}
	case 12:
		//line grammar.y:49
		{
			yyVAL.val = Mul{yyS[yypt-2].val, yyS[yypt-0].val}
		}
	case 13:
		//line grammar.y:50
		{
			yyVAL.val = Div{yyS[yypt-2].val, yyS[yypt-0].val}
		}
	case 14:
		//line grammar.y:51
		{
			yyVAL.val = Mul{Num{big.NewRat(-1, 1)}, yyS[yypt-0].val}
		}
	case 15:
		//line grammar.y:52
		{
			yyVAL.val = Pow{yyS[yypt-2].val, yyS[yypt-0].val}
		}
	case 16:
		//line grammar.y:55
		{
			yyVAL.val = List{}
		}
	case 17:
		//line grammar.y:56
		{
			yyVAL.val = yyS[yypt-0].val
		}
	case 18:
		//line grammar.y:57
		{
			yyVAL.val = List{yyS[yypt-0].val}
		}
	case 19:
		//line grammar.y:58
		{
			yyVAL.val = append(yyS[yypt-2].val.(List), yyS[yypt-0].val)
		}
	}
	goto yystack /* stack new state and value */
}
