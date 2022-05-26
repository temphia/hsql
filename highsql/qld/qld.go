package qld

import (
	"github.com/araddon/qlbridge/rel"
)

type VendorTransformer func(tenantId string, groupId string, query string) (string, error)

var transformers map[string]VendorTransformer = map[string]VendorTransformer{
	"genericsql": pgTransform,
}

type SqlCtx interface {
	Transform() (string, error)
	PreCheck() error
	PostCheck() error
}

type Bridge interface {
	NewCtx(tenantId string, groupId string, query string) (SqlCtx, error)
}

type bridge struct {
	vendor string
}

func (b *bridge) NewCtx(tenantId string, groupId string, query string) (SqlCtx, error) {

	st, err := rel.ParseSqlSelect(query)
	if err != nil {
		return nil, err
	}

	return &sqlCtx{
		bridge:   b,
		stmt:     st,
		tenantId: tenantId,
		groupId:  groupId,
	}, nil
}

type sqlCtx struct {
	bridge   *bridge
	stmt     *rel.SqlSelect
	tenantId string
	groupId  string
}

func (s *sqlCtx) Transform() (string, error) {
	s.stmt.WriteDialect(rel.NewSqlDialect())
	return s.stmt.String(), nil
}

func (s *sqlCtx) PreCheck() error { return nil }

func (s *sqlCtx) PostCheck() error { return nil }

type Env interface {
	Is(string) bool // tags select, with,
	HasTableAccess(...string) bool
	HasColumnAccess(string, ...string) bool
}

///////////////////

type SqlEnv struct {
	tenantId string
	groupId  string
}
