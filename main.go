package main

import (
	"log"
	"net/http"
)

const port = "3000"

func main() {
	http.HandleFunc("/ping", PongHandler)
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
