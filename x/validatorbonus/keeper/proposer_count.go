package keeper

import (
	"context"

	"blockmazechain/x/validatorbonus/types"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetProposerCount set a specific proposerCount in the store from its index
func (k Keeper) SetProposerCount(ctx context.Context, proposerCount types.ProposerCount) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProposerCountKeyPrefix))
	b := k.cdc.MustMarshal(&proposerCount)
	store.Set(types.ProposerCountKey(
		proposerCount.Id,
	), b)
}

// GetProposerCount returns a proposerCount from its index
func (k Keeper) GetProposerCount(
	ctx context.Context,
	index string,

) (val types.ProposerCount, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProposerCountKeyPrefix))

	b := store.Get(types.ProposerCountKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveProposerCount removes a proposerCount from the store
func (k Keeper) RemoveProposerCount(
	ctx context.Context,
	index string,

) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProposerCountKeyPrefix))
	store.Delete(types.ProposerCountKey(
		index,
	))
}

// GetAllProposerCount returns all proposerCount
func (k Keeper) GetAllProposerCount(ctx context.Context) (list []types.ProposerCount) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProposerCountKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ProposerCount
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
