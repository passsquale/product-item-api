package sender

import "github.com/passsquale/product-item-api/internal/model"

type EventSender interface {
	Send(subdomain *model.ItemEvent) error
}
