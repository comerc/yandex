## Локальная установка модуля

`go install` компилирует и устанавливает пакеты в `$(go env GOPATH)/bin`.

Особенности `go install ./...`:
- Устанавливает дочерние пакеты с `main()` 
- Имя исполняемого файла = имя директории
- `cmd/main.go` → исполняемый файл `cmd`

Команды:
- `go install ./...` - установит дочерние пакеты
- `go install path/to/package` - конкретный пакет
- `go install github.com/user/other@latest` - внешний пакет

