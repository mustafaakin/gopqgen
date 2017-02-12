package main

import (
	"strings"

	"fmt"

	"log"

	"github.com/jmoiron/sqlx"
)

type InOutType struct {
	Name string
	Type string
}

type summary struct {
	Enums     []Enum
	Tables    []Table
	Views     []Table
	Functions []Function
}

func NewSummaryFromDB(db *sqlx.DB) (*summary, error) {
	s := &summary{}
	// Populate the enums
	es, err := GetEnums(db)
	if err != nil {
		return nil, err
	}
	s.Enums = es

	// Populate the tables
	ts, err := GetTables(db)
	if err != nil {
		return nil, err
	}

	s.Tables = make([]Table, 0)
	s.Views = make([]Table, 0)

	// Get the indexes for each table
	for _, t := range ts {
		switch t.Type {
		case "table":
			inds, err := GetIndices(db, t.OID)
			if err != nil {
				return nil, err
			}
			t.Indices = inds
			s.Tables = append(s.Tables, t)
		case "view":
			s.Views = append(s.Views, t)
		}
	}

	// Create functions
	// -- The indices from functions
	s.Functions = make([]Function, 0)
	for _, table := range s.Tables {
		for _, index := range table.Indices {

			var name string
			if index.IsUnique {
				name = "Get"
			} else {
				name = "List"
			}

			name += capitalize(table.Name)

			if len(index.Columns) > 0 {
				name += "By"
				if len(index.Columns) == 1 {
					name += capitalize(index.Columns[0].Name)
				} else {
					name += capitalize(index.Name)
				}
			}

			fn := Function{
				Name:       name,
				Inputs:     make([]InOutType, len(index.Columns)),
				Outputs:    []InOutType{{Name: "output", Type: table.Name}},
				Type:       "SELECT",
				IsOutArray: !index.IsUnique,
			}

			// Starting point
			fn.Query = "SELECT "

			// The projected fields, in index case, all fields
			var fs = []string{}
			for _, field := range table.Fields {
				fs = append(fs, strings.ToLower(field.Name))
			}

			fn.Query += strings.Join(fs, ", ") + " FROM " + strings.ToLower(table.Name) + " WHERE "
			whereCasues := []string{}
			for i, column := range index.Columns {
				// Add to function struct
				fn.Inputs[i] = InOutType{Name: column.Name, Type: column.Typename}
				whereCasues = append(whereCasues, fmt.Sprintf("`%s` = $%d", column.Name, i+1))
			}

			fn.Query += strings.Join(whereCasues, " AND ")
			s.Functions = append(s.Functions, fn)
		}
	}

	// -- The functions from table lists only
	for _, table := range s.Tables {
		fn := Function{
			Name:       "List" + capitalize(table.Name),
			Outputs:    []InOutType{{Name: "output", Type: table.Name}},
			Type:       "SELECT",
			IsOutArray: true,
		}
		fn.Query = "SELECT "

		// The projected fields, in index case, all fields
		var fs = []string{}
		for _, field := range table.Fields {
			fs = append(fs, strings.ToLower(field.Name))
		}

		fn.Query += strings.Join(fs, ", ") + " FROM " + strings.ToLower(table.Name)
		s.Functions = append(s.Functions, fn)
	}

	// -- The user-defined functions
	fs, err := GetUserFunctions(db)

	if err != nil {
		log.Fatal("haydaaa", err)
	}
	s.Functions = append(s.Functions, fs...)

	return s, nil
}

func capitalize(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func (s summary) generateProtoSummary() protoSummary {
	ps := newProtoSummary()

	// Enums are independent, they are first
	for _, enum := range s.Enums {
		ps.addEnum(protoEnum(enum))
	}

	// TODO: Other composite types

	// Table types
	for _, table := range s.Tables {
		msg := newProtoMessage(table.Name)
		for _, field := range table.Fields {
			msg.addField(field.Name, field.Type, field.Num)
		}
		ps.addMessage(msg)
	}

	// Views
	for _, view := range s.Views {
		msg := newProtoMessage(view.Name)
		for _, field := range view.Fields {
			msg.addField(field.Name, field.Type, field.Num)
		}
		ps.addMessage(msg)
	}

	for _, fn := range s.Functions {
		var inputName, outputName string
		if fn.Inputs == nil {
			inputName = "VoidRequest"
		} else if len(fn.Inputs) == 1 {
			inputName = fn.Inputs[0].Type
		} else {
			inputName = fn.Name + "Arg"
			msg := newProtoMessage(inputName)
			for idx, inp := range fn.Inputs {
				if inp.Name == "" {
					inp.Name = fmt.Sprintf("var%d", idx+1)
				}
				msg.addField(inp.Name, inp.Type, idx+1)
			}
			ps.addMessage(msg)
		}

		if len(fn.Outputs) == 1 {
			if fn.IsOutArray {
				outputName = "stream " + fn.Outputs[0].Type
			} else {
				outputName = fn.Outputs[0].Type
			}
		} else {
			outputName = fn.Name + "Out"
			msg := newProtoMessage(outputName)
			for idx, outp := range fn.Outputs {
				if outp.Name == "" {
					outp.Name = fmt.Sprintf("var%d", idx+1)
				}
				msg.addField(outp.Name, outp.Type, idx+1)
			}
			ps.addMessage(msg)
		}

		rpc := protoRpc{
			Function:   fn,
			Name:       fn.Name,
			InputName:  inputName,
			OutputName: outputName,
		}

		ps.addFunc(rpc)
	}

	return *ps
}
