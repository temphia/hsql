package highsql

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/araddon/qlbridge/expr"
	"github.com/araddon/qlbridge/rel"
)

var errNotSupported = errors.New("not supported")

func (q *queryctx) transformPQ() (string, error) {

	w := expr.NewDefaultWriter()

	err := transform(q.stmt, 0, w)
	if err != nil {
		return "", err
	}
	return w.String(), nil
}

func transform(qast *rel.SqlSelect, depth int, w expr.DialectWriter) error {

	io.WriteString(w, "SELECT ")
	if qast.Distinct {
		io.WriteString(w, "DISTINCT ")
	}
	qast.Columns.WriteDialect(w)
	if qast.Into != nil {
		return errNotSupported
	}

	if qast.From != nil {
		io.WriteString(w, " FROM")
		for i, from := range qast.From {
			if i == 0 {
				io.WriteString(w, " ")
			} else {
				if from.SubQuery != nil {
					io.WriteString(w, "\n")
					io.WriteString(w, strings.Repeat("\t", depth+1))
				} else {
					io.WriteString(w, "\n")
					io.WriteString(w, strings.Repeat("\t", depth+1))
				}
			}

			transformFrom(from, depth+1, w)
		}
	}
	if qast.Where != nil {
		io.WriteString(w, " WHERE ")
		transformWhere(qast.Where, depth, w)
	}
	if len(qast.GroupBy) > 0 {
		io.WriteString(w, " GROUP BY ")
		qast.GroupBy.WriteDialect(w)
	}
	if qast.Having != nil {
		io.WriteString(w, " HAVING ")
		qast.Having.WriteDialect(w)
	}
	if len(qast.OrderBy) > 0 {
		io.WriteString(w, " ORDER BY ")
		qast.OrderBy.WriteDialect(w)
	}
	if qast.Limit > 0 {
		io.WriteString(w, fmt.Sprintf(" LIMIT %d", qast.Limit))
	}
	if qast.Offset > 0 {
		io.WriteString(w, fmt.Sprintf(" OFFSET %d", qast.Offset))
	}

	return nil
}

func transformFrom(qast *rel.SqlSource, depth int, w expr.DialectWriter) error {
	if int(qast.Op) == 0 && int(qast.LeftOrRight) == 0 && int(qast.JoinType) == 0 {
		if qast.Alias != "" {
			w.WriteIdentity(qast.Name)
			io.WriteString(w, " AS ")
			w.WriteIdentity(qast.Alias)
			return nil
		}
		if qast.Schema == "" {
			w.WriteIdentity(qast.Name)
		} else {
			w.WriteIdentity(qast.Schema)
			io.WriteString(w, ".")
			w.WriteIdentity(qast.Name)
		}
		return nil
	}

	//   Jointype                Op
	//  INNER JOIN orders AS o 	ON
	if int(qast.JoinType) != 0 {
		io.WriteString(w, strings.ToTitle(qast.JoinType.String())) // inner/outer
		io.WriteString(w, " ")
	}
	io.WriteString(w, "JOIN ")

	if qast.SubQuery != nil {
		io.WriteString(w, "(\n"+strings.Repeat("\t", depth+1))
		transform(qast.SubQuery, depth+1, w)
		io.WriteString(w, "\n"+strings.Repeat("\t", depth)+")")
	} else {
		if qast.Schema == "" {
			w.WriteIdentity(qast.Name)
		} else {
			w.WriteIdentity(qast.Schema)
			io.WriteString(w, ".")
			w.WriteIdentity(qast.Name)
		}

	}
	if qast.Alias != "" {
		io.WriteString(w, " AS ")
		w.WriteIdentity(qast.Alias)
	}

	io.WriteString(w, " ")
	io.WriteString(w, strings.ToTitle(qast.Op.String()))

	if qast.JoinExpr != nil {
		w.Write([]byte{' '})
		qast.JoinExpr.WriteDialect(w)
	}

	return nil
}

func transformWhere(qast *rel.SqlWhere, depth int, w expr.DialectWriter) error {
	if int(qast.Op) == 0 && qast.Source == nil && qast.Expr != nil {
		qast.Expr.WriteDialect(w)
		return nil
	}
	// Op = subselect or in etc
	//  SELECT ... WHERE IN (SELECT ...)
	if int(qast.Op) != 0 && qast.Source != nil {
		io.WriteString(w, qast.Op.String())
		io.WriteString(w, " (")
		transform(qast.Source, depth+1, w)
		io.WriteString(w, ")")
		return nil
	}

	return errNotSupported
}
