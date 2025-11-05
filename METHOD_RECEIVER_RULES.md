# Методы типов и указателей

Методы, определённые для значения типа `T`, автоматически доступны и для указателей `*T`. Однако методы, определённые только для указателей `*T`, не могут быть вызваны для значений типа `T`.

```go
type MyStruct struct{}

func (t MyStruct) ValueMethod() {
    println("Метод для значения")
}

func (t *MyStruct) PointerMethod() {
    println("Метод для указателя")
}

func main() {
    s := MyStruct{}
    ps := &s

    s.ValueMethod()     // ✓ Работает
    ps.ValueMethod()    // ✓ Работает
    
    ps.PointerMethod() // ✓ Работает
    s.PointerMethod()  // ✗ Не скомпилируется
}
```