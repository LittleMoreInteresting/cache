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
	tests := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
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
func BenchmarkShardSingleFlightBadSetParallel(b *testing.B) {
	g := NewShardSingleFlight()
	tests := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	rand.Seed(time.Now().Unix())
	b.RunParallel(func(pb *testing.PB) {
		_ = rand.Intn(len(tests))
		for pb.Next() {
			_, err := g.Do(makeKey(tests, 0), func() (interface{}, error) {
				return "Val", nil
			})
			if err != nil {
				return
			}
		}
	})
}

// go test -bench=SetParallel -benchmem -count=10 -cpuprofile profile.out ./. | tee old.txt
// benchstat old.txt
//go tool pprof profile.out
