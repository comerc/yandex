# Nil интерфейсных типов

Интерфейс не предоставляет типизированный nil, в отличие от указателей конкретных типов. Nil интерфейса не имеет конкретного типа.

```go
var w io.Writer         // Nil интерфейс без типа
var f *os.File = nil    // Типизированный nil указателя

func main() {
    // Nil интерфейса
    fmt.Printf("%T %v\n", w, w)  // <nil> <nil>

    // Типизированный nil указателя
    fmt.Printf("%T %v\n", f, f)  // *os.File <nil>

    // Сравнение nil интерфейса и nil указателя
    fmt.Println(w == nil)   // true
    fmt.Println(f == nil)   // true
    fmt.Println(w == f)     // Ошибка компиляции: нельзя сравнить
}
```