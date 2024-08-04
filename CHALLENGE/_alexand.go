package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Item struct {
	ID int64
}

type IStore interface {
	Set(ctx context.Context, id int64, item *Item) error
	Get(ctx context.Context, id int64) (*Item, error)
	GetBatch(ctx context.Context, ids []int64) ([]*Item, error)
}

type Example struct {
	storeV1 IStore
	storeV2 IStore
	logger  *slog.Logger
}

// Get извлекает Item по заданному id, используя два хранилища (storeV1 и storeV2)
// параллельно. Возвращает первый успешный результат или ошибку, если оба хранилища
// не смогли вернуть результат.
func (e *Example) Get(ctx context.Context, id int64) (*Item, error) {
	// Определяем структуру для хранения результата операции
	type result struct {
		item *Item
		err  error
	}

	// Создаем новый контекст, который можно отменить
	innerCtx, cancel := context.WithCancel(ctx)
	defer cancel() // Гарантируем освобождение ресурсов в конце выполнения функции

	// Создаем каналы для получения результатов от каждого хранилища
	v1Chan := make(chan result, 1)
	v2Chan := make(chan result, 1)

	// Запускаем горутину для получения данных из первого хранилища
	go func() {
		item, err := e.storeV1.Get(innerCtx, id)
		select {
		case v1Chan <- result{item: item, err: err}: // Отправляем результат в канал
		case <-innerCtx.Done(): // Если контекст отменен, логируем сообщение
			e.logger.Info("storeV1: контекст отменен")
		}
	}()

	// Запускаем горутину для получения данных из второго хранилища
	go func() {
		item, err := e.storeV2.Get(innerCtx, id)
		select {
		case v2Chan <- result{item: item, err: err}: // Отправляем результат в канал
		case <-innerCtx.Done(): // Если контекст отменен, логируем сообщение
			e.logger.Info("storeV2: контекст отменен")
		}
	}()

	var errors []error
	for i := 0; i < 2; i++ { // Ожидаем максимум 2 результата
		select {
		case r := <-v1Chan: // Получаем результат из первого хранилища
			if r.err == nil {
				e.logger.Info("Получен результат из storeV1")
				return r.item, nil // Возвращаем успешный результат
			}
			e.logger.Error("Ошибка из storeV1", "error", r.err)
			errors = append(errors, fmt.Errorf("\nstoreV1: %w", r.err))
		case r := <-v2Chan: // Получаем результат из второго хранилища
			if r.err == nil {
				e.logger.Info("Получен результат из storeV2")
				return r.item, nil // Возвращаем успешный результат
			}
			e.logger.Error("Ошибка из storeV2", "error", r.err)
			errors = append(errors, fmt.Errorf("\nstoreV2: %w", r.err))
		case <-ctx.Done(): // Если внешний контекст отменен, возвращаем ошибку
			return nil, fmt.Errorf("превышено время ожидания контекста: %w", ctx.Err())
		}
	}

	// Если мы дошли до этой точки, значит обе попытки завершились с ошибкой
	return nil, fmt.Errorf("все хранилища завершились с ошибкой: \n%v\n", errors)
}

type Example2 struct {
	storeV1 IStore
	storeV2 IStore
	logger  *slog.Logger
}

func (e *Example2) Get(ctx context.Context, id int64) (*Item, error) {
	// Создаем новый контекст с возможностью отмены
	ctx, cancel := context.WithCancel(ctx)
	// Гарантируем, что функция отмены будет вызвана при выходе из Get
	defer cancel()

	// Структура для хранения результата и возможной ошибки
	type result struct {
		item *Item
		err  error
	}

	// Буферизованный канал для результатов. Размер 2 соответствует количеству хранилищ
	resultCh := make(chan result, 2)

	// Слайс с информацией о хранилищах
	stores := []struct {
		source string
		store  IStore
	}{
		{"storeV1", e.storeV1},
		{"storeV2", e.storeV2},
	}

	// Запускаем горутину для каждого хранилища
	for _, storeInfo := range stores {
		go func(source string, store IStore) {
			// Запрашиваем элемент из хранилища
			item, err := store.Get(ctx, id)
			// Пытаемся отправить результат в канал
			select {
			case resultCh <- result{item: item, err: err}:
				// Результат успешно отправлен
			case <-ctx.Done():
				// Контекст был отменен, завершаем горутину
			}
		}(storeInfo.source, storeInfo.store)
	}

	var errs []string
	// Ожидаем результаты от всех запущенных горутин
	for i := 0; i < len(stores); i++ {
		select {
		case <-ctx.Done():
			// Если контекст отменен, немедленно возвращаем ошибку
			return nil, ctx.Err()
		case res := <-resultCh:
			if res.err != nil {
				// Если получили ошибку, добавляем ее в список
				errs = append(errs, res.err.Error())
			} else {
				// Если получили успешный результат, немедленно возвращаем его
				return res.item, nil
			}
		}
	}

	// Если все запросы завершились с ошибками, возвращаем их все
	if len(errs) > 0 {
		return nil, fmt.Errorf("%s", strings.Join(errs, ", "))
	}

	// Если мы дошли до этой точки, значит ни одно хранилище не вернуло результат
	return nil, fmt.Errorf("ни одно хранилище не вернуло результат")
}

type Example3 struct {
	storeV1 IStore
	storeV2 IStore
	logger  *slog.Logger
}

func (e *Example3) Get(ctx context.Context, id int64) (*Item, error) {
	// Создаем каналы для обработки результатов и ошибок
	errCh := make(chan error, 2)
	resultCh := make(chan *Item, 1)

	// Создаем новый контекст с возможностью отмены
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // Гарантируем отмену контекста при выходе из функции

	// Определяем список хранилищ для параллельного запроса
	stores := []struct {
		source string
		store  IStore
	}{
		{"storeV1", e.storeV1},
		{"storeV2", e.storeV2},
	}

	// Запускаем горутины для каждого хранилища
	for _, storeInfo := range stores {
		go func(source string, store IStore) {
			item, err := store.Get(ctx, id)
			e.logger.Info("Запрос к хранилищу", "source", source)

			if err != nil {
				select {
				case errCh <- err:
				case <-ctx.Done():
				}
				return
			}

			// Пытаемся отправить результат, если контекст не отменен
			select {
			case resultCh <- item:
				cancel() // Отменяем контекст, так как получили успешный результат
			case <-ctx.Done():
			}
		}(storeInfo.source, storeInfo.store)
	}

	errors := make([]string, 0, len(stores))
	// Ожидаем результаты или ошибки от горутин
	for i := 0; i < len(stores); i++ {
		select {
		case <-ctx.Done():
			if len(errors) > 0 {
				return nil, fmt.Errorf("%s", strings.Join(errors, ", "))
			}
			return nil, ctx.Err()
		case item := <-resultCh:
			return item, nil // Возвращаем первый успешный результат
		case err := <-errCh:
			errors = append(errors, err.Error())
		}
	}

	// Если получили ошибки от всех хранилищ, возвращаем их
	return nil, fmt.Errorf("%s", strings.Join(errors, ", "))
}

type Example4 struct {
	storeV1 IStore
	storeV2 IStore
	logger  *slog.Logger
}

func (e *Example4) Get(ctx context.Context, id int64) (*Item, error) {
	logger := e.logger.With("id", id)

	type result struct {
		item *Item
		err  error
	}

	var wg sync.WaitGroup
	results := make(chan result, 2)

	// Функция для запуска запроса к хранилищу
	queryStore := func(name string, store interface {
		Get(context.Context, int64) (*Item, error)
	}) {
		defer wg.Done()
		item, err := store.Get(ctx, id)
		if err != nil {
			logger.Error("Не удалось получить элемент из хранилища", "store", name, "error", err)
			results <- result{err: fmt.Errorf("%s: %w", name, err)}
			return
		}
		logger.Info("Получен результат из хранилища", "store", name)
		results <- result{item: item}
	}

	wg.Add(2)
	go queryStore("storeV1", e.storeV1)
	go queryStore("storeV2", e.storeV2)

	go func() {
		wg.Wait()
		close(results)
	}()

	var errs []error
	for r := range results {
		if r.err == nil {
			return r.item, nil
		}
		errs = append(errs, r.err)
	}

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("ошибка контекста: %w", err)
	}

	return nil, fmt.Errorf("все хранилища завершились с ошибкой: %w", errors.Join(errs...))
}

// MockStore реализует интерфейс IStore для тестирования
type MockStore struct {
	name     string
	delay    time.Duration
	failMode bool
	logger   *slog.Logger
}

func (m *MockStore) Set(ctx context.Context, id int64, item *Item) error {
	return nil
}

func (m *MockStore) Get(ctx context.Context, id int64) (*Item, error) {
	select {
	case <-time.After(m.delay):
		if m.failMode {
			m.logger.Error("Не удалось получить элемент", "store", m.name)
			return nil, fmt.Errorf("%s не удалось получить элемент", m.name)
		}
		return &Item{ID: id}, nil
	case <-ctx.Done():
		m.logger.Info("Запрос отменен", "store", m.name)
		return nil, ctx.Err()
	}
}

func (m *MockStore) GetBatch(ctx context.Context, ids []int64) ([]*Item, error) {
	return nil, fmt.Errorf("не реализовано")
}

// Добавляем функцию для тестирования производительности
func runPerformanceTest(name string, getter func(context.Context, int64) (*Item, error), iterations int, logger *slog.Logger) {
	logger.Info("Производительность", "name", name)
	initialGoroutines := runtime.NumGoroutine()
	logger.Info("Количество горутин перед тестом", "goroutines", initialGoroutines)

	totalDuration := time.Duration(0)
	for i := 0; i < iterations; i++ {
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := getter(ctx, int64(i+1))
		cancel()
		duration := time.Since(start)
		totalDuration += duration

		if err != nil {
			logger.Error("Ошибка в итерации", "iteration", i+1, "error", err, "duration", duration)
		} else {
			logger.Info("Успешное выполнение итерации", "iteration", i+1, "duration", duration)
		}
	}
	time.Sleep(1 * time.Millisecond)
	currentGoroutines := runtime.NumGoroutine()
	logger.Info("Количество горутин после теста", "goroutines", currentGoroutines)

	if currentGoroutines > initialGoroutines {
		logger.Warn("Обнаружена утечка горутин", "before", initialGoroutines, "after", currentGoroutines)
	}

	avgDuration := totalDuration / time.Duration(iterations)
	logger.Info("Среднее время выполнения", "name", name, "avgDuration", avgDuration)
}

func runTest(e interface {
	Get(context.Context, int64) (*Item, error)
}, i int, timeout time.Duration, logger *slog.Logger) {
	initialGoroutines := runtime.NumGoroutine()
	logger.Info("Количество горутин перед Get", "goroutines", initialGoroutines)

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	item, err := e.Get(ctx, int64(i+1))
	cancel()
	duration := time.Since(start)

	// Даем немного времени на завершение горутин
	time.Sleep(1 * time.Millisecond)

	currentGoroutines := runtime.NumGoroutine()
	logger.Info("Количество горутин после Get", "goroutines", currentGoroutines)

	if err != nil {
		logger.Error("Ошибка при выполнении Get", "error", err, "duration", duration)
		// Добавляем дополнительную обработку для вывода отдельных ошибок
		if errors, ok := err.(interface{ Unwrap() []error }); ok {
			for _, e := range errors.Unwrap() {
				logger.Error("Детали ошибки", "error", e)
			}
		}
	} else {
		logger.Info("Получен элемент", "id", item.ID, "duration", duration)
	}

	if currentGoroutines > initialGoroutines {
		logger.Warn("Обнаружена утечка горутин", "before", initialGoroutines, "after", currentGoroutines)
	}
}

func main() {
	// Настройка логгера
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   a.Key,
					Value: slog.StringValue(time.Now().Format("02.01.06 15:04:05")),
				}
			}
			return a
		},
	})
	logger := slog.New(handler)

	// Создаем два тестовых хранилища
	storeV1 := &MockStore{name: "StoreV1", delay: 0 * time.Millisecond, failMode: false, logger: logger}
	storeV2 := &MockStore{name: "StoreV2", delay: 0 * time.Millisecond, failMode: false, logger: logger}

	// Инициализируем Example и Example2 с двумя хранилищами
	example := &Example{
		storeV1: storeV1,
		storeV2: storeV2,
		logger:  logger,
	}

	example2 := &Example2{
		storeV1: storeV1,
		storeV2: storeV2,
		logger:  logger,
	}

	example3 := &Example3{
		storeV1: storeV1,
		storeV2: storeV2,
		logger:  logger,
	}

	example4 := &Example4{
		storeV1: storeV1,
		storeV2: storeV2,
		logger:  logger,
	}

	// Тестирование производительности
	// logger.Info("Тестирование производительности")
	// iterations := 10

	// runPerformanceTest("Example", example.Get, iterations, logger)
	// runPerformanceTest("Example2", example2.Get, iterations, logger)
	// runPerformanceTest("Example3", example3.Get, iterations, logger)
	// runPerformanceTest("Example4", example4.Get, iterations, logger)

	// Определяем параметры для четырех разных тестов
	tests := []struct {
		name         string
		storeV1Delay time.Duration
		storeV2Delay time.Duration
		storeV1Fail  bool
		storeV2Fail  bool
		timeout      time.Duration
	}{
		{"Тест 1: StoreV1 быстрее", 1000 * time.Millisecond, 5000 * time.Millisecond, false, false, 5 * time.Second},
		{"Тест 2: StoreV1 медленнее", 5000 * time.Millisecond, 1000 * time.Millisecond, false, false, 5 * time.Second},
		{"Тест 3: Оба хранилища возвращают ошибку", 1000 * time.Millisecond, 1000 * time.Millisecond, true, true, 5 * time.Second},
		{"Тест 4: Оба отваливаются по таймауту", 6000 * time.Millisecond, 6000 * time.Millisecond, false, false, 5 * time.Second},
	}

	for i, test := range tests {
		logger.Info("-> Начало теста", "name", test.name)
		storeV1.delay = test.storeV1Delay
		storeV2.delay = test.storeV2Delay
		storeV1.failMode = test.storeV1Fail
		storeV2.failMode = test.storeV2Fail

		fmt.Println()
		// Тестирование Example
		logger.Info("Тестирование Example")
		runTest(example, i, test.timeout, logger)

		fmt.Println()
		// Тестирование Example2
		logger.Info("Тестирование Example2")
		runTest(example2, i, test.timeout, logger)

		fmt.Println()
		// Тестирование Example3
		logger.Info("Тестирование Example3")
		runTest(example3, i, test.timeout, logger)

		fmt.Println()
		// Тестирование Example4
		logger.Info("Тестирование Example4")
		runTest(example4, i, test.timeout, logger)
		fmt.Println()
	}
}
