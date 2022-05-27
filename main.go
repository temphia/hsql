package main

import (
	"github.com/k0kubun/pp"
	"github.com/upper/db/v4/adapter/sqlite"
)

var testConds = []string{
	`SELECT * FROM table_name;`, `SELECT * FROM \"table_name\"`,
	`SELECT a, b from (select a,b,c from tableb);`,
	`SELECT distinct(mno) FROM table_name;`,
	`SELECT count(mno) FROM table_name;`,
	`SELECT COUNT(mno) FROM table_name;`,
	`SELECT Customers.customer_id, Customers.first_name, Orders.amount
	FROM Customers
	INNER JOIN Orders
	ON Customers.customer_id = Orders.customer;`,
}

func main() {

	sess, err := sqlite.Open(sqlite.ConnectionURL{
		Database: `example.db`,
	})

	handlerErr(err)

	h := NewHSQL(sess, testConds[0])
	err = h.Parse()
	handlerErr(err)
	h.Transform()

}

func handlerErr(err error) {
	if err != nil {
		pp.Println(err)
		panic(err)
	}
}
