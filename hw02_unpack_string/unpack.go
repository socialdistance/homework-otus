package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var result, res, letter string

	var b strings.Builder

	for i, r := range s {
		switch {
		case unicode.IsDigit(r) && i == 0:
			return "", ErrInvalidString

		case unicode.IsDigit(rune(s[i])) && unicode.IsDigit(rune(s[i+1])):
			return "", ErrInvalidString

		case unicode.IsDigit(r):
			count, err := strconv.Atoi(string(r))
			if err != nil {
				return "", err
			}
			if count == 0 {
				res = strings.Replace(s, string(s[i-1]), "", 1)
				letter = strings.Replace(res, string(s[i]), "", 1)
				return letter, nil
			}

			if count != 0 {
				res = strings.Repeat(string(s[i-1]), count-1)
				b.WriteString(res)
			}

		case r == '\\' && string(s[i+2]) == "\\":
			res = strings.Replace(s, string(s[i]), "", 1)
			letter = strings.Replace(res, string(s[i+2]), "", 1)
			return letter, nil

		case r == '\\' && unicode.IsDigit(rune(s[i+2])):
			count, err := strconv.Atoi(string(s[i+2]))
			if err != nil {
				return "", err
			}
			res = strings.Repeat(string(s[i+1]), count)
			b.WriteString(res)
			result = b.String()
			return result, nil

		case r == '"':
			return "", ErrInvalidString

		case string(r) == "\n":
			return "", ErrInvalidString

		default:
			b.WriteRune(r)
		}
	}

	result = b.String()

	return result, nil
}
