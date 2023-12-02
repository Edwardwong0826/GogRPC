package service

import (
	"errors"
	"fmt"
	wongProto "gogrpc/pb"
	"sync"

	"github.com/jinzhu/copier"
)

var ErrAlreadyExists = errors.New("record already exists")

// this is an interface to store laptop
type LaptopStore interface {
	Save(laptop *wongProto.Laptop) error
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
