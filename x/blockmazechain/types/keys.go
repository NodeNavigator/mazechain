package types

const (
	// ModuleName defines the module name
	ModuleName = "blockmazechain"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_blockmazechain"
)

var (
	ParamsKey = []byte("p_blockmazechain")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
