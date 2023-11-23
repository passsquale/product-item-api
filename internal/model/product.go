package model

import (
	"fmt"
	"time"
)

type Item struct {
	ID        uint64
	OwnerID   uint64
	ProductID uint64
	Title     string
	Created   time.Time
	Updated   *time.Time
}

func (i *Item) String() string {
	return fmt.Sprintf("Item{id:%v, ownerId:%v, pruductId:%v, title:%v}",
		i.ID, i.OwnerID, i.ProductID, i.Title)
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
