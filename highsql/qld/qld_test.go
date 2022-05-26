package qld

import (
	"log"
	"testing"

	"github.com/k0kubun/pp"
)

var testCases = map[string]string{
	"simple":     "select 1;",
	"s2":         "select `xyz aaa`, mno from pqr;",
	"subquery 1": "select * from test where c in (select c from test2);",

	//"subquery 1": "select * from test where c in (select c from test2 where c<3 limit 1);",
}

func TestQld(t *testing.T) {

	vendor := "genericsql"

	for idx, tc := range testCases {

		tsql, err := transformers[vendor]("ten1", "grp1", tc)
		if err != nil {
			log.Fatal(idx, err)
		}

		pp.Println(tsql)

		err = executeECPGCommand(tsql)
		if err != nil {
			pp.Println("here")
			log.Fatal(err)
		}
	}

}
