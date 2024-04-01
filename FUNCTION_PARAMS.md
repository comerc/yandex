```go
package main

import "fmt"

type Person struct {
	Name string
}

func changeName(person *Person) {
	person = &Person{
		Name: "Alice",
	}
}

func main() {
	person := Person{
		Name: "Bob",
	}
	fmt.Println(person.Name)
	changeName(&person)
	fmt.Println(person.Name)
}
```

***

В языке программирования Go параметры функций передаются по значению. Это означает, что когда вы передаете переменную в функцию, Go создает копию этой переменной. Все изменения, которые вы делаете с копией, не влияют на исходную переменную.

В вашем коде вы передаете указатель на структуру `Person` в функцию `changeName`. Указатель - это переменная, которая хранит адрес другой переменной. В данном случае, указатель `person` в функции `changeName` является копией исходного указателя, и когда вы присваиваете ему новое значение (`&Person{Name: "Alice"}`), это не влияет на исходный указатель.

Если вы хотите изменить исходную структуру `Person`, вы должны изменить поле `Name` напрямую, как показано ниже:

```go
func changeName(person *Person) {
	person.Name = "Alice"
}
```

Теперь функция `changeName` изменяет поле `Name` исходной структуры `Person`, на которую указывает `person`, и вы увидите "Alice" при печати `person.Name` после вызова `changeName(person)`. 

