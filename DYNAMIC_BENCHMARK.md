Как получить бенчмарки динамически? (применение: для выбора между параллельным и последовательным поведением программы в зависимости от конфигурации конкретной вычислительной среды и рабочих нагрузок CPU-bound + I/O-bound + mem-bound)

```go
package main

import (
	"fmt"
	"testing"
)

func BenchmarkExample(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Здесь ваш код
	}
}

func main() {
	b := testing.Benchmark(BenchmarkExample)
	fmt.Println(b.String())
}
```
