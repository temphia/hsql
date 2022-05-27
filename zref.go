package main

import (
	"github.com/araddon/qlbridge/expr"
	"github.com/araddon/qlbridge/lex"
	"github.com/k0kubun/pp"
)

func a() {
	x := &expr.BinaryNode{

		Paren: false,
		Args: []expr.Node{
			&expr.BinaryNode{

				Paren: false,
				Args: []expr.Node{
					&expr.BinaryNode{

						Paren: false,
						Args: []expr.Node{
							&expr.IdentityNode{
								Quote: 0x00,
								Text:  "mpr",
							},
							&expr.NumberNode{
								IsInt:   true,
								IsFloat: true,
								Int64:   12,
								Float64: 12.000000,
								Text:    "12",
							},
						},
						Operator: lex.Token{
							T:      0x0043,
							V:      "=",
							Quote:  0x00,
							Line:   1,
							Column: 35,
							Pos:    35,
						},
					},
					&expr.BinaryNode{
						Paren: false,
						Args: []expr.Node{
							&expr.IdentityNode{
								Quote: 0x00,
								Text:  "mno",
							},
							&expr.StringNode{
								Quote: 0x27,
								Text:  "12",
							},
						},
						Operator: lex.Token{
							T:      0x0049,
							V:      "<",
							Quote:  0x00,
							Line:   1,
							Column: 47,
							Pos:    47,
						},
					},
				},
				Operator: lex.Token{
					T:      0x004f,
					V:      "AND",
					Quote:  0x00,
					Line:   1,
					Column: 41,
					Pos:    41,
				},
			},
			&expr.BinaryNode{
				Paren: true,
				Args: []expr.Node{
					&expr.BinaryNode{
						Paren: false,
						Args: []expr.Node{
							&expr.IdentityNode{
								Quote: 0x00,
								Text:  "pqr1",
							},
							&expr.NumberNode{
								IsInt:   true,
								IsFloat: true,
								Int64:   23,
								Float64: 23.000000,
								Text:    "23",
							},
						},
						Operator: lex.Token{
							T:      0x0043,
							V:      "=",
							Quote:  0x00,
							Line:   1,
							Column: 64,
							Pos:    64,
						},
					},
					&expr.UnaryNode{
						Arg: &expr.NullNode{},
						Operator: lex.Token{
							T:      0x0057,
							V:      "IS",
							Quote:  0x00,
							Line:   1,
							Column: 78,
							Pos:    78,
						},
					},
				},
				Operator: lex.Token{
					T:      0x004e,
					V:      "OR",
					Quote:  0x00,
					Line:   1,
					Column: 70,
					Pos:    70,
				},
			},
		},
		Operator: lex.Token{
			T:      0x004f,
			V:      "AND",
			Quote:  0x00,
			Line:   1,
			Column: 56,
			Pos:    56,
		},
	}

	pp.Println(x)

}
