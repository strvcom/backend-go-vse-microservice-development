package main

import (
	"fmt"
	"net/http"

	"auth/custom/api/rest"
	"auth/custom/service"

	"go.uber.org/zap"
)

const (
	port = 8080
)

func main() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	controller := rest.NewController(service.NewService(), &service.TokenParser{}, l)

	if err = http.ListenAndServe(fmt.Sprintf(":%d", port), controller); err != nil {
		l.Panic(err.Error())
	}
}
