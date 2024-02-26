package services

import (
	"fmt"
	"strconv"
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
	// Преобразуем в целочисленный массив
	var numbers []int64
	for i := 0; i < len(inpStr); i++ {
		num, err := strconv.ParseInt(string(inpStr[i]), 10, 32)
		if err != nil {
			return fmt.Errorf("parseint fail.. %w", err)
		}
		numbers = append(numbers, num)
	}

	// Далее проверяем на четность.
	if len(numbers)%2 == 0 {
		// Четное количество цифр.
		serv.l.ZL.Debug("it's even")
		err = forEven(numbers)
		return err
	} else {
		// Нечетное.
		serv.l.ZL.Debug("it's odd")
		err = forOdd(numbers)
		return err
	}
}

func forEven(numbers []int64) (err error) {
	// Если четное количество цифр в номере заказа.
	return nil
}

func forOdd(numbers []int64) (err error) {
	// Если нечетное количество цифр в номере заказа.
	return nil
}
