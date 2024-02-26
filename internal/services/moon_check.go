package services

import (
	"fmt"
	"unicode"
)

func (serv *Services) MoonCheck(inpStr string) (err error) {
	// Сначала нам нужно убедиться что в строке только цифры
	digits := []rune{}
	for _, char := range inpStr {
		if unicode.IsDigit(char) {
			digits = append(digits, char)
		}
	}
	if len(inpStr) != len(digits) {
		return fmt.Errorf("not only digits in input string %w", err)
	}
	return nil
}
