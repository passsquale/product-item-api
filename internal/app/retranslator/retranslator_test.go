package restranslator

import (
	"fmt"
	"github.com/passsquale/product-item-api/internal/mocks"
	"github.com/passsquale/product-item-api/internal/model"
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestStart(t *testing.T) {

	t.Run("correct run and stop", func(t *testing.T) {
		repo, _, retranslator := setup(t, 2)

		repo.EXPECT().Lock(gomock.Any()).AnyTimes()
		retranslator.Start()
		retranslator.Close()
	})

	t.Run("correctly read all events and send", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			t.Run(fmt.Sprintf("attempt %d", i), func(t *testing.T) {

				batchSize := 32
				repo, sender, retranslator := setup(t, uint64(batchSize))

				controlChannel := make(chan struct{}, 1)

				eventsCount := 1000
				db := generate(eventsCount)

				offsetMutex := sync.Mutex{}
				offset := uint64(0)

				sendCount := int32(0)

				removedIDsMutex := sync.Mutex{}
				removedIDs := make([]uint64, 0, eventsCount)

				repo.EXPECT().Lock(gomock.Any()).DoAndReturn(func(size uint64) ([]model.ItemEvent, error) {
					offsetMutex.Lock()
					defer offsetMutex.Unlock()
					if offset >= uint64(len(db)) {
						return make([]model.ItemEvent, 0), nil
					}
					// chunk := db[offset : offset+size : offset+size]
					if offset+size >= uint64(len(db)) {
						chunk := db[offset:]
						offset += size
						return chunk, nil
					}
					chunk := db[offset : offset+size]
					offset += size
					return chunk, nil
				}).AnyTimes()

				repo.EXPECT().Remove(gomock.Any()).AnyTimes().Do(func(arr []uint64) {
					removedIDsMutex.Lock()
					removedIDs = append(removedIDs, arr...)
					if len(removedIDs) == eventsCount {
						controlChannel <- struct{}{}
					}
					removedIDsMutex.Unlock()
				})

				sender.EXPECT().Send(gomock.Any()).Times(eventsCount).Do(func(ptr *model.ItemEvent) {
					atomic.AddInt32(&sendCount, 1)
				}).Return(nil)

				retranslator.Start()
				go func() {
					time.Sleep(10 * time.Second)
					controlChannel <- struct{}{}
				}()
				<-controlChannel

				retranslator.Close()

				assert.Equal(t, int32(eventsCount), sendCount)
				assert.Equal(t, eventsCount, len(removedIDs), "len of removedIDs")
				assert.Equal(t, true, checkArr(removedIDs, uint64(eventsCount)))
			})
		}
	})
}

func generate(count int) []model.ItemEvent {
	result := make([]model.ItemEvent, 0, count)
	for i := 0; i < count; i++ {
		event := model.ItemEvent{
			ID:     uint64(i),
			Type:   model.Created,
			Status: 0,
			Item: &model.Item{
				ID:        uint64(i),
				OwnerID:   uint64(1),
				ProductID: uint64(2),
				Title:     "",
			},
		}
		result = append(result, event)
	}

	return result
}

func setup(t *testing.T, batchSize uint64) (*mocks.MockEventRepo, *mocks.MockEventSender, Retranslator) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	cfg := RetranslatorConfig{
		ChannelSize:     512,
		ConsumerCount:   16,
		BatchSize:       batchSize,
		ConsumeInterval: 100 * time.Millisecond,
		ProducerCount:   8,
		WorkerCount:     4,
		Repo:            repo,
		Sender:          sender,
	}

	retranslator := NewRetranslator(cfg)

	return repo, sender, retranslator
}

func checkArr(arr []uint64, n uint64) bool {
e:
	for i := uint64(0); i < n; i++ {
		for _, elem := range arr {
			if elem == i {
				continue e
			}
		}
		log.Printf("arr missed %d", i)
		return false
	}

	return true
}
