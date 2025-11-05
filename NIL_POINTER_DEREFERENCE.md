# Разименование nil-указателя

Разименование nil-указателя вызывает немедленную остановку программы с panic.

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
	// case 1: nil-указатель как аргумент интерфейса
	var buf *bytes.Buffer = nil
	fn(buf) // panic

	// case 2: прямое разименование nil-указателя
	var a *int
	b := *a // panic
	println(b)
}
```

> nilaway помогает выявлять второй случай, но не первый.