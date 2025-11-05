# Broadcast Channels

## Обзор

Broadcast каналы позволяют отправлять сообщения множественным получателям. Существуют три основные модели: **Push**, **Pull** и **Fan-out**.

## 1. Push модель (PushBroadcaster)

**Принцип:** Отправитель активно "толкает" сообщения всем получателям.

### Реализация

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// PushBroadcaster - простой broadcaster
type PushBroadcaster struct {
    mu        sync.RWMutex
    listeners map[string]chan interface{}
    closed    bool
}

// NewPushBroadcaster создает новый broadcaster
func NewPushBroadcaster() *PushBroadcaster {
    return &PushBroadcaster{
        listeners: make(map[string]chan interface{}),
    }
}

// Subscribe добавляет нового слушателя
func (b *PushBroadcaster) Subscribe(id string, bufferSize int) (chan interface{}, error) {
    b.mu.Lock()
    defer b.mu.Unlock()

    if b.closed {
        return nil, fmt.Errorf("broadcaster закрыт")
    }

    if _, exists := b.listeners[id]; exists {
        return nil, fmt.Errorf("слушатель с ID %s уже существует", id)
    }

    listener := make(chan interface{}, bufferSize)
    b.listeners[id] = listener
    return listener, nil
}

// Broadcast отправляет сообщение всем слушателям
func (b *PushBroadcaster) Broadcast(msg interface{}) {
    b.mu.RLock()
    defer b.mu.RUnlock()

    if b.closed {
        return
    }

    for id, listener := range b.listeners {
        select {
        case listener <- msg:
            // Сообщение отправлено успешно
        default:
            // Канал заполнен, сообщение пропущено
            fmt.Printf("Listener %s: канал заполнен, сообщение пропущено\n", id)
        }
    }
}

// Close закрывает broadcaster
func (b *PushBroadcaster) Close() {
    b.mu.Lock()
    defer b.mu.Unlock()

    if b.closed {
        return
    }

    b.closed = true
    for _, listener := range b.listeners {
        close(listener)
    }
    b.listeners = nil
}
```

### Пример использования Push модели

```go
func ExamplePushModel() {
    broadcaster := NewPushBroadcaster()
    defer broadcaster.Close()

    // Подписываем слушателей
    ch1, _ := broadcaster.Subscribe("user1", 5)
    ch2, _ := broadcaster.Subscribe("user2", 3)

    // Запускаем горутины для чтения
    go func() {
        for msg := range ch1 {
            fmt.Printf("User1 получил: %s\n", msg)
        }
    }()

    go func() {
        for msg := range ch2 {
            fmt.Printf("User2 получил: %s\n", msg)
        }
    }()

    // Отправляем сообщения
    broadcaster.Broadcast("Привет всем!")
    broadcaster.Broadcast("Как дела?")
}
```

### Преимущества Push модели

**Преимущества:**
- ✅ Простота реализации
- ✅ Немедленная доставка сообщений
- ✅ Нет блокировок получателей
- ✅ Лучше для активных получателей

**Недостатки:**
- ❌ Потеря сообщений при заполненных каналах
- ❌ Нет гарантии доставки
- ❌ Тратит ресурсы на отправку всем

**Когда использовать:**
- Потеря сообщений допустима
- Получатели активно читают
- Нужна простота и производительность

## 2. Pull модель (PullBroadcaster)

**Принцип:** Получатели "вытягивают" сообщения, когда готовы их обработать.

### Реализация

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// PullBroadcaster - реализация через sync.Cond
type PullBroadcaster struct {
    mu        sync.RWMutex
    cond      *sync.Cond
    listeners map[int]chan interface{}
    nextID    int
    closed    bool
    lastMsg   interface{}
    hasMsg    bool
}

// NewPullBroadcaster создает новый broadcaster
func NewPullBroadcaster() *PullBroadcaster {
    b := &PullBroadcaster{
        listeners: make(map[int]chan interface{}),
    }
    b.cond = sync.NewCond(&b.mu)
    return b
}

// Subscribe добавляет нового слушателя
func (b *PullBroadcaster) Subscribe(ctx context.Context, bufferSize int) (int, chan interface{}) {
    b.mu.Lock()
    defer b.mu.Unlock()

    if b.closed {
        ch := make(chan interface{})
        close(ch)
        return -1, ch
    }

    id := b.nextID
    b.nextID++

    listener := make(chan interface{}, bufferSize)
    b.listeners[id] = listener

    // Запускаем горутину для слушания
    go b.listen(ctx, id, listener)

    return id, listener
}

// listen горутина для каждого слушателя
func (b *PullBroadcaster) listen(ctx context.Context, id int, ch chan interface{}) {
    defer func() {
        b.mu.Lock()
        delete(b.listeners, id)
        close(ch)
        b.mu.Unlock()
    }()

    for {
        b.mu.Lock()

        // Ждем сообщения или закрытия
        for !b.hasMsg && !b.closed {
            select {
            case <-ctx.Done():
                b.mu.Unlock()
                return
            default:
                b.cond.Wait() // Блокируемся до получения сигнала
            }
        }

        if b.closed {
            b.mu.Unlock()
            return
        }

        msg := b.lastMsg
        b.mu.Unlock()

        // Отправляем сообщение в канал
        select {
        case ch <- msg:
        case <-ctx.Done():
            return
        default:
            fmt.Printf("Listener %d: канал заполнен, сообщение пропущено\n", id)
        }
    }
}

// Broadcast отправляет сообщение всем слушателям
func (b *PullBroadcaster) Broadcast(msg interface{}) {
    b.mu.Lock()
    defer b.mu.Unlock()

    if b.closed {
        return
    }

    b.lastMsg = msg
    b.hasMsg = true

    // Уведомляем всех ожидающих горутин
    b.cond.Broadcast()

    // Уведомляем одну горутину (альтернатива Broadcast)
    // b.cond.Signal()

    // Сбрасываем флаг после небольшой задержки
    go func() {
        time.Sleep(time.Millisecond)
        b.mu.Lock()
        b.hasMsg = false
        b.mu.Unlock()
    }()
}


// Close закрывает broadcaster
func (b *PullBroadcaster) Close() {
    b.mu.Lock()
    defer b.mu.Unlock()

    if b.closed {
        return
    }

    b.closed = true
    b.cond.Broadcast()
}
```

### Пример использования Pull модели

```go
func ExamplePullModel() {
    broadcaster := NewPullBroadcaster()
    defer broadcaster.Close()

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    var wg sync.WaitGroup

    // Создаем несколько слушателей
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(listenerNum int) {
            defer wg.Done()

            id, ch := broadcaster.Subscribe(ctx, 5)
            fmt.Printf("Listener %d подписался с ID %d\n", listenerNum, id)

            for msg := range ch {
                fmt.Printf("Listener %d получил: %v\n", listenerNum, msg)
                if listenerNum == 2 {
                    time.Sleep(20 * time.Millisecond) // Медленная обработка
                }
            }
        }(i)
    }

    // Даем время на подписку
    time.Sleep(50 * time.Millisecond)

    // Отправляем сообщения
    for i := 0; i < 5; i++ {
        broadcaster.Broadcast(fmt.Sprintf("Сообщение %d", i+1))
        time.Sleep(30 * time.Millisecond)
    }

    cancel()
    wg.Wait()
}
```

### Преимущества Pull модели

**Преимущества:**
- ✅ Гарантия доставки сообщений
- ✅ Эффективность - сообщения только тем, кто ждет
- ✅ Контроль потока сообщений
- ✅ Лучше для пассивных получателей

**Недостатки:**
- ❌ Сложность реализации
- ❌ Потенциальные deadlock'и при неправильном использовании
- ❌ Получатели могут блокироваться

**Когда использовать:**
- Важна гарантия доставки
- Получатели работают "по требованию"
- Нужен контроль над потоком

## 3. Fan-out паттерн

**Принцип:** Один источник сообщений распределяет их между множественными получателями через отдельные каналы.

### Реализация Fan-out (FanOutBroadcaster)

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// FanOutBroadcaster - реализация через fan-out паттерн
type FanOutBroadcaster struct {
    mu        sync.RWMutex
    listeners map[int]chan interface{}
    nextID    int
    closed    bool
    input     chan interface{}
}

// NewFanOutBroadcaster создает новый fan-out broadcaster
func NewFanOutBroadcaster(bufferSize int) *FanOutBroadcaster {
    b := &FanOutBroadcaster{
        listeners: make(map[int]chan interface{}),
        input:     make(chan interface{}, bufferSize),
    }
    
    // Запускаем горутину для распределения сообщений
    go b.distribute()
    
    return b
}

// distribute распределяет сообщения между слушателями
func (b *FanOutBroadcaster) distribute() {
    for msg := range b.input {
        b.mu.RLock()
        
        // Отправляем сообщение всем активным слушателям
        for id, listener := range b.listeners {
            select {
            case listener <- msg:
                // Сообщение отправлено успешно
            default:
                // Канал заполнен, пропускаем
                fmt.Printf("Listener %d: канал заполнен, сообщение пропущено\n", id)
            }
        }
        
        b.mu.RUnlock()
    }
}

// Subscribe добавляет нового слушателя
func (b *FanOutBroadcaster) Subscribe(bufferSize int) (int, chan interface{}) {
    b.mu.Lock()
    defer b.mu.Unlock()

    if b.closed {
        ch := make(chan interface{})
        close(ch)
        return -1, ch
    }

    id := b.nextID
    b.nextID++

    listener := make(chan interface{}, bufferSize)
    b.listeners[id] = listener

    return id, listener
}

// Send отправляет сообщение в input канал
func (b *FanOutBroadcaster) Send(msg interface{}) {
    if !b.closed {
        b.input <- msg
    }
}

// Close закрывает broadcaster
func (b *FanOutBroadcaster) Close() {
    b.mu.Lock()
    defer b.mu.Unlock()

    if b.closed {
        return
    }

    b.closed = true
    close(b.input)
    
    // Закрываем все каналы слушателей
    for _, listener := range b.listeners {
        close(listener)
    }
    b.listeners = nil
}
```

### Пример использования Fan-out

```go
func ExampleFanOut() {
    broadcaster := NewFanOutBroadcaster(100)
    defer broadcaster.Close()

    // Подписываем слушателей
    ch1, _ := broadcaster.Subscribe(5)
    ch2, _ := broadcaster.Subscribe(3)
    ch3, _ := broadcaster.Subscribe(10)

    // Запускаем горутины для чтения
    var wg sync.WaitGroup
    wg.Add(3)

    go func() {
        defer wg.Done()
        for msg := range ch1 {
            fmt.Printf("Worker1 получил: %v\n", msg)
            time.Sleep(10 * time.Millisecond)
        }
    }()

    go func() {
        defer wg.Done()
        for msg := range ch2 {
            fmt.Printf("Worker2 получил: %v\n", msg)
            time.Sleep(20 * time.Millisecond)
        }
    }()

    go func() {
        defer wg.Done()
        for msg := range ch3 {
            fmt.Printf("Worker3 получил: %v\n", msg)
            time.Sleep(15 * time.Millisecond)
        }
    }()

    // Отправляем сообщения
    for i := 0; i < 10; i++ {
        broadcaster.Send(fmt.Sprintf("Задача %d", i+1))
        time.Sleep(50 * time.Millisecond)
    }

    time.Sleep(200 * time.Millisecond)
}
```

### Преимущества Fan-out паттерна

**Преимущества:**
- ✅ **Асинхронность** - отправитель не блокируется
- ✅ **Буферизация** - сообщения накапливаются в input канале
- ✅ **Масштабируемость** - легко добавлять/удалять слушателей
- ✅ **Изоляция** - каждый слушатель работает независимо
- ✅ **Производительность** - нет блокировок при отправке

**Недостатки:**
- ❌ **Потеря сообщений** при заполненных каналах
- ❌ **Потребление памяти** - input канал может накапливать сообщения
- ❌ **Нет гарантии доставки** - как и в Push модели

**Когда использовать:**
- Нужна асинхронная отправка сообщений
- Важна производительность отправителя
- Слушатели работают с разной скоростью
- Нужна буферизация сообщений

## Сравнение всех трех моделей

```
| Аспект              | Push модель      | Pull модель      | Fan-out модель   |
|---------------------|------------------|------------------|------------------|
| Доставка            | Немедленная      | По требованию    | Асинхронная      |
| Гарантия            | Нет              | Есть             | Нет              |
| Сложность           | Простая          | Сложная          | Средняя          |
| Блокировки          | Нет              | Есть             | Нет              |
| Ресурсы             | Тратит на всех   | Эффективна       | Буферизует       |
| Производительность  | Высокая          | Средняя          | Очень высокая    |
| Масштабируемость    | Средняя          | Низкая           | Высокая          |
```

## Заключение

Все три модели имеют свои преимущества и недостатки. Выбор зависит от требований к надежности, простоте и производительности. Push модель подходит для большинства случаев, Pull модель - для критически важных систем, где потеря сообщений недопустима, а Fan-out модель - для высокопроизводительных асинхронных систем.
