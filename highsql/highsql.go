package highsql

import (
	"github.com/araddon/qlbridge/rel"
)

// highsql is a mysql specific sql format that maps to different db vendors
type highsql struct {
	vendor         string
	queryTransform func(*rel.SqlSelect) (string, error)
}
