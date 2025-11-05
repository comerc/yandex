Краткие ответы на вопросы для собеседования.

## Go Questions

**Слайсы**
- Динамические массивы с указателем на массив, длиной и вместимостью
- `make([]int, len, cap)`, append может вызвать перевыделение памяти
- Срезы ссылаются на один массив, изменения видны во всех срезах

**Мапы**
- Хеш-таблицы, не потокобезопасны, требуют sync.RWMutex для конкурентного доступа
- `make(map[string]int)`, проверка существования: `val, ok := m[key]`
- Итерация не гарантирует порядок

**Горутины**
- Легковесные потоки выполнения, управляются runtime Go
- Запуск: `go func(){}()`, стек растет динамически (~2KB начально)
- M:N модель - M горутин на N OS потоков

**sync.\***
- `sync.Mutex/RWMutex` - мьютексы для синхронизации
- `sync.WaitGroup` - ожидание завершения горутин
- `sync.Once` - однократное выполнение
- `sync.Pool` - пул объектов для переиспользования
- `sync.Cond` - условная синхронизация горутин
- `sync.Atomic` - атомарные операции

**Шедулинг**
- Кооперативный планировщик, G-M-P модель
- P (процессоры) = GOMAXPROCS, M (потоки ОС), G (горутины)
- Work stealing между P, блокирующие операции передают P другому M

**Сборка мусора**
- Concurrent mark-and-sweep, три цвета (белый/серый/черный)
- Stop-the-world паузы минимальны (<1мс)
- Настройка через GOGC, принудительный вызов `runtime.GC()`

**Контекст**
- `context.Context` для отмены операций и передачи значений
- `WithCancel/WithTimeout/WithDeadline/WithValue`
- Всегда первый параметр в функциях

**Graceful shutdown**
- Обработка сигналов `os.Signal`, обычно SIGTERM/SIGINT
- Закрытие серверов с таймаутом: `server.Shutdown(ctx)`
- Ожидание завершения горутин через WaitGroup

**Каналы**
- CSP модель, типизированные `chan T`, `chan<- T`, `<-chan T`
- Буферизованные и небуферизованные
- `close(ch)` закрывает канал, `select` для неблокирующих операций

**Error handling**
- Panic только для критических ошибок, recover только в defer
- `errors.Is/As` для проверки и извлечения ошибок
- Кастомные ошибки через интерфейс error или errors.New

**JSON unmarshalling**
- `json.Marshal/Unmarshal`, теги `json:"field_name,omitempty"`
- Для больших JSON: `json.Decoder` с потоковым чтением
- Интерфейсы `json.Marshaler/Unmarshaler` для кастомной логики

**nil и typed nil interface**
- `var i interface{} = (*int)(nil)` - typed nil, `i != nil` = true
- Интерфейс nil только если и тип и значение nil
- Проверка: `reflect.ValueOf(i).IsNil()`

**Testing**
- `testing` пакет, функции `TestXxx(t *testing.T)`
- Бенчмарки `BenchmarkXxx(b *testing.B)`, `go test -bench=.`
- Моки через интерфейсы, table-driven tests

**Time**
- `time.Time`, `time.Duration`, часовые пояса `time.Location`
- `time.After/Sleep/Tick` для таймеров
- Парсинг: `time.Parse("2006-01-02", "2023-12-25")`

**Select**
- Неблокирующий выбор между каналами
- `default` для неблокирующей проверки
- Случайный выбор при нескольких готовых каналах

## Other

**SQL DB**
- Реляционные БД с ACID, JOIN операции, транзакции, уровни изоляции
- Индексы для ускорения запросов, нормализация схемы
- ORM vs чистый SQL, uptrace/bun

**Docker**
- Контейнеризация приложений, изоляция через namespaces/cgroups
- Dockerfile для сборки образов, docker-compose для многоконтейнерных приложений
- Слои образов, registry для хранения

**gRPC**
- HTTP/2 протокол с Protocol Buffers
- Типы вызовов: unary, server/client/bidirectional streaming
- Автогенерация кода, перехватчики (interceptors)

**Git**
- Распределенная система контроля версий
- Ветвление, слияние, конфликты
- `git rebase/merge/cherry-pick`, gitflow

**Kubernetes**
- Оркестрация контейнеров, Pods/Services/Deployments
- ConfigMaps/Secrets для конфигурации
- Автоскалирование, health checks

**REST**
- HTTP методы (GET/POST/PUT/DELETE/OPTION), stateless
- Ресурсы как URL, JSON/XML для данных
- Статус коды, idempotency

**NoSQL DB**
- Document (MongoDB), Key-Value (Redis), Column (Cassandra), Graph (Neo4j)
- Eventual consistency, CAP теорема
- Горизонтальное масштабирование

**Брокеры очередей**
- Асинхронная обработка сообщений (RabbitMQ, Kafka, NATS)
- Паттерны: pub/sub, work queues, routing
- Гарантии доставки: at-least-once, exactly-once

**Unit-тесты**
- Изолированное тестирование компонентов
- Моки и стабы, test coverage
- TDD, AAA паттерн (Arrange-Act-Assert)

**Linux**
- Процессы, файловая система, права доступа
- Shell команды, pipes, environment variables
- Сигналы, демоны, systemd

**WebSocket**
- Полнодуплексная связь поверх HTTP
- Upgrade handshake, фреймы данных
- Real-time приложения, чаты

**GraphQL**
- Язык запросов для API, single endpoint
- Схема типов, resolvers, mutations/subscriptions
- Over-fetching/under-fetching решение

**dlv (Delve)**
- Отладчик для Go
- Breakpoints, step debugging, variable inspection
- Интеграция с IDE

**Grafana**
- Визуализация метрик и логов
- Дашборды, алерты, data sources
- PromQL для Prometheus

**pprof**
- Профилирование Go приложений
- CPU/memory/goroutine профили
- `go tool pprof`, flame graphs

**Prometheus**
- Система мониторинга с pull моделью
- Метрики: counter/gauge/histogram/summary
- PromQL для запросов, alerting

**gin**
- HTTP web framework для Go
- Middleware, роутинг, JSON binding
- Быстрая производительность
