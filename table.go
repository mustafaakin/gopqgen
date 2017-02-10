package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type table struct {
	OID  int
	Name string
	Type string
}

func getTableNames(db *sqlx.DB) ([]table, error) {
	var sql = `
    SELECT
      c.oid as "oid",
      c.relname as "name",
      CASE c.relkind
        WHEN 'r' THEN 'table'
        WHEN 'v' THEN 'view'
        WHEN 'm' THEN 'materialized view'
        WHEN 'i' THEN 'index'
        WHEN 'S' THEN 'sequence'
        WHEN 's' THEN 'special'
        WHEN 'f' THEN 'foreign table'
        WHEN 'c' THEN 'composite type'
      END as "type"
    FROM pg_catalog.pg_class c
         LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
    WHERE n.nspname <> 'pg_catalog'
          AND n.nspname <> 'information_schema'
          AND n.nspname !~ '^pg_toast'
      AND pg_catalog.pg_table_is_visible(c.oid)
    ORDER BY 2;
  `

	t := make([]table, 0)
	err := db.Select(&t, sql)
	return t, err
}

type Table struct {
	Name   string
	Type   string
	Fields []field
}

func GetTables(db *sqlx.DB) ([]Table, error) {
	tblNames, err := getTableNames(db)
	if err != nil {
		return nil, err
	}

	tbls := make([]Table, len(tblNames))
	for i, tblName := range tblNames {
		fields, err := getFields(db, tblName.OID)
		if err != nil {
			return nil, err
		}

		tbls[i] = Table{
			Name:   tblName.Name,
			Type:   tblName.Type,
			Fields: fields,
		}
	}

	return tbls, nil
}

// This will require a better mechanism, how about varchar(a).. etc?
// how are they stored? we might need regexes
var protoFieldLookup = map[string]string{
	"uuid":             "string",
	"text":             "string",
	"bytea":            "bytes",
	"boolean":          "bool",
	"integer":          "int32",
	"bigint":           "int64",
	"serial":           "int32",
	"bigserial":        "int64",
	"smallint":         "int32",
	"real":             "float",
	"double precision": "double",
}

func (t Table) toProto() string {
	var s = "message " + t.Name + " {\n"
	for _, field := range t.Fields {
		protoType, ok := protoFieldLookup[field.Type]
		if !ok {
			protoType = field.Type
		}

		s += fmt.Sprintf("  %s %s = %d;\n", protoType, field.Name, field.Num)
	}
	s += "}\n"
	return s
}
