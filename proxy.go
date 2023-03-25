package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var subsystemUrlPrefix map[string]string = map[string]string{
	// system A
	"/a/foo/bar": "/v1/foo/bar",
	// system B
	"/b/baz": "/baz",
}

const (
	systemARoutePrefix = "/a"
	systemBRoutePrefix = "/b"
)

func NewProxy(systemAURL, systemBURL string) (*httputil.ReverseProxy, error) {
	urlA, urlErr := url.Parse(systemAURL)
	if urlErr != nil {
		return nil, fmt.Errorf("cannot parse system A URL: %w", urlErr)
	}

	urlB, urlErr := url.Parse(systemBURL)
	if urlErr != nil {
		return nil, fmt.Errorf("cannot parse system B URL: %w", urlErr)
	}

	director := func(req *http.Request) {
		originalURL := req.URL

		if strings.HasPrefix(originalURL.Path, systemARoutePrefix) {
			req.URL = urlA
		} else if strings.HasPrefix(originalURL.Path, systemBRoutePrefix) {
			req.URL = urlB
		} else {
			return
		}

		req.URL.Fragment = originalURL.Fragment
		req.URL.RawQuery = originalURL.RawQuery
		req.URL.Path = mapPath(originalURL.Path)
	}

	errResponse := struct{ Message string }{
		Message: "Something went wrong while handling your request. Our team has been notified. Please try again later.",
	}
	resJson, err := json.Marshal(errResponse)
	if err != nil {
		log.Panicf("error marshalling error response to JSON")
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("error proxying request from '%s' to '%s' because %v", r.RequestURI, r.URL, err)

		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write(resJson)
		if err != nil {
			log.Printf("error writing proxy error response")
			return
		}
	}

	return &httputil.ReverseProxy{Director: director, ErrorHandler: errorHandler}, nil
}

func mapPath(path string) string {
	for apiPrefix, subsystemPrefix := range subsystemUrlPrefix {
		if strings.HasPrefix(path, apiPrefix) {
			return strings.Replace(path, apiPrefix, subsystemPrefix, 1)
		}
	}

	return path
}
