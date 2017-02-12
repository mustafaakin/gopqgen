package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db := sqlx.MustOpen("postgres", "user=postgres dbname=gopqgen sslmode=disable")
	s, err := NewSummaryFromDB(db)

	if err != nil {
		log.Fatal("Could not get an overview of the database")
	}

	//	ps := s.getProtoSummary()

	/*
		for _, v := range ps.Enums {
			println(v.String())
		}

		for _, v := range ps.Msgs {
			println(v.String())
		}
	*/

	ps := s.generateProtoSummary()
	// json.NewEncoder(os.Stdout).Encode(s)
	//	fmt.Println(ps.ToJSON())
	fmt.Println(ps)
}
