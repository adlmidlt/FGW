package http_web

import (
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"github.com/google/uuid"
	"net/http"
)

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
