# Структуры без имени

Механизм создания структур непосредственно в точке использования позволяет описывать сложные данные "на лету", не объявляя заранее их тип. Это удобно для временных конструкций, когда нужна быстрая и компактная структуризация информации без overhead объявления именованного типа.

## Примеры

### Простое создание
```go
myCar := struct { 
    make, model string 
}{"tesla", "model3"}
```

### Вложенная структура
```go
person := struct {
    name string
    car struct {
        make, model string
    }
}{
    name: "Алиса",
    car: struct {
        make, model string
    }{"tesla", "model3"}
}
```

### Использование в тестах
```go
testCases := []struct {
    input    int
    expected bool
}{
    {5, true},
    {10, false},
}
```