package convert

import (
	"FGW/internal/handler"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io"
	"net/http"
	"strconv"
	"strings"
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
	i, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	return i
}

// ConvStrToFloat конвертировать строку в вещественное число.
func ConvStrToFloat(str string) float64 {
	i, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	return i
}

// ConvStrToBool конвертировать строку в bool.
func ConvStrToBool(str string) bool {
	str = strings.ToLower(strings.TrimSpace(str))

	return str == "true" || str == "on" || str == "1"
}

func ParseUUIDUnsafe(str string) uuid.UUID {
	UUID, err := uuid.Parse(str)
	if err != nil {
		fmt.Println(err.Error())

		return uuid.Nil
	}

	return UUID
}

// ParseStrToUUID парсит строку в UUID и пишет ошибку в HTTP-ответ при неудаче.
func ParseStrToUUID(formValue string, w http.ResponseWriter, r *http.Request, wLogg *wlogger.CustomWLogg) (uuid.UUID, error) {
	idUUID, err := uuid.Parse(formValue)
	if err != nil {
		handler.WriteBadRequest(w, r, wLogg, msg.H7004, err)

		return uuid.Nil, err
	}

	return idUUID, nil
}

// ParseStrToID парсит строку в UUID и пишет ошибку в HTTP-ответ при неудаче.
func ParseStrToID(formValue string, w http.ResponseWriter, r *http.Request, wLogg *wlogger.CustomWLogg) (int, error) {
	id, err := strconv.Atoi(formValue)
	if err != nil {
		handler.WriteBadRequest(w, r, wLogg, msg.H7013, err)

		return 0, err
	}
	return id, nil
}
