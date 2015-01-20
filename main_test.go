package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var redirections = map[string]string{
	"https://appsdeck.eu":                                 "https://scalingo.com",
	"https://appsdeck.eu/home/apps":                       "https://my.scalingo.com/apps",
	"https://appsdeck.eu/home/apps/app":                   "https://my.scalingo.com/apps/app",
	"https://appsdeck.eu/home/apps/app#containers-tab":    "https://my.scalingo.com/apps/app",
	"https://appsdeck.eu/home/apps/app#env-tab":           "https://my.scalingo.com/apps/app/variables",
	"https://appsdeck.eu/home/apps/app#addons-tab":        "https://my.scalingo.com/apps/app/addons",
	"https://appsdeck.eu/home/apps/app#collaborators-tab": "https://my.scalingo.com/apps/app/collaborators",
	"https://appsdeck.eu/home/apps/app#dns-tab":           "https://my.scalingo.com/apps/app/domains",
	"https://appsdeck.eu/home/apps/app#management-tab":    "https://my.scalingo.com/apps/app/settings",
	"https://appsdeck.eu/home/apps/app#deployments-tab":   "https://my.scalingo.com/apps/app/timeline",
	"https://appsdeck.eu/contact":                         "https://scalingo.com/contact",
	"https://appsdeck.eu/pricing":                         "https://scalingo.com/pricing",
	"https://appsdeck.eu/roadmap":                         "http://changelog.scalingo.com",
	"https://staging.appsdeck.eu":                         "https://staging.scalingo.com",
	"https://staging.appsdeck.eu/home/apps":               "https://scalingo-dashboard.staging.scalingo.io/apps",
}

func TestRedirectAppsdeck(t *testing.T) {
	for src, dst := range redirections {
		testRedirection(t, src, dst)
	}
}

func testRedirection(t *testing.T, path, location string) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	redirectAppsdeck(res, req)

	if res.Code != 301 {
		t.Fatal(res.Code, "should be 301")
	}

	if res.Header().Get("Location") != location {
		t.Error("for", path, "expected", location, "got", res.HeaderMap.Get("Location"))
	}
}
