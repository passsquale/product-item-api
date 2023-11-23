package consumer

import (
	"github.com/passsquale/product-item-api/internal/mocks"
	"github.com/passsquale/product-item-api/internal/model"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
)

func TestConsumer(t *testing.T) {
	cases := []struct {
		name           string
		batchSize      uint64
		correctDbLocks int
	}{
		{
			name:           "40",
			batchSize:      10,
			correctDbLocks: 4,
		},
	}

	for _, testcase := range cases {
		t.Run(testcase.name, func(t *testing.T) {
			eventsChannel := make(chan<- model.ItemEvent, 100)

			mockController := gomock.NewController(t)
			defer mockController.Finish()
			mockRepo := mocks.NewMockEventRepo(mockController)

			itemID := uint64(0)
			itemIDMutex := sync.Mutex{}

			gomock.InOrder(
				mockRepo.EXPECT().Lock(testcase.batchSize).Times(testcase.correctDbLocks).DoAndReturn(func(batchSize uint64) ([]model.ItemEvent, error) {
					itemIDMutex.Lock()
					output := make([]model.ItemEvent, 0, batchSize)
					for i := uint64(0); i < batchSize; i++ {
						e := model.ItemEvent{
							ID:     itemID,
							Type:   model.Created,
							Status: model.Processed,
							Item: &model.Item{
								ID: itemID,
							},
						}
						itemID++
						output = append(output, e)
					}
					itemIDMutex.Unlock()
					time.Sleep(100 * time.Millisecond)
					return output, nil
				}),
				mockRepo.EXPECT().Lock(testcase.batchSize).AnyTimes().Return([]model.ItemEvent{}, nil),
			)

			cosumerCfg := ConsumerConfig{
				ConsumeCount:    2,
				EventsChannel:   eventsChannel,
				Repo:            mockRepo,
				BatchSize:       testcase.batchSize,
				ConsumeInterval: 250 * time.Millisecond,
			}

			testingConsumer := NewDbConsumer(cosumerCfg)

			testingConsumer.Start()
			time.Sleep(2 * time.Second)

			assert.Equal(t, int(testcase.batchSize)*testcase.correctDbLocks, len(eventsChannel))

		})
	}
}
