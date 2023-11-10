пример, когда -race не обнаруживает гонку данных

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
