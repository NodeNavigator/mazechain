package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// EligibleValidatorKeyPrefix is the prefix to retrieve all EligibleValidator
	EligibleValidatorKeyPrefix = "EligibleValidator/value/"
)

// EligibleValidatorKey returns the store key to retrieve a EligibleValidator from the index fields
func EligibleValidatorKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
