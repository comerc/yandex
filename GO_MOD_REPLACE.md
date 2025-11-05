## Replace в go.mod

### Локальная замена
```go.mod
replace mylib v0.0.0 => ../mylib
```

### Внешняя замена
```go.mod
replace github.com/original v1.0.0 => github.com/fork v1.1.0
```
