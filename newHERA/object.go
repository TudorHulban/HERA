package main

type Person struct {
	ID             uint `hera:"pk"`
	Name           string
	Age            int16
	AllowedToDrive bool
}
