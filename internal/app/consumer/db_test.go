package consumer

import (
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hablof/logistic-package-api/internal/mocks"
	"github.com/hablof/logistic-package-api/internal/model"
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
			eventsChannel := make(chan<- model.PackageEvent, 100)

			mockController := gomock.NewController(t)
			defer mockController.Finish()
			mockRepo := mocks.NewMockEventRepo(mockController)

			entityID := uint64(0)
			entityIDMutex := sync.Mutex{}

			gomock.InOrder(
				mockRepo.EXPECT().Lock(testcase.batchSize).Times(testcase.correctDbLocks).DoAndReturn(func(batchSize uint64) ([]model.PackageEvent, error) {
					entityIDMutex.Lock()
					output := make([]model.PackageEvent, 0, batchSize)
					for i := uint64(0); i < batchSize; i++ {
						e := model.PackageEvent{
							ID:      entityID,
							Type:    model.Created,
							Status:  model.Processed,
							Defered: 0,
							Entity: &model.Package{
								ID: entityID,
							},
						}
						entityID++
						output = append(output, e)
					}
					entityIDMutex.Unlock()
					time.Sleep(100 * time.Millisecond)
					return output, nil
				}),
				mockRepo.EXPECT().Lock(testcase.batchSize).AnyTimes().Return([]model.PackageEvent{}, nil),
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
