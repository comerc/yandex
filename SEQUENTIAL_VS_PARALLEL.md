компромисс между последовательной и параллельной обрбаботкой - требуется подбирать оптимальное значение max для конкретного железа

```go
package main_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func sequentialMergesort(s []int) {
	if len(s) <= 1 {
		return
	}
	middle := len(s) / 2
	sequentialMergesort(s[:middle])
	sequentialMergesort(s[middle:]) // Вторая половина
	merge(s, middle)                // Объединение двух половин
}

func parallelMergesortV1(s []int) {
	if len(s) <= 1 {
		return
	}
	middle := len(s) / 2
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { // Запускается первая половина работы в горутине
		defer wg.Done()
		parallelMergesortV1(s[:middle])
	}()
	go func() { // Запускается вторая половина работы в горутине
		defer wg.Done()
		parallelMergesortV1(s[middle:])
	}()
	wg.Wait()
	merge(s, middle) // Объединение этих половин
}

const max = 100 // Задание величины порога

func parallelMergesortV2(s []int) {
	if len(s) <= 1 {
		return
	}
	if len(s) <= max {
		sequentialMergesort(s) // Вызов последовательной версии
	} else { // Если порог превышен, то выполняется параллельная версия
		middle := len(s) / 2
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			parallelMergesortV2(s[:middle])
		}()
		go func() {
			defer wg.Done()
			parallelMergesortV2(s[middle:])
		}()
		wg.Wait()
		merge(s, middle)
	}
}

func merge(s []int, middle int) {
	left := make([]int, middle)
	right := make([]int, len(s)-middle)

	copy(left, s[:middle])
	copy(right, s[middle:])

	i := 0
	j := 0
	k := 0

	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			s[k] = left[i]
			i++
		} else {
			s[k] = right[j]
			j++
		}
		k++
	}

	for i < len(left) {
		s[k] = left[i]
		i++
		k++
	}

	for j < len(right) {
		s[k] = right[j]
		j++
		k++
	}
}

func generateRandomSlice(size int) []int {
	rand.Seed(time.Now().UnixNano())
	slice := make([]int, size)
	for i := range slice {
		slice[i] = rand.Intn(100) // Генерирует случайное число от 0 до 99
	}
	return slice
}

const ln = 1_000

// func BenchmarkSequentialMergesort(b *testing.B) {
// 	s := generateRandomSlice(ln)
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		sequentialMergesort(s)
// 	}
// }

func BenchmarkParallelMergesortV1(b *testing.B) {
	s := generateRandomSlice(ln)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parallelMergesortV1(s)
	}
}

func BenchmarkParallelMergesortV2(b *testing.B) {
	s := generateRandomSlice(ln)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parallelMergesortV2(s)
	}
}
```
