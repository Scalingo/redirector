package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		url := *req.URL
		url.Host = strings.Replace(url.Host, "appsdeck.eu", "scalingo.com", 1)
		http.Redirect(res, req, url.String(), 301)
	})

	log.Println("listenening on", os.Getenv("PORT"))
	log.Fatalln(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
