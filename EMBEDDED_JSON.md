как обойти вопрос встроенного поля, в котором реализован Marshaler (Unmrshaler):

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	ID        int
	time.Time // "embedded"
}

// type Marshaler interface {
// 	MarshalJSON() ([]byte, error)
// }

func (e Event) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			ID   int
			Time time.Time
		}{
			ID:   e.ID,
			Time: e.Time,
		},
	)
}

func main() {
	event := Event{ID: 1234, Time: time.Now()}
	fmt.Printf("%+v", event)
	b, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
```

такое же решение может быть с форматированием строк:

```go
func (e Event) String() (string, error) {
	return fmt.Sprint(
		struct {
			ID   int
			Time time.Time
		}{
			ID:   e.ID,
			Time: e.Time,
		},
	), nil
}
```

или можно просто добавить имя, чтобы поле time.Time больше не было встроенным:

```go
type Event struct {
	ID        int
	Time time.Time // !!!
}
```
