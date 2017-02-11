package main

import (
	"log"

	"github.com/jmoiron/sqlx"
)

/*
NOTE TO ME.
https://www.postgresql.org/docs/9.6/static/catalog-pg-proc.html
* divide inputs and outputs
* make use of proretset
* if output only 1, it uses prorettype
* otherwise, everything is on proallargtypes
* inputs are always proargtypes
* make sure we dont get trigger funcs
* weird but not hard i'm hungry now

SELECT
  ArrayOrder as arrayorder,
  proargnames[ArrayOrder] as name,
  proretset,
  CASE proargmodes[ArrayOrder]
    WHEN 'i' THEN 'input'
    WHEN 'o' THEN 'output'
    WHEN 'v' THEN 'variadic'
    WHEN 't' THEN 'table'
  END AS kind,
  typname as typename
FROM
  pg_proc,
  pg_type,
  UNNEST(proargtypes) WITH ORDINALITY AS T(tag, ArrayOrder)
WHERE
  pg_proc.oid = 17615 AND
  pg_type.oid = tag
ORDER BY
  arrayorder ASC
*/

type definedFunction struct {
	Oid  int
	Name string
	Tip  string
}

type definedFunctionArg struct {
	ArrayOrder int
	Name       string
	Kind       string
	TypeName   string
}

func getDefinedFunctions(db *sqlx.DB) ([]definedFunction, error) {
	var sql = `
SELECT
  p.oid,
  p.proname as "name",
 CASE
  WHEN p.proisagg THEN 'agg'
  WHEN p.proiswindow THEN 'window'
  WHEN p.prorettype = 'pg_catalog.trigger'::pg_catalog.regtype THEN 'trigger'
  ELSE 'normal'
 END as "tip"
FROM pg_catalog.pg_proc p
     LEFT JOIN pg_catalog.pg_namespace n ON n.oid = p.pronamespace
WHERE pg_catalog.pg_function_is_visible(p.oid)
      AND n.nspname <> 'pg_catalog'
      AND n.nspname <> 'information_schema'
ORDER BY 2
`
	args := make([]definedFunction, 0)
	err := db.Select(&args, sql)
	if err != nil {
		panic(err)
	}

	return args, err
}

func getDefinedFunctionArgs(db *sqlx.DB, oid int) ([]definedFunctionArg, error) {
	var sql = `
SELECT
  ArrayOrder as arrayorder,
  proargnames[ArrayOrder] as name,
  CASE proargmodes[ArrayOrder]
    WHEN 'i' THEN 'input'
    WHEN 'o' THEN 'output'
    WHEN 'v' THEN 'variadic'
    WHEN 't' THEN 'table'
  END AS kind,
  typname as typename
FROM
  pg_proc,
  pg_type,
  UNNEST(proallargtypes) WITH ORDINALITY AS T(tag, ArrayOrder)
WHERE
  pg_proc.oid = $1 AND
  pg_type.oid = tag
ORDER BY
 arrayorder
`
	args := make([]definedFunctionArg, 0)
	err := db.Select(&args, sql, oid)

	return args, err
}

type Function struct {
	Name       string
	Type       string // SELECT, UPDATE, INSERT
	Query      string // The SELECT * FROM procedure
	Inputs     []InOutType
	Outputs    []InOutType
	IsOutArray bool
}

func GetUserFunctions(db *sqlx.DB) ([]Function, error) {
	dfns, err := getDefinedFunctions(db)
	if err != nil {
		return nil, err
	}

	fns := make([]Function, 0)
	for _, dfn := range dfns {
		if dfn.Tip == "normal" {
			args, err := getDefinedFunctionArgs(db, dfn.Oid)
			if err != nil {
				return nil, err
			}

			fn := Function{
				Name:       dfn.Name,
				Type:       "SELECT",
				Query:      "SELECT",
				Inputs:     make([]InOutType, 0),
				Outputs:    make([]InOutType, 0),
				IsOutArray: false,
			}

			ins := []string{}
			outs := []string{}

			log.Println(dfn)
			log.Println(args)
			log.Println("---")

			for _, arg := range args {
				if arg.Kind == "input" {
					ins = append(ins, arg.Name)
					fn.Inputs = append(fn.Inputs, InOutType{
						Name: arg.Name,
						Type: arg.TypeName,
					})
				} else if arg.Kind == "output" {
					outs = append(outs, arg.Name)
					fn.Outputs = append(fn.Outputs, InOutType{
						Name: arg.Name,
						Type: arg.TypeName,
					})
				} else {
					// TODO: İS İT THOU?
					fn.IsOutArray = true
					fn.Outputs = append(fn.Outputs, InOutType{
						Name: arg.Name,
						Type: arg.TypeName,
					})
				}
			}

			fns = append(fns, fn)
		}
		// TODO: Others are not cool, don't need triggers, but may need agg and window?
	}

	return fns, nil
}
