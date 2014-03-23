// Copyright (c) 2014 by Christoph Hack <christoph@tux21b.org>
// All rights reserved. Distributed under the Simplified BSD License.

package main

import (
	"testing"
)

var brunoTests = []struct {
	input  string
	output string
}{
	{
		"3",
		"3",
	},
	{
		"identifier",
		"identifier",
	},
	{
		"3*x + 6 * y ^ 2",
		"((3 * x) + (6 * (y ^ 2)))",
	},
	{
		"[1, 2, a, b, c]",
		"[1 2 a b c]",
	},
	{
		"q = 3",
		"q = 3",
	},
	{
		"multicoeff(18*x^2*y + y*z, [x, y, z], [2, 1, 0])",
		"18",
	},
	{
		"multicoeff(18*x^2*y + 19*x^2*y*z + y*z, [x, y], [2, 1])",
		"19*z + 18",
	},
	{
		"multicoeff(18*x^2*y + 19*x^2*y*z + y*z, [x, y, z], [2, 1, 0])",
		"18",
	},
	{
		"f = p(3*x*y^2 + 8*x^2 + 7 + 20*x*y + 3*y^10)",
		"f = 8*x^2 + 3*x*y^2 + 20*x*y + 3*y^10 + 7",
	},
	{
		"totalorder(f)",
		"3*y^10 + 3*x*y^2 + 8*x^2 + 20*x*y + 7",
	},
	{
		"lexorder(f)",
		"8*x^2 + 3*x*y^2 + 20*x*y + 3*y^10 + 7",
	},
	{
		"lpp(f)",
		"1*x^2",
	},
	{
		"lc(f)",
		"8",
	},
	{
		"lm(f)",
		"8*x^2",
	},
	{
		"lm(totalorder(f))",
		"3*y^10",
	},
	{
		"support(lexorder(f), [x, y])",
		"[[2 0] [1 2] [1 1] [0 10] [0 0]]",
	},
	{
		"support(f, [y])",
		"[[0] [2] [1] [10] [0]]",
	},
}

func TestBruno(t *testing.T) {
	bruno := NewBruno()
	for i := range brunoTests {
		result, err := bruno.Exec(brunoTests[i].input)
		if err != nil {
			t.Errorf("test %q: unexpected error %v.", brunoTests[i].input, err)
			continue
		}
		if output := result.String(); output != brunoTests[i].output {
			t.Errorf("test %q: expected output %q, got %q.",
				brunoTests[i].input, brunoTests[i].output, output)
		}
	}
}
