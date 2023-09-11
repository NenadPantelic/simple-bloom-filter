package main

import (
	"fmt"
	"hash"
	"time"

	"github.com/google/uuid"
	"github.com/spaolacci/murmur3"
)

// https://en.wikipedia.org/wiki/MurmurHash
var murmurHasher hash.Hash32

func init() {
	murmurHasher = murmur3.New32WithSeed(uint32(time.Now().Unix()))
}

type BloomFilter struct {
	filter []bool
	size   int32
}

func murmurhash(key string, size int32) int32 {
	murmurHasher.Write([]byte(key))
	result := murmurHasher.Sum32() % uint32(size)
	murmurHasher.Reset()
	return int32(result)
}

func NewBloomFilter(size int32) *BloomFilter {
	return &BloomFilter{
		filter: make([]bool, size),
		size:   size,
	}
}

func (b *BloomFilter) Add(key string) {
	idx := murmurhash(key, b.size)
	b.filter[idx] = true
}

func (b *BloomFilter) Exists(key string) bool {
	return b.filter[murmurhash(key, b.size)]
}

func (b *BloomFilter) Print() {
	fmt.Print(b.filter)
}

func main() {
	dataset := make([]string, 0)
	datasetExists := make(map[string]bool)
	datasetNotExist := make(map[string]bool)

	numOfElements := 5000

	for i := 0; i < numOfElements/2; i++ {
		u := uuid.New().String()
		dataset = append(dataset, u)
		datasetExists[u] = true
	}

	for i := 0; i < numOfElements/2; i++ {
		u := uuid.New().String()
		dataset = append(dataset, u)
		datasetNotExist[u] = true
	}

	for j := 100; j <= 50_000; j += 100 {
		bloom := NewBloomFilter(int32(j))

		for key, _ := range datasetExists {
			bloom.Add(key)
		}

		falsePositive := 0

		for _, key := range dataset {
			exists := bloom.Exists(key)
			_, ok := datasetNotExist[key]

			if exists && ok {
				falsePositive++
			}

		}

		fmt.Println(float64(falsePositive) / float64(len(dataset)))
	}

}
