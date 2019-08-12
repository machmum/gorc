package request

import (
	"testing"
)

func BenchmarkRequestID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RequestID()
	}
}
