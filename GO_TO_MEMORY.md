[Go To Memory](https://habr.com/ru/companies/oleg-bunin/articles/676332/)

```go
package main_test

import "testing"

func BenchmarkAssignmentIndirect(b *testing.B) {
	type X struct {
		p *int
	}
	for i := 0; i < b.N; i++ {
		// go test -bench=. >>> 0.2588 ns/op
		var i1 int
		x1 := &X{
			p: &i1,
		}
		_ = x1

		// go test -bench=. >>> 12.36 ns/op
		// var i2 int
		// x2 := &X{}
		// x2.p = &i2
	}
}
```
