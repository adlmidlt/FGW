package convert

import (
	"FGW/internal/handler"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"github.com/google/uuid"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Win1251ToUTF8 конвертирует из win1251 в utf8.
func Win1251ToUTF8(str string) (string, error) {
	tr := transform.NewReader(strings.NewReader(str), charmap.Windows1251.NewDecoder())
	buf, err := io.ReadAll(tr)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

// ConvStrToInt конвертировать строку в число.
func ConvStrToInt(str string) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("Ошибка: [%s] --- ссылка на код: [ %s ] --- значение: [%v]", err.Error(), wlogger.FileWithLineNum(), value)

		return 0
	}

	return value
}

// ConvStrToFloat конвертировать строку в вещественное число.
func ConvStrToFloat(str string) float64 {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Printf("Ошибка: [%s] --- ссылка на код: [ %s ] --- значение: [%v]", err.Error(), wlogger.FileWithLineNum(), value)

		return 0
	}

	return value
}

// ParseFormFieldInt преобразует поле в целое число, полученное из HTTP запроса.
func ParseFormFieldInt(r *http.Request, fieldName string) int {
	formValue := r.FormValue(fieldName)
	if formValue == "" {
		formValue = "0"

		return 0
	}
	value, err := strconv.Atoi(formValue)
	if err != nil {
		log.Printf("Ошибка: [%s] --- ссылка на код: [ %s ] --- поле: [%s] --- значение: [%v]", err.Error(), wlogger.FileWithLineNum(), fieldName, value)

		return 0
	}

	return value
}

// ParseFormFieldFloat преобразует поле в вещественное число, полученное из HTTP запроса.
func ParseFormFieldFloat(r *http.Request, fieldName string) float64 {
	formValue := r.FormValue(fieldName)
	if formValue == "" {
		formValue = "0"

		return 0
	}
	value, err := strconv.ParseFloat(formValue, 64)
	if err != nil {
		log.Printf("Ошибка: [%s] --- ссылка на код: [ %s ] --- поле: [%s] --- значение: [%v]", err.Error(), wlogger.FileWithLineNum(), fieldName, value)

		return 0
	}

	return value
}

// ParseFormFieldBool преобразует поле в булево значение, полученное из HTTP запроса.
func ParseFormFieldBool(r *http.Request, fieldName string) bool {
	formValue := r.FormValue(fieldName)
	if formValue == "" {
		formValue = "false"
	} else {
		formValue = "true"
	}

	value, err := strconv.ParseBool(formValue)

	if err != nil {
		log.Printf("Ошибка: [%s] --- ссылка на код: [ %s ] --- поле: [%s] --- значение: [%v]", err.Error(), wlogger.FileWithLineNum(), fieldName, value)

		return false
	}

	return value
}

// ParseUUIDUnsafe пытается разобрать строку в UUID.
func ParseUUIDUnsafe(str string) uuid.UUID {
	value, err := uuid.Parse(str)
	if err != nil {
		log.Printf("Ошибка: [%s] --- ссылка на код: [ %s ] --- значение: [%v]", err.Error(), wlogger.FileWithLineNum(), value)

		return uuid.Nil
	}

	return value
}

// ParseStrToUUID пытается разобрать строку в UUID и пишет ошибку в HTTP-ответ при неудаче.
func ParseStrToUUID(fieldName string, w http.ResponseWriter, r *http.Request, wLogg *wlogger.CustomWLogg) (uuid.UUID, error) {
	value, err := uuid.Parse(fieldName)
	if err != nil {
		log.Printf("Ошибка: [%s] --- ссылка на код: [ %s ] --- значение: [%v]", err.Error(), wlogger.FileWithLineNum(), value)
		handler.WriteBadRequest(w, r, wLogg, msg.H7004, err)

		return uuid.Nil, err
	}

	return value, nil
}

// ParseStrToID пытается разобрать строку в UUID и пишет ошибку в HTTP-ответ при неудаче.
func ParseStrToID(fieldName string, w http.ResponseWriter, r *http.Request, wLogg *wlogger.CustomWLogg) (int, error) {
	value, err := strconv.Atoi(fieldName)
	if err != nil {
		log.Printf("Ошибка: [%s] --- ссылка на код: [ %s ] --- значение: [%v]", err.Error(), wlogger.FileWithLineNum(), value)
		handler.WriteBadRequest(w, r, wLogg, msg.H7013, err)

		return 0, err
	}

	return value, nil
}

// FormatDateTime - функция форматирования даты в формате ДД.ММ.ГГГГ ЧЧ:ММ
func FormatDateTime(dateTime string) string {
	t, err := time.Parse(time.RFC3339, dateTime)
	if err != nil {
		return dateTime
	}
	return t.Format("02.01.2006 15:04:05")
}
