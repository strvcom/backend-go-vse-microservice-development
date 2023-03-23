package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func main() {
	john := Person{
		FirstName: "John",
		LastName:  "Doe",
		Age:       26,
	}

	// JSON marshalling
	jsonOutput, err := json.Marshal(john)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonOutput))

	// JSON unmarshalling
	var copyOfJohnJSON Person
	err = json.Unmarshal(jsonOutput, &copyOfJohnJSON)
	if err != nil {
		panic(err)
	}
	fmt.Println(copyOfJohnJSON)

	// xml encoding
	outputBuffer := &bytes.Buffer{}
	err = xml.NewEncoder(outputBuffer).Encode(john)
	if err != nil {
		panic(err)
	}
	fmt.Println(outputBuffer.String())

	// xml decoding
	var copyOfJohnXML Person
	err = xml.NewDecoder(outputBuffer).Decode(&copyOfJohnXML)
	if err != nil {
		panic(err)
	}
	fmt.Println(copyOfJohnXML)
}
