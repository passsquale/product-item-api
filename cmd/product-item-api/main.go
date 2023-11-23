package main

import (
	restranslator "github.com/passsquale/product-item-api/internal/app/retranslator"
	"github.com/passsquale/product-item-api/internal/mocks"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	cfg := restranslator.RetranslatorConfig{
		ChannelSize:     512,
		ConsumerCount:   2,
		BatchSize:       10,
		ConsumeInterval: 2 * time.Second,
		ProducerCount:   28,
		WorkerCount:     2,
		Repo:            &mocks.MockEventRepo{},
		Sender:          &mocks.MockEventSender{},
	}

	retranslator := restranslator.NewRetranslator(cfg)
	retranslator.Start()

	<-sigs

	retranslator.Close()
}
