package main

import (
	"fmt"
	"net/http"
	"os"
)

func handleHTTP() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	handleHTTP()
}
