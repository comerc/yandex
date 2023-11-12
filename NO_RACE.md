пример, когда -race не обнаруживает гонку данных (если убрать for {}, то обнаруживает)

```go
package main

func main() {
	i := 0
	for {
		go func() {
			i++
		}()
		go func() {
			i++
		}()
		println(i)
	}
}
```
