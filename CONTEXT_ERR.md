```go
  for {
    select {
    case <-ctx.Done():
      return ctx.Err()
    case <-ticker.C:
      // ...
    }
  }
```
