package common

import (
	"strings"
	"sync"
	"testing"
)

const iterCount = 10000
const coroutinesCount = 100

func Test_RndGen_GenString01(t *testing.T) {
	possible := "01"
	rndgen := NewRandomGenerator(possible)
	for iter := 0; iter < iterCount; iter++ {
		result := rndgen.NextString(rndgen.NextIntN(10) + 1)
		for index := range result {
			if strings.Index(possible, string(result[index])) == -1 {
				t.Fatalf("generated string containes impossible symbols")
			}
		}
	}
}

func Test_RndGen_GenStringABCD01(t *testing.T) {
	possible := "abcdABCD01"
	rndgen := NewRandomGenerator(possible)
	for iter := 0; iter < iterCount; iter++ {
		result := rndgen.NextString(rndgen.NextIntN(10) + 1)
		for index := range result {
			if strings.Index(possible, string(result[index])) == -1 {
				t.Fatalf("generated string containes impossible symbols")
			}
		}
	}
}

func Test_RndGen_GenStringParallel(t *testing.T) {
	possible := "abcdABCD01"
	rndgen := NewRandomGenerator(possible)
	var wg sync.WaitGroup
	for coroutine := 0; coroutine < coroutinesCount; coroutine++ {
		wg.Add(1)
		go func() {
			for iter := 0; iter < iterCount; iter++ {
				result := rndgen.NextString(rndgen.NextIntN(10) + 1)
				for index := range result {
					if strings.Index(possible, string(result[index])) == -1 {
						t.Fatalf("generated string containes impossible symbols")
					}
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
