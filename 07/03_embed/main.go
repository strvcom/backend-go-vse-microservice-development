package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
)

// content holds our static web server content.
//
//go:embed images/*
//go:embed notes/*
var content embed.FS

//go:embed notes/hello.txt
var greeting string

func main() {
	fmt.Println("There is a greeting from embed:", greeting)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(content))))
	fmt.Println("Check out http://localhost:8080/static/image/gopher.png.")
	fmt.Println("Do you want to know what is embedded in content? See http://localhost:8080/static/.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
