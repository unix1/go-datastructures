/*
Copyright 2014 Workiva, LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package augmentedtree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func constructMultiDimensionQueryTestTree() (
	*tree, Interval, Interval, Interval) {

	it := newTree(2)

	iv1 := constructMultiDimensionInterval(
		0, constructDimension(5, 10), constructDimension(5, 10),
	)
	it.Add(iv1)

	iv2 := constructMultiDimensionInterval(
		1, constructDimension(4, 5), constructDimension(4, 5),
	)
	it.Add(iv2)

	iv3 := constructMultiDimensionInterval(
		2, constructDimension(7, 12), constructDimension(7, 12),
	)
	it.Add(iv3)

	return it, iv1, iv2, iv3
}

func TestRootAddMultipleDimensions(t *testing.T) {
	it := newTree(2)
	iv := constructMultiDimensionInterval(
		1, constructDimension(0, 5), constructDimension(1, 6),
	)

	it.Add(iv)

	checkRedBlack(t, it.root, 1)
	result := it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(0, 10), constructDimension(0, 10),
		),
	)
	assert.Equal(t, Intervals{iv}, result)

	result = it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(100, 200), constructDimension(100, 200),
		),
	)
	assert.Len(t, result, 0)
}

func TestMultipleAddMultipleDimensions(t *testing.T) {
	it, iv1, iv2, iv3 := constructMultiDimensionQueryTestTree()

	checkRedBlack(t, it.root, 1)

	result := it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(0, 100), constructDimension(0, 100),
		),
	)
	assert.Equal(t, Intervals{iv2, iv1, iv3}, result)

	result = it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(3, 5), constructDimension(3, 5),
		),
	)
	assert.Equal(t, Intervals{iv2}, result)

	result = it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(5, 8), constructDimension(5, 8),
		),
	)
	assert.Equal(t, Intervals{iv1, iv3}, result)

	result = it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(11, 15), constructDimension(11, 15),
		),
	)
	assert.Equal(t, Intervals{iv3}, result)

	result = it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(15, 20), constructDimension(15, 20),
		),
	)
	assert.Len(t, result, 0)
}

func TestAddRebalanceInOrderMultiDimensions(t *testing.T) {
	it := newTree(2)

	for i := int64(0); i < 10; i++ {
		iv := constructMultiDimensionInterval(
			uint64(i), constructDimension(i, i+1), constructDimension(i, i+1),
		)
		it.Add(iv)
	}

	checkRedBlack(t, it.root, 1)
	result := it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(0, 10), constructDimension(0, 10),
		),
	)
	assert.Len(t, result, 10)
	assert.Equal(t, uint64(10), it.Len())
}

func TestAddRebalanceReverseOrderMultiDimensions(t *testing.T) {
	it := newTree(2)

	for i := int64(9); i >= 0; i-- {
		iv := constructMultiDimensionInterval(
			uint64(i), constructDimension(i, i+1), constructDimension(i, i+1),
		)
		it.Add(iv)
	}

	checkRedBlack(t, it.root, 1)
	result := it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(0, 10), constructDimension(0, 10),
		),
	)
	assert.Len(t, result, 10)
	assert.Equal(t, uint64(10), it.Len())
}

func TestAddRebalanceRandomOrderMultiDimensions(t *testing.T) {
	it := newTree(2)

	starts := []int64{0, 4, 2, 1, 3}

	for i, start := range starts {
		iv := constructMultiDimensionInterval(
			uint64(i), constructDimension(start, start+1), constructDimension(start, start+1),
		)
		it.Add(iv)
	}

	checkRedBlack(t, it.root, 1)
	result := it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(0, 10), constructDimension(0, 10),
		),
	)
	assert.Len(t, result, 5)
	assert.Equal(t, uint64(5), it.Len())
}

func TestAddLargeNumbersMultiDimensions(t *testing.T) {
	numItems := int64(1000)
	it := newTree(2)

	for i := int64(0); i < numItems; i++ {
		iv := constructMultiDimensionInterval(
			uint64(i), constructDimension(i, i+1), constructDimension(i, i+1),
		)
		it.Add(iv)
	}

	checkRedBlack(t, it.root, 1)
	result := it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(0, numItems), constructDimension(0, numItems),
		),
	)
	assert.Len(t, result, int(numItems))
	assert.Equal(t, uint64(numItems), it.Len())
}

func BenchmarkAddItemsMultiDimensions(b *testing.B) {
	numItems := int64(b.N)
	intervals := make(Intervals, 0, numItems)

	for i := int64(0); i < numItems; i++ {
		iv := constructMultiDimensionInterval(
			uint64(i), constructDimension(i, i+1), constructDimension(i, i+1),
		)
		intervals = append(intervals, iv)
	}

	it := newTree(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		it.Add(intervals[int64(i)%numItems])
	}
}

func BenchmarkQueryItemsMultiDimensions(b *testing.B) {
	numItems := int64(1000)
	intervals := make(Intervals, 0, numItems)

	for i := int64(0); i < numItems; i++ {
		iv := constructMultiDimensionInterval(
			uint64(i), constructDimension(i, i+1), constructDimension(i, i+1),
		)
		intervals = append(intervals, iv)
	}

	it := newTree(2)
	it.Add(intervals...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it.Query(
			constructMultiDimensionInterval(
				0, constructDimension(0, numItems), constructDimension(0, numItems),
			),
		)
	}
}

func TestRootDeleteMultiDimensions(t *testing.T) {
	it := newTree(2)
	iv := constructMultiDimensionInterval(
		0, constructDimension(5, 10), constructDimension(5, 10),
	)
	it.Add(iv)

	it.Delete(iv)

	checkRedBlack(t, it.root, 1)
	result := it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(0, 100), constructDimension(0, 100),
		),
	)
	assert.Len(t, result, 0)
	assert.Equal(t, uint64(0), it.Len())
}

func TestDeleteMultiDimensions(t *testing.T) {
	it, iv1, iv2, iv3 := constructMultiDimensionQueryTestTree()

	checkRedBlack(t, it.root, 1)

	it.Delete(iv1)

	result := it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(0, 100), constructDimension(0, 100),
		),
	)
	assert.Equal(t, Intervals{iv2, iv3}, result)

	result = it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(3, 5), constructDimension(3, 5),
		),
	)
	assert.Equal(t, Intervals{iv2}, result)

	result = it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(5, 8), constructDimension(5, 8),
		),
	)
	assert.Equal(t, Intervals{iv3}, result)

	result = it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(11, 15), constructDimension(11, 15),
		),
	)
	assert.Equal(t, Intervals{iv3}, result)

	result = it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(15, 20), constructDimension(15, 20),
		),
	)
	assert.Len(t, result, 0)
}

func TestDeleteRebalanceInOrderMultiDimensions(t *testing.T) {
	it := newTree(2)

	var toDelete *mockInterval

	for i := int64(0); i < 10; i++ {
		iv := constructMultiDimensionInterval(
			uint64(i), constructDimension(i, i+1), constructDimension(i, i+1),
		)
		it.Add(iv)
		if i == 5 {
			toDelete = iv
		}
	}

	it.Delete(toDelete)

	checkRedBlack(t, it.root, 1)
	result := it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(0, 10), constructDimension(0, 10),
		),
	)
	assert.Len(t, result, 9)
	assert.Equal(t, uint64(9), it.Len())
}

func TestDeleteRebalanceReverseOrderMultiDimensions(t *testing.T) {
	it := newTree(2)

	var toDelete *mockInterval

	for i := int64(9); i >= 0; i-- {
		iv := constructMultiDimensionInterval(
			uint64(i), constructDimension(i, i+1), constructDimension(i, i+1),
		)
		it.Add(iv)
		if i == 5 {
			toDelete = iv
		}
	}

	it.Delete(toDelete)

	checkRedBlack(t, it.root, 1)
	result := it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(0, 10), constructDimension(0, 10),
		),
	)
	assert.Len(t, result, 9)
	assert.Equal(t, uint64(9), it.Len())
}

func TestDeleteRebalanceRandomOrderMultiDimensions(t *testing.T) {
	it := newTree(2)

	starts := []int64{0, 4, 2, 1, 3}

	var toDelete *mockInterval

	for i, start := range starts {
		iv := constructMultiDimensionInterval(
			uint64(i), constructDimension(start, start+1), constructDimension(start, start+1),
		)
		it.Add(iv)
		if start == 1 {
			toDelete = iv
		}
	}

	it.Delete(toDelete)

	checkRedBlack(t, it.root, 1)
	result := it.Query(
		constructMultiDimensionInterval(
			0, constructDimension(0, 10), constructDimension(0, 10),
		),
	)
	assert.Len(t, result, 4)
	assert.Equal(t, uint64(4), it.Len())
}

func TestDeleteEmptyTreeMultiDimensions(t *testing.T) {
	it := newTree(2)

	it.Delete(
		constructMultiDimensionInterval(
			0, constructDimension(0, 10), constructDimension(0, 10),
		),
	)
	assert.Equal(t, uint64(0), it.Len())
}

func BenchmarkDeleteItemsMultiDimensions(b *testing.B) {
	numItems := int64(1000)
	intervals := make(Intervals, 0, numItems)

	for i := int64(0); i < numItems; i++ {
		iv := constructMultiDimensionInterval(
			uint64(i), constructDimension(i, i+1), constructDimension(i, i+1),
		)
		intervals = append(intervals, iv)
	}

	trees := make([]*tree, 0, b.N)
	for i := 0; i < b.N; i++ {
		it := newTree(2)
		it.Add(intervals...)
		trees = append(trees, it)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		trees[i].Delete(intervals...)
	}
}

func TestAddDeleteDuplicatesRebalanceInOrderMultiDimensions(t *testing.T) {
	it := newTree(2)

	intervals := make(Intervals, 0, 10)

	for i := 0; i < 10; i++ {
		iv := constructMultiDimensionInterval(
			uint64(i), constructDimension(0, 10), constructDimension(0, 10),
		)
		intervals = append(intervals, iv)
	}

	it.Add(intervals...)
	it.Delete(intervals...)
	assert.Equal(t, uint64(0), it.Len())
}

func TestAddDeleteDuplicatesRebalanceReverseOrderMultiDimensions(t *testing.T) {
	it := newTree(2)

	intervals := make(Intervals, 0, 10)

	for i := 9; i >= 0; i-- {
		iv := constructMultiDimensionInterval(
			uint64(i), constructDimension(0, 10), constructDimension(0, 10),
		)
		intervals = append(intervals, iv)
	}

	it.Add(intervals...)
	it.Delete(intervals...)
	assert.Equal(t, uint64(0), it.Len())
}

func TestAddDeleteDuplicatesRebalanceRandomOrderMultiDimensions(t *testing.T) {
	it := newTree(2)

	intervals := make(Intervals, 0, 5)
	starts := []int{0, 4, 2, 1, 3}

	for _, start := range starts {
		iv := constructMultiDimensionInterval(
			uint64(start), constructDimension(0, 10), constructDimension(0, 10),
		)
		intervals = append(intervals, iv)
	}

	it.Add(intervals...)
	it.Delete(intervals...)
	assert.Equal(t, uint64(0), it.Len())
}
