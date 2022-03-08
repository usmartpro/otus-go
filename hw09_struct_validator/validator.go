package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	text := ""
	if len(v) > 0 {
		text = "ошибка валидации: \n"
		for _, err := range v {
			text += fmt.Sprintf("поле %s: %v\n", err.Field, err.Err)
		}
	}
	return text
}

func (v *ValidationErrors) Add(field string, err error) {
	vError := ValidationError{Field: field, Err: err}
	*v = append(*v, vError)
}

var ErrFormatData = errors.New("неверный формат входных данных")

var (
	ErrValidationLessThanMinimalValue = errors.New("значение меньше минимального")
	ErrValidationMoreThanMaximalValue = errors.New("значение больше максимального")
	ErrValidationNotInAgreeValues     = errors.New("значение не соответствует ожидаемому")
	ErrValidationLengthValue          = errors.New("некорректное значение длины поля")
	ErrValidationRegexpValue          = errors.New("значение некорректно")
)

func Validate(v interface{}) (ValidationErrors, error) {
	// рефлексия значения
	input := reflect.ValueOf(v)

	// проверка входного параметра на структуру
	if input.Kind() != reflect.Struct {
		return nil, ErrFormatData
	}

	// список ошибок валидации
	vErrors := ValidationErrors{}

	// рефлексия типа структуры
	structType := input.Type()

	// проходим по всем полям
	for i := 0; i < structType.NumField(); i++ {
		// рефлексия конкретного поля
		field := structType.Field(i)
		// получение тега валидации
		validateRules := field.Tag.Get("validate")

		// разделение на правила
		rules := strings.Split(validateRules, "|")
		for _, rule := range rules {
			ruleItem := strings.Split(rule, ":")

			// правило должно состоять из самого правила и ограничения по нему
			if len(ruleItem) == 2 {
				// правило
				ruleType := ruleItem[0]
				// ограничение
				ruleRestrict := ruleItem[1]
				// исследуем поле
				inspectField(&vErrors, input.Field(i), field, ruleType, ruleRestrict)
			}
		}
	}
	return vErrors, nil
}

func inspectField(vErrors *ValidationErrors, inputField reflect.Value, field reflect.StructField,
	ruleType string, ruleRestrict string) {
	// проверяем тип поля
	switch field.Type.Kind() { //nolint:exhaustive
	case reflect.Int:
		inspectInt(vErrors, inputField, field, ruleType, ruleRestrict)
	case reflect.String:
		inspectString(vErrors, inputField, field, ruleType, ruleRestrict)
	case reflect.Slice:
		// берем значение поля в нужном типе
		slice := inputField.Slice(0, inputField.Len())
		// проверяем правила
		switch inputField.Type().String() {
		case "[]string":
			for i := 0; i < inputField.Len(); i++ {
				inspectString(vErrors, slice.Index(i), field, ruleType, ruleRestrict)
			}
		case "[]int":
			for i := 0; i < inputField.Len(); i++ {
				inspectInt(vErrors, slice.Index(i), field, ruleType, ruleRestrict)
			}
		}
	}
}

func inspectInt(vErrors *ValidationErrors, inputField reflect.Value, field reflect.StructField,
	ruleType string, ruleRestrict string) {
	// берем значение поля в нужном типе
	intVal := inputField.Int()
	// проверяем правила
	switch ruleType {
	case "min":
		checkMinInt(vErrors, field.Name, intVal, ruleRestrict)
	case "max":
		checkMaxInt(vErrors, field.Name, intVal, ruleRestrict)
	case "in":
		checkInInt(vErrors, field.Name, intVal, ruleRestrict)
	}
}

func inspectString(vErrors *ValidationErrors, inputField reflect.Value, field reflect.StructField,
	ruleType string, ruleRestrict string) {
	// берем значение поля в нужном типе
	stringVal := inputField.String()
	// проверяем правила
	switch ruleType {
	case "len":
		checkLenString(vErrors, field.Name, stringVal, ruleRestrict)
	case "regexp":
		checkRegexpString(vErrors, field.Name, stringVal, ruleRestrict)
	case "in":
		checkInString(vErrors, field.Name, stringVal, ruleRestrict)
	}
}

// проверка на минимальное значение.
func checkMinInt(vErrors *ValidationErrors, fieldName string, value int64, ruleRestrict string) {
	ruleRestrictInt, _ := strconv.ParseInt(ruleRestrict, 10, 64)
	if value < ruleRestrictInt {
		vErrors.Add(fieldName, ErrValidationLessThanMinimalValue)
	}
}

// проверка на максимальное значение.
func checkMaxInt(vErrors *ValidationErrors, fieldName string, value int64, ruleRestrict string) {
	ruleRestrictInt, _ := strconv.ParseInt(ruleRestrict, 10, 64)
	if value > ruleRestrictInt {
		vErrors.Add(fieldName, ErrValidationMoreThanMaximalValue)
	}
}

// проверка на вхождение значения int.
func checkInInt(vErrors *ValidationErrors, fieldName string, value int64, ruleRestrict string) {
	restrict := strings.Split(ruleRestrict, ",")
	inArray := false
	for _, v := range restrict {
		ruleRestrictItemInt, _ := strconv.ParseInt(v, 10, 64)
		if value == ruleRestrictItemInt {
			inArray = true
			break
		}
	}
	if !inArray {
		vErrors.Add(fieldName, ErrValidationNotInAgreeValues)
	}
}

// проверка на вхождение значения string.
func checkInString(vErrors *ValidationErrors, fieldName, value, ruleRestrict string) {
	restrict := strings.Split(ruleRestrict, ",")
	inArray := false
	for _, v := range restrict {
		if value == v {
			inArray = true
			break
		}
	}
	if !inArray {
		vErrors.Add(fieldName, ErrValidationNotInAgreeValues)
	}
}

// проверка на длину строки.
func checkLenString(vErrors *ValidationErrors, fieldName, value, ruleRestrict string) {
	length := utf8.RuneCount([]byte(value))
	ruleRestrictItem, _ := strconv.ParseInt(ruleRestrict, 10, 64)
	if length != int(ruleRestrictItem) {
		vErrors.Add(fieldName, ErrValidationLengthValue)
	}
}

// проверка на регулярное выражение.
func checkRegexpString(vErrors *ValidationErrors, fieldName, value, ruleRestrict string) {
	if ok, _ := regexp.Match(ruleRestrict, []byte(value)); !ok {
		vErrors.Add(fieldName, ErrValidationRegexpValue)
	}
}
