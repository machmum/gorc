package requestc

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"sync/atomic"

	gorc "github.com/machmum/gorc/string"
)

// Ported from Goji's middleware, source:
// https://github.com/zenazn/goji/tree/master/web/middlewar

// var prefix string
var reqid uint64

// RequestID is a middleware that injects a request ID into the context of each
// request. A request ID is a string of the form "host.example.com/random-0001",
// where "random" is a base62 random string that uniquely identifies this go
// process, and where the last number is an atomically incremented request
// counter.
func RequestID() string {
	var prefix string

	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}
	prefix = gorc.StringBuilder(hostname, b64[0:10])

	return fmt.Sprintf("%s-%06d", prefix, atomic.AddUint64(&reqid, 1))
}
