//go:build ignore

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const TotalSpots = 3

// Garage implements a parking garage with a semaphore (spots) to limit capacity
// and a mutex (carsMu) to protect the parkedCars slice (critical section).
type Garage struct {
	spots      chan struct{} // buffered channel used a semaphore
	carsMu     sync.Mutex
	parkedCars []string
}

func NewGarage(totalSpots int) *Garage {
	return &Garage{
		spots:      make(chan struct{}, totalSpots),
		parkedCars: make([]string, 0),
	}
}

func (g *Garage) CountParkedCars() int {
	g.carsMu.Lock()
	defer g.carsMu.Unlock()
	return len(g.parkedCars)
}

// Enter acquires a spot permit (may block if full),
// then enters the critical section to update shared state.
func (g *Garage) Enter(carName string) {
	g.spots <- struct{}{} // acquire one permit

	g.carsMu.Lock()

	g.parkedCars = append(g.parkedCars, carName)

	fmt.Printf("%s parked\n", carName)
	g.carsMu.Unlock()
}

// Exit removes the car from shared state, then releases a spot permit.
func (g *Garage) Exit(carName string) {
	g.carsMu.Lock()

	for i, name := range g.parkedCars {
		if name == carName {
			last := len(g.parkedCars) - 1
			g.parkedCars[i] = g.parkedCars[last]
			g.parkedCars = g.parkedCars[:last]
			break
		}
	}

	fmt.Printf("%s leaving\n", carName)
	g.carsMu.Unlock()

	<-g.spots // release one permit
}

func parkCar(garage *Garage, carName string) {
	garage.Enter(carName)

	sleep := time.Duration(1000+rand.Intn(2000)) * time.Millisecond
	time.Sleep(sleep) // simulate parking duration

	garage.Exit(carName)
}

func testGarage(garage *Garage, numberOfCars int) {
	var wg sync.WaitGroup
	wg.Add(numberOfCars)

	for carNum := range numberOfCars {
		carName := fmt.Sprintf("Car-%d", carNum)
		go func(name string) {
			defer wg.Done()
			parkCar(garage, name)
		}(carName)
	}

	wg.Wait()
}

func main() {
	numberOfCars := 10
	garage := NewGarage(TotalSpots)
	// test garage by concurrently arriving cars
	testGarage(garage, numberOfCars)

	fmt.Println("Number of parked cars after a busy day:")
	fmt.Printf("Actual: %d\nExpected: 0\n", garage.CountParkedCars())
}
