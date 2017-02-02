package main

import "github.com/jmoiron/sqlx"

type field struct {
	Name    string
	Type    string
	NotNull bool
	Num     int
}

func getFields(db *sqlx.DB, oid int) ([]field, error) {
	var sql = `
  SELECT a.attname as "name",
    pg_catalog.format_type(a.atttypid, a.atttypmod) as "type",
    a.attnotnull as "notnull",
    a.attnum as "num"
  FROM pg_catalog.pg_attribute a
  WHERE a.attrelid = $1 AND a.attnum > 0 AND NOT a.attisdropped
  ORDER BY a.attnum;
  `

	fs := make([]field, 0)
	err := db.Select(&fs, sql, oid)
	return fs, err
}
