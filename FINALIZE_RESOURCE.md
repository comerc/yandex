```go
package main

import (
	"fmt"
	"runtime"
	"time"
)

type Resource struct {
	name string
}

func finalizeResource(r *Resource) {
	fmt.Printf("Finalizing resource: %s\n", r.name)
}

func main() {
	r := &Resource{name: "MyResource"}

	// Устанавливаем финализатор для объекта r
	runtime.SetFinalizer(r, finalizeResource)

	// Симулируем выход за пределы области видимости
	r = nil

	// Форсируем сборку мусора (не рекомендуется использовать в продакшене)
	runtime.GC()

	// Дадим время для выполнения финализации
	// В реальной программе это не нужно
	time.Sleep(time.Second)
}
```