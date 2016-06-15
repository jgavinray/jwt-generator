package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/token", ValidateRequest(generateToken))
	http.HandleFunc("/validateToken", TokenAuth(resource))
	fmt.Println("Token Generator listening on Port 8000\n")
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}

func resource(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!\n")
}
