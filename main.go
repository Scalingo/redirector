package main

import (
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
	"/pricing":           "https://scalingo.com/pricing",
	"/roadmap":           "http://changelog.scalingo.com",
	"/contact":           "https://scalingo.com/contact",
}

func redirectAppsdeck(res http.ResponseWriter, req *http.Request) {
	path := strings.TrimSpace(strings.TrimPrefix(req.URL.Path, "/home"))
	pathFragment := path
	if len(req.URL.Fragment) > 0 {
		pathFragment += "#" + req.URL.Fragment
	}

	u := &url.URL{}
	u.Scheme = "https"
	for f, p := range appsdeckRedirections {
		if req.URL.Path == f && strings.HasPrefix(p, "http") {
			pParsed, _ := url.Parse(p)
			u.Scheme = pParsed.Scheme
			u.Host = pParsed.Host
			u.Path = pParsed.Path
			break
		} else if strings.HasPrefix(pathFragment, "/apps") {
			if req.Host == "staging.appsdeck.eu" {
				u.Host = "scalingo-dashboard.staging.scalingo.io"
			} else {
				u.Host = "my.scalingo.com"
			}
			pathFragment = strings.Replace(pathFragment, f, p, 1)
			u.Path = pathFragment
		}
	}

	http.Redirect(res, req, u.String(), 301)
}
