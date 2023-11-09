В чём разница errors.Is() и errors.As()

```go
package main

import "errors"

type MyError struct {
  message string
}

func (e MyError) Error() string {
  return e.message
}

func main() {
  err := MyError{"My custom error"}
  // сравнение со значением:
  println(errors.Is(err, MyError{"My custom error"})) // true
  // сравнение с типом:
  println(errors.As(err, &MyError{})) // true
}
```
