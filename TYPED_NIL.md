```go
var data *int
println(data == nil)
```

`data` объявлен как указатель на `int`, и поскольку вы не присваиваете ему значение, его значение по умолчанию будет `nil`. Однако, это `nil` является "типизированным" `nil`, что означает, что он имеет тип `*int`.

Когда вы сравниваете `data` с `nil` в `println(data == nil)`, вы сравниваете `*int nil` (значение `data`) с `nil`. В Go, `nil` без типа может быть сравним с `nil` любого типа, и поэтому `println(data == nil)` вернет `true`.

Это отличается от случая с интерфейсами, где `nil` интерфейс не равен `nil` указателю. Это связано с тем, как интерфейсы реализованы в Go.

```go
	var data *int
	var try any
	println(data == nil)         // true
	println(data == (*int)(nil)) // true
	try = data
	println(try == nil)         // false
	println(try == (*int)(nil)) // true
```
