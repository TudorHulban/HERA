package main

type Person struct {
	ID   uint `hera:"pk"`
	Name string
	Age  uint
}
