package repo

import "github.com/passsquale/product-item-api/internal/model"

type EventRepo interface {
	Lock(n uint64) ([]model.ItemEvent, error)
	Unlock(eventIDs []uint64) error

	Add(event []model.ItemEvent) error
	Remove(eventIDs []uint64) error
}
