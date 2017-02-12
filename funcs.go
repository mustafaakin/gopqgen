package main

import (
	"fmt"
	"log"

	"strings"

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
	fns := make([]definedFunction, 0)
	err := db.Select(&fns, sql)
	if err != nil {
		panic(err)
	}

	return fns, err
}

type definedFunctionInputArg struct {
	ArrayOrder int
	Name       *string
	TypeName   string
}

func getDefinedFuncArgs(db *sqlx.DB, oid int) ([]definedFunctionInputArg, error) {
	var sql = `
SELECT
  ArrayOrder as arrayorder,
  proargnames[ArrayOrder] as name,
  typname as typename
FROM
  pg_proc,
  pg_type,
  UNNEST(proargtypes) WITH ORDINALITY AS T(tag, ArrayOrder)
WHERE
  pg_proc.oid = $1 AND
  pg_type.oid = tag
ORDER BY
  arrayorder ASC
`

	args := make([]definedFunctionInputArg, 0)
	err := db.Select(&args, sql, oid)
	if err != nil {
		return nil, err
	}

	return args, err
}

type definedFuncOutputArg struct {
	ArrayOrder  int
	Name        *string
	Argmode     string
	IsReturnSet bool
	TypeName    string
}

func getDefinedFuncOutputArgs(db *sqlx.DB, oid int) ([]definedFuncOutputArg, error) {
	var sql = `
  SELECT
    ArrayOrder as arrayorder,
    proargnames[ArrayOrder] as name,
    proargmodes[ArrayOrder] as argmode,
    proretset as isreturnset,
    typname as typename
  FROM
    pg_proc,
    pg_type,
    UNNEST(proallargtypes) WITH ORDINALITY AS T(tag, ArrayOrder)
  WHERE
    pg_proc.oid = $1 AND
    pg_type.oid = tag AND
    proargmodes[ArrayOrder] IN ('t','o')
  ORDER BY
    arrayorder ASC`

	outs := make([]definedFuncOutputArg, 0)
	err := db.Select(&outs, sql, oid)
	if err != nil {
		return nil, err
	}

	// Is really no output?
	if len(outs) == 0 {
		sql = `
  SELECT
    1 as arrayorder,
    proretset as isreturnset,
    typname as typename
  FROM
    pg_proc,
    pg_type
  WHERE
    pg_proc.oid = $1 AND
    pg_type.oid = prorettype
    `

		err = db.Select(&outs, sql, oid)
		if err != nil {
			return nil, err
		}
	}

	return outs, err

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
			fn := Function{
				Name:       capitalize(dfn.Name),
				Type:       "SELECT",
				Inputs:     make([]InOutType, 0),
				Outputs:    make([]InOutType, 0),
				IsOutArray: false,
			}

			// Process Input args
			args, err := getDefinedFuncArgs(db, dfn.Oid)
			if err != nil {
				return nil, err
			}

			argsStr := []string{}
			for i, arg := range args {
				argsStr = append(argsStr, fmt.Sprintf("$%d", i+1))
				var name string
				if arg.Name != nil {
					name = *(arg.Name)
				} else {
					name = ""
				}

				fn.Inputs = append(fn.Inputs, InOutType{
					Name: name,
					Type: arg.TypeName,
				})
			}

			// Process Output args
			outs, err := getDefinedFuncOutputArgs(db, dfn.Oid)
			if err != nil {
				log.Fatal(err)
				return nil, nil
			}

			for _, out := range outs {
				var name string
				if out.Name != nil {
					name = *(out.Name)
				} else {
					name = ""
				}

				// All outputs same, set or not-set, its enough to check only one
				if out.IsReturnSet {
					fn.IsOutArray = true
				}

				fn.Outputs = append(fn.Outputs, InOutType{
					Name: name,
					Type: out.TypeName,
				})
			}

			fn.Query = fmt.Sprintf(
				"SELECT * FROM %s(%s)",
				dfn.Name,
				strings.Join(argsStr, ", "),
			)

			fns = append(fns, fn)
		}
		// TODO: Others are not cool, don't need triggers, but may need agg and window?
	}

	return fns, nil
}
