package main

import (
	"fmt"
	"log"
	"net/http"

	"vse-course/config"
	netchi "vse-course/transport/chi"
)

func main() {
	cfg, err := config.ReadConfigFromFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	h := netchi.Initialize(cfg.Port)

	fmt.Println(fmt.Sprintf("server is running on port: %d", cfg.Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", h.Port), h.Mux); err != nil {
		log.Fatal(err)
	}
}
