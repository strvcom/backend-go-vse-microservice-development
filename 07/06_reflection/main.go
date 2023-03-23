package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

func main() {
	person := Person{
		FirstName: "John",
		LastName:  "Doe",
		Age:       26,
	}

	typ := reflect.TypeOf(person)
	fmt.Println(typ)

	kind := typ.Kind()
	fmt.Println(kind)

	value := reflect.ValueOf(person)
	if value.Kind() == reflect.Struct {
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			if field.CanInt() {
				fmt.Println(field.Int())
			} else {
				fmt.Println(field.String())
			}
		}
	}

	// json marshal
	fmt.Println(marshalJSON(person))
	fmt.Println(marshalJSONWithTags(person))
}

func marshalJSON(input any) string {
	output := &strings.Builder{}
	// opening bracket
	output.WriteRune('{')

	typ := reflect.TypeOf(input)
	value := reflect.ValueOf(input)
	if value.Kind() == reflect.Struct {
		for i := 0; i < value.NumField(); i++ {
			fieldValue := value.Field(i)
			fieldType := typ.Field(i)

			// write field name
			output.WriteRune('"')
			output.WriteString(fieldType.Name)
			output.WriteRune('"')
			output.WriteRune(':')

			// write field value
			if fieldValue.CanInt() {
				output.WriteString(strconv.FormatInt(fieldValue.Int(), 10))
			} else {
				output.WriteRune('"')
				output.WriteString(fieldValue.String())
				output.WriteRune('"')
			}

			// write comma if the field is not last
			if i != value.NumField()-1 {
				output.WriteRune(',')
			}
		}
	}
	// closing bracket
	output.WriteRune('}')

	return output.String()
}

func marshalJSONWithTags(input any) string {
	output := &strings.Builder{}
	// opening bracket
	output.WriteRune('{')

	typ := reflect.TypeOf(input)
	value := reflect.ValueOf(input)
	if value.Kind() == reflect.Struct {
		for i := 0; i < value.NumField(); i++ {
			fieldValue := value.Field(i)
			fieldType := typ.Field(i)

			// write field name
			output.WriteRune('"')
			if tag, found := fieldType.Tag.Lookup("json"); found {
				output.WriteString(tag)
			} else {
				output.WriteString(fieldType.Name)
			}
			output.WriteRune('"')
			output.WriteRune(':')

			// write field value
			if fieldValue.CanInt() {
				output.WriteString(strconv.FormatInt(fieldValue.Int(), 10))
			} else {
				output.WriteRune('"')
				output.WriteString(fieldValue.String())
				output.WriteRune('"')
			}

			// write comma if the field is not last
			if i != value.NumField()-1 {
				output.WriteRune(',')
			}
		}
	}
	// closing bracket
	output.WriteRune('}')

	return output.String()
}
