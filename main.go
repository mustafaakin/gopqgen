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

	prepareEnums(db)
	/*
	   	enums, err := getEnums(db)
	   	if err != nil {
	   		log.Fatal("Could not list ENUMs: ", err)
	   	}

	   	for _, enum := range enums {
	   		log.Println(enum)
	   	}

	   	tables, err := getTables(db)
	   	if err != nil {
	   		log.Println("Could not list TABLEs: ", err)
	   	}

	   	for _, table := range tables {
	       if table.Name ==
	     	log.Println(table.Type, table.Name)
	   		fields, err := getFields(db, table.OID)
	   		if err != nil {
	   			log.Fatal("Could not get fields for", table, err)
	   		}
	   		for _, field := range fields {
	   			log.Println("\t", field.Name, field.Type)
	   		}
	   	}

	*/
}
