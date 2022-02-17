package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

type Validator interface {
	Validate(v reflect.Value) error
}

type StrLenVal struct {
	len int
}

type MinNumVal struct {
	min int64
}

type maxNumValidator struct {
	max int64
}

type RegExpVal struct {
	re *regexp.Regexp
}

type EnumVal struct {
	enum []string
}

func (r *RegExpVal) Validate(v reflect.Value) error {
	val := v.String()

	if !r.re.Match([]byte(val)) {
		return ErrRegexp
	}

	return nil
}

func (m *maxNumValidator) Validate(v reflect.Value) error {
	if intValue := v.Int(); intValue > m.max {
		return ErrNumRange
	}

	return nil
}

func (m *MinNumVal) Validate(v reflect.Value) error {
	if intValue := v.Int(); intValue < m.min {
		return ErrNumRange
	}

	return nil
}

func (s *StrLenVal) Validate(v reflect.Value) error {
	strValue := v.String()
	if len(strValue) != s.len {
		return ErrStrLen
	}

	return nil
}

func (e *EnumVal) Validate(v reflect.Value) error {
	var val string

	if v.Kind() == reflect.Int {
		val = strconv.Itoa(int(v.Int()))
	} else {
		val = v.String()
	}

	for _, r := range e.enum {
		if val == r {
			return nil
		}
	}

	return ErrStrEnum
}

var (
	ErrStrLen   = errors.New("string length error")
	ErrNumRange = errors.New("number is out of range")
	ErrStrEnum  = errors.New("enum error")
	ErrRegexp   = errors.New("regexp error")
)

func (v ValidationErrors) Error() string {
	m := make([]string, 0, len(v))
	for _, err := range v {
		m = append(m, err.Err.Error())
	}

	return strings.Join(m, "\n")
}

func Validate(v interface{}) error {
	var validationErr ValidationErrors

	rType := reflect.TypeOf(v)
	rValue := reflect.ValueOf(v)

	for rType.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < rType.NumField(); i++ {
		rTypeField := rType.Field(i)
		rValueField := rValue.Field(i)

		tag := rTypeField.Tag.Get("validate")

		if tag == "" {
			continue
		}

		if !TypeAllow(rTypeField.Type) {
			continue
		}

		validators, err := parseTag(tag)
		if err != nil {
			return err
		}

		for _, validator := range validators {
			if rTypeField.Type.Kind() == reflect.Slice {
				for j := 0; j < rValue.Field(i).Len(); j++ {
					fieldFull := strings.Join([]string{rTypeField.Name, strconv.Itoa(j)}, ".")
					validateField(fieldFull, rValue.Field(i).Index(j), validator, &validationErr)
				}
			} else {
				validateField(rTypeField.Name, rValueField, validator, &validationErr)
			}
		}
	}

	if len(validationErr) == 0 {
		return nil
	}

	return validationErr
}

func validateField(name string, r reflect.Value, validator Validator, errorVal *ValidationErrors) {
	var validErr ValidationErrors

	valError := validator.Validate(r)
	if valError != nil {
		if errors.As(valError, &validErr) {
			*errorVal = append(*errorVal, validErr...)
		} else {
			*errorVal = append(*errorVal, ValidationError{name, valError})
		}
	}
}

func TypeAllow(r reflect.Type) bool {
	return r.Kind() == reflect.Int || r.Kind() == reflect.String ||
		r == reflect.SliceOf(reflect.TypeOf("")) ||
		r == reflect.SliceOf(reflect.TypeOf(1)) ||
		r.Kind() == reflect.Struct
}

func parseTag(tag string) ([]Validator, error) {
	var validators []Validator
	var valName, valOption string
	var validator Validator
	var err error

	splitString := strings.Split(tag, "|")

	for _, resParts := range splitString {
		splitStringSec := strings.Split(resParts, ":")

		if len(splitStringSec) == 0 {
			return nil, fmt.Errorf("failed parse parts %s", splitString)
		}

		if len(splitStringSec) > 0 {
			valName = splitStringSec[0]
		}

		if len(splitStringSec) > 1 {
			valOption = strings.TrimSpace(splitStringSec[1])
		}

		switch valName {
		case "len":
			validator, err = StrLenValidator(valOption)
		case "regexp":
			validator, err = RegExpValidator(valOption)
		case "in":
			validator, err = StrNumValidator(valOption)
		case "min":
			validator, err = MinNumValidator(valOption)
		case "max":
			validator, err = MaxNumValidator(valOption)
		default:
			validator = nil
		}

		if err != nil {
			return nil, fmt.Errorf("cant create validator %s: %w", valName, err)
		}

		if validator != nil {
			validators = append(validators, validator)
		}
	}

	return validators, nil
}

func StrNumValidator(valOption string) (*EnumVal, error) {
	res := strings.Split(valOption, ",")

	return &EnumVal{res}, nil
}

func RegExpValidator(valOption string) (*RegExpVal, error) {
	re, err := regexp.Compile(valOption)
	if err != nil {
		return nil, fmt.Errorf("failed parse %w", err)
	}

	return &RegExpVal{re}, nil
}

func MaxNumValidator(valOption string) (Validator, error) {
	var num int
	var err error

	if num, err = strconv.Atoi(valOption); err != nil {
		return nil, fmt.Errorf("%s is not number", valOption)
	}

	return &maxNumValidator{int64(num)}, nil
}

func MinNumValidator(valOption string) (*MinNumVal, error) {
	var num int
	var err error

	if num, err = strconv.Atoi(valOption); err != nil {
		return nil, fmt.Errorf("%s is not number", valOption)
	}

	return &MinNumVal{int64(num)}, nil
}

func StrLenValidator(valOption string) (*StrLenVal, error) {
	var length int
	var err error

	if length, err = strconv.Atoi(valOption); err != nil {
		return nil, fmt.Errorf("%s is not number", valOption)
	}

	return &StrLenVal{length}, nil
}
