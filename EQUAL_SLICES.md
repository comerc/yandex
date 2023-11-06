getInternalArray:

```go
package main

import (
	"unsafe"
)

func getInternalArray[T unsafe.Pointer](nums ...int) T {
	return *(*T)(unsafe.Pointer(&nums[0]))
}

func main() {
	a := []int{1, 2, 3}
	b := []int{1, 2, 3}
	c := []int{3, 2, 1}
	println(getInternalArray(a...) == getInternalArray(b...)) // true
	println(getInternalArray(b...) == getInternalArray(c...)) // false
}
```

проигрывает в произволительности на порядок встроенной функции `slices.Equal`:

```go
package main

import (
	"slices"
	"testing"
	"unsafe"
)

func getInternalArray[T any](nums ...int) T {
	return *(*T)(unsafe.Pointer(&nums[0]))
}

var (
	a, c []int
)

func init() {
	a = make([]int, 10_000_000)
	c = make([]int, 10_000_000)
	for i := 0; i < 10_000_000; i++ {
		a[i] = i
		c[i] = i
	}
}

func BenchmarkOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = getInternalArray[[10_000_000]int](a...) == getInternalArray[[10000000]int](c...)
	}
}

func BenchmarkTwo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		slices.Equal(a, c)
	}
}
```

```
cpu: Intel(R) Core(TM) i5-10400 CPU @ 2.90GHz
BenchmarkOne-12    	      13	  86404167 ns/op	160006320 B/op	       2 allocs/op
BenchmarkTwo-12    	     120	   9667527 ns/op	       0 B/op	       0 allocs/op
```

но можно ускориться:

```go
package main

import (
	"slices"
	"testing"
	"unsafe"
)

func EqualSlices[T comparable](s1, s2 []int) bool {
	return *(*T)(unsafe.Pointer(&s1[0])) == *(*T)(unsafe.Pointer(&s2[0]))
}

var (
	a, c []int
)

func init() {
	a = make([]int, 10_000_000)
	c = make([]int, 10_000_000)
	for i := 0; i < 10_000_000; i++ {
		a[i] = i
		c[i] = i
	}
}

func BenchmarkOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EqualSlices[[10_000_000]int](a, c)
	}
}

func BenchmarkTwo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		slices.Equal(a, c)
	}
}
```

```
cpu: Intel(R) Core(TM) i5-10400 CPU @ 2.90GHz
BenchmarkOne-12    	     157	   7760157 ns/op	       0 B/op	       0 allocs/op
BenchmarkTwo-12    	     124	  10045822 ns/op	       0 B/op	       0 allocs/op
```
