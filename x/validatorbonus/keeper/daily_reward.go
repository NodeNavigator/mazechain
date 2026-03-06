package keeper

import (
	"context"

	"blockmazechain/x/validatorbonus/types"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetDailyReward set a specific dailyReward in the store from its index
func (k Keeper) SetDailyReward(ctx context.Context, dailyReward types.DailyReward) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.DailyRewardKeyPrefix))
	b := k.cdc.MustMarshal(&dailyReward)
	store.Set(types.DailyRewardKey(
		dailyReward.Index,
	), b)
}

// GetDailyReward returns a dailyReward from its index
func (k Keeper) GetDailyReward(
	ctx context.Context,
	index string,

) (val types.DailyReward, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.DailyRewardKeyPrefix))

	b := store.Get(types.DailyRewardKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveDailyReward removes a dailyReward from the store
func (k Keeper) RemoveDailyReward(
	ctx context.Context,
	index string,

) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.DailyRewardKeyPrefix))
	store.Delete(types.DailyRewardKey(
		index,
	))
}

// GetAllDailyReward returns all dailyReward
func (k Keeper) GetAllDailyReward(ctx context.Context) (list []types.DailyReward) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.DailyRewardKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DailyReward
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
