package handler

import (
	"FGW/pkg/wlogger"
	"net/http"
)

func WriteNotFound(w http.ResponseWriter, r *http.Request, wLogg *wlogger.CustomWLogg, message string, err error) {
	wLogg.LogHttpW(http.StatusNotFound, r.Method, r.URL.Path, message, err)
	http.Error(w, message, http.StatusNotFound)
}

func WriteServerError(w http.ResponseWriter, r *http.Request, wLogg *wlogger.CustomWLogg, message string, err error) {
	wLogg.LogHttpW(http.StatusInternalServerError, r.Method, r.URL.Path, message, err)
	http.Error(w, message, http.StatusInternalServerError)
}

func WriteBadRequest(w http.ResponseWriter, r *http.Request, wLogg *wlogger.CustomWLogg, message string, err error) {
	wLogg.LogHttpE(http.StatusBadRequest, r.Method, r.URL.Path, message, err)
	http.Error(w, message, http.StatusBadRequest)
}

func WriteMethodNotAllowed(w http.ResponseWriter, r *http.Request, wLogg *wlogger.CustomWLogg, message string, err error) {
	wLogg.LogHttpE(http.StatusMethodNotAllowed, r.Method, r.URL.Path, message, err)
	http.Error(w, message, http.StatusMethodNotAllowed)
}

func WriteUnauthorized(w http.ResponseWriter, r *http.Request, wLogg *wlogger.CustomWLogg, message string, err error) {
	wLogg.LogHttpE(http.StatusUnauthorized, r.Method, r.URL.Path, message, err)
	http.Error(w, message, http.StatusUnauthorized)
}
