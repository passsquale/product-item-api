package producer

import (
	"context"
	"github.com/passsquale/product-item-api/internal/app/cleaner"
	"github.com/passsquale/product-item-api/internal/app/repo"
	"github.com/passsquale/product-item-api/internal/app/sender"
	"github.com/passsquale/product-item-api/internal/model"
	"log"
	"sync"
	"time"
)

type Producer interface {
	Start()
	Close()
}

type producer struct {
	producerCount  uint64
	cleanerChannel chan<- cleaner.PackageCleanerEvent
	sender         sender.EventSender
	eventsChannel  chan model.ItemEvent

	wg     *sync.WaitGroup
	cancel context.CancelFunc
}

type ProducerConfig struct {
	ProducerCount  uint64
	Repo           repo.EventRepo
	Sender         sender.EventSender
	CleanerChannel chan<- cleaner.PackageCleanerEvent
	EventsChannel  chan model.ItemEvent
}

func NewKafkaProducer(cfg ProducerConfig) Producer {

	wg := &sync.WaitGroup{}

	return &producer{
		producerCount:  cfg.ProducerCount,
		cleanerChannel: cfg.CleanerChannel,
		sender:         cfg.Sender,
		eventsChannel:  cfg.EventsChannel,
		wg:             wg,
		cancel: func() {
		},
	}
}

func (p *producer) runHandler(ctx context.Context) {
	for {
		select {
		case event := <-p.eventsChannel:
			switch err := p.sender.Send(&event); err {
			case nil:
				p.cleanerChannel <- cleaner.PackageCleanerEvent{
					Status:  cleaner.Ok,
					EventID: event.ID,
				}
			default:
				log.Println(event, err)
				p.cleanerChannel <- cleaner.PackageCleanerEvent{
					Status:  cleaner.Fail,
					EventID: event.ID,
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

func (p *producer) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	for i := uint64(0); i < p.producerCount; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			p.runHandler(ctx)
		}()
	}

	log.Printf("producer started with %d workers", p.producerCount)
}

func (p *producer) Close() {

	for len(p.eventsChannel) != 0 {
		time.Sleep(250 * time.Millisecond)
	}

	p.cancel()
	p.wg.Wait()
}
