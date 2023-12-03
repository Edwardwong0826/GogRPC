package service

import (
	"context"
	"errors"
	"fmt"
	wongProto "gogrpc/pb"
	"log"
	"sync"

	"github.com/jinzhu/copier"
)

var ErrAlreadyExists = errors.New("record already exists")

// this is an interface to store laptop
type LaptopStore interface {
	Save(laptop *wongProto.Laptop) error
	Find(id string) (*wongProto.Laptop, error)
	Search(ctx context.Context, filter *wongProto.Filter, found func(laptop *wongProto.Laptop) error) error
}

// stores laptop in db
// type DBLaptopStore struct{
// }

// stores laptop in memory
type InMemoryLaptopStore struct {
	// this read-write mutex is to handle concurrency
	mutex sync.RWMutex
	data  map[string]*wongProto.Laptop
}

// function to returns a new InMemoryLaptopStore
func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*wongProto.Laptop),
	}
}

func (store *InMemoryLaptopStore) Save(laptop *wongProto.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	// deep copy
	// go get "github.com/jinzhu/copier"
	other := &wongProto.Laptop{}
	err := copier.Copy(other, laptop)

	if err != nil {
		return fmt.Errorf("cannot copy laptop data: %v", err)
	}

	store.data[other.Id] = other
	return nil
}

// Find a laptop by Ud
func (store *InMemoryLaptopStore) Find(id string) (*wongProto.Laptop, error) {
	store.mutex.RLock() // require read lock
	defer store.mutex.RUnlock()

	laptop := store.data[id]
	if laptop == nil {
		return nil, nil
	}

	return deepCopy(laptop)
}

// Sarch a laptop
func (store *InMemoryLaptopStore) Search(ctx context.Context, filter *wongProto.Filter, found func(laptop *wongProto.Laptop) error) error {

	store.mutex.RLock() // require read lock
	defer store.mutex.RUnlock()

	for _, laptop := range store.data {

		// time.Sleep(time.Second)
		log.Print("checking laptop id: ", laptop.GetId())

		if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
			log.Print("context is cancelled")
			return nil
		}

		if isQualified(filter, laptop) {
			// deep copy
			other, err := deepCopy(laptop)
			if err != nil {
				return err
			}

			err = found(other)
			if err != nil {
				return err
			}
		}
	}

	return nil

}

func isQualified(filter *wongProto.Filter, laptop *wongProto.Laptop) bool {
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}

	if laptop.GetCpu().GetNumberCores() < filter.GetMinCpuCores() {
		return false
	}

	if laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz() {
		return false
	}

	if toBit(laptop.GetRam()) < toBit(filter.GetMinRam()) {
		return false
	}

	return true
}

func toBit(memory *wongProto.Memory) uint64 {
	value := memory.GetValue()

	switch memory.GetUnit() {
	case wongProto.Memory_BIT:
		return value
	case wongProto.Memory_BYTE:
		return value << 3 // 8 = 2^3
	case wongProto.Memory_KILOBYTE:
		return value << 13 // 1024 * 8 = 2^10 * 2^3 = 2^13
	case wongProto.Memory_MEGABYTE:
		return value << 23
	case wongProto.Memory_GIGABYTE:
		return value << 33
	case wongProto.Memory_TERABYTE:
		return value << 43
	default:
		return 0
	}
}

func deepCopy(laptop *wongProto.Laptop) (*wongProto.Laptop, error) {
	other := &wongProto.Laptop{}

	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}

	return other, nil
}
