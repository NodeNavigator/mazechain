package keeper

import (
	"context"

	"blockmazechain/x/validatorbonus/types"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetEligibleValidator set a specific eligibleValidator in the store from its index
func (k Keeper) SetEligibleValidator(ctx context.Context, eligibleValidator types.EligibleValidator) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.EligibleValidatorKeyPrefix))
	b := k.cdc.MustMarshal(&eligibleValidator)
	store.Set(types.EligibleValidatorKey(
		eligibleValidator.Index,
	), b)
}

// GetEligibleValidator returns a eligibleValidator from its index
func (k Keeper) GetEligibleValidator(
	ctx context.Context,
	index string,

) (val types.EligibleValidator, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.EligibleValidatorKeyPrefix))

	b := store.Get(types.EligibleValidatorKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveEligibleValidator removes a eligibleValidator from the store
func (k Keeper) RemoveEligibleValidator(
	ctx context.Context,
	index string,

) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.EligibleValidatorKeyPrefix))
	store.Delete(types.EligibleValidatorKey(
		index,
	))
}

// GetAllEligibleValidator returns all eligibleValidator
func (k Keeper) GetAllEligibleValidator(ctx context.Context) (list []types.EligibleValidator) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.EligibleValidatorKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.EligibleValidator
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
