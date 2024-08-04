// Кирилл Жухарев

package example

import (
	"context"
	"fmt"
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
	// enable priority read from new storage
	// ...
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	resultChan := make(chan struct {
		item *Item
		err  error
	})

	for _, store := range []IStore{e.storeV1, e.storeV2} {
		go func(s IStore) {
			item, err := s.Get(ctx, id)
			select {
			case resultChan <- struct {
				item *Item
				err  error
			}{item, err}:
			case <-ctx.Done():
			}
		}(store)
	}

	select {
	case result := <-resultChan:
		if result.err != nil {
			// wait second result
			select {
			case result := <-resultChan:
				if result.err != nil {
					return nil, fmt.Errorf("failed to get item: %w", result.err)
				}

				return result.item, nil
			case <-ctx.Done():
				return nil, fmt.Errorf("failed to get item: %w", ctx.Err())
			}
		}
		return result.item, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
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
