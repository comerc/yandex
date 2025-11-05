```go
package article

import (
	"sort"
)

// Оба решения имеют свои преимущества и недостатки, как вы уже упомянули.
// Решение с использованием словарей быстрее и не меняет входные данные,
// но требует дополнительной памяти для хранения словарей.
// Решение с сортировкой работает медленнее, но не требует дополнительной памяти.
// В зависимости от конкретных требований проекта, можно выбрать наиболее подходящее решение.

func compareMaps(map1, map2 map[rune]int) bool {
	if len(map1) != len(map2) {
		return false
	}

	for key, val1 := range map1 {
		val2, ok := map2[key]
		if !ok || val1 != val2 {
			return false
		}
	}

	return true
}

func CompareStrings1(s1, s2 string) int {
	dict1 := make(map[rune]int)
	dict2 := make(map[rune]int)

	for _, char := range s1 {
		dict1[char]++
	}

	for _, char := range s2 {
		dict2[char]++
	}

	if compareMaps(dict1, dict2) {
		return 1
	} else {
		return 0
	}
}

func sortString1(s string) string {
	chars := []rune(s)
	sort.Slice(chars, func(i, j int) bool {
		return chars[i] < chars[j]
	})
	return string(chars)
}

func sortString2(s string) string {
	chars := []rune(s)
	n := len(chars)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if chars[j] > chars[j+1] {
				chars[j], chars[j+1] = chars[j+1], chars[j]
			}
		}
	}
	return string(chars)
}

func CompareStrings2(s1, s2 string) int {
	sorted1 := sortString1(s1)
	sorted2 := sortString2(s2)

	if sorted1 == sorted2 {
		return 1
	} else {
		return 0
	}
}
```