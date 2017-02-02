package main

import "github.com/jmoiron/sqlx"

type table struct {
	OID  int
	Name string
	Type string
}

func getTables(db *sqlx.DB) ([]table, error) {
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
    WHERE /* c.relkind IN ('r','')
          AND */ n.nspname <> 'pg_catalog'
          AND n.nspname <> 'information_schema'
          AND n.nspname !~ '^pg_toast'
      AND pg_catalog.pg_table_is_visible(c.oid)
    ORDER BY 2;
  `

	t := make([]table, 0)
	err := db.Select(&t, sql)
	return t, err
}
