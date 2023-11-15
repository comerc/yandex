Ложное совместное использование происходит, когда какая-то кэш-линия совместно используется двумя ядрами и при этом хотя бы одна горутина что-то записывает в память. Можно предотвратить его либо с помощью заполнения (padding), либо с помощью коммуницирования через каналы.

```go
package main

import (
	"sync"
	"testing"
)

type Input struct {
	a int64
	b int64
}

type Result struct {
	sumA int64
	// _    [56]byte // улучшение до 20%
	sumB int64
}

func count(inputs []Input) Result {
	wg := sync.WaitGroup{}
	wg.Add(2)
	result := Result{}
	go func() {
		for i := 0; i < len(inputs); i++ {
			result.sumA += inputs[i].a
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < len(inputs); i++ {
			result.sumB += inputs[i].b
		}
		wg.Done()
	}()
	wg.Wait()
	return result
}

const ln = 10_000_000

func BenchmarkCount(b *testing.B) {
	inputs := make([]Input, ln)
	for i := 0; i < ln; i++ {
		inputs[i] = Input{a: int64(i), b: int64(i)}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count(inputs)
	}
}
```
