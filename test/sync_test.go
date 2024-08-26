package demo

import (
	"math/rand"
	"sync"
	"testing"
)

// sharing memory
var i int

// sync primitive
var mu sync.Mutex

// safe for concurrent use by multiple goroutines
var m2 sync.Map

var m1 = make(map[int]struct{})

var c = make(chan struct{}, 1)
var c1 = make(chan int)
var c2 = make(chan int)

func init() {
	go func() {
		for {
			mu.Lock()
			k := rand.Intn(100)
			m1[k] = struct{}{}
			mu.Unlock()
			c1 <- k
		}
	}()
	go func() {
		for {
			k := rand.Intn(100)
			m2.Store(k, struct{}{})
			c2 <- k
		}
	}()
}

func criticalSectionSyncByMutex() {
	mu.Lock()
	i++
	mu.Unlock()
}

func criticalSectionSyncByChannel() {
	c <- struct{}{}
	i++
	<-c
}

func BenchmarkCriticalSectionSyncByMutex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		criticalSectionSyncByMutex()
	}
}

func BenchmarkCriticalSectionSyncByMutexInParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			criticalSectionSyncByMutex()
		}
	})
}

func BenchmarkCriticalSectionSyncByChannel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		criticalSectionSyncByChannel()
	}
}

func BenchmarkCriticalSectionSyncByChannelInParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			criticalSectionSyncByChannel()
		}
	})
}

func recv1() {
	k := <-c1
	mu.Lock()
	_ = m1[k]
	delete(m1, k)
	mu.Unlock()
}

func recv2() {
	k := <-c2
	m2.LoadAndDelete(k)
}

func BenchmarkMapWithMutex(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			recv1()
		}
	})
}

func BenchmarkWithSyncMap(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			recv2()
		}
	})
}
