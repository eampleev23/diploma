package services

import (
	"fmt"
	"strconv"
	"unicode"

	"go.uber.org/zap"
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
		err = serv.forEven(numbers)
		return err
	} else {
		// Нечетное.
		serv.l.ZL.Debug("it's odd")
		err = serv.forOdd(numbers)
		return err
	}
}

func (serv *Services) forEven(numbers []int64) (err error) {
	// Если четное количество цифр в номере заказа.
	// Перебираем нечетные номера индексов в массиве, последняя цифра контрольная
	// Перебираем четные цифры в номере
	for i := 0; i < len(numbers); i += 2 {
		// serv.l.ZL.Debug("", zap.Int64("v", numbers[i]))
		// В резалт записываем произвидение четной цифры и 2
		result := numbers[i] * 2 //nolint:gomnd //такой алгоритм луны
		// Если произведение больш 9
		if result > 9 { //nolint:gomnd //такой алгоритм луны
			// Значит складывам 2 разряда результата и записываем вместо текущего значения
			numbers[i] = serv.sum2Places(result)
			// Иначе просто записываем резалт
		} else {
			// Значит просто заносим результат умножения
			numbers[i] = result
		}
	}
	var sum int64
	sum = 0
	for _, v := range numbers {
		sum += v
	}
	serv.l.ZL.Debug("forEven / Сумма значений всех цифр в результате:", zap.Int64("sum", sum))
	if sum%10 != 0 {
		return fmt.Errorf("moon test fail")
	}
	return nil
}

func (serv *Services) forOdd(numbers []int64) (err error) {
	// Если нечетное количество цифр в номере заказа.
	// Значит перебираем четные номера
	serv.l.ZL.Debug("forOdd start..")
	for i := 1; i < len(numbers); i += 2 {
		// serv.l.ZL.Debug("", zap.Int64("v", numbers[i]))
		result := numbers[i] * 2 //nolint:gomnd //такой алгоритм луны
		if result > 9 {          //nolint:gomnd //такой алгоритм луны
			// Значит складывам 2 разряда результата и записываем вместо текущего значения
			numbers[i] = serv.sum2Places(result)
		} else {
			// Значит просто заносим результат умножения
			numbers[i] = result
		}
	}
	var sum int64
	sum = 0
	for _, v := range numbers {
		sum += v
	}
	if sum%10 != 0 {
		return fmt.Errorf("moon test fail")
	}
	return nil
}

func (serv *Services) sum2Places(number int64) (result int64) {
	numberString := strconv.FormatInt(number, 10)
	result = 0
	for _, digit := range numberString {
		// Преобразуем руну в int64
		digitNumber, _ := strconv.ParseInt(string(digit), 10, 64)
		// Прибавляем в result
		result += digitNumber
	}
	return result
}
