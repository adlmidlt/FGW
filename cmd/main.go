package main

import (
	"FGW/pkg/wlogger"
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, new project \"FGW\"!!!")
	})

	logger, _ := wlogger.NewCustomWLogger()
	defer logger.Close()

	http.ListenAndServe(":7777", nil)
}
