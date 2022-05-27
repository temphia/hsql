package main

import (
	"database/sql"
	"time"

	"github.com/k0kubun/pp"
	"github.com/upper/db/v4/adapter/sqlite"
)

const schema = `
DROP TABLE IF EXISTS "pp_birthday";

CREATE TABLE "pp_birthday" (
  "pp_name" varchar(50) DEFAULT NULL,
  "pp_born" DATETIME DEFAULT CURRENT_TIMESTAMP
);`

var testConds = []string{
	`SELECT * FROM table_name;`,
	`SELECT a, b from (select a,b,c from tableb);`,
	`SELECT distinct(mno) FROM table_name;`,
	`SELECT count(mno) FROM table_name;`,
	`SELECT COUNT(mno) FROM table_name;`,
	`SELECT Customers.customer_id, Customers.first_name, Orders.amount
	FROM Customers
	INNER JOIN Orders
	ON Customers.customer_id = Orders.customer;`,
}

type Birthday struct {
	Name string    `db:"pp_name"`
	Born time.Time `db:"pp_born"`
}

func main() {

	// Attempt to open the 'example.db' database file
	sess, err := sqlite.Open(sqlite.ConnectionURL{
		Database: `example.db`,
	})

	handlerErr(err)

	driver := sess.Driver().(*sql.DB)
	_, err = driver.Exec(schema)
	handlerErr(err)

	defer sess.Close()

	birthdayCollection := sess.Collection("pp_birthday")

	birthdayCollection.Insert(Birthday{
		Name: "Hayao Miyazaki",
		Born: time.Date(1941, time.January, 5, 0, 0, 0, 0, time.Local),
	})

	birthdayCollection.Insert(Birthday{
		Name: "Nobuo Uematsu",
		Born: time.Date(1959, time.March, 21, 0, 0, 0, 0, time.Local),
	})

	birthdayCollection.Insert(Birthday{
		Name: "Hironobu Sakaguchi",
		Born: time.Date(1962, time.November, 25, 0, 0, 0, 0, time.Local),
	})

	h := NewHSQL(sess, testConds[1])
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
