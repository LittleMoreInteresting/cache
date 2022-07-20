package singleflight

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func makeKey(tests []string, idx int) string {
	return fmt.Sprintf("key-%s", tests[idx])
}
func BenchmarkSingleFlightSetParallel(b *testing.B) {
	g := NewSingleFlight()
	tests := []string{"e", "a", "e", "a", "b", "c", "b", "a", "c", "d", "b", "c", "d"}
	rand.Seed(time.Now().Unix())
	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(len(tests))
		for pb.Next() {
			_, err := g.Do(makeKey(tests, id), func() (interface{}, error) {
				return "Val", nil
			})
			if err != nil {
				return
			}
		}
	})
}

func BenchmarkShardSingleFlightSetParallel(b *testing.B) {
	g := NewShardSingleFlight()
	tests := []string{"e", "a", "e", "a", "b", "c", "b", "a", "c", "d", "b", "c", "d"}
	rand.Seed(time.Now().Unix())
	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(len(tests))
		for pb.Next() {
			_, err := g.Do(makeKey(tests, id), func() (interface{}, error) {
				return "Val", nil
			})
			if err != nil {
				return
			}
		}
	})
}

// go test -bench=SetParallel -benchmem -count=10 -cpuprofile profile.out ./. | tee old.txt
//go tool pprof profile.out
