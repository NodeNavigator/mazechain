package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// CycleRewardKeyPrefix is the prefix to retrieve all CycleReward
	CycleRewardKeyPrefix = "CycleReward/value/"
)

// CycleRewardKey returns the store key to retrieve a CycleReward from the index fields
func CycleRewardKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
