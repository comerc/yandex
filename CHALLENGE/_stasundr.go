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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	type Result struct {
		Item *Item
		Err  error
	}

	v1Chan := make(chan Result, 1)
	v2Chan := make(chan Result, 1)

	go func() {
		item, err := e.storeV1.Get(ctx, id)
		v1Chan <- Result{Item: item, Err: err}
	}()

	go func() {
		item, err := e.storeV2.Get(ctx, id)
		v2Chan <- Result{Item: item, Err: err}
	}()

	select {
	case result := <-v1Chan:
		if result.Err == nil {
			return result.Item, nil
		}
		select {
		case result := <-v2Chan:
			if result.Err == nil {
				return result.Item, nil
			}
			return nil, fmt.Errorf("failed to get item from both stores: v1(%v), v2(%v)", result.Err, result.Err)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	case result := <-v2Chan:
		if result.Err == nil {
			return result.Item, nil
		}
		select {
		case result := <-v1Chan:
			if result.Err == nil {
				return result.Item, nil
			}
			return nil, fmt.Errorf("failed to get item from both stores: v1(%v), v2(%v)", result.Err, result.Err)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
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
