с помощью `runtime.Memstats` мы можем записывать статистику распределителя памяти, например количество байтов, выделенных в куче:

```go
func printAlloc() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%d KB\n", m.Alloc/1024)
}
```

функция полезна для изучения поведения `runtime.GC()`

```go
package main

import (
	"fmt"
	"runtime"
)

type Foo struct {
	v []byte
}

func main() {
	foos := make([]Foo, 1_000)
	printAlloc()
	for i := 0; i < len(foos); i++ {
		foos[i] = Foo{
			v: make([]byte, 1024*1024),
		}
	}
	printAlloc()
	two := keepFirstTwoElementsOnly(foos)
	runtime.GC()           // Вызывается GC, чтобы принудительно вызвать очистку кучи printAlloc()
	runtime.KeepAlive(two) // Сохраняется ссылка на переменную two
	printAlloc()
}

func keepFirstTwoElementsOnly(foos []Foo) []Foo {
  // утечка
  return foos[:2]
  // решение с расходом mem
	// res := make([]Foo, 2)
	// copy(res, foos)
	// return res
  // решение с расходом cpu
	// for i := 2; i < len(foos); i++ {
	// 	foos[i].v = nil
	// }
	// return foos[:2]
}
```
