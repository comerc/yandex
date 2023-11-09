```go
package main

import (
	"errors"
	"fmt"
)

func main() {
	err1 := errors.New("foo")
	err2 := fmt.Errorf("baz %w", err1)
	fmt.Println(err2)                // baz foo
	fmt.Println(err1)                // foo
	fmt.Println(errors.Unwrap(err2)) // foo
}
```
