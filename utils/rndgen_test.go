package utils

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

const iterCount = 10000
const coroutinesCount = 100

func Test_RndGen_GenString01(t *testing.T) {
	t.Parallel()
	possible := "01"
	rndgen := NewRandomGenerator(possible)
	for iter := 0; iter < iterCount; iter++ {
		result := rndgen.NextString(rndgen.NextIntN(10) + 1)
		for index := range result {
			if !strings.Contains(possible, string(result[index])) {
				t.Fatalf("generated string contains impossible symbols")
			}
		}
	}
}

func Test_RndGen_GenStringABCD01(t *testing.T) {
	t.Parallel()
	possible := "abcdABCD01"
	rndgen := NewRandomGenerator(possible)
	for iter := 0; iter < iterCount; iter++ {
		result := rndgen.NextString(rndgen.NextIntN(10) + 1)
		for index := range result {
			if !strings.Contains(possible, string(result[index])) {
				t.Fatalf("generated string contains impossible symbols")
			}
		}
	}
}

func Test_RndGen_GenStringParallel(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip()
	}
	possible := "abcdABCD012345"
	rndgen := NewRandomGenerator(possible)
	var wg sync.WaitGroup
	var errMutex sync.Mutex
	var err error = nil
	for coroutine := 0; coroutine < coroutinesCount; coroutine++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for iter := 0; iter < iterCount; iter++ {
				result := rndgen.NextString(rndgen.NextIntN(10) + 1)
				for index := range result {
					if !strings.Contains(possible, string(result[index])) {
						errMutex.Lock()
						err = fmt.Errorf("generated string contains impossible symbols")
						errMutex.Unlock()
						return
					}
				}
			}
		}()
	}
	wg.Wait()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
