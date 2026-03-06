package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// ProposerCountKeyPrefix is the prefix to retrieve all ProposerCount
	ProposerCountKeyPrefix = "ProposerCount/value/"
)

// ProposerCountKey returns the store key to retrieve a ProposerCount from the index fields
func ProposerCountKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
