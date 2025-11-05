# Иерархический контекст

**Ссылка**: https://github.com/mcfly722/context

Альтернативная реализация пакета `context` для Go, решающая проблему с порядком закрытия контекстов в иерархических системах.

## Проблема стандартного context

В стандартном пакете `context` родительский контекст может завершиться раньше дочерних, что приводит к непредсказуемому поведению при попытке использовать уже закрытые ресурсы родителей.

### Механизм отмены в стандартном context

```go
import (
    "context"
    "fmt"
    "time"
)

func main() {
    parentCtx, cancel := context.WithCancel(context.Background())
    
    // Дочерний контекст наследует механизм отмены
    childCtx, _ := context.WithCancel(parentCtx)
    grandChildCtx, _ := context.WithTimeout(childCtx, time.Hour)
    
    // Все разделяют один Done() канал
    go func() { <-parentCtx.Done(); fmt.Println("parent done") }()
    go func() { <-childCtx.Done(); fmt.Println("child done") }()
    go func() { <-grandChildCtx.Done(); fmt.Println("grandchild done") }()
    
    cancel() // Все три горутины проснутся одновременно!
}
```

**Проблема**: Контексты сигнализируют о закрытии, но не контролируют порядок завершения горутин.

## Решение: Hierarchical Context

Библиотека гарантирует, что контексты закрываются в обратном порядке (дети → родители), и родитель ждет завершения всех дочерних контекстов.

### Основные особенности

1. **Контролируемый порядок закрытия**: Родитель не завершится пока дети не закончат работу
2. **Защита от race conditions**: Гарантированная последовательность завершения
3. **Иерархическая структура**: Явное дерево зависимостей контекстов

### Интерфейс ContextedInstance

```go
import (
    "fmt"
    "time"
    context "github.com/mcfly722/context"
)

type node struct {
    name string
}

func (n *node) Go(current context.Context) {
    for {
        select {
        case <-time.After(time.Second):
            fmt.Printf("Node %s working...\n", n.name)
        case _, isOpened := <-current.Context():
            if !isOpened {
                fmt.Printf("Node %s shutting down\n", n.name)
                return
            }
        }
    }
}
```

### Создание иерархии контекстов

```go
// Создание корневого контекста
node0 := &node{name: "root"}
ctx0 := context.NewRootContext(node0)

// Создание дочерних контекстов
node1 := &node{name: "child1"}
ctx1, err := ctx0.NewContextFor(node1)
if err != nil {
    // Обработка ошибки - родитель может быть в состоянии закрытия
    return
}

// Создание поддочернего контекста
node2 := &node{name: "grandchild"}
ctx2, err := ctx1.NewContextFor(node2)
if err != nil {
    return
}

// Закрытие всей иерархии
ctx0.Close() // Закроется в порядке: grandchild → child1 → root
```

## Ограничения

1. **Обязательная проверка канала**: Нельзя выходить из горутины без проверки `current.Context()` - защита от deadlock'ов
2. **Обработка ошибок**: Всегда проверяйте ошибку при создании дочерних контекстов
3. **Отсутствие встроенной поддержки**: Нет deadline/timeout/value - реализуйте самостоятельно:

```go
// Вместо context.WithTimeout()
select {
case <-time.After(timeout):
    return
case <-ctx.Context():
    return
}
```

## Сравнение подходов

| Аспект | Стандартный context | Hierarchical context |
|--------|---------------------|----------------------|
| Отмена | ✅ Отлично работает | ✅ Работает |
| Порядок завершения | ❌ Не контролируется | ✅ Гарантирован |
| Защита от ошибок | ❌ Требует внимания | ✅ Принудительная |
| Deadline/Timeout | ✅ Встроенная поддержка | ❌ Реализуйте сами |
| Сложность | Простой | Более строгий |

## Когда использовать

**Выбирайте hierarchical context когда:**
- Нужен строгий контроль порядка завершения горутин
- Важна предсказуемость поведения при остановке системы
- Есть иерархические ресурсы с зависимостями

**Выбирайте стандартный context когда:**
- Нужна простая отмена операций
- Требуются встроенные deadline и timeout
- Нет сложных зависимостей между горутинами

Так что стандартный context отлично отменяет, но не контролирует порядок завершения. А hierarchical context - наоборот: жертвует удобством ради надёжности.