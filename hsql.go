package main

import (
	"github.com/upper/db/v4"

	"github.com/araddon/qlbridge/rel"
)

type HSQL struct {
	db   db.Session
	qstr string
	qast *rel.SqlSelect

	tqstr  string
	tqargs []interface{}
}

func NewHSQL(db db.Session, qstr string) *HSQL {

	return &HSQL{
		db:   db,
		qstr: qstr,
		qast: nil,
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
	s, err := transform(&HqlCtx{
		group:      "xyz",
		groupAlias: "mno",
		tables:     []string{"xy1", "xy2", "xy3", "xy4"},
		sess:       h.db,
		maxDepth:   10,
	}, h.qast, 0)

	if err != nil {
		panic(err)
	}

	h.tqstr = s.String()
	h.tqargs = s.Arguments()

}
