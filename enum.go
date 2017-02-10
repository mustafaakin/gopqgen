package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type enumVal struct {
	SortOrder int
	Label     string
}

type Enum struct {
	Name   string
	Values []enumVal
}

func getEnumTypes(db *sqlx.DB) ([]string, error) {
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

func GetEnums(db *sqlx.DB) ([]Enum, error) {
	enumTypes, err := getEnumTypes(db)
	if err != nil {
		return nil, nil
	}

	enums := make([]Enum, len(enumTypes))
	for i, enumType := range enumTypes {
		values, err := getEnumVals(db, enumType)
		if err != nil {
			return nil, nil
		}

		enums[i] = Enum{
			Name:   enumType,
			Values: values,
		}
	}

	return enums, nil
}

func (e Enum) toProto() string {
	s := "enum " + e.Name + " {\n"
	s += "  UNKNOWN = 0;\n"
	for _, value := range e.Values {
		s += fmt.Sprintf("  %s = %d;\n", value.Label, value.SortOrder)
	}
	s += "}\n"

	return s
}
