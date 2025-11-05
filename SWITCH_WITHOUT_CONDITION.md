# Switch без условия

Switch statement без условия (без переменной после `switch`) эквивалентен `switch true` и работает как цепочка `if-else-if`, где каждый `case` содержит логическое выражение, что позволяет более читаемо организовывать проверки нескольких условий с ранним выходом и предотвращает глубокую вложенность if'ов.

```go
// Вместо вложенных if'ов:
if x > 10 {
    fmt.Println("x больше 10")
} else if x > 5 {
    fmt.Println("x больше 5")
} else if x > 0 {
    fmt.Println("x положительное")
} else {
    fmt.Println("x неположительное")
}

// Можно использовать switch без условия:
switch { // нет переменной - каждый case содержит условие
case x > 10:
    fmt.Println("x больше 10")
case x > 5:
    fmt.Println("x больше 5")
case x > 0:
    fmt.Println("x положительное")
default:
    fmt.Println("x неположительное")
}

// Эквивалентно switch true:
switch true { // каждый case сравнивается с true
case x > 10:
    fmt.Println("x больше 10")
case x > 5:
    fmt.Println("x больше 5")
case x > 0:
    fmt.Println("x положительное")
default:
    fmt.Println("x неположительное")
}
```
