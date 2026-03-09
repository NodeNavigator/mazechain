package keeper

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"blockmazechain/x/validatorbonus/types"
)

// BeginBlocker executes the BeginBlocker logic for the validator bonus module.
// It tracks the block proposer and increments their proposer count for the current day.
//
// Logic:
// 1. Get the block proposer from the context
// 2. Get current block time
// 3. Convert time → day index since genesis
// 4. Check if proposer exists in eligibleValidator store
// 5. If eligible: increment proposerCount/{validatorAddress}/{day} += 1
func (k Keeper) BeginBlocker(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get the proposer from the block header
	proposerConsAddr := sdkCtx.BlockHeader().ProposerAddress
	if len(proposerConsAddr) == 0 {
		k.Logger().Debug("no proposer in block header")
		return nil
	}

	// Convert consensus address to validator address using staking keeper
	validator, err := k.stakingKeeper.ValidatorByConsAddr(ctx, proposerConsAddr)
	if err != nil {
		k.Logger().Debug("validator not found by consensus address", "error", err)
		return nil
	}

	// Get the validator's operator address
	validatorAddr := validator.GetOperator()
	if validatorAddr == "" {
		k.Logger().Error("validator operator address is empty")
		return nil
	}

	// Get current block time for eligibility check
	blockTime := sdkCtx.BlockTime()

	// Check if validator is eligible (exists in eligibleValidator store by address)
	_, found := k.GetEligibleValidatorByAddress(ctx, validatorAddr)
	if !found {
		// Check if we're still within the 30-day eligibility window from genesis
		isWithin30Days, err := k.IsWithin30DayEligibilityWindow(ctx, blockTime)
		if err != nil {
			// Log error but continue - don't add validator if we can't determine eligibility
			k.Logger().Error("failed to check eligibility window", "validator", validatorAddr, "error", err)
			return nil
		}

		if isWithin30Days {
			// Auto-add validator to eligible list
			// Assign a sequential numeric id (1, 2, 3, ...) using a stored counter.
			nextNumericID := k.getNextEligibleValidatorID(ctx)
			eligibleValidator := types.EligibleValidator{
				Id:               strconv.FormatUint(nextNumericID, 10), // "1", "2", "3", ...
				ValidatorAddress: validatorAddr,                         // actual validator operator address
				JoinTime:         int32(sdkCtx.BlockTime().Unix()),      // current block time (unix)
				Creator:          k.authority,                           // module authority
			}
			k.SetEligibleValidator(ctx, eligibleValidator)
			k.setNextEligibleValidatorID(ctx, nextNumericID+1)
			k.Logger().Info("auto-added validator to eligible list", "validator", validatorAddr, "height", sdkCtx.BlockHeight())
		} else {
			// We're past the 30-day window, validator not eligible
			k.Logger().Debug("validator joined after 30-day eligibility window", "validator", validatorAddr)
			return nil
		}
	}

	// Get genesis time from context
	// Retrieve genesis time (stored in store during genesis)
	genesisTime, err := k.GetGenesisTime(ctx)
	if err != nil {
		k.Logger().Debug("genesis time not yet initialized", "validator", validatorAddr, "error", err)
		return nil
	}

	// Calculate day from genesis
	secondsElapsed := blockTime.Unix() - genesisTime
	if secondsElapsed < 0 {
		k.Logger().Warn("block time is before genesis time", "validator", validatorAddr, "block_time", blockTime.Unix(), "genesis_time", genesisTime)
		return nil
	}
	day := uint64(secondsElapsed) / 86400

	// Increment proposer count for the current day
	err = k.IncrementProposerCount(ctx, validatorAddr, day)
	if err != nil {
		k.Logger().Error("failed to increment proposer count", "validator", validatorAddr, "day", day, "error", err)
		return err
	}

	return nil
}

// IncrementProposerCount increments the proposer count for a validator on a specific day.
// This function avoids heavy KV writes by only updating when necessary.
func (k Keeper) IncrementProposerCount(ctx context.Context, validatorAddr string, day uint64) error {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProposerCountKeyPrefix))

	// Construct composite key: validatorAddress:day
	index := fmt.Sprintf("%s:%d", validatorAddr, day)

	// Get existing count
	b := store.Get(types.ProposerCountKey(index))

	var count uint64 = 0
	if b != nil {
		var pc types.ProposerCount
		k.cdc.MustUnmarshal(b, &pc)
		count = pc.Count
	}

	// Create new ProposerCount entry with incremented count
	proposerCount := types.ProposerCount{
		Id:               index,
		ValidatorAddress: validatorAddr,
		Day:              day,
		Count:            count + 1,
		Creator:          k.authority,
	}

	// Store the updated count
	newB := k.cdc.MustMarshal(&proposerCount)
	store.Set(types.ProposerCountKey(index), newB)

	return nil
}

// GetProposerCountForDay retrieves the proposer count for a specific validator on a specific day.
func (k Keeper) GetProposerCountForDay(ctx context.Context, validatorAddr string, day uint64) uint64 {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProposerCountKeyPrefix))

	index := fmt.Sprintf("%s:%d", validatorAddr, day)
	b := store.Get(types.ProposerCountKey(index))

	if b == nil {
		return 0
	}

	var pc types.ProposerCount
	k.cdc.MustUnmarshal(b, &pc)
	return pc.Count
}

// GetAllProposerCountsForDay retrieves all proposer counts for a specific day.
// Used during end-of-day reward calculation.
func (k Keeper) GetAllProposerCountsForDay(ctx context.Context, day uint64) []types.ProposerCount {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProposerCountKeyPrefix))

	var counts []types.ProposerCount

	// Iterate all and filter by day
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var pc types.ProposerCount
		k.cdc.MustUnmarshal(iterator.Value(), &pc)
		if pc.Day == day {
			counts = append(counts, pc)
		}
	}

	return counts
}

// GetTotalBlocksForDay calculates the total number of blocks proposed on a specific day.
func (k Keeper) GetTotalBlocksForDay(ctx context.Context, day uint64) uint64 {
	counts := k.GetAllProposerCountsForDay(ctx, day)
	var total uint64 = 0
	for _, count := range counts {
		total += count.Count
	}
	return total
}

// IsWithin30DayEligibilityWindow checks if we're still within the 30-day window from genesis
// where new validators can automatically join the eligible list
func (k Keeper) IsWithin30DayEligibilityWindow(ctx context.Context, currentTime time.Time) (bool, error) {
	// Get genesis time
	genesisTime, err := k.GetGenesisTime(ctx)
	if err != nil {
		return false, err
	}

	// Convert genesis time to Time
	genesisTimeObj := time.Unix(genesisTime, 0)

	// Calculate 30 days from genesis
	thirtyDaysFromGenesis := genesisTimeObj.Add(30 * 24 * time.Hour)

	// Check if current time is still within 30 days of genesis
	if currentTime.After(thirtyDaysFromGenesis) {
		// We're past the 30-day window, so no new validators can join automatically
		return false, nil
	}

	// We're within the 30-day window, so new validators can join
	return true, nil
}
