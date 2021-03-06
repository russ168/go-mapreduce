// Package mapreduce provides a simple abstraction for the general Map/Reduce
// pattern.
package mapreduce

import (
	"sync"
)

// In order to utilize this package you must create a struct that implements
// the following interface.
type MapReduce interface {
	Map(in chan interface{}, out chan interface{})
	Reduce(in chan interface{}) interface{}
}

// Configuration used by the Map Reducer.
type Configuration struct {
	MapperCount int
	InChan      chan interface{}
	OutChan     chan interface{}
}

// NewMapReduceConfig returns a MapReduce Configuration struct with sensible
// defaults.
func NewMapReduceConfig() *Configuration {
	inChan := make(chan interface{})
	outChan := make(chan interface{})

	return &Configuration{
		MapperCount: 1,
		InChan:      inChan,
		OutChan:     outChan,
	}
}

// Run executes the MapReduce process.
func Run(mr MapReduce, c *Configuration) (interface{}, error) {

	var wg sync.WaitGroup

	// Map
	wg.Add(c.MapperCount)
	for i := 0; i < c.MapperCount; i++ {
		go func(wg *sync.WaitGroup) {
			mr.Map(c.InChan, c.OutChan)
			wg.Done()
		}(&wg)
	}

	go func(w *sync.WaitGroup) {
		w.Wait()
		close(c.OutChan)
	}(&wg)

	// Reduce
	res := mr.Reduce(c.OutChan)
	return res, nil
}
