package main

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
)

var (
	log = logrus.WithField("cmd", "web")
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.WithField("PORT", port).Fatal("$PORT must be set")
	}

	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println(http.ListenAndServe(":"+port, nil))
}
