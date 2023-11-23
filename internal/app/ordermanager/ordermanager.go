package ordermanager

import (
	"errors"
	"github.com/passsquale/product-item-api/internal/model"
	"sync"
)

var (
	ErrIncorrectOrder = errors.New("has registred with incorrect order")
)

type OrderManager interface {
	ApproveOrder(incomingEvent model.ItemEvent) bool
	RegisterEvent(incomingEvent model.ItemEvent) error
}

func NewOrderManager() OrderManager {
	return &orderManager{
		mu:       &sync.Mutex{},
		ordermap: make(map[uint64]model.EventType),
	}
}

type orderManager struct {
	mu       *sync.Mutex
	ordermap map[uint64]model.EventType
}

func (o *orderManager) ApproveOrder(incomingEvent model.ItemEvent) bool {
	o.mu.Lock()
	defer o.mu.Unlock()

	prevEventType, ok := o.ordermap[incomingEvent.Item.ID]

	switch incomingEvent.Type {
	case model.Created:
		if !ok {
			return true
		}

	case model.Updated, model.Removed:
		if ok && (prevEventType == model.Created || prevEventType == model.Updated) {
			return true
		}
	}

	return false
}

func (o *orderManager) RegisterEvent(incomingEvent model.ItemEvent) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	prevEventType, ok := o.ordermap[incomingEvent.Item.ID]

	switch incomingEvent.Type {
	case model.Created:
		o.ordermap[incomingEvent.Item.ID] = model.Created
		if !ok {
			return nil
		}

	case model.Updated:
		o.ordermap[incomingEvent.Item.ID] = model.Updated
		if ok && (prevEventType == model.Created || prevEventType == model.Updated) {
			return nil
		}

	case model.Removed:
		delete(o.ordermap, incomingEvent.Item.ID)
		if ok && (prevEventType == model.Created || prevEventType == model.Updated) {
			return nil
		}
	}
	return ErrIncorrectOrder
}
