package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type enumVal struct {
	SortOrder int
	Label     string
}

func getEnums(db *sqlx.DB) ([]string, error) {
	var sql = `
  SELECT pg_type.typname FROM pg_type WHERE pg_type.typcategory = 'E'
  `
	enums := make([]string, 0)
	err := db.Select(&enums, sql)
	return enums, err
}

func getEnumVals(db *sqlx.DB, typname string) ([]enumVal, error) {
	var sql = `
    SELECT
    	pg_enum.enumsortorder AS "sortorder",
    	pg_enum.enumlabel AS "label"
    FROM pg_type
    JOIN pg_enum
    ON pg_enum.enumtypid = pg_type.oid
    WHERE pg_type.typname = $1
    ORDER BY "sortorder" ASC;
  `

	enums := make([]enumVal, 0)
	err := db.Select(&enums, sql, typname)
	return enums, err
}

func prepareEnums(db *sqlx.DB) []string {
	enums, err := getEnums(db)
	if err != nil {
		log.Fatal("Could not fetch ENUM list:", err)
	}

	for _, enum := range enums {
		src := newSrc(fmt.Sprintf("%s.gen.go", enum))
		// Make it template-wise
		src.addLine("enum hi")
	}

	return nil
}
