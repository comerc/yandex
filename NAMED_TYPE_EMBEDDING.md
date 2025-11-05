# Неявное встраивание типов

Именованный тип как анонимное поле структуры предоставляет прямой доступ к его методам и значениям, упрощая композицию и расширение поведения с сохранением строгой типизации.

```go
type MyInt int

func (m MyInt) Double() int { return int(m * 2) }

type MyObj struct { MyInt }

func main() {
    obj := MyObj{MyInt(5)}
    println(obj.MyInt, obj.Double())
}
```
