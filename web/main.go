package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("index.html")
	defer f.Close()
	if err != nil {
		fmt.Fprint(w, "index.html not found")
		return
	}

	w.Header().Set("Content-Type", "text/html")
	io.Copy(w, f)
}

func handleHTTP() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	handleHTTP()
}
