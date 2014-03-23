%{
// Copyright (c) 2014 by Christoph Hack <christoph@tux21b.org>
// All rights reserved. Distributed under the Simplified BSD License.

package main

import (
	"fmt"
	"math/big"
	"unicode"
	"unicode/utf8"
	"strings"
)

%}

%start all

%union{
	val Expr
}

%token <val> NUM ID
%type <val> stmt expr items items2

%right '='
%left '+'  '-'
%left '*'  '/'
%left NEG
%right '^'

%%

all : /* empty */ { yylex.(*Lexer).result = nil }
	| stmt { yylex.(*Lexer).result = $1
}

stmt
	: expr { $$ = $1 }
	| ID '=' expr { $$ = Assign{$1.(Ident), $3}}

expr : NUM { $$ = $1 }
	 | ID { $$ = $1 }
	 | ID '(' items ')' { $$ = Call{$1.(Ident), $3.(List)} }
	 | '(' expr ')' { $$ = $2 }
	 | '[' items ']' { $$ = $2 }
     | expr '+' expr { $$ = Add{$1, $3} }
     | expr '-' expr { $$ = Sub{$1, $3} }
	 | expr '*' expr { $$ = Mul{$1, $3} }
	 | expr '/' expr { $$ = Div{$1, $3} }
	 | '-' expr %prec NEG { $$ = Mul{Num{big.NewRat(-1, 1)}, $2} }
	 | expr '^' expr { $$ = Pow{$1, $3} }
	 ;

items : /* empty */ { $$ = List{} }
	  | items2 { $$ = $1 }
items2 : expr { $$ = List{$1}}
	  | items2 ',' expr { $$ = append($1.(List), $3) }

%%

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