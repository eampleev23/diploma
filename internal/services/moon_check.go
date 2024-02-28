package services

import (
	"fmt"
	"log"
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
		serv.l.ZL.Debug("it's even")
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
	serv.l.ZL.Debug("forEven start..")
	// Перебираем четные цифры в номере
	for i := 0; i < len(numbers); i += 2 {
		//serv.l.ZL.Debug("", zap.Int64("v", numbers[i]))
		// В резалт записываем произвидение четной цифры и 2
		result := numbers[i] * 2
		// Если произведение больш 9
		if result > 9 {
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
	serv.l.ZL.Debug("forEven / Записали в sum int64 0..")
	log.Println("numbers = ", numbers)
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
		//serv.l.ZL.Debug("", zap.Int64("v", numbers[i]))
		result := numbers[i] * 2
		if result > 9 {
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
	serv.l.ZL.Debug("sum2Places/ got a number", zap.Int64("number", number))
	numberString := strconv.FormatInt(number, 10)
	//serv.l.ZL.Debug("sum2Places/ преобразовали в строку", zap.String("numberString", numberString))
	numberRunes := []rune(numberString)
	//serv.l.ZL.Debug("sum2Places/ преобразовали в массив рун")
	result = 0
	//serv.l.ZL.Debug("sum2Places/ в result int64 присвоили 0", zap.Int64("result", result))
	//serv.l.ZL.Debug("sum2Places/ начнаем перебирать каждую руну..")
	for i, digit := range numberRunes {
		serv.l.ZL.Debug("sum2Places/ итерация ", zap.Int("iter #", i+1))
		//serv.l.ZL.Debug("", zap.String("sum2Places/ iter/ берем руну:", string(digit)))
		// Преобразуем руну в int64
		digitNumber, _ := strconv.ParseInt(string(digit), 10, 64)
		// Прибавляем в result
		result += digitNumber
		//serv.l.ZL.Debug("sum2Places/ в result int64 после итерации", zap.Int64("result", result))
	}
	serv.l.ZL.Debug("sum2Places/ отправляем результат", zap.Int64("result", result))
	return result
}