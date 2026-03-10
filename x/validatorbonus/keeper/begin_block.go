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

	// Get current block time
	blockTime := sdkCtx.BlockTime()

	// SINGLE initialization point for ActivationTime (fix loophole #3).
	// This must run before either eligibility check or day calculation so
	// both code paths always reference the same stored timestamp.
	activationTime, err := k.EnsureActivationTimeInitialized(ctx, blockTime)
	if err != nil {
		k.Logger().Error("failed to ensure activation time", "error", err)
		return nil
	}

	// Check if validator is eligible (exists in eligibleValidator store by address)
	_, found := k.GetEligibleValidatorByAddress(ctx, validatorAddr)
	if !found {
		// Check if we're still within the eligibility window from activation
		isWithinWindow, err := k.IsWithinEligibilityWindow(ctx, blockTime)
		if err != nil {
			// Log error but continue - don't add validator if we can't determine eligibility
			k.Logger().Error("failed to check eligibility window", "validator", validatorAddr, "error", err)
			return nil
		}

		if isWithinWindow {
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
			// We're past the eligibility window, validator not eligible
			k.Logger().Debug("validator joined after eligibility window", "validator", validatorAddr)
			return nil
		}
	}

	// Calculate day from activation
	secondsElapsed := blockTime.Unix() - activationTime
	if secondsElapsed < 0 {
		k.Logger().Warn("block time is before activation time", "validator", validatorAddr, "block_time", blockTime.Unix(), "activation_time", activationTime)
		return nil
	}
	params := k.GetParams(ctx)
	secondsPerDay := params.SecondsPerDay
	if secondsPerDay == 0 {
		secondsPerDay = 86400
	}
	day := uint64(secondsElapsed) / secondsPerDay

	// Increment proposer count for the current day
	err = k.IncrementProposerCount(ctx, validatorAddr, day)
	if err != nil {
		k.Logger().Error("failed to increment proposer count", "validator", validatorAddr, "day", day, "error", err)
		return err
	}

	return nil
}

// EnsureActivationTimeInitialized is the single authoritative point for reading ActivationTime.
// If the ActivationTime is not yet in the store it is set to blockTime, which correctly handles
// both brand-new chains (Block 1) and Cosmovisor upgrades (first block of the new binary).
// All other code should call this instead of GetActivationTime directly.
func (k Keeper) EnsureActivationTimeInitialized(ctx context.Context, blockTime time.Time) (int64, error) {
	activationTime, err := k.GetActivationTime(ctx)
	if err != nil {
		activationTimeUnix := blockTime.Unix()
		if setErr := k.SetActivationTime(ctx, activationTimeUnix); setErr != nil {
			return 0, setErr
		}
		k.Logger().Info("activation time initialized", "activation_time", activationTimeUnix)
		return activationTimeUnix, nil
	}
	return activationTime, nil
}

// IncrementProposerCount increments the proposer count for a validator on a specific day.
// This function avoids heavy KV writes by only updating when necessary.
func (k Keeper) IncrementProposerCount(ctx context.Context, validatorAddr string, day uint64) error {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProposerCountKeyPrefix))

	// Construct composite key: validatorAddress:day
	// Use %08d to zero-pad the day for correct lexicographical sorting in pagination
	index := fmt.Sprintf("%s:%08d", validatorAddr, day)

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

	index := fmt.Sprintf("%s:%08d", validatorAddr, day)
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

// IsWithinEligibilityWindow checks if we're still within the validator onboarding window
// from activation. The window duration equals exactly one cycle (CycleDays * SecondsPerDay),
// so it stays consistent with whatever cycle parameters are set in params.
// This replaces the old hard-coded "30 * 24 * time.Hour" window (fix loophole #2).
func (k Keeper) IsWithinEligibilityWindow(ctx context.Context, currentTime time.Time) (bool, error) {
	// Activation time is guaranteed to be initialized before this point
	// (EnsureActivationTimeInitialized is always called first in BeginBlocker).
	activationTime, err := k.GetActivationTime(ctx)
	if err != nil {
		return false, err
	}

	params := k.GetParams(ctx)
	cycleDays := params.CycleDays
	if cycleDays == 0 {
		cycleDays = 30
	}
	secondsPerDay := params.SecondsPerDay
	if secondsPerDay == 0 {
		secondsPerDay = 86400
	}

	// Window = CycleDays * SecondsPerDay (e.g. 30 days * 86400 s/day = 2,592,000 s)
	windowSeconds := int64(cycleDays * secondsPerDay)
	windowEnd := time.Unix(activationTime+windowSeconds, 0)

	if currentTime.After(windowEnd) {
		// Past the eligibility window
		return false, nil
	}
	return true, nil
}
