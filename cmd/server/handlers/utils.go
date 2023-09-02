package handlers

import (
	"net/http"
	"net/url"
	"strings"
)

func ValidateHTTPMethod(r *http.Request, method string) bool {
	return r.Method == method
}

func ValidateContentType(r *http.Request, contentType string) bool {
	requestContentType := r.Header.Get("Content-Type")
	if requestContentType == "" {
		return true
	}

	return requestContentType == contentType
}

type ParsedUrlParams []string

func ParseUrlParams(path string, url *url.URL) ParsedUrlParams {
	urlPath := strings.Replace(url.Path, path, "", 1)
	split := strings.Split(urlPath, "/")

	params := make(ParsedUrlParams, 0)
	for _, param := range split {
		if param == "" {
			continue
		}

		params = append(params, param)
	}

	return params
}
