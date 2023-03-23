package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	FirstName string
	LastName  string
	Age       int
	Password  string
}

type PersonWithTags struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Age       int    `json:"age"`
	Password  string `json:"-"`
}

func main() {
	person := Person{
		FirstName: "John",
		LastName:  "",
		Age:       26,
		Password:  "HelloKitty",
	}
	personWithTags := PersonWithTags{
		FirstName: "John",
		LastName:  "",
		Age:       26,
		Password:  "HelloKitty",
	}

	output, err := json.MarshalIndent(person, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(output))

	outputWithTags, err := json.MarshalIndent(personWithTags, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(outputWithTags))
}
