1. Почему false?
2. Что надо поправить в реализации Comparator чтобы код перестал паниковать?

```go
package main

import (
  "fmt"
)

// Comparator - интерфейс компаратора, модифицировать нельзя
type Comparator interface {
  Compare(a, b interface{}) bool
}

func wrapper(c Comparator, a, b interface{}) bool {
  verdict := c.Compare(a, b)
  fmt.Println(a, b, verdict)

  return verdict
}

type ci struct{}

// Compare - реализация компаратора, которую надо поправить
func (c *ci) Compare(a, b interface{}) bool {
  return a == b
}

type T0 int
type T1 int

func main() {
  var value0 T0
  var value1 T1

  ptr0a := (*T0)(&value0)
  ptr1a := (*T1)(ptr0a)

  ptr0b := (*T0)(&value1)
  ptr1b := (*T1)(ptr0b)

  f := &ci{}

  v0 := wrapper(f, ptr0a, ptr1a) // true
  v1 := wrapper(f, ptr0b, ptr1b) // true
  v2 := wrapper(f, ptr0a, ptr0b) // false
  v3 := wrapper(f, ptr1a, ptr1b) // false

  if !(v0 || v1) || (v2 || v3) {
    panic("failed")
  }
}
```

Давайте разберем каждый из случаев:

v0: Сравнивает ptr0a и ptr1a, которые оба указывают на value0. Однако, ptr0a является указателем на T0, а ptr1a является указателем на T1. В Go, указатели на разные типы считаются разными, даже если они указывают на одно и то же значение. Поэтому v0 возвращает false.

v1: Сравнивает ptr0b и ptr1b, которые оба указывают на value1. Опять же, ptr0b является указателем на T0, а ptr1b является указателем на T1. Поэтому v1 также возвращает false.

v2: Сравнивает ptr0a и ptr0b. Оба являются указателями на T0, но они указывают на разные значения (value0 и value1 соответственно). Поэтому v2 возвращает false.

v3: Сравнивает ptr1a и ptr1b. Оба являются указателями на T1, но они указывают на разные значения (value0 и value1 соответственно). Поэтому v3 возвращает false.

Ваш код вызывает панику, если (v0 v1) возвращает false или (v2 v3) возвращает true. Поскольку все четыре значения возвращают false, ваш код вызывает панику.

- решение с рефлексией:

```go
func (c *ci) Compare(a, b interface{}) bool {
  return fmt.Sprintf("%p", a) == fmt.Sprintf("%p", b)
}
```

- решение без рефлексии:

```go
func (c *ci) Compare(a, b interface{}) bool {
  return getAddr(a) == getAddr(b)
}

func getAddr(a interface{}) uintptr {
  return (*[2]uintptr)(unsafe.Pointer(&a))[1]
}
```
