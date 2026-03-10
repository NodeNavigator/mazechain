package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker executes the EndBlocker logic for the validator bonus module.
// It handles end-of-day reward calculations and end-of-cycle aggregations.
//
// Note: This function should be called from the module's EndBlock handler.
// It checks if a day or cycle boundary has been reached and performs the necessary calculations.
func (k Keeper) EndBlocker(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get params for cycle days
	params := k.GetParams(ctx)
	cycleDays := params.CycleDays

	if cycleDays == 0 {
		k.Logger().Error("cycle_days is not set in params, skipping reward calculations")
		return nil
	}

	// Get current block time and calculate current day
	blockTime := sdkCtx.BlockTime()

	// Retrieve activation time from store (should be set during genesis or upgrade)
	activationTime, err := k.GetActivationTime(ctx)
	if err != nil {
		// Activation time not set yet, set it now using current block time
		// This assumes the first block time (whether block 1 or upgrade block) represents activation time
		activationTimeUnix := blockTime.Unix()
		err := k.SetActivationTime(ctx, activationTimeUnix)
		if err != nil {
			k.Logger().Error("failed to set activation time", "error", err)
			return nil
		}
		k.Logger().Info("activation time initialized", "activation_time", activationTimeUnix)
		activationTime = activationTimeUnix
	}

	// Calculate current day
	secondsElapsed := blockTime.Unix() - activationTime
	if secondsElapsed < 0 {
		return nil
	}
	secondsPerDay := params.SecondsPerDay
	if secondsPerDay == 0 {
		secondsPerDay = 86400
	}
	currentDay := uint64(secondsElapsed) / secondsPerDay

	// Get the last processed day from store
	lastProcessedDay, err := k.GetLastProcessedDay(ctx)
	if err != nil {
		// First time initialization
		k.SetLastProcessedDay(ctx, currentDay)
		return nil
	}

	// Check if we've moved to a new day
	if currentDay > lastProcessedDay {
		// End-of-day logic: calculate daily rewards for the previous day
		previousDay := lastProcessedDay

		k.Logger().Info("End of day reached", "previous_day", previousDay)

		// Calculate and store daily rewards for the previous day
		if err := k.CalculateAndStoreDailyRewards(ctx, previousDay); err != nil {
			k.Logger().Error("failed to calculate daily rewards", "day", previousDay, "error", err)
		}

		// Check if we've completed a cycle
		currentCycle := k.GetCycleFromDay(currentDay, cycleDays)
		previousCycle := k.GetCycleFromDay(previousDay, cycleDays)

		if currentCycle > previousCycle {
			// End-of-cycle logic: aggregate daily rewards into cycle rewards
			k.Logger().Info("End of cycle reached", "previous_cycle", previousCycle)

			// Calculate cycle end day
			cycleEndDay := (previousCycle+1)*cycleDays - 1

			// Calculate and store cycle rewards
			if err := k.CalculateAndStoreCycleRewards(ctx, cycleEndDay+1); err != nil {
				k.Logger().Error("failed to calculate cycle rewards", "cycle", previousCycle, "error", err)
			}
		}

		// Update the last processed day
		k.SetLastProcessedDay(ctx, currentDay)
	}

	return nil
}

// GetActivationTime retrieves the activation time from the module store.
// This is set automatically on the very first block the module runs (genesis or upgrade).
func (k Keeper) GetActivationTime(ctx context.Context) (int64, error) {
	storeAdapter := k.storeService.OpenKVStore(ctx)

	activationTimeKey := []byte("activation_time")
	value, _ := storeAdapter.Get(activationTimeKey)

	if value == nil {
		return 0, fmt.Errorf("activation time not found in store")
	}

	if len(value) != 8 {
		return 0, fmt.Errorf("invalid activation time format")
	}

	// Parse int64 from bytes (big-endian)
	activationTime := int64(0)
	for i := 0; i < 8; i++ {
		activationTime = (activationTime << 8) | int64(value[i])
	}

	return activationTime, nil
}

// SetActivationTime stores the activation time in the module store.
func (k Keeper) SetActivationTime(ctx context.Context, activationTime int64) error {
	storeAdapter := k.storeService.OpenKVStore(ctx)

	// Encode int64 as big-endian bytes
	value := make([]byte, 8)
	for i := 0; i < 8; i++ {
		value[7-i] = byte(activationTime >> (i * 8))
	}

	storeAdapter.Set([]byte("activation_time"), value)
	return nil
}

// GetLastProcessedDay retrieves the last day that was processed for rewards.
func (k Keeper) GetLastProcessedDay(ctx context.Context) (uint64, error) {
	storeAdapter := k.storeService.OpenKVStore(ctx)

	lastProcessedDayKey := []byte("last_processed_day")
	value, _ := storeAdapter.Get(lastProcessedDayKey)

	if value == nil {
		return 0, fmt.Errorf("last processed day not found in store")
	}

	if len(value) != 8 {
		return 0, fmt.Errorf("invalid last processed day format")
	}

	// Parse uint64 from bytes (big-endian)
	lastProcessedDay := uint64(0)
	for i := 0; i < 8; i++ {
		lastProcessedDay = (lastProcessedDay << 8) | uint64(value[i])
	}

	return lastProcessedDay, nil
}

// SetLastProcessedDay stores the last processed day in the module store.
func (k Keeper) SetLastProcessedDay(ctx context.Context, day uint64) error {
	storeAdapter := k.storeService.OpenKVStore(ctx)

	// Encode uint64 as big-endian bytes
	value := make([]byte, 8)
	for i := 0; i < 8; i++ {
		value[7-i] = byte(day >> (i * 8))
	}

	storeAdapter.Set([]byte("last_processed_day"), value)
	return nil
}
