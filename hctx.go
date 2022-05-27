package main

import (
	"fmt"

	"github.com/upper/db/v4"
)

type TNS interface {
	TableName(tenant, group, table string) string
}

type HqlCtx struct {
	group      string
	groupAlias string
	tables     []string
	sess       db.Session
	tns        TNS
	maxDepth   uint8
}

type tns struct{}

func (t *tns) TableName(tenant, group, table string) string {
	return fmt.Sprintf("%s_%s_%s", tenant, group, table)
}
