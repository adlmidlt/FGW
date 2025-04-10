package convert

import (
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io"
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
