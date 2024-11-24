```go
package main

import "fmt"

type Key struct {
	ID int
}

func main() {
	m := make(map[Key]int)
	m[Key{ID: 1}] = 123
	fmt.Printf("%v\n", m[Key{ID: 1}]) // 123
}
```

ключ в виде структуры работает по значениям этой структуры. 