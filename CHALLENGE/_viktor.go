package example

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
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
	var item Item
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		item, err := e.storeV1.Get(egCtx, id)
		if err != nil {
			return fmt.Errorf("failed to get item (id: %d) from storeV1: %w", id, err)
		}

		return nil
	})
	eg.Go(func() error {
		item, err := e.storeV2.Get(egCtx, id)
		if err != nil {
			return fmt.Errorf("failed to get item (id: %d) from storeV2: %w", id, err)
		}
	})

	return &item, eg.Wait()
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
