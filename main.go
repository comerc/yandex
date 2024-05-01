package main

import (
	"fmt"
	"strconv"
	"strings"
)

func convStrToInt(v string) int {
	num, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}
	return num
}

func fn(predicted, actual string) int {
	if predicted == actual {
		return 2
	}

	predictedScores := strings.Split(predicted, ":")
	actualScores := strings.Split(actual, ":")

	predictedHome := convStrToInt(predictedScores[0])
	predictedAway := convStrToInt(predictedScores[1])

	actualHome := convStrToInt(actualScores[0])
	actualAway := convStrToInt(actualScores[1])

	if (predictedHome > predictedAway && actualHome > actualAway) ||
		(predictedHome < predictedAway && actualHome < actualAway) ||
		(predictedHome == predictedAway && actualHome == actualAway) {
		return 1
	}

	return 0
}

func main() {
	fmt.Println(fn("2:3", "2:3") == 2) // если полностью угадал
	fmt.Println(fn("2:1", "2:0") == 1) // хозяева
	fmt.Println(fn("1:2", "0:2") == 1) // гости
	fmt.Println(fn("2:2", "1:1") == 1) // ничья
	fmt.Println(fn("2:3", "3:2") == 0) // не угадал
}

// Представь, что ты работаешь в букмекерской конторе программистом
// и тебя попросили написать функцию, которая на вход принимает два футбольных счета
// - тот который загадал клиент когда делал ставку и реальный результат футбольного матча
// (то как сыграли команды на самом деле). На выходе нужно получить:
// 2 - если клиент полностью угадал счет
// 1 - если клиент угадал победившую команду (в том числе ничью)
// 0 - если не угадал ничего
// Счет задается строкой вида “2:3”.
// Первое число называется счет хозяев, второе - счет гостей.
