package main

import (
	"fmt"
	"hash"
	"math/rand"

	"github.com/google/uuid"
	"github.com/spaolacci/murmur3"
)

// https://en.wikipedia.org/wiki/MurmurHash
var hashFunctions []hash.Hash32

func init() {
	// hashFunctions = []hash.Hash32{
	// 	murmur3.New32WithSeed(uint32(11)),
	// 	murmur3.New32WithSeed(uint32(4214315325)),
	// 	murmur3.New32WithSeed(uint32(442343243)),
	// 	murmur3.New32WithSeed(uint32(87988798)),
	// 	murmur3.New32WithSeed(uint32(768797855)),
	// 	murmur3.New32WithSeed(uint32(465765768)),
	// 	murmur3.New32WithSeed(uint32(889807453)),
	// 	murmur3.New32WithSeed(uint32(456363)),
	// 	murmur3.New32WithSeed(uint32(44353)),
	// }

	hashFunctions = make([]hash.Hash32, 0)
	for i := 0; i < 100; i++ {
		hashFunctions = append(hashFunctions, murmur3.New32WithSeed(rand.Uint32()))
	}
}

type BloomFilter struct {
	filter []uint8
	size   int32
}

func murmurhash(key string, size int32, hashFunctionsIdx int) int32 {
	hashFunctions[hashFunctionsIdx].Write([]byte(key))
	result := hashFunctions[hashFunctionsIdx].Sum32() % uint32(size)
	hashFunctions[hashFunctionsIdx].Reset()
	return int32(result)
}

func NewBloomFilter(size int32) *BloomFilter {
	return &BloomFilter{
		filter: make([]uint8, size),
		size:   size,
	}
}

func (b *BloomFilter) Add(key string, numOfHashFunctions int) {
	for i := 0; i <= numOfHashFunctions; i++ {
		idx := murmurhash(key, b.size, numOfHashFunctions)
		// xth byte - idx / 8
		byteIndex := idx / 8
		// yth bit in that byte - idx % 8
		bitPosition := idx % 8
		b.filter[byteIndex] = b.filter[byteIndex] | (1 << bitPosition)
	}
}

func (b *BloomFilter) Exists(key string, numOfFunctions int) bool {

	for i := 0; i < numOfFunctions; i++ {
		idx := murmurhash(key, b.size, i)
		// xth byte - idx / 8
		byteIndex := idx / 8
		// yth bit in that byte - idx % 8
		bitPosition := idx % 8

		exist := b.filter[byteIndex]&(1<<bitPosition) == 1
		if !exist {
			return false
		}

	}

	return true
}

func (b *BloomFilter) Print() {
	fmt.Print(b.filter)
}

func main() {
	dataset := make([]string, 0)
	datasetExists := make(map[string]bool)
	datasetNotExist := make(map[string]bool)

	numOfElements := 8000

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

	for i := 0; i < len(hashFunctions); i++ {
		bloom := NewBloomFilter(int32(10_00))

		for key, _ := range datasetExists {
			bloom.Add(key, i)
		}

		falsePositive := 0

		for _, key := range dataset {
			exists := bloom.Exists(key, i)
			_, ok := datasetNotExist[key]

			if exists && ok {
				falsePositive++
			}

		}

		fmt.Println(float64(falsePositive) / float64(len(dataset)))
	}

}
