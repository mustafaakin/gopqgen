package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Open("postgres", "user=postgres dbname=armada sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	println("// Enums")
	es, _ := GetEnums(db)
	for _, e := range es {
		println(e.toProto())
	}

	println("// Tables")
	ts, _ := GetTables(db)

	// TOOD: Make it an interface so that we can call index, view, or user defined funcs as well
	funcs := make([]string, 0)

	for _, t := range ts {
		println(t.toProto())

		inds, err := GetIndices(db, t.OID)
		if err != nil {
			log.Fatal(err)
		}

		for _, ind := range inds {
			funcs = append(funcs, ind.toProto(t.Name))
		}
	}

	println("service DBService {")
	for _, fn := range funcs {
		print(fn)
	}
	println("}")
}
