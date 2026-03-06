package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// DailyRewardKeyPrefix is the prefix to retrieve all DailyReward
	DailyRewardKeyPrefix = "DailyReward/value/"
)

// DailyRewardKey returns the store key to retrieve a DailyReward from the index fields
func DailyRewardKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
