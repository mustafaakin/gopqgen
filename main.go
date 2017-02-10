package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Open("postgres", "user=armada dbname=armada sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	es, _ := GetEnums(db)
	for _, e := range es {
		println(e.toProto())
	}

	ts, _ := GetTables(db)
	for _, t := range ts {
		println(t.toProto())
	}
}
