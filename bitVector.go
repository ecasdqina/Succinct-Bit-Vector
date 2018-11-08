package bitVector

import "errors"

const (
	bitLength = 64
	maskFF    = uint64(0xffffffffffffffff)
	mask55    = uint64(0x5555555555555555)
	mask33    = uint64(0x3333333333333333)
	mask0F    = uint64(0x0f0f0f0f0f0f0f0f)
	mask01    = uint64(0x0101010101010101)
)

var (
	// ErrorOutOfRange indicates out of range access.
	ErrorOutOfRange = errors.New("Out of range access")
)

type BitVector struct {
	size int
	rank []int    // the vector of the number of 1s in the bit vector pers BitLength.
	v    []uint64 // the bit vector
}

// Len returns the size of the bit vector.
func (b BitVector) Len() int {
	return b.size
}

// Get returns true or false, the value of the i-th bit in the bit vector.
func (b BitVector) Get(i int) (bool, error) {
	if i > b.size {
		return false, ErrorOutOfRange
	}
	return ((b.v[i/bitLength] >> uint(i%bitLength)) & 1) == 1, nil
}

// Rank returns the count of 1s or 0s before the i-th bit.
func (b BitVector) Rank(i int, x bool) (int, error) {
	if x {
		val, err := b.Rank1(i)
		if err != nil {
			return 0, err
		}
		return val, nil
	} else {
		val, err := b.Rank0(i)
		if err != nil {
			return 0, err
		}
		return val, nil
	}
}

// Rank1 returns the count of 1s before the i-th bit.
func (b BitVector) Rank1(i int) (int, error) {
	if i > b.size {
		return 0, ErrorOutOfRange
	}
	offset := uint(i % bitLength)
	return int(b.rank[i/bitLength] + popcount(b.v[i/bitLength] & ^(maskFF<<offset))), nil
}

// Rank0 return the count of 0s before the i-th bit.
func (b BitVector) Rank0(i int) (int, error) {
	val, err := b.Rank1(i)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// Builder is a builder of BitVector.
type Builder struct {
	size int
	v    []uint64
}

// NewBuilder makes a new builder of BitVector of the specified size.
func NewBuilder(size int) *Builder {
	bufsize := size/bitLength + 1

	return &Builder{
		size: size,
		v:    make([]uint64, bufsize),
	}
}

// Len returns the size of the bit vector.
func (b Builder) Len() int {
	return b.size
}

// Set sets i-th bit in the bit vector to v.
func (b *Builder) Set(i int, v bool) {
	if v {
		b.v[i/bitLength] |= uint64(1) << uint(i%bitLength)
	} else {
		b.v[i/bitLength] &^= (uint64(1) << uint(i%bitLength))
	}
}

// Set1 sets i-th bit in the bit vector to 1.
func (b *Builder) Set1(i int) {
	b.Set(i, true)
}

// Set0 sets i-th bit in the bit vector to 0.
func (b *Builder) Set0(i int) {
	b.Set(i, false)
}

// Get returns true or false, i-th bit in the bit vector.
func (b Builder) Get(i int) bool {
	return (b.v[i/bitLength] << uint(i%bitLength) & 1) == 1
}

// Build builds a BitVector from the builder.
func (b Builder) Build() *BitVector {
	rank := make([]int, len(b.v))
	count := 0

	for i, x := range b.v {
		rank[i] = count
		count += popcount(x)
	}

	return &BitVector{
		size: b.size,
		v:    b.v,
		rank: rank,
	}
}

func popcount(x uint64) int {
	x = (x & mask55) + (x >> 1 & mask55)
	x = (x & mask33) + (x >> 2 & mask33)
	x = (x + (x >> 4)) & mask0F
	return int(x * mask01 >> 56 & uint64(0x7f))
}
