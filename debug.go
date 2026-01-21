package tuak

import (
	"encoding/hex"
	"strings"
)

// DebugHook receives labeled copies of intermediate buffers.
type DebugHook func(label string, data []byte)

// DebugHexBytes formats bytes as space-separated lowercase hex (e.g. "01 02 03").
func DebugHexBytes(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(len(data)*3 - 1)
	for i := 0; i < len(data); i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		var buf [2]byte
		hex.Encode(buf[:], []byte{data[i]})
		b.Write(buf[:])
	}
	return b.String()
}
