package main

import (
	"database/sql"
	"time"

	"github.com/upper/db/v4/adapter/sqlite"
)

const schema = `
DROP TABLE IF EXISTS "pp_birthday";

CREATE TABLE "pp_birthday" (
  "pp_name" varchar(50) DEFAULT NULL,
  "pp_born" DATETIME DEFAULT CURRENT_TIMESTAMP
);`

type Birthday struct {
	Name string    `db:"pp_name"`
	Born time.Time `db:"pp_born"`
}

func mainold() {

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

}
