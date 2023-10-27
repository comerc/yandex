https://youtu.be/MHn-taXfQ8o?t=3748

```go
package main

import (
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	runtime.MemProfileRate = 0
	defer func() {
		w, _ := os.Create("mem.pb.gz")
		runtime.GC()
		_ = pprof.WriteHeapProfile(w)
		_ = w.Close()
	}()

	w, _ := os.Create("cpu.pb.gz")
	_ = pprof.StartCPUProfile(w)
	defer func() {
		pprof.StopCPUProfile()
		_ = w.Close()
	}()

	// ...
}
```
