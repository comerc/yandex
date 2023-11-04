прикольная задачка для собеса:

```go
package main

func try(nums ...int) {
	nums[0] = 123
}

func main() {
	a := [...]int{1, 2, 3}
	try(a[:]...)
	println(a[0]) // ?
}
```

термины: sum(nums ...int) - вариативная функция (variadic function), где nums - слайс, ... - spread operator; sum(nums...) - при вызове тоже можно использовать ... - slice unpacking
