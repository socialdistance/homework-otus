package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var b strings.Builder
	runes := []rune(s)

	for i, r := range runes {
		switch {
		case unicode.IsDigit(r):
			if unicode.IsDigit(r) && i == 0 {
				return "", ErrInvalidString
			}

			count, err := strconv.Atoi(string(r))
			if err != nil {
				return "", err
			}

			if count == 0 {
				if unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i-1]) {
					return "", ErrInvalidString
				}
				letter := strings.Replace(string(runes), string(runes[i]), "", 1)
				res := strings.Replace(letter, string(runes[i-1]), "", 1)
				return res, nil
			}

			for idx := 0; idx < count-1; idx++ {
				b.WriteRune(runes[i-1])
			}

		case r == '\\' && string(runes[i+2]) == "\\":
			res := strings.Replace(string(runes), string(runes[i]), "", 1)
			letter := strings.Replace(res, string(runes[i+2]), "", 1)
			return letter, nil

		case r == '\\' && unicode.IsDigit(runes[i+2]):
			count, err := strconv.Atoi(string(runes[i+2]))
			if err != nil {
				return "", err
			}
			res := strings.Repeat(string(runes[i+1]), count)
			b.WriteString(res)
			result := b.String()
			return result, nil

		case r == '"':
			return "", ErrInvalidString

		default:
			b.WriteRune(r)
		}
	}

	return b.String(), nil
}
