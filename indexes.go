package main

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Index struct {
	OID       int
	Name      string
	Columns   []*Column
	IsUnique  bool
	IsPrimary bool
}

func (i Index) toProto(table string) string {
	var s = "  rpc "
	if i.IsUnique {
		s += "Get"
	} else {
		s += "List"
	}

	cols := make([]string, len(i.Columns))
	for x, col := range i.Columns {
		cols[x] = col.Name
	}

	var by string
	by = strings.Join(cols, "")

	s += fmt.Sprintf("%sBy%s ( %s ) returns (%s) {}\n",
		table, by, strings.Join(cols, ", "), table,
	)

	return s
}

type Column struct {
	Name     string
	Typename string
}

func getIndexNames(db *sqlx.DB, oid int) ([]*Index, error) {
	var sql = `
SELECT
	idx.indexrelid as oid,
	i.relname as name,
    idx.indisunique as isunique,
    idx.indisprimary as isprimary
FROM   pg_index as idx
JOIN   pg_class as i
ON     i.oid = idx.indexrelid
JOIN   pg_am as am
ON     i.relam = am.oid
JOIN   pg_namespace as ns
ON     ns.oid = i.relnamespace
AND    ns.nspname = ANY(current_schemas(false))
AND    idx.indrelid = $1
	`
	is := make([]*Index, 0)
	err := db.Select(&is, sql, oid)
	return is, err
}

func getIndexColumns(db *sqlx.DB, oid int) ([]*Column, error) {
	var sql = `
SELECT
  a.attname as name,
  pg_catalog.format_type(a.atttypid, a.atttypmod) as typename
FROM pg_catalog.pg_attribute a
WHERE a.attrelid = $1 AND a.attnum > 0 AND NOT a.attisdropped
ORDER BY a.attnum;
`
	cols := make([]*Column, 0)
	err := db.Select(&cols, sql, oid)
	return cols, err
}

func GetIndices(db *sqlx.DB, tableOid int) ([]*Index, error) {
	inds, err := getIndexNames(db, tableOid)
	if err != nil {
		return inds, err
	}

	for _, ind := range inds {
		cols, err := getIndexColumns(db, ind.OID)
		if err != nil {
			return nil, err
		}
		ind.Columns = cols
	}

	return inds, nil
}
