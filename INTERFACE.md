```go
package main

import "fmt"

type Celsius float64

func (Celsius) String() string {
	return "C"
}

type Fahrenheit float64

func (*Fahrenheit) String() string {
	return "F"
}

func main() {
	var _ fmt.Stringer = (*Celsius)(nil)
	var _ fmt.Stringer = (*Fahrenheit)(nil) // компиляция проходит

	var _ fmt.Stringer = Celsius(20.0)
	// var _ fmt.Stringer = Fahrenheit(20.0) // компиляция не проходит
}
```

Методы, определенные для типа T, также доступны для указателей этого типа (*T). Однако обратное не верно: методы, определенные для *T, не доступны для T.
