package parameters

import (
	"fmt"
	"reflect"
	"testing"
)

// TestUniqueUint64 test unique uint64
func TestUniqueUint64(t *testing.T) {
	one := []uint64{3, 2, 1}
	if !reflect.DeepEqual(UniqueUint64(one), one) {
		t.Errorf("slice with no dupes is different")
	}

	two := []uint64{3, 2, 1, 3, 3, 3, 3}

	if !reflect.DeepEqual(UniqueUint64(two), one) {
		t.Errorf("slice with dupes is not what we expected")
	}
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
