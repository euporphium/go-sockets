package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/echo", handleConnection)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	go handleMessages()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}
