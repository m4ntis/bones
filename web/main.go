package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello bones!")
}

func handleHTTP() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe("localhost:80", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	handleHTTP()
}
