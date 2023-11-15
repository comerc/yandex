аргументы функций **всегда** копируются...

...и для слайса:

```go
func fn(in []int) {
	println(&in) // 0xc00004c718
}

func main() {
	out := []int{1}
	println(&out) // 0xc00004c700
	fn(out)
}
```

...и для указателя:

```go
func fn(p *int) {
	println(&p) // 0xc00004c720
}

func main() {
	v := 1
	p := &v
	println(&p) // 0xc00004c728
	fn(p)
}
```
