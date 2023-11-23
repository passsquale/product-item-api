package model

import "fmt"

type Item struct {
	ID      uint64
	OwnerId uint64
	ProductId uint64
	Title     string
}

func (i *Item) String() string {
	return fmt.Sprintf("Item{id:%v, ownerId:%v, pruductId:%v, title:%v}",
		i.ID, i.OwnerId, i.ProductId, i.Title)
}

type EventType uint8

type EventStatus uint8

const (
	Created EventType = iota
	Updated
	Removed

	Deferred EventStatus = iota
	Processed
)

type ItemEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Item   *Item
}
