# Go Workspace

Структура проекта:
```
project/
├── mylib/
│   └── go.mod
└── myapp/
    └── go.mod
```

Создание workspace:
```bash
cd project
go work init myapp mylib
```

Содержимое `go.work`:
```
go 1.18

use (
    ./myapp
    ./mylib
)
```