package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Print("Hello, new project \"FGW\"!!!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, new project \"FGW123\"!!!")
	})

	http.ListenAndServe(":7000", nil)
}
