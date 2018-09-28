package trackers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtilsStringSliceContains(t *testing.T) {
	var sliceA = []string{"a", "b", "c"}
	var contains bool

	contains = stringSliceContains(sliceA, "b")
	assert.True(t, contains)

	contains = stringSliceContains(sliceA, "d")
	assert.False(t, contains)
}

func TestUtilsStringSliceRemove(t *testing.T) {
	var sliceA = []string{"a", "b", "c"}
	var sliceB = stringSliceRemove(sliceA, "c")
	var contains bool

	contains = stringSliceContains(sliceB, "c")
	assert.False(t, contains)
}

func TestUtilsIntSliceContains(t *testing.T) {
	var sliceA = []int{1, 2, 3}
	var contains bool

	contains = intSliceContains(sliceA, 1)
	assert.True(t, contains)

	contains = intSliceContains(sliceA, 4)
	assert.False(t, contains)
}

func TestUtilsIntSliceEq(t *testing.T) {
	var sliceA = []int{1, 2, 3}
	var sliceB = []int{1, 2, 3}
	var equals bool

	equals = intSliceEq(sliceA, sliceB)
	assert.True(t, equals)

	equals = intSliceEq(sliceA, []int{1})
	assert.False(t, equals)
}

func TestUtilsStringSliceToInt(t *testing.T) {
	var sliceA = []string{"1", "2", "3"}
	var sliceB = []int{1, 2, 3}
	var slice []int
	var err error

	slice, err = stringSliceToInt(sliceA)
	assert.Nil(t, err)
	assert.Equal(t, sliceB, slice)
}
