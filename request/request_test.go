package gorc_test

import (
	"testing"

	"github.com/machmum/gorc"
)

func BenchmarkRequestID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		gorc.RequestID()
	}
}
