`type Number interface { ~int | ~int8 }` - допускает "кастомный" тип для дженериков: `func Fn[T Number](a T) {};` "~" нужна для наследников `int`, например: `type MyInt int`


```go
type MyBase interface {
	// int32 | int64 // так не работает sum() для MyCustomType
  ~int32 | ~int64
}

type MyCustomType int32

func sum[T MyBase](a, b T) T {
	return a + b
}

func main() {
	var a, b MyCustomType = 1, 2 // применяю MyCustomType для sum()

	println(sum(a, b))
}

```