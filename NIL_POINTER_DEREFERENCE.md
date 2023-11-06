код приведёт к ошибке, какой линтер может это предупредить? ChatGPT говорит:

> К сожалению, на данный момент нет линтера в Go, который мог бы предупредить об этой конкретной ошибке. Это связано с тем, что в Go nil может иметь тип, и проверка на nil не гарантирует отсутствие ошибки nil pointer dereference.

```go
package main

import (
  "bytes"
  "io"
)

func fn(out io.Writer) {
  if out != nil {
    out.Write([]byte("OK\n"))
  }
}

func main() {
  var buf *bytes.Buffer // вместо io.Writer
  fn(buf)
}
```
