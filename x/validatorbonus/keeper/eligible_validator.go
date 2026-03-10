package keeper

import (
	"context"
	"encoding/binary"

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
		eligibleValidator.Id,
	), b)

	// Update secondary index: address -> id
	if eligibleValidator.ValidatorAddress != "" {
		addressStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.EligibleValidatorByAddressKeyPrefix))
		addressStore.Set([]byte(eligibleValidator.ValidatorAddress), []byte(eligibleValidator.Id))
	}
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
	// Get the validator to find its address for secondary index removal
	val, found := k.GetEligibleValidator(ctx, index)
	if found {
		addressStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.EligibleValidatorByAddressKeyPrefix))
		addressStore.Delete([]byte(val.ValidatorAddress))
	}

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

// internal key for auto-incrementing eligible validator IDs.
const eligibleValidatorSeqKey = "EligibleValidatorSeq"

// getNextEligibleValidatorID reads the next numeric ID to assign; defaults to 1 if unset.
func (k Keeper) getNextEligibleValidatorID(ctx context.Context) uint64 {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := storeAdapter

	bz := store.Get([]byte(eligibleValidatorSeqKey))
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz)
}

// setNextEligibleValidatorID persists the next numeric ID to assign.
func (k Keeper) setNextEligibleValidatorID(ctx context.Context, nextID uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := storeAdapter

	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, nextID)
	store.Set([]byte(eligibleValidatorSeqKey), bz)
}

// GetEligibleValidatorByAddress returns an eligibleValidator by its validatorAddress field.
// This is used internally for eligibility checks where the natural key is the validator address.
func (k Keeper) GetEligibleValidatorByAddress(ctx context.Context, validatorAddr string) (val types.EligibleValidator, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	addressStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.EligibleValidatorByAddressKeyPrefix))

	// Look up the ID from the secondary index
	idBytes := addressStore.Get([]byte(validatorAddr))
	if idBytes == nil {
		return types.EligibleValidator{}, false
	}

	// Fetch the full record using the ID
	return k.GetEligibleValidator(ctx, string(idBytes))
}
