package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

const (
	SLASH  = "slash"
	SYMBOL = "symbol"
	DIGIT  = "digit"
)

func Unpack(str string) (res string, err error) {
	var b strings.Builder
	var lastSymbol, symbol symbolState

	strRune := []rune(str)
	resetIndex := true
	strLen := utf8.RuneCountInString(str)

	for index, item := range strRune {
		isLast := strLen == index+1
		isFirst := index == 0
		// текущий символ
		symbol = create(item)

		if lastSymbol.State == SLASH {
			// экранированный символ ( или цифра или слэш), переопределяем как символ
			symbol.State = SYMBOL
		}

		switch {
		case isFirst && symbol.State == DIGIT:
			// первый символ - цифра
			return "", ErrInvalidString
		case resetIndex: // состояние после мульти-добавления символа (символ + цифра)
			if isLast && symbol.State == SYMBOL {
				b.WriteString(string(symbol.Symbol))
				symbol.IsAdded = true
				break
			}
			if symbol.State == DIGIT {
				// повторно цифра быть не может
				return "", ErrInvalidString
			}
			lastSymbol = symbol
			resetIndex = false
			continue
		case symbol.State == DIGIT && lastSymbol.State == SYMBOL:
			// цифра за символом
			count, _ := strconv.Atoi(string(item))
			b.WriteString(strings.Repeat(string(lastSymbol.Symbol), count)) // мультидобавление символа
			lastSymbol.IsAdded = true
			resetIndex = true
		case isLast || symbol.State == SYMBOL:
			// одиночный символ
			if lastSymbol.State != SLASH {
				// слэши не добавляем
				b.WriteString(string(lastSymbol.Symbol))
				lastSymbol.IsAdded = true
			}
			if isLast {
				// последний символ
				b.WriteString(string(symbol.Symbol))
				symbol.IsAdded = true
			}
		case symbol.State == SLASH && !lastSymbol.IsAdded:
			b.WriteString(string(lastSymbol.Symbol))
			lastSymbol.IsAdded = true
		}
		// сохраняем предыдущий символ
		lastSymbol = symbol
	}

	return b.String(), nil
}

type symbolState struct { // структура символа
	Symbol  rune   // его руна
	State   string // статус (символ, цифра, слэш)
	IsAdded bool   // был ли добавлен символ
}

func create(symbol rune) symbolState { // конструктор
	return symbolState{symbol, getState(symbol), false}
}

func getState(symbol rune) string { // функция расчета статуса руны
	switch {
	case string(symbol) == `\`:
		return SLASH
	case unicode.IsDigit(symbol):
		return DIGIT
	default:
		return SYMBOL
	}
}
