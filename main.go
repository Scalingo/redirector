package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		url := *req.URL
		url.Scheme = "http"
		if req.Host == "appsdeck.eu" || req.Host == "staging.appsdeck.eu" {
			redirectAppsdeck(res, req)
			return
		}
		url.Host = strings.Replace(req.Host, "appsdeck.eu", "scalingo.com", 1)
		http.Redirect(res, req, url.String(), 301)
	})

	log.Println("listenening on", os.Getenv("PORT"))
	log.Fatalln(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

var appsdeckRedirections = map[string]string{
	"#containers-tab":    "",
	"#logs-tab":          "/logs",
	"#dns-tab":           "/domains",
	"#env-tab":           "/variables",
	"#addons-tab":        "/addons",
	"#collaborators-tab": "/collaborators",
	"#deployments-tab":   "/timeline",
	"#management-tab":    "/settings",
	"/pricing":           "/pricing",
	"/roadmap":           "http://changelog.scalingo.com",
	"/contact":           "/contact",
	"":                   "",
}

func redirectAppsdeck(res http.ResponseWriter, req *http.Request) {
	path := strings.TrimSpace(req.URL.Path)
	u := &url.URL{Scheme: "https"}

	if strings.Contains(path, "/home") {
		switch req.Host {
		case "appsdeck.eu":
			u.Host = "my.scalingo.com"
		case "staging.appsdeck.eu":
			u.Host = "scalingo-dashboard.staging.scalingo.io"
		}
	} else {
		switch req.Host {
		case "appsdeck.eu":
			u.Host = "scalingo.com"
		case "staging.appsdeck.eu":
			u.Host = "staging.scalingo.com"
		}
	}

	if strings.HasPrefix(path, "/home") {
		path = strings.TrimPrefix(path, "/home")
		fragment := "#" + req.URL.Fragment
		for src, dst := range appsdeckRedirections {
			if src == fragment {
				u.Path = path + dst
				break
			}
		}
		if u.Path == "" {
			u.Path = path
		}
	} else {
		for src, dst := range appsdeckRedirections {
			if path == src {
				dstURL, err := url.Parse(dst)
				if err != nil {
					fmt.Println(err)
				} else {
					if len(dstURL.Host) != 0 {
						u.Scheme = dstURL.Scheme
						u.Host = dstURL.Host
						u.Path = dstURL.Path
					} else {
						u.Path = dst
					}
				}
				break
			}
		}
	}

	http.Redirect(res, req, u.String(), 301)
}
