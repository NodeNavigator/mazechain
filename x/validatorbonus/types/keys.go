package types

const (
	// ModuleName defines the module name
	ModuleName = "validatorbonus"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_validatorbonus"
)

var (
	ParamsKey = []byte("p_validatorbonus")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
