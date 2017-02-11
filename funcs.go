package main

type Function struct {
	Name       string
	Type       string // SELECT, UPDATE, INSERT
	Query      string // The SELECT * FROM procedure
	Inputs     []InOutType
	Outputs    []InOutType
	IsOutArray bool
}


