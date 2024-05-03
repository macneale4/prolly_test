package main

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProllyBinSearchUneven(t *testing.T) {
	// We construct a prefix list which is not well distributed to ensure that the search still works, even if not
	// optimal.
	pf := make([]int, 1000)
	for i := 0; i < 900; i++ {
		pf[i] = i
	}
	target := 12345
	pf[900] = target
	for i := 901; i < 1000; i++ {
		pf[i] = 10000000 + i
	}
	// In normal circumstances, a value of 12345 would be far to the left side of the list
	found := prollyBinSearch(pf, target)
	assert.Equal(t, 900, found)

	// Same test, but from something on the right side of the list.
	for i := 999; i > 100; i-- {
		pf[i] = math.MaxInt - i
	}
	target = math.MaxInt - 12345
	pf[100] = target
	for i := 99; i >= 0; i-- {
		pf[i] = 10000000 - i
	}
	found = prollyBinSearch(pf, target)
	assert.Equal(t, 100, found)
}

func TestProllyBinSearch(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	curVal := r.Int()
	pf := make([]int, 10000)
	for i := 0; i < 10000; i++ {
		pf[i] = curVal
		curVal += r.Intn(10)
	}

	for i := 0; i < 10000; i++ {
		idx := prollyBinSearch(pf, pf[i])
		// There are dupes in the list, so we don't always end up with the same index.
		assert.Equal(t, pf[i], pf[idx])
	}

	idx := prollyBinSearch(pf, pf[0]-1)
	assert.Equal(t, -1, idx)
	idx = prollyBinSearch(pf, pf[9999]+1)
	assert.Equal(t, -1, idx)

	// 23 is not a dupe, and neighbors don't match. stable due to seed.
	idx = prollyBinSearch(pf, pf[23]+1)
	assert.Equal(t, -1, idx)
	idx = prollyBinSearch(pf, pf[23]-1)
	assert.Equal(t, -1, idx)
}

func TestAaronSearch(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	curVal := r.Int()
	pf := make([]int, 10000)
	for i := 0; i < 10000; i++ {
		pf[i] = curVal
		curVal += r.Intn(10)
	}

	for i := 0; i < 10000; i++ {
		idx := aaronSearch(pf, pf[i])
		// There are dupes in the list, so we don't always end up with the same index.
		assert.Equal(t, pf[i], pf[idx])
	}

	idx := aaronSearch(pf, pf[0]-1)
	assert.Equal(t, 0, idx)
	idx = aaronSearch(pf, pf[9999]+1)
	assert.Equal(t, 10000, idx)

	// 23 is not a dupe, and neighbors don't match. stable due to seed.
	idx = aaronSearch(pf, pf[23]+1)
	assert.Equal(t, 10000, idx)
	idx = aaronSearch(pf, pf[23]-1)
	assert.Equal(t, 10000, idx)
}
