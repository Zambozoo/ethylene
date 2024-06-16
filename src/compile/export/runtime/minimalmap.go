package runtime

import (
	"slices"
	"sort"
)

type Serializable interface {
	Serialize() []byte
}

type UInt64 uint64

func (u UInt64) Serialize() []byte {
	return LongBytes(uint64(u))
}

type UInt64s []uint64

func (is UInt64s) Serialize() []byte {
	bytes := make([]byte, 0, len(is)*8)
	for _, i := range is {
		bytes = append(bytes, LongBytes(i)...)
	}

	if len(bytes)%8 != 0 {
		bytes = append(bytes, []byte{0, 0, 0, 0}...)
	}

	return bytes
}

func MinimalMapSize(numItems, itemByteSize uint64) uint64 {
	size := (itemByteSize+4)*numItems + 8
	if size%8 != 0 {
		return size + 4
	}

	return size
}

// based on : http://stevehanov.ca/blog/?id=119
type MinimalMap[T Serializable] struct {
	redirects []int
	values    []T
}

func (mm *MinimalMap[T]) Bytes() []byte {
	bytes := append(
		IntBytes(len(mm.redirects)),
		IntBytes(len(mm.values))...,
	)
	for _, i := range mm.redirects {
		bytes = append(bytes, IntBytes(i)...)
	}
	for _, v := range mm.values {
		bytes = append(bytes, v.Serialize()...)
	}

	return bytes
}

// isthe.com/chongo/tech/comp/fnv/
func fnv(d int, key uint64) int {
	if d == 0 {
		d = 0x01000193
	}

	for _, b := range LongBytes(key) {
		d = ((d * 0x01000193) ^ int(b)) & 0xffffffff
	}

	return d
}

func NewMinimalMap[T Serializable](m map[uint64]T, eq func(a T, b T) bool) *MinimalMap[T] {
	var zero T
	size := len(m)

	// Step 1: Place all of the keys into buckets
	buckets := make([][]uint64, len(m))
	mm := MinimalMap[T]{
		redirects: make([]int, size),
		values:    make([]T, size),
	}

	for key := range m {
		i := fnv(0, key) % size
		buckets[i] = append(buckets[i], key)
	}

	// Step 2: Sort the buckets and process the ones with the most items first.
	sort.Slice(buckets, func(i, j int) bool {
		return len(buckets[i]) > len(buckets[j])
	})

	b := 0
	for ; b < size; b++ {
		bucket := buckets[b]
		if len(bucket) <= 1 {
			break
		}

		d := 1
		item := 0
		slots := []int{}

		// Repeatedly try different values of d until we find a hash function
		// that places all items in the bucket into free slots
		for item < len(bucket) {
			slot := fnv(d, bucket[item]) % size
			if !eq(mm.values[slot], zero) || slices.Contains(slots, slot) {
				d += 1
				item = 0
				slots = []int{}
			} else {
				slots = append(slots, slot)
				item += 1
			}
		}

		mm.redirects[fnv(0, bucket[0])%size] = d
		for i := range bucket {
			mm.values[slots[i]] = m[bucket[i]]
		}
	}

	//Only buckets with 1 item remain. Process them more quickly by directly
	//placing them into a free slot. Use a negative value of d to indicate
	//this.
	freelist := []int{}
	for i := range mm.values {
		if eq(mm.values[i], zero) {
			freelist = append(freelist, i)
		}
	}

	for ; b < size; b++ {
		bucket := buckets[b]
		if len(bucket) == 0 {
			break
		}
		slot := freelist[len(freelist)-1]
		// We subtract one to ensure it's negative even if the zeroeth slot was
		// used.
		mm.redirects[fnv(0, bucket[0])%size] = -slot - 1
		mm.values[slot] = m[bucket[0]]
	}

	return &mm
}

func (mm *MinimalMap[T]) Get(key uint64) T {
	d := mm.redirects[fnv(0, key)%len(mm.redirects)]
	if d < 0 {
		return mm.values[len(mm.values)-int(d)-1]
	}

	return mm.values[int(fnv(d, key))%len(mm.values)]
}
