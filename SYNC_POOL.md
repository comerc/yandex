Вот пример использования `sync.Pool`:

```go
package main

import (
	"sync"
)

// Создаем пул
var bytePool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 1024)
		return &b
	},
}

func Hi() {
	// Получаем байтовый срез из пула в буфер
	a := bytePool.Get().(*[]byte)
	(*a)[0] = 'H'
	(*a)[1] = 'i'
	*a = (*a)[:2]       // Сброс буфера // ?? почему "ignoring self-assignment" (-gcflags "-m=1")
	println(string(*a)) // Выводит "Hi"
	bytePool.Put(a)     // Возвращаем буфер в пул
}

func main() {
  Hi()
}
```

В этом примере мы создаем пул байтовых срезов размером 1024. Когда нам нужен байтовый срез, мы можем получить его из пула, используя `bytePool.Get()`. После использования мы возвращаем срез обратно в пул, вызывая `bytePool.Put(a)`. Это позволяет нам повторно использовать уже выделенную память, что может быть более эффективным, чем каждый раз выделять и освобождать память. Однако стоит отметить, что `sync.Pool` предназначен для объектов, которые могут быть дорогими для создания, и его использование может быть избыточным, если объекты дешевы для создания или если они не создаются часто.

---

Второй вариант смотрится проще без россыпи звёздочек, но смущает "moved to heap: a" - если Hi() вызвать 100500 раз. Как быть?

```go
package main

import (
	"sync"
)

// Создаем пул
var bytePool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024)
	},
}

func Hi() {
	// Получаем байтовый срез из пула в буфер
	a := bytePool.Get().([]byte)
	a[0] = 'H'
	a[1] = 'i'
	a = a[:2]          // Сброс буфера
	println(string(a)) // Выводит "Hi"
	bytePool.Put(&a)   // Возвращаем буфер в пул // ?? почему &a, иначе go-staticcheck: "argument should be pointer-like to avoid allocations"
}

func main() {
	Hi()
}
```

Вариант 1

```
./main.go:10:3: moved to heap: b
./main.go:10:12: make([]byte, 1024) escapes to heap
./main.go:20:5: Hi ignoring self-assignment in *a = (*a)[:2]
./main.go:21:17: string(*a) does not escape
```

Вариант 2

```
./main.go:43:14: make([]byte, 1024) escapes to heap
./main.go:43:14: make([]byte, 1024) escapes to heap
./main.go:49:2: moved to heap: a
./main.go:53:17: string(a) does not escape
```
