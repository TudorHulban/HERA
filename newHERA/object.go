package main

type Person struct {
	Persons struct{} `hera:"tablename"`

	ID             uint `hera:"pk"`
	Name           string
	Age            int16
	AllowedToDrive bool `hera:"default:false, columnname:driving,"`
}
