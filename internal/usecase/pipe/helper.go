package pipe

import (
	"strings"
	"unicode"
)

type Helper struct {
}

func (h Helper) GetKeyValue(line string, delimeter rune) (k, v string, err error) {
	cleanLine := strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return unicode.ToLower(r)
		}

		return -1
	}, line)

	if pos := strings.IndexRune(cleanLine, delimeter); pos != -1 {
		return strings.Trim(cleanLine[:pos], " "), strings.Trim(cleanLine[pos+1:], " "), nil
	}

	return "", "", ErrNotFound
}
