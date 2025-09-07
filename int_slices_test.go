package parameters

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUniqueUint64 test unique uint64
func TestUniqueUint64(t *testing.T) {
	one := []uint64{3, 2, 1}
	assert.True(t, reflect.DeepEqual(UniqueUint64(one), one))

	two := []uint64{3, 2, 1, 3, 3, 3, 3}

	assert.True(t, reflect.DeepEqual(UniqueUint64(two), one))
}

// ExampleUniqueUint64 shows an example using the method
func ExampleUniqueUint64() {
	one := []uint64{3, 2, 1, 3, 3, 3, 3}
	unique := UniqueUint64(one)

	fmt.Println(unique)
	// Output: [3 2 1]
}

// BenchmarkUniqueUint64 benchmarks the method
func BenchmarkUniqueUint64(b *testing.B) {
	one := []uint64{3, 2, 1, 3, 3, 3, 3}
	for i := 0; i < b.N; i++ {
		_ = UniqueUint64(one)
	}
}
