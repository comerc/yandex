```go
var i int // инициализирует значение 0
var s string // инициализирует значение ""
var o struct{} // инициализирует структуру o
var a [3]int // инициализирует массив a
var v []int // объявляет nil-слайс v
var m map[string]int // объявляет nil-карту m
var ch chan string // объявляет nil-канал ch
```

точно так же это работает для полей структуры:

```go
type MyType struct {
  i int // инициализирует значение 0
  s string // инициализирует значение ""
  o struct{} // инициализирует структуру o
  a [3]int // инициализирует массив a
  v []int // объявляет nil-слайс v
  m map[string]int // объявляет nil-карту m
  ch chan string // объявляет nil-канал ch
}
var data MyType
fmt.Printf("%#v", data) 
```

***

как получить типизированный `nil` для типа `MyType`

```go
o := (*MyType)(nil) 
```

