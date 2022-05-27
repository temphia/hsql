package main

import (
	"fmt"
	"strings"

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
		cols, err := transformColumns(hctx, qast.Columns)
		if err != nil {
			return nil, err
		}
		selecter = sql.Select().Distinct(cols...)
	} else {
		if qast.Star && len(qast.Columns) == 0 {
			selecter = sql.Select("*")
		} else {
			cols, err := transformColumns(hctx, qast.Columns)
			if err != nil {
				return nil, err
			}
			selecter = sql.Select(cols...)
		}
	}

	tables := make([]interface{}, 0, len(qast.From))

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

	if qast.Where != nil {
		if int(qast.Where.Op) == 0 && qast.Where.Source == nil && qast.Where.Expr != nil {
			estr, args, err := extractExpr(hctx, qast.Where.Expr)
			if err != nil {
				return nil, err
			}
			wargs := make([]interface{}, 0, len(args)+1)
			wargs = append(wargs, estr)
			wargs = append(wargs, args...)
			selecter = selecter.Where(wargs...)

		} else {
			panic("handle select a from b where in (select x from subqueryx);")
		}

	}

	if len(qast.GroupBy) != 0 {
		// fixme => can groub by cond can have '?' ? then we also need to pass args([]interface{})

		cols, err := transformColumns(hctx, qast.GroupBy)
		if err != nil {
			return nil, err
		}

		selecter = selecter.GroupBy(cols...)
	}

	if len(qast.OrderBy) != 0 {
		// fixme => can order by cond can have '?' ? then we also need to pass args([]interface{})

		cols, err := transformColumns(hctx, qast.OrderBy)
		if err != nil {
			return nil, err
		}

		selecter = selecter.OrderBy(cols...)
	}

	if qast.Having != nil {
		panic("Having not Supported")
	}

	return selecter, nil
}

func extractExpr(hctx *HqlCtx, e expr.Node) (string, []interface{}, error) {

	switch n := e.(type) {

	case *expr.StringNode:
		return " ? ", []interface{}{n.Text}, nil

	case *expr.NumberNode:
		if n.IsFloat {
			return " ? ", []interface{}{n.Float64}, nil
		} else {
			return " ? ", []interface{}{n.Int64}, nil
		}
	case *expr.FuncNode:

		switch n.Name {
		case "tolower", "count", "now":
			fallthrough
		default:
			pp.Println("FIXME => properly whitelist allowed func")
		}

		var buf strings.Builder

		buf.WriteString(n.Name)
		buf.WriteByte('(')

		args := make([]interface{}, 0)

		for _, arg := range n.Args {
			estr, eargs, err := extractExpr(hctx, arg)
			if err != nil {
				return "", nil, err
			}

			buf.WriteString(estr)
			args = append(args, eargs...)
		}

		return buf.String(), args, nil
	case *expr.IdentityNode:
		left, right, ok := n.LeftRight()
		if ok {
			return fmt.Sprintf("%s.%s", left, right), []interface{}{}, nil
		}
		return n.Text, []interface{}{}, nil
	case *expr.BooleanNode:

	case *expr.NullNode:
		return " NULL ", []interface{}{}, nil
	case *expr.UnaryNode:
		estr, eargs, err := extractExpr(hctx, n.Arg)
		if err != nil {
			return "", nil, err
		}

		return fmt.Sprintf("IS %s", estr), eargs, nil
	case *expr.TriNode:

	case *expr.ArrayNode:
		// fixme = ?
		return " ? ", nil, nil

	case *expr.BinaryNode:

		leftstr, leftargs, err := extractExpr(hctx, n.Args[0])
		if err != nil {
			return "", nil, err
		}
		rightstr, rightargs, err := extractExpr(hctx, n.Args[1])
		if err != nil {
			return "", nil, err
		}

		var fmtstr string

		switch n.Operator.T {
		case lex.TokenMinus: // -
			fmtstr = "%s - %s"
		case lex.TokenPlus: // +
			fmtstr = "%s + %s"
		case lex.TokenDivide: // /
			fmtstr = "%s / %s"
		case lex.TokenMultiply: // *
			fmtstr = "%s * %s"
		case lex.TokenModulus: // %
			fmtstr = "%s % %s"
		case lex.TokenEqual: // =
			fmtstr = "%s = %s"
		case lex.TokenEqualEqual: // ==
			fmtstr = "%s == %s"
		case lex.TokenNE: // !=
			fmtstr = "%s != %s"
		case lex.TokenGE: // >=
			fmtstr = "%s >= %s"
		case lex.TokenLE: // <=
			fmtstr = "%s <= %s"
		case lex.TokenGT: // >
			fmtstr = "%s > %s"
		case lex.TokenLT: // <
			fmtstr = "%s < %s"
		case lex.TokenOr: // ||
			fmtstr = "%s || %s"
		case lex.TokenAnd: // &&
			fmtstr = "%s && %s"
		case lex.TokenBetween: // between
			fmtstr = "%s BETWEEN %s"
		case lex.TokenLogicOr: // OR
			fmtstr = "%s OR %s"
		case lex.TokenLogicAnd: // AND
			fmtstr = "%s AND %s"
		case lex.TokenIN: // IN
			fmtstr = "%s IN %s"
		case lex.TokenLike: // LIKE
			fmtstr = "%s LIKE %s"
		default:
			panic("not implemented")
		}

		farsgs := make([]interface{}, 0, len(leftargs)+len(rightargs))
		farsgs = append(farsgs, leftargs...)
		farsgs = append(farsgs, rightargs...)
		return fmt.Sprintf(fmtstr, leftstr, rightstr), farsgs, nil

	default:
		pp.Println(n)
		panic("Not implemented ")
	}

	return "", nil, nil

}

func transformColumns(hctx *HqlCtx, cols rel.Columns) ([]interface{}, error) {
	columns := make([]interface{}, 0, len(cols))

	for _, col := range cols {
		if col.CountStar() {
			columns = append(columns, "COUNT(*)")
			continue
		}
		if col.Expr == nil {
			continue
		}

		// fixme => donot throw args  (look up)^
		cstr, _, err := extractExpr(hctx, col.Expr)
		if err != nil {
			return nil, err
		}

		columns = append(columns, cstr)
	}

	return columns, nil

}
