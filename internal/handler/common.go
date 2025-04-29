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

// EntityExistByUUIDChecker — универсальный интерфейс для проверки существования сущности по её UUID.
type EntityExistByUUIDChecker interface {
	ExistsByUUID(ctx context.Context, idEntity uuid.UUID) (bool, error)
}

// EntityExistByIDChecker — универсальный интерфейс для проверки существования сущности по её ID.
type EntityExistByIDChecker interface {
	ExistsByID(ctx context.Context, idEntity int) (bool, error)
}

// EntityExistsByUUID проверяет на существование сущности по её UUID.
// Использует интерфейс EntityExistByUUIDChecker, реализующий метод ExistsByUUID.
func EntityExistsByUUID(ctx context.Context, idEntity uuid.UUID, w http.ResponseWriter, r *http.Request, wLogg *wlogger.CustomWLogg, entityChecker EntityExistByUUIDChecker) bool {
	exist, err := entityChecker.ExistsByUUID(ctx, idEntity)

	return checkEntityExistence(err, w, r, wLogg, exist)
}

// EntityExistsByID проверяет на существование сущности по её ID.
// Использует интерфейс EntityExistByIDChecker, реализующий метод ExistsByID.
func EntityExistsByID(ctx context.Context, idEntity int, w http.ResponseWriter, r *http.Request, wLogg *wlogger.CustomWLogg, entityChecker EntityExistByIDChecker) bool {
	exist, err := entityChecker.ExistsByID(ctx, idEntity)

	return checkEntityExistence(err, w, r, wLogg, exist)
}

// checkEntityExistence проверяет, существует ли сущность.
func checkEntityExistence(err error, w http.ResponseWriter, r *http.Request, wLogg *wlogger.CustomWLogg, exist bool) bool {
	if err != nil {
		WriteServerError(w, r, wLogg, msg.H7003, err)
		WriteJSON(w, map[string]string{"message": msg.H7003}, wLogg)

		return false
	}

	if !exist {
		WriteNotFound(w, r, wLogg, msg.H7005, err)
		WriteJSON(w, map[string]string{"message": msg.W1002}, wLogg)

		return false
	}

	return true
}

// MethodNotAllowed проверяет, соответствует ли HTTP-метод запроса ожидаемому.
func MethodNotAllowed(w http.ResponseWriter, r *http.Request, expected string, wLogg *wlogger.CustomWLogg) bool {
	if r.Method != expected {
		WriteMethodNotAllowed(w, r, wLogg, msg.H7002, nil)

		return true
	}

	return false
}

// WriteJSON сериализует переданный объект в JSON и отправляет его в HTTP-ответ.
func WriteJSON(w http.ResponseWriter, entity interface{}, wLogg *wlogger.CustomWLogg) {
	w.Header().Set("Content-Type", "application/json_api; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(entity); err != nil {
		WriteServerError(w, nil, wLogg, msg.E3105, err)

		return
	}
}

// ParseTemplateHTML загружает и парсит HTML-шаблон по указанному пути.
func ParseTemplateHTML(pathToHTML string, w http.ResponseWriter, r *http.Request, wLogg *wlogger.CustomWLogg) (*template.Template, bool) {
	tmpl, err := template.ParseFiles(pathToHTML)
	if err != nil {
		WriteServerError(w, r, wLogg, msg.H7006, err)

		return nil, false
	}

	return tmpl, true
}

// ExecuteTemplate выполняет шаблон и записывает результат в http.ResponseWriter.
func ExecuteTemplate(tmpl *template.Template, data interface{}, w http.ResponseWriter, r *http.Request, wLogg *wlogger.CustomWLogg) bool {
	if err := tmpl.Execute(w, data); err != nil {
		WriteServerError(w, r, wLogg, msg.H7007, err)

		return false
	}

	return true
}
