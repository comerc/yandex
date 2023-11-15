У примера абсолютно одинаковые результаты замеров для обоих бенчмарков; как изменить, чтобы было наглядно видна выгода от выравнивания заполнением?

Поскольку каждая кэш-линия содержит в себе большее число переменных i, итерация по срезу Foo2 требует меньшего общего количества кэш-линий

```go
package main

import (
	"testing"
)

type Foo1 struct {
	b1 byte
	i  int64
	b2 byte
}

func sum1(foos []Foo1) int64 {
	var s int64
	for i := 0; i < len(foos); i++ {
		s += foos[i].i
	}
	return s
}

type Foo2 struct {
	i  int64
	b1 byte
	b2 byte
}

func sum2(foos []Foo2) int64 {
	var s int64
	for i := 0; i < len(foos); i++ {
		s += foos[i].i
	return s
}

var ln = 10_000_000

func BenchmarkSum1(b *testing.B) {
	foos := make([]Foo1, ln)
	for i := 0; i < ln; i++ {
		foos[i] = Foo1{i: int64(i)}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum1(foos)
	}
}

func BenchmarkSum2(b *testing.B) {
	foos := make([]Foo2, ln)
	for i := 0; i < ln; i++ {
		foos[i] = Foo2{i: int64(i)}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum2(foos)
	}
}

func main() {
  println(unsafe.Sizeof(Foo1{})) // 24
  println(unsafe.Sizeof(Foo2{})) // 16
  //
	println("alignof(Foo1.i) =", unsafe.Alignof(Foo1{}.i), "offset =", unsafe.Offsetof(Foo1{}.i))
	println("alignof(Foo2.i) =", unsafe.Alignof(Foo2{}.i), "offset =", unsafe.Offsetof(Foo2{}.i))
}
```
