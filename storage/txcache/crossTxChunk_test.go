package txcache

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCrossTxChunk_ImmunizeKeys(t *testing.T) {
	chunk := newUnconstrainedCrossTxChunkToTest()

	chunk.addTestTxs("x", "y", "z")
	require.Equal(t, 3, chunk.CountItems())

	// No immune items, all removed
	numRemoved := chunk.RemoveOldest(42)
	require.Equal(t, 3, numRemoved)
	require.Equal(t, 0, chunk.CountItems())

	chunk.addTestTxs("x", "y", "z")
	require.Equal(t, 3, chunk.CountItems())

	// Immunize some items
	numNow, numFuture := chunk.ImmunizeKeys(hashesAsBytes([]string{"x", "z"}))
	require.Equal(t, 2, numNow)
	require.Equal(t, 0, numFuture)

	numRemoved = chunk.RemoveOldest(42)
	require.Equal(t, 1, numRemoved)
	require.Equal(t, 2, chunk.CountItems())
	require.Equal(t, []string{"x", "z"}, hashesAsStrings(chunk.KeysInOrder()))
}

func TestCrossTxChunk_AddItemIgnoresDuplicates(t *testing.T) {
	chunk := newUnconstrainedCrossTxChunkToTest()
	chunk.addTestTxs("x", "y", "z")
	require.Equal(t, 3, chunk.CountItems())

	ok, added := chunk.addTestTx("a")
	require.True(t, ok)
	require.True(t, added)
	require.Equal(t, 4, chunk.CountItems())

	ok, added = chunk.addTestTx("x")
	require.True(t, ok)
	require.False(t, added)
	require.Equal(t, 4, chunk.CountItems())
}

func TestCrossTxChunk_AddItemEvictsWhenTooMany(t *testing.T) {
	chunk := newCrossTxChunkToTest(3, math.MaxUint32)
	chunk.addTestTxs("x", "y", "z")
	require.Equal(t, 3, chunk.CountItems())

	chunk.addTestTxs("a", "b")
	require.Equal(t, []string{"z", "a", "b"}, hashesAsStrings(chunk.KeysInOrder()))
}

func TestCrossTxChunk_AddItemDoesNotEvictImmuneItems(t *testing.T) {
	chunk := newCrossTxChunkToTest(3, math.MaxUint32)
	chunk.addTestTxs("x", "y", "z")
	require.Equal(t, 3, chunk.CountItems())

	_, _ = chunk.ImmunizeKeys(hashesAsBytes([]string{"x", "y"}))

	chunk.addTestTxs("a")
	require.Equal(t, []string{"x", "y", "a"}, hashesAsStrings(chunk.KeysInOrder()))
	chunk.addTestTxs("b")
	require.Equal(t, []string{"x", "y", "b"}, hashesAsStrings(chunk.KeysInOrder()))

	_, _ = chunk.ImmunizeKeys(hashesAsBytes([]string{"b"}))
	ok, added := chunk.addTestTx("c")
	require.False(t, ok)
	require.False(t, added)
	require.Equal(t, []string{"x", "y", "b"}, hashesAsStrings(chunk.KeysInOrder()))
}

func newUnconstrainedCrossTxChunkToTest() *crossTxChunk {
	chunk := newCrossTxChunk(crossTxChunkConfig{
		maxNumItems:                 math.MaxUint32,
		maxNumBytes:                 math.MaxUint32,
		numItemsToPreemptivelyEvict: math.MaxUint32,
	})

	return chunk
}

func newCrossTxChunkToTest(maxNumItems uint32, numMaxBytes uint32) *crossTxChunk {
	chunk := newCrossTxChunk(crossTxChunkConfig{
		maxNumItems:                 maxNumItems,
		maxNumBytes:                 numMaxBytes,
		numItemsToPreemptivelyEvict: 1,
	})

	return chunk
}

func (chunk *crossTxChunk) addTestTxs(hashes ...string) {
	for _, hash := range hashes {
		_, _ = chunk.addTestTx(hash)
	}
}

func (chunk *crossTxChunk) addTestTx(hash string) (ok, added bool) {
	return chunk.addItem(createTx([]byte(hash), ".", uint64(42)))
}
