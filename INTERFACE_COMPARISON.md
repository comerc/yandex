# Сравнение интерфейсов: семантические особенности

Сравнивайте значения интерфейсов, только если вы уверены, что они содержат динамические значения сравниваемых типов.

```go
func main() {
    var x any = []int{1, 2, 3}
    println(x == x)  // Паника во время выполнения

    // Безопасные способы сравнения
    switch v := x.(type) {
    case []int:
        // Корректное сравнение среза
        fmt.Println(reflect.DeepEqual(x, v))
    }
}
```

## Способ проверки

```go
// gorules/gorules.go
//go:build ruleguard

package gorules

import "github.com/quasilyte/go-ruleguard/dsl"

// Запрещает сравнение interface значений между собой через ==/!= (кроме сравнения с nil).
// Причина: возможна panic при равенстве динамического типа и его несравнимости (slice/map/func).
func forbidIfaceEq(m dsl.Matcher) {
	m.Match(`$x == $y`, `$x != $y`).
		Where(
			// Поддерживает как interface{}, так и именованные интерфейсы, через Underlying().
			m["x"].Type.Underlying().Is(`interface{}`) &&
				m["y"].Type.Underlying().Is(`interface{}`) &&
				!m["x"].Text.Matches(`^nil$`) &&
				!m["y"].Text.Matches(`^nil$`),
		).
		Report(`comparison of interfaces via ==/!= is forbidden: may panic on incomparable dynamic type; use type switch + slices.Equal/maps.Equal/reflect.DeepEqual`).
		At(m["x"])
}
```

```bash
$ go install github.com/quasilyte/go-ruleguard/cmd/ruleguard@latest
$ ruleguard -rules ./gorules/gorules.go .
```

```yml
# .golangci.yml
version: "2"
linters:
  enable:
    - gocritic
  settings:
    gocritic:
      enabled-checks:
        - ruleguard
      settings:
        ruleguard:
          rules: ./gorules/gorules.go
```