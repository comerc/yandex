Наименее противный вариант, наверно:

```go
package main

import (
	"fmt"
)

type MyTest struct {
	items any
}

func main() {
	tests := []MyTest{
		{
			items: []int{1, 2, 3},
		},
		{
			items: []string{"1", "2", "3"},
		},
	}

	for _, test := range tests {
		switch v := test.items.(type) {
		case []any:
			process(v)
		default:
			fmt.Printf("unknown type: %T\n", test.items)
		}
	}
}

func process(items []any) {
	for _, item := range items {
		switch v := item.(type) {
		case string, int:
			fmt.Printf("%#v", v)
		default:
			panic("invalid type")
		}
	}
}
```

Хочется явно декларировать тип, вроде этого:

```go
type MyTest struct {
  items []int | []string
}
```

Но так сделать в Go нельзя и начинаются пляски с дженериками, которые только загромождают код.

---

Лаконичный вариант на дженериках:

```go
package main

import (
	"fmt"
)

type ItemsProcessor interface {
	Process()
}

type Items[T any] []T

func (items Items[T]) Process() {
	for _, item := range items {
		fmt.Printf("%#v ", item)
	}
	fmt.Println()
}

type MyTest struct {
	items ItemsProcessor
}

func main() {
	tests := []MyTest{
		{
			items: Items[int]{1, 2, 3},
		},
		{
			items: Items[string]{"1", "2", "3"},
		},
	}

	for _, test := range tests {
		test.items.Process()
	}
}
```
