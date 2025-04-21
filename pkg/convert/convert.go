package convert

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io"
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
	res, err := strconv.ParseBool(str)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return res
}

func ParseStrToUUID(str string) uuid.UUID {
	UUID, err := uuid.Parse(str)
	if err != nil {
		fmt.Println(err.Error())

		return uuid.Nil
	}

	return UUID
}
