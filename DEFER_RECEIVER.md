Если изменить получателя на указатель, то напечатает "bar" (иначе "foo")

```go
package main

import "fmt"

func main() {
	s := Struct{id: "foo"}
	defer s.print() // s вычисляется немедленно
	s.id = "bar"    // обновление s.id (невидимое)
}

type Struct struct {
	id string
}

func (s Struct) print() {
	fmt.Println(s.id) // foo
}
```
