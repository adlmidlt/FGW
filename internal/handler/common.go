package handler

import (
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"encoding/json"
	"github.com/google/uuid"
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

// ParseStrToUUID парсит строку в UUID и пишет ошибку в HTTP-ответ при неудаче.
func ParseStrToUUID(formValue string, w http.ResponseWriter, r *http.Request, wLogg *wlogger.CustomWLogg) (uuid.UUID, error) {
	idUUID, err := uuid.Parse(formValue)
	if err != nil {
		wLogg.LogHttpE(http.StatusBadRequest, r.Method, r.URL.Path, msg.H7004, err)
		http.Error(w, msg.H7004, http.StatusBadRequest)

		return uuid.Nil, err
	}

	return idUUID, nil
}

func WriteJSON(writer http.ResponseWriter, obj interface{}, wLogg *wlogger.CustomWLogg) {
	writer.Header().Set("Content-Type", "application/json_api; charset=UTF-8")
	if err := json.NewEncoder(writer).Encode(obj); err != nil {
		wLogg.LogE(msg.E3105, err)
		http.Error(writer, msg.E3105, http.StatusInternalServerError)

		return
	}
}
