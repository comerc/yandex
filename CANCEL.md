Одновременно запускаю две горутины. Если первая горутина закончила работу, нужно как-то сообщить второй горутине, чтобы она прервала свою работу.

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Первая горутина
	go func() {
		defer cancel() // Сообщить второй горутине при завершении работы
		// Симуляция работы
		time.Sleep(2 * time.Second)
		fmt.Println("Первая горутина завершила работу")
	}()

	// Вторая горутина
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Вторая горутина прервана")
				return
			default:
				// Симуляция работы
				fmt.Println("Вторая горутина работает")
				time.Sleep(1 * time.Second)
			}
		}
	}(ctx)

	// Ожидание завершения работы горутин
	time.Sleep(5 * time.Second)
}
```
