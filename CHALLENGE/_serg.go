package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// MockStore реализует интерфейс IStore для тестирования
type MockStore struct {
	name     string
	delay    time.Duration
	failMode bool
}

func (m *MockStore) Set(ctx context.Context, id int64, item *Item) error {
	return nil
}

func (m *MockStore) Get(ctx context.Context, id int64) (*Item, error) {
	time.Sleep(m.delay)
	if m.failMode {
		return nil, fmt.Errorf("%s failed to get item", m.name)
	}
	return &Item{ID: id}, nil

	// select {
	// case <-time.After(m.delay):
	// 	if m.failMode {
	// 		return nil, fmt.Errorf("%s failed to get item", m.name)
	// 	}
	// 	return &Item{ID: id}, nil
	// case <-ctx.Done():
	// 	fmt.Printf("%s: request cancelled\n", m.name)
	// 	return nil, ctx.Err()
	// }

	// если отменится контекст, функция выдет не выполнив основное тело - поясни
	// select {
	// case <-time.After(m.delay):
	// case <-ctx.Done():
	// 	// fmt.Printf("%s: request cancelled\n", m.name)
	// 	return nil, ctx.Err()
	// }
	// if m.failMode {
	// 	return nil, fmt.Errorf("%s failed to get item", m.name)
	// }
	// return &Item{ID: id}, nil
}

func (m *MockStore) GetBatch(ctx context.Context, ids []int64) ([]*Item, error) {
	return nil, nil
}

func main() {
	// Создаем два тестовых хранилища
	storeV1 := &MockStore{name: "StoreV1"}
	storeV2 := &MockStore{name: "StoreV2"}

	// Инициализируем Example с двумя хранилищами
	example := &Example{
		storeV1: storeV1,
		storeV2: storeV2,
	}

	ctx := context.Background()

	// // Тест 1: StoreV1 быстрее
	// {
	// 	storeV1.delay = 100 * time.Millisecond
	// 	storeV2.delay = 200 * time.Millisecond
	// 	start := time.Now()
	// 	item, err := example.Get(ctx, 1)
	// 	duration := time.Since(start)
	// 	if err != nil {
	// 		fmt.Printf("Error: %v\n", err)
	// 	} else {
	// 		fmt.Printf("Got item with ID %d in %v\n", item.ID, duration)
	// 	}
	// }

	// Тест 2: StoreV1 медленнее
	{
		storeV1.delay = 200 * time.Millisecond
		storeV2.delay = 100 * time.Millisecond
		start := time.Now()
		item, err := example.Get(ctx, 2)
		duration := time.Since(start)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Got item with ID %d in %v\n", item.ID, duration)
		}
		time.Sleep(1 * time.Second)
		println(runtime.NumGoroutine())
	}

	// // Тест 3: Оба хранилища возвращают ошибку
	// {
	// 	storeV1.failMode = true
	// 	storeV2.failMode = true
	// 	start := time.Now()
	// 	item, err := example.Get(ctx, 3)
	// 	duration := time.Since(start)
	// 	if err != nil {
	// 		fmt.Printf("Error: %v (duration: %v)\n", err, duration)
	// 	} else {
	// 		fmt.Printf("Got item with ID %d in %v\n", item.ID, duration)
	// 	}
	// }
}

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
}

func (e *Example) Set(ctx context.Context, id int64, item *Item) error {
	// enable write to new storage
	// ...
	if err := e.storeV1.Set(ctx, id, item); err != nil {
		return fmt.Errorf("failed to store item: %w", err)
	}
	return nil
}

func (e *Example) Get(ctx context.Context, id int64) (*Item, error) {
	// 1. enable shadow read from new storage
	// ...
	// 2. enable priority read from new storage - тоже не совсем корректно,
	// т.к. мы ждём старый сторадж, даже если уже выполнится новый
	// ...

	errCh := make(chan error)
	resultCh := make(chan *Item)
	// resultCh1 := make(chan *Item)
	// resultCh2 := make(chan *Item)

	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
		// close(errCh)
		// close(resultCh1)
		// close(resultCh2)
	}()

	stores := []struct {
		source string
		store  IStore
	}{
		{"storeV1", e.storeV1},
		{"storeV2", e.storeV2},
	}

	for _, storeInfo := range stores {
		go func(source string, store IStore) {
			item, err := store.Get(ctx, id)
			println(source)
			if err != nil {
				errCh <- err
				return
			}
			select {
			case resultCh <- item:
			// case resultCh1 <- item:
			// case resultCh2 <- item:
			case <-ctx.Done():
			}
			// resultCh <- item

			// switch source {
			// case "storeV1":
			// 	resultCh1 <- item
			// case "storeV2":
			// 	resultCh2 <- item
			// }
		}(storeInfo.source, storeInfo.store)
	}

	// errs := make([]string, 0, len(stores))
	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return nil, ctx.Err()
	// 	case item := <-resultCh:
	// 		return item, nil
	// 	// case item := <-resultCh1:
	// 	// 	return item, nil
	// 	// case item := <-resultCh2:
	// 	// 	return item, nil
	// 	case err := <-errCh:
	// 		errs = append(errs, err.Error())
	// 		if len(errs) == len(stores) {
	// 			return nil, fmt.Errorf("%s", strings.Join(errs, ", "))
	// 		}
	// 	}
	// }

}

func (e *Example) GetBatch(ctx context.Context, ids []int64) ([]*Item, error) {
	// enable shadow read from new storage
	// ...
	// enable priority read from new storage
	// ...
	items, err := e.storeV1.GetBatch(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	return items, nil
}
