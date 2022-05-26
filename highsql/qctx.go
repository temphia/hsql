package highsql

import (
	"github.com/araddon/qlbridge/rel"
)

type queryctx struct {
	parent   *highsql
	tenantId string
	groupId  string
	stmt     *rel.SqlSelect
}

func (q *queryctx) transform() (string, error) {
	return "", nil
}

func (q *queryctx) Is(qtype string) bool {
	// tags select, nested=(sub_query | joins ) ?

	switch qtype {
	case "select":
		return true
	default:
		notImplemented()
	}

	return false
}

func (q *queryctx) GetTables(tables ...string) ([]string, error) {
	notImplemented()
	return nil, nil
}

func (q *queryctx) TouchesTable(table string) bool {
	notImplemented()
	return true
}

func (q *queryctx) GetTableColumns(table string) ([]string, error) {
	notImplemented()
	return nil, nil
}

func (q *queryctx) TouchesTableColumns(table string) bool {
	notImplemented()
	return true
}

func (q *queryctx) ForceFilter(table string, filers map[string]interface{}) bool {
	notImplemented()
	return false
}

func notImplemented() {
	panic("not implemented")
}
