package handlers

import (
	"github.com/k-orolevsk-y/go-metricts-tpl/cmd/server/storage"
	"net/http"
	"strconv"
)

func UpdateGauge(storage stor.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !ValidateHTTPMethod(r, http.MethodPost) {
			handleBadRequest(w)
			return
		} else if !ValidateContentType(r, "text/plain") {
			handleBadRequest(w)
			return
		}

		params := ParseUrlParams("/update/gauge/", r.URL)
		if len(params) != 2 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		value, err := strconv.ParseFloat(params[1], 64)
		if err != nil {
			handleBadRequest(w)
			return
		}

		storage.SetGauge(params[1], value)
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateCounter(storage stor.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !ValidateHTTPMethod(r, http.MethodPost) {
			handleBadRequest(w)
			return
		} else if !ValidateContentType(r, "text/plain") {
			handleBadRequest(w)
			return
		}

		params := ParseUrlParams("/update/counter/", r.URL)
		if len(params) != 2 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		value, err := strconv.ParseInt(params[1], 10, 64)
		if err != nil {
			handleBadRequest(w)
			return
		}

		storage.AddCounter(params[1], value)
		w.WriteHeader(http.StatusOK)
	}
}
