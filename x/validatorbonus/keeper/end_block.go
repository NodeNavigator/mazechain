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

	// Use hardcoded defaults for cycle days (will use params after proto regeneration)
	cycleDays := uint64(30)

	if cycleDays == 0 {
		k.Logger().Error("cycle_days is not set, skipping reward calculations")
		return nil
	}

	// Get current block time and calculate current day
	blockTime := sdkCtx.BlockTime()

	// Retrieve genesis time from store (should be set during genesis)
	genesisTime, err := k.GetGenesisTime(ctx)
	if err != nil {
		// Genesis time not set yet, set it now using current block time
		// This assumes the first block time represents genesis time
		genesisTimeUnix := blockTime.Unix()
		err := k.SetGenesisTime(ctx, genesisTimeUnix)
		if err != nil {
			k.Logger().Error("failed to set genesis time", "error", err)
			return nil
		}
		k.Logger().Info("genesis time initialized", "genesis_time", genesisTimeUnix)
		genesisTime = genesisTimeUnix
	}

	// Calculate current day
	secondsElapsed := blockTime.Unix() - genesisTime
	if secondsElapsed < 0 {
		return nil
	}
	currentDay := uint64(secondsElapsed) / 86400

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

// GetGenesisTime retrieves the genesis time from the module store.
// This should be set during genesis initialization.
func (k Keeper) GetGenesisTime(ctx context.Context) (int64, error) {
	storeAdapter := k.storeService.OpenKVStore(ctx)

	genesisTimeKey := []byte("genesis_time")
	value, _ := storeAdapter.Get(genesisTimeKey)

	if value == nil {
		return 0, fmt.Errorf("genesis time not found in store")
	}

	if len(value) != 8 {
		return 0, fmt.Errorf("invalid genesis time format")
	}

	// Parse int64 from bytes (big-endian)
	genesisTime := int64(0)
	for i := 0; i < 8; i++ {
		genesisTime = (genesisTime << 8) | int64(value[i])
	}

	return genesisTime, nil
}

// SetGenesisTime stores the genesis time in the module store.
func (k Keeper) SetGenesisTime(ctx context.Context, genesisTime int64) error {
	storeAdapter := k.storeService.OpenKVStore(ctx)

	// Encode int64 as big-endian bytes
	value := make([]byte, 8)
	for i := 0; i < 8; i++ {
		value[7-i] = byte(genesisTime >> (i * 8))
	}

	storeAdapter.Set([]byte("genesis_time"), value)
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
