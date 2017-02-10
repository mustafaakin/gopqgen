package main

import "github.com/jmoiron/sqlx"

type Index struct {
	OID      string
	Name     string
	Columns  []string
	IsUnique bool
}

func getIndexNames(db *sqlx.DB, oid int) ([]Index, error) {
	var sql = `
   SELECT
    i.OID,
    i.relname AS name
   FROM pg_index x
     JOIN pg_class c ON c.oid = x.indrelid
     JOIN pg_class i ON i.oid = x.indexrelid
     LEFT JOIN pg_namespace n ON n.oid = c.relnamespace
     LEFT JOIN pg_tablespace t ON t.oid = i.reltablespace
  WHERE (c.relkind = ANY (ARRAY['r'::"char", 'm'::"char"])) AND i.relkind = 'i'::"char";
	`
	is := make([]Index, 0)
	err := db.Select(&is, sql, oid)
	return is, err
}

func (i *Index) getColumns(db *sqlx.DB) error {

}
