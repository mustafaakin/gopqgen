package main

import (
	"encoding/json"
	"fmt"
)

type messageField struct {
	Name     string
	Typename string
	Index    int
}

var protoFieldLookup = map[string]string{
	"int4":                        "int32",
	"int8":                        "int64",
	"uuid":                        "string",
	"text":                        "string",
	"bytea":                       "bytes",
	"boolean":                     "bool",
	"integer":                     "int32",
	"bigint":                      "int64",
	"serial":                      "int32",
	"bigserial":                   "int64",
	"smallint":                    "int32",
	"real":                        "float",
	"double precision":            "double",
	"timestamp without time zone": "google.protobuf.Timestamp",
}

func getProtoFriendlyFieldName(s string) string {
	// TODO: This will require a better mechanism, how about varchar(a).. etc? how are they stored? we might need regexes
	if val, ok := protoFieldLookup[s]; ok {
		return val
	}
	return s
}

type protoMessage struct {
	Name   string
	Fields []messageField
}

const indentation = "    "

func (m protoMessage) String() string {
	s := fmt.Sprintf("message %s {\n", m.Name)
	for _, f := range m.Fields {
		s += fmt.Sprintf("%s%s %s = %d;\n",
			indentation, getProtoFriendlyFieldName(f.Typename), f.Name, f.Index,
		)
	}
	s += "}\n"
	return s
}

func (p *protoMessage) addField(name, typename string, index int) {
	p.Fields = append(p.Fields, messageField{
		Name:     capitalize(name),
		Typename: capitalize(typename),
		Index:    index,
	})
}

func newProtoMessage(name string) protoMessage {
	return protoMessage{
		Name:   capitalize(name),
		Fields: make([]messageField, 0),
	}
}

type protoEnum Enum

func (p protoEnum) String() string {
	var s = fmt.Sprintf("enum %s {\n", p.Name)
	s += indentation + "UNKNOWN = 0; \n"
	for _, val := range p.Values {
		s += fmt.Sprintf("%s%s = %d; \n", indentation, val.Label, val.SortOrder)
	}
	s += "}\n"
	return s
}

type protoRpc struct {
	Name       string
	InputName  string
	OutputName string
	Function   Function
}

func (pr protoRpc) String() string {
	s := fmt.Sprintf("%s// %s\n", indentation, pr.Function.Query)
	s += fmt.Sprintf("%srpc %s(%s) returns (%s) {}\n",
		indentation, pr.Name, getProtoFriendlyFieldName(pr.InputName), getProtoFriendlyFieldName(pr.OutputName),
	)
	return s
}

type protoSummary struct {
	Fns   map[string]protoRpc
	Enums map[string]protoEnum
	Msgs  map[string]protoMessage
}

func newProtoSummary() *protoSummary {
	return &protoSummary{
		Fns:   make(map[string]protoRpc, 0),
		Enums: make(map[string]protoEnum, 0),
		Msgs:  make(map[string]protoMessage, 0),
	}
}

func (ps *protoSummary) addFunc(fn protoRpc) bool {
	if _, ok := ps.Fns[fn.Name]; ok {
		return false
	}
	ps.Fns[fn.Name] = fn
	return true
}

func (ps *protoSummary) addEnum(enum protoEnum) bool {
	if _, ok := ps.Enums[enum.Name]; ok {
		return false
	}
	ps.Enums[enum.Name] = enum
	return true
}

func (ps *protoSummary) containsMessageOrEnum(key string) bool {
	if _, ok := ps.Msgs[key]; ok {
		return true
	}
	if _, ok := ps.Enums[key]; ok {
		return true
	}
	return false
}

func (ps *protoSummary) addMessage(msg protoMessage) bool {
	if _, ok := ps.Msgs[msg.Name]; ok {
		return false
	}
	ps.Msgs[msg.Name] = msg
	return true
}

func (ps protoSummary) String() string {
	s := `syntax = "proto3";` + "\n\n"

	s += "// Enums\n"
	for _, v := range ps.Enums {
		s += v.String() + "\n"
	}

	s += "// Messages, Field Types\n"

	for _, v := range ps.Msgs {
		s += v.String() + "\n"
	}

	s += "// Service Definition\n"

	s += "service DatabaseService { \n"
	for _, v := range ps.Fns {
		s += v.String()
	}
	s += "}\n"

	return s
}

func (ps protoSummary) ToJSON() string {
	b, _ := json.MarshalIndent(ps, "", "  ")
	return string(b)
}
