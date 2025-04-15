package handler

import (
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"net/http"
)

func MethodNotAllowed(w http.ResponseWriter, r *http.Request, expected string, wLogg *wlogger.CustomWLogg) bool {
	if r.Method != expected {
		wLogg.LogHttpE(http.StatusMethodNotAllowed, r.Method, r.URL.Path, msg.H7002, nil)
		http.Error(w, msg.H7002, http.StatusMethodNotAllowed)

		return true
	}

	return false
}
