# Nil интерфейса: семантика сравнения

Проверка `if out != nil` может не сработать для типизированного nil, когда интерфейс содержит nil-указатель. Такой nil не равен `nil`, что может привести к неожиданному поведению.

```go
func fn(out io.Writer) {
    if out != nil {
        out.Write([]byte("data"))  // Выполнится, хотя buf == nil
    }
}

func main() {
    var buf *bytes.Buffer
    fn(buf)
}
```