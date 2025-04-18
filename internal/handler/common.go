package handler

import (
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"html/template"
	"net/http"
)

// MethodNotAllowed проверяет, соответствует ли HTTP-метод запроса ожидаемому.
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

// WriteJSON сериализует переданный объект в JSON и отправляет его в HTTP-ответ.
func WriteJSON(w http.ResponseWriter, entity interface{}, wLogg *wlogger.CustomWLogg) {
	w.Header().Set("Content-Type", "application/json_api; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(entity); err != nil {
		wLogg.LogE(msg.E3105, err)
		http.Error(w, msg.E3105, http.StatusInternalServerError)

		return
	}
}

// EntityExistChecker — универсальный интерфейс для проверки существования сущности по её UUID.
type EntityExistChecker interface {
	Exists(ctx context.Context, idEntity uuid.UUID) (bool, error)
}

// EntityExists проверяет на существование сущности по её UUID.
// Использует интерфейс EntityExistChecker, реализующий метод Exists.
func EntityExists(ctx context.Context, idEntity uuid.UUID, w http.ResponseWriter, r *http.Request, wLogg *wlogger.CustomWLogg, entityChecker EntityExistChecker) bool {
	_, err := entityChecker.Exists(ctx, idEntity)
	if err != nil {
		wLogg.LogHttpW(http.StatusNotFound, r.Method, r.URL.Path, msg.H7005, err)
		http.Error(w, msg.H7005, http.StatusNotFound)
		WriteJSON(w, map[string]string{"message": msg.W1002}, wLogg)

		return false
	}

	return true
}

// ParseTemplateHTML загружает и парсит HTML-шаблон по указанному пути.
func ParseTemplateHTML(pathToHTML string, w http.ResponseWriter, r *http.Request, wLogg *wlogger.CustomWLogg) (*template.Template, bool) {
	tmpl, err := template.ParseFiles(pathToHTML)
	if err != nil {
		wLogg.LogHttpE(http.StatusInternalServerError, r.Method, r.URL.Path, msg.H7006, err)
		http.Error(w, msg.H7006, http.StatusInternalServerError)

		return nil, false
	}

	return tmpl, true
}

// ExecuteTemplate выполняет шаблон и записывает результат в http.ResponseWriter.
func ExecuteTemplate(tmpl *template.Template, data interface{}, w http.ResponseWriter, r *http.Request, wLogg *wlogger.CustomWLogg) bool {
	if err := tmpl.Execute(w, data); err != nil {
		wLogg.LogHttpE(http.StatusInternalServerError, r.Method, r.URL.Path, msg.H7007, err)
		http.Error(w, msg.H7007, http.StatusInternalServerError)

		return false
	}

	return true
}
