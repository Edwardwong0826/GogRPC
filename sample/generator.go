package sample

import (
	proto "gogrpc/pb"

	"github.com/golang/protobuf/ptypes"
)

func NewKeyboard() *proto.Keyboard {
	keyboard := &proto.Keyboard{
		Layout:  randomKeyboardLayout(),
		Backlit: randomBool(),
	}

	return keyboard
}

func NewCPU() *proto.CPU {
	brand := randomCPUBrand()
	name := randomCPUName(brand)

	numberCores := randomInt(2, 8)
	numberThreads := randomInt(numberCores, 12)

	minGhz := randomFloat64(2.0, 3.5)
	maxGhz := randomFloat64(minGhz, 5.0)

	cpu := &proto.CPU{
		Brand:         brand,
		Name:          name,
		NumberCores:   uint32(numberCores),
		NumberThreads: uint32(numberThreads),
		MinGhz:        minGhz,
		MaxGhz:        maxGhz,
	}

	return cpu
}

// NewGPU returns a new sample GPU
func NewGPU() *proto.GPU {
	brand := randomGPUBrand()
	name := randomGPUName(brand)

	minGhz := randomFloat64(1.0, 1.5)
	maxGhz := randomFloat64(minGhz, 2.0)
	memGB := randomInt(2, 6)

	gpu := &proto.GPU{
		Brand:  brand,
		Name:   name,
		MinGhz: minGhz,
		MaxGhz: maxGhz,
		Memory: &proto.Memory{
			Value: uint64(memGB),
			Unit:  proto.Memory_GIGABYTE,
		},
	}

	return gpu
}

// NewRAM returns a new sample RAM
func NewRAM() *proto.Memory {
	memGB := randomInt(4, 64)

	ram := &proto.Memory{
		Value: uint64(memGB),
		Unit:  proto.Memory_GIGABYTE,
	}

	return ram
}

// NewSSD returns a new sample SSD
func NewSSD() *proto.Storage {
	memGB := randomInt(128, 1024)

	ssd := &proto.Storage{
		Driver: proto.Storage_SSD,
		Memory: &proto.Memory{
			Value: uint64(memGB),
			Unit:  proto.Memory_GIGABYTE,
		},
	}

	return ssd
}

// NewHDD returns a new sample HDD
func NewHDD() *proto.Storage {
	memTB := randomInt(1, 6)

	hdd := &proto.Storage{
		Driver: proto.Storage_HDD,
		Memory: &proto.Memory{
			Value: uint64(memTB),
			Unit:  proto.Memory_TERABYTE,
		},
	}

	return hdd
}

// NewScreen returns a new sample Screen
func NewScreen() *proto.Screen {
	screen := &proto.Screen{
		SizeInch:   randomFloat32(13, 17),
		Resolution: randomScreenResolution(),
		Panel:      randomScreenPanel(),
		Multitouch: randomBool(),
	}

	return screen
}

// NewLaptop returns a new sample Laptop
func NewLaptop() *proto.Laptop {
	brand := randomLaptopBrand()
	name := randomLaptopName(brand)

	laptop := &proto.Laptop{
		Id:       randomID(),
		Brand:    brand,
		Name:     name,
		Cpu:      NewCPU(),
		Ram:      NewRAM(),
		Gpus:     []*proto.GPU{NewGPU()},
		Storages: []*proto.Storage{NewSSD(), NewHDD()},
		Screen:   NewScreen(),
		Keyboard: NewKeyboard(),
		Weight: &proto.Laptop_WeightKg{
			WeightKg: randomFloat64(1.0, 3.0),
		},
		PriceUsd:    randomFloat64(1500, 3500),
		ReleaseYear: uint32(randomInt(2015, 2019)),
		UpdatedAt:   ptypes.TimestampNow(),
	}

	return laptop
}

// RandomLaptopScore returns a random laptop score
func RandomLaptopScore() float64 {
	return float64(randomInt(1, 10))
}
