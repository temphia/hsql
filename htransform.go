package main

import (
	"fmt"

	"github.com/araddon/qlbridge/expr"
	"github.com/araddon/qlbridge/lex"
	"github.com/araddon/qlbridge/rel"
	"github.com/k0kubun/pp"
	"github.com/upper/db/v4"
)

func transform(hctx *HqlCtx, qast *rel.SqlSelect, depth uint8) (db.Selector, error) {
	if depth > hctx.maxDepth {
		panic("Max depth reached")
	}

	var selecter db.Selector
	sql := hctx.sess.SQL()

	if qast.Into != nil {
		panic("Not supported")
	}

	if qast.Distinct {
		cols, err := transformColumns(qast.Columns)
		if err != nil {
			return nil, err
		}
		selecter = sql.Select().Distinct(cols...)
	} else {
		if qast.Star && len(qast.Columns) == 0 {
			selecter = sql.Select("*")
		} else {
			cols, err := transformColumns(qast.Columns)
			if err != nil {
				return nil, err
			}
			selecter = sql.Select(cols...)
		}
	}

	tables := make([]interface{}, 0, len(qast.From))

	pp.Println(qast.From)

	for _, frm := range qast.From {

		if frm.SubQuery != nil {
			s2, err := transform(hctx, frm.SubQuery, depth+1)
			if err != nil {
				return nil, err
			}
			tables = append(tables, s2)
			continue
		}

		if int(frm.Op) == 0 && int(frm.LeftOrRight) == 0 && int(frm.JoinType) == 0 {
			if frm.Alias != "" {
				tables = append(tables, fmt.Sprintf("%s AS %s", frm.Name, frm.Alias))

			} else if frm.Schema == "" {
				tables = append(tables, frm.Name)
			} else {
				tables = append(tables, fmt.Sprintf("%s.%s", frm.Schema, frm.Name))
			}
		}

	}

	selecter = selecter.From(tables...)

	for _, frm := range qast.From {
		if int(frm.JoinType) != 0 {

			switch frm.JoinType {
			case lex.TokenInner:
				nexpr := frm.JoinExpr.(*expr.BinaryNode)

				first := nexpr.Args[0].(*expr.IdentityNode)
				second := nexpr.Args[1].(*expr.IdentityNode)

				fleft, fright, _ := first.LeftRight()
				sleft, sright, _ := second.LeftRight()
				selecter = selecter.Join(fleft).On(fmt.Sprintf("%s.%s = %s.%s", fleft, fright, sleft, sright))

			// case lex.TokenCross:
			// case lex.TokenOuter:
			// case lex.TokenLeft:
			// case lex.TokenRight:

			default:
				panic("Unknown join type")
			}

			pp.Println(frm.JoinType)
		}

	}

	return selecter.From(tables...), nil
}

func transformColumns(cols rel.Columns) ([]interface{}, error) {
	columns := make([]interface{}, 0, len(cols))

	for _, col := range cols {
		if col.CountStar() {
			columns = append(columns, "COUNT(*)")
			continue
		}
		if col.Expr == nil {
			continue
		}
		switch n := col.Expr.(type) {
		case *expr.FuncNode:

			switch n.Name {
			case "tolower":
				pp.Println(n.Args)

			case "count":
				fallthrough
			case "now":
				fallthrough
			default:
				// count(*)
				// now()
				columns = append(columns, n.String())
			}

			// tolower(field_name)
		case *expr.IdentityNode:

			left, right, ok := col.LeftRight()
			if ok {
				columns = append(columns, fmt.Sprintf("%s.%s", left, right))
			} else {
				columns = append(columns, n.String())
			}
		default:
			columns = append(columns, n.String())
		}

	}

	return columns, nil

}
