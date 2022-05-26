package main

import (
	"fmt"

	"github.com/k0kubun/pp"
	"github.com/upper/db/v4"

	"github.com/araddon/qlbridge/expr"
	"github.com/araddon/qlbridge/rel"
)

type HSQL struct {
	db    db.Session
	depth int
	qstr  string
	qast  *rel.SqlSelect
	w     expr.DialectWriter // fixme => remove this && emit upper sql type instead ?
}

func NewHSQL(db db.Session, qstr string) *HSQL {

	return &HSQL{
		db:    db,
		depth: 10,
		qstr:  qstr,
		qast:  nil,
		w:     nil,
	}
}

func (h *HSQL) Parse() error {
	st, err := rel.ParseSqlSelect(h.qstr)
	if err != nil {
		return err
	}

	h.qast = st

	return nil
}

func (h *HSQL) Transform() {
	s, err := transform(h.qast, 0, h.db)
	if err != nil {
		panic(err)
	}
	pp.Println(s.String())

}

func transform(qast *rel.SqlSelect, depth int, sess db.Session) (db.Selector, error) {

	var selecter db.Selector
	sql := sess.SQL()

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

	if len(qast.From) == 1 {
		frm := qast.From[0]

		if int(frm.Op) == 0 && int(frm.LeftOrRight) == 0 && int(frm.JoinType) == 0 {
			if frm.Alias != "" {
				selecter.From(fmt.Sprintf("%s AS %s", frm.Name, frm.Alias))
			} else if frm.Schema == "" {
				selecter.From(frm.Name)
			} else {
				selecter.From(fmt.Sprintf("%s.%s", frm.Schema, frm.Name))
			}
		} else if int(frm.JoinType) != 0 {
			pp.Println(frm.JoinType.String())
		}

	} else {
		pp.Println(qast.From)
		panic("Does not handle multiple from sources.")
	}

	return selecter, nil
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

			n.IdentityPb()

			left, right, ok := col.LeftRight()
			if ok {
				pp.Println("LEFT | RIGHT", left, right)
				// fixme => check if table is in dtable etc and perm
				// and map to tns format
			} else {
				columns = append(columns, n.String())
			}
		default:
			columns = append(columns, n.String())
		}

	}

	return columns, nil

}
