package ppln

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestIDTrieSanity(t *testing.T) {
	tr := NewIDTrie(1)

	pipeId := uint64(1)

	tr.Insert(pipeId, Value{Lineage: []uint64{1, 2, 3}, Value: "hello1"})
	tr.Insert(pipeId, Value{Lineage: []uint64{1, 2, 4}, Value: "hello2"})

	v, ok := tr.Get(pipeId, []uint64{1, 2, 3}, 0)
	assert.True(t, ok)
	assert.Equal(t, Value{Lineage: []uint64{1, 2, 3}, Value: "hello1"}, v)

	v, ok = tr.Get(pipeId, []uint64{1, 2, 4}, 0)
	assert.True(t, ok)
	assert.Equal(t, Value{Lineage: []uint64{1, 2, 4}, Value: "hello2"}, v)

	values := make([]Value, 0)
	tr.Walk(func(value Value) {
		values = append(values, value)
	})
	assert.EqualValues(t, []Value{
		{Lineage: []uint64{1, 2, 3}, Value: "hello1"},
		{Lineage: []uint64{1, 2, 4}, Value: "hello2"},
	}, values)

	{
		tr, ok := tr.GetTrie(pipeId, nil)
		assert.True(t, ok)

		values = make([]Value, 0)
		tr.Walk(func(value Value) {
			values = append(values, value)
		})
		assert.EqualValues(t, []Value{
			{Lineage: []uint64{1, 2, 3}, Value: "hello1"},
			{Lineage: []uint64{1, 2, 4}, Value: "hello2"},
		}, values)
	}

	tr.Remove(pipeId, []uint64{1, 2, 3}, 0)

	values = make([]Value, 0)
	tr.Walk(func(value Value) {
		values = append(values, value)
	})
	assert.EqualValues(t, []Value{
		{Lineage: []uint64{1, 2, 4}, Value: "hello2"},
	}, values)
}

func TestIDTrieSanityLarge(t *testing.T) {
	tr := NewIDTrie(1)

	pipeId := uint64(1)

	tr.Insert(pipeId, Value{Lineage: []uint64{1, numBuckets, numBuckets * 2}, Value: "hello"})

	v, ok := tr.Get(pipeId, []uint64{1, numBuckets, numBuckets * 2}, 0)

	assert.True(t, ok)
	assert.Equal(t, Value{Lineage: []uint64{1, numBuckets, numBuckets * 2}, Value: "hello"}, v)
}

func BenchmarkIDTrieInsert3_1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		benchmarkIDTrieInsert(b, 3, 1)
	}
}
func BenchmarkIDTrieInsert3_10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		benchmarkIDTrieInsert(b, 3, 10)
	}
}
func BenchmarkIDTrieInsert3_1000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		benchmarkIDTrieInsert(b, 3, 1000)
	}
}

func BenchmarkIDTrieInsert1000_1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		benchmarkIDTrieInsert(b, 1000, 1)
	}
}
func BenchmarkIDTrieInsert1000_10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		benchmarkIDTrieInsert(b, 1000, 10)
	}
}
func BenchmarkIDTrieInsert1000_1000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		benchmarkIDTrieInsert(b, 1000, 1000)
	}
}

func benchmarkIDTrieInsert(b *testing.B, levels, n int) {
	b.Helper()
	b.ReportAllocs()

	ls := make([]LineageID, 0, n)
	for i := 0; i < n; i++ {
		l := make(LineageID, levels)
		for j := 0; j < levels; j++ {
			l[j] = rand.Uint64()
		}
		ls = append(ls, l)
	}

	// this makes it hang for whatever reason
	//b.ResetTimer()

	tr := NewIDTrie(1)

	for i := 0; i < n; i++ {
		tr.Insert(0, Value{Lineage: ls[i], Value: "hello"})
	}
}
