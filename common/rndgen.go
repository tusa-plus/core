package common

import (
	"math/rand"
	"sync"
	"time"
)

type RandomGenerator struct {
	symbols string
	lock    *sync.Mutex
	gen     *rand.Rand
}

func NewRandomGenerator(symbols string) *RandomGenerator {
	source := rand.NewSource(time.Now().UnixNano())
	return &RandomGenerator{
		symbols: symbols,
		lock:    &sync.Mutex{},
		gen:     rand.New(source),
	}
}

func (rnd *RandomGenerator) NextString(cnt int) string {
	result := make([]byte, cnt)
	rnd.lock.Lock()
	defer rnd.lock.Unlock()
	for index := range result {
		result[index] = rnd.symbols[rnd.gen.Intn(len(rnd.symbols))]
	}
	return string(result)
}

func (rnd *RandomGenerator) NextIntN(n int) int {
	rnd.lock.Lock()
	defer rnd.lock.Unlock()
	return rnd.gen.Intn(n)
}
