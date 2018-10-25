package main

import (
	"encoding/binary"
	"fmt"
	"unicode/utf8"
)

// stringify ensures that we can print only valid characters.
// It's wrong to assume that everything is a string, since BoltDB is typeless.
func stringify(v []byte) string {
	if utf8.Valid(v) {
		ok := true
		for _, r := range string(v) {
			if r < 0x20 {
				ok = false
				break
			} else if r >= 0x7f && r <= 0x9f {
				ok = false
				break
			}
		}
		if ok {
			return string(v)
		}
	}
	if len(v) == 8 {
		return fmt.Sprintf("%v", binary.BigEndian.Uint64(v))
	}

	return fmt.Sprintf("%x", v)
}
