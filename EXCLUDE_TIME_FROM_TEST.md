Правильный ответ - добавить контекст и вынести `period`. Тогда мы сможем управлять временем в тестах без добавления каких-то зависимостей от этих тестов.

```go
func (p *Processor) Run(ctx context.Context, period time.Duration) {

	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Контекст отменен, выходим из цикла
			return
		case <-ticker.C:
			// Продолжаем цикл
      // ...
		}
	}
}
```

Если тебе "что-то" мешает написать тест, значит это "что-то" не на своём месте.