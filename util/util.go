package util

import (
	"fmt"
	"hash"
)

func HexHash(hash hash.Hash) string {
	return fmt.Sprintf("%x", hash.Sum(nil))
}
