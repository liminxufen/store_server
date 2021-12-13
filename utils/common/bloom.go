package common

import (
	"math"
	"sync"

	"github.com/spaolacci/murmur3"
)

const (
	mod7       = 1<<3 - 1
	bitPerByte = 8
)

//bloom filter class
type BloomFilter interface {
	Put([]byte) error
	PutString(string) error
	Has([]byte) (bool, error)
	HasString(string) (bool, error)
	Reset()
}

//filter define
type Filter struct {
	m          uint64 //bit array of m bits
	n          uint64 //number of inserted elements
	k          uint64 //the number of hash function
	keys       []byte //byte array to store hash value
	lock       sync.RWMutex
	concurrent bool
}

func (f *Filter) Put(data []byte) error {
	h := HashData(data)
	for i := uint64(0); i < f.k; i++ {
		loc := Location(h, i)
		slot, mod := f.location(loc)
		f.keys[slot] |= 1 << mod
	}
	f.n++
	return nil
}

func (f *Filter) PutString(data string) error {
	d := String2ByteSlice(data)
	return f.Put(d)
}

func (f *Filter) Has(data []byte) (bool, error) {
	h := HashData(data)
	for i := uint64(0); i < f.k; i++ {
		loc := Location(h, i)
		slot, mod := f.location(loc)
		if f.keys[slot]&(1<<mod) == 0 {
			return false, nil
		}
	}
	return true, nil
}

func (f *Filter) HasString(data string) (bool, error) {
	d := String2ByteSlice(data)
	return f.Has(d)
}

func (f *Filter) location(h uint64) (uint64, uint64) {
	slot := (h / bitPerByte) & (f.m - 1)
	mod := h & mod7
	return slot, mod
}

func (f *Filter) Reset() {
	for i := 0; i < len(f.keys); i++ {
		f.keys[i] &= 0
	}
	f.n = 0
}

func NewFilter(size uint64, k uint64, race bool) BloomFilter {
	log2 := uint64(math.Ceil(math.Log2(float64(size))))
	filter := &Filter{
		m:          1 << log2,
		k:          k,
		keys:       make([]byte, 1<<log2),
		concurrent: race,
	}
	if filter.concurrent {
		filter.lock = sync.RWMutex{}
	}
	return filter
}

func HashData(data []byte) []uint64 {
	a1 := []byte{1} //to grab another bit of data
	hasher := murmur3.New128()
	hasher.Write(data)
	v1, v2 := hasher.Sum128()
	hasher.Write(a1)
	v3, v4 := hasher.Sum128()
	return []uint64{
		v1, v2, v3, v4,
	}
}

func Location(h []uint64, i uint64) uint64 {
	return h[i&1] + i*h[2+(((i+(i&1))&3)/2)]
}
