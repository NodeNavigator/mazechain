package keeper

import (
	"context"

	"blockmazechain/x/validatorbonus/types"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetCycleReward set a specific cycleReward in the store from its index
func (k Keeper) SetCycleReward(ctx context.Context, cycleReward types.CycleReward) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.CycleRewardKeyPrefix))
	b := k.cdc.MustMarshal(&cycleReward)
	store.Set(types.CycleRewardKey(
		cycleReward.Index,
	), b)
}

// GetCycleReward returns a cycleReward from its index
func (k Keeper) GetCycleReward(
	ctx context.Context,
	index string,

) (val types.CycleReward, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.CycleRewardKeyPrefix))

	b := store.Get(types.CycleRewardKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveCycleReward removes a cycleReward from the store
func (k Keeper) RemoveCycleReward(
	ctx context.Context,
	index string,

) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.CycleRewardKeyPrefix))
	store.Delete(types.CycleRewardKey(
		index,
	))
}

// GetAllCycleReward returns all cycleReward
func (k Keeper) GetAllCycleReward(ctx context.Context) (list []types.CycleReward) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.CycleRewardKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.CycleReward
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
