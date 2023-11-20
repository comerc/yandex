Пример функции-генератора на Go, которая возвращает канал, через который передаются значения:

```go
func genOneThruThree() chan int {
    c := make(chan int)
    go func() {
        for i := 1; i <= 3; i++ {
            c <- i
        }
        close(c)
    }()
    return c
}

func main() {
    generator := genOneThruThree()
    for value := range generator {
        fmt.Println(value)
    }
}
```
