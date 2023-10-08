```go
// var o struct{} - инициализирует структуру o
// var a [3]int - инициализирует массив a
// var v []int - объявляет nil-слайс v
// var m map[string]int - объявляет nil-карту m
// var ch chan string - объявляет nil-канал ch
// if i := 0; i == 0 {} - присваивание; условие
// switch t := v.(type) - .(type) работает только внутри switch
// fallthrough - переход в следующий case (без его проверки), или в default
// a := [5]any{"1", 2, true}
// l := t[:] - не меняет unsafe.Pointer
// l := t[1:2:3] - третий член определяет ёмкость
// cap() - для массивов, как алиас на len()
// cap() - применим для слайсов и каналов
// for i, ch := range s {} "go" - ch имеет тип rune, s имеет тип string
// ch := make(chan string) - блокирует передатчик, пока не готов приёмник
// bufferedCh := make(chan string, 1) - неблокирует передатчик, пока не готов приёмник
// make(<-chan bool) / make(chan<- bool) - зачем создавать каналы только для чтения / записи?
// j, ok := <-jobs - проверка, что канал закрыт
// for t := range ticker.C {} - чтение из канала ticker.C
// wg.Add(1) запускать в том же потоке, где и wg.Wait()
// time.Tick() - лучше не использовать, it "leaks"; вместо него NewTicker() + Stop()
// myType(val) - приведение типа работает, если известен тип для val
// val.(myType) - утверждение типа работает, если тип для val неопределён (т.е. interface{})
// var i int; defer func(i int) { println(i) }(i); i = 1 - выведет 0
// var i int; defer func() { println(i) }(); i = 1 - выведет 1
// есть ли в стандартной библиотеке пакет для работы с коллекциями разных типов через дженерики?
// [Лучший regexp для Go](https://habr.com/ru/articles/756222/)
// regexp.MustCompile() для глобальных переменных вместо regexp.Compile()
// type MyType struct { Page int `json:"page"` } - поле публичное (с большой буквы) + "json-тег структуры"
// diff := time.Now().Sub(other) - разница между двумя временами
// time.Now().Add(-diff) - продвинуть время на заданную продолжительность
// u, _ := url.Parse(s); net.SplitHostPort(u.Host) - как вытащить порт
// ew := &errWriter{Writer: w} - где errWriter реализует свой метод Write с предварительной обработкой сохранённой ошибки https://habr.com/ru/articles/759920/
// os.RemoveAll("dir") - аналог "rm -rf"
// go test -v - флаг -v отображает прохождение тестов
// go test -run="VectorA.*$|TestVectorMag" -v - применение regex для фильтра тестов
// [fuzzing-тесты](https://habr.com/ru/companies/oleg-bunin/articles/709248/)
// testify - в помощь по тестированию к стандартной библиотеке
// mockery - для тестов
// [migrate](https://github.com/golang-migrate/migrate)
// flag - стандартный пакет, плюс выбор: flaggy | go-flags | pflag
// exec.Command("bash", "-c", "ls -a -l -h") - как создать полную команду в одну строку вместо exec.Command("ls", "-a", "-l", "-h")
// os.Exit(1) - игнорирует defer
// цветные логи: https://github.com/GolangLessons/url-shortener/blob/c3987f66469a8d0769add18521adb9023520be95/internal/lib/logger/handlers/slogpretty/slogpretty.go
// vegeta, wrk - для стресс-тестов
// что не нравиться в go? импорт без алиасов
// когда нужен сервер httpserver - https://github.com/evrone/go-clean-template/tree/master/pkg/httpserver
// allegro/bigcache - когда нужен просто кеш (рекомендации лучших собаководов из Avito)
// go-playground/validator - правильный валидатор
// func New(ctx context.Context, connectionString string, opts ...Option) (*Storage, error) - паттерн опций для конструктора в функциях
// tdlibClient.GetMessage(&client.GetMessageRequest{}) паттерн опций для методов в структуре
// благодатное выключение - Graceful Shutdown
// func New() или func NewSubscriber() для пакета subscriber ?
//	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM) - неправильно, signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM) - правильно. os.Interrupt == syscall.SIGINT
// "божественный конфиг" - https://youtu.be/0Fhsgmz-Gig?list=PLZvfMc-lVSSO2zhyyxQLFmio8NxvQqZoN&t=906
// ilyakaznacheev/cleanenv - yaml & env в одном флаконе + godotenv для чтения .env
// если функция кидает панику, то у неё должен быть префикс Must*, например MustLoad()
// TEST EXPLORER внутри VS Code
// jackc/pgx/v5/pgxpool / go-pg + pool - PG Pool
// Masterminds/squirrel - SQL Builder (by Avito)
// [Dependency Injection](https://youtu.be/0Fhsgmz-Gig?list=PLZvfMc-lVSSO2zhyyxQLFmio8NxvQqZoN&t=1001)
// [Dependency Injection на примере Uber fx](https://www.youtube.com/watch?v=KRdrH9a98HQ)
// внутри interface ненужно прописывать ключевое слово func
// type Number interface { ~int | ~int8 } - тип для дженериков: func Fn[T Number](a T) {}; "~" нужна для наследников int, например: type MyInt int
// er := errgroup.Group{}; eg.SetLimit(limit) - ещё один примитив синхронизации
// math.Pow() - возведение в степень
// механизм эвакуации в map
// RWMutex - читаем без блокировок на чтение, но с блокировкой на запись при чтении(!), или записи
// map в Go не гарантирует порядок ключей, ES6 - гарантирует, а Dart - нет (hash map vs b-tree map)
// log.Fatal(http.ListenAndServe(":8080", httpserver.NewHandler())) - как вариант обработки ошибок
// "божественный" main.go
// string - это тоже структура и лежит в куче; при передаче аргументом, что копируется?
// func (Bear) Speak() - можно не указывать "this" при реализации метода структуры
// Интерфейсы - способ, как следать программу SOLIDной? (Dependency Inversion)
// где лучше объявлять интерфейсы: где применяются или где реализуются?
// type (A struct {}; B struct {}) - типы можно объявлять группой
// go run . | ts '%.Ss' - не работает с println(), только с fmt.Println()
```
