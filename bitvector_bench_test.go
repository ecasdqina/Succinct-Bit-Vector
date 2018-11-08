package bitvector

import (
	"math/rand"
	"testing"
)

const (
	bigSize = 1e6
)

func BenchmarkRank(b *testing.B) {
	_, bv := random(bigSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bv.Rank(rand.Intn(bigSize), itoB(rand.Intn(2)))
	}
	b.StopTimer()
}

func BenchmarkSelect(b *testing.B) {
	_, bv := random(bigSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bv.Select(rand.Intn(bigSize), itoB(rand.Intn(2)))
	}
	b.StopTimer()
}

func BenchmarkGet(b *testing.B) {
	_, bv := random(bigSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bv.Get(rand.Intn(bigSize))
	}
	b.StopTimer()
}

func itoB(i int) bool {
	return i != 0
}

func random(size int) (string, *BitVector) {
	var bs []byte
	b := NewBuilder(size)
	for i := 0; i < size; i++ {
		if rand.Intn(2) == 1 {
			bs = append(bs, '1')
			b.Set1(i)
		} else {
			bs = append(bs, '0')
		}
	}

	return string(bs), b.Build()
}
