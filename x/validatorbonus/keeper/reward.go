package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"

	"blockmazechain/x/validatorbonus/types"
)

// CalculateAndStoreDailyRewards calculates and stores daily rewards for all validators
// at the end of a day. This should be called when transitioning to the next day.
//
// Logic:
// For each validator in proposerCount store for that day:
// - Calculate reward share: validatorShare = proposerBlocks / totalBlocksThatDay
// - Calculate daily reward: dailyReward = validatorShare * (total_reward_pool / (cycle_days * total_cycles))
// - Store in dailyReward/{validator}/{day}
func (k Keeper) CalculateAndStoreDailyRewards(ctx context.Context, day uint64) error {
	// Get total blocks for the day
	totalBlocks := k.GetTotalBlocksForDay(ctx, day)
	if totalBlocks == 0 {
		// No proposers for this day
		return nil
	}

	// Note: Params will be populated after proto regeneration
	// For now, use hardcoded defaults that match the requirement
	totalRewardPool := math.LegacyNewDec(1000000000) // 1 billion
	cycleDays := uint64(30)
	totalCycles := uint64(15)

	totalBlocksDec := math.NewInt(int64(totalBlocks))

	// Get all proposer counts for this day
	proposerCounts := k.GetAllProposerCountsForDay(ctx, day)

	// Iterate through each proposer and calculate their reward
	for _, proposerCount := range proposerCounts {
		validatorAddr := proposerCount.ValidatorAddress
		proposerBlocks := math.NewInt(int64(proposerCount.Count))

		// Calculate daily reward using helper function
		dailyReward := k.CalculateDailyRewardShare(
			proposerBlocks,
			totalBlocksDec,
			totalRewardPool,
			cycleDays,
			totalCycles,
		)

		// Store the daily reward
		err := k.StoreDailyRewardInternal(ctx, validatorAddr, day, dailyReward.String())
		if err != nil {
			// Log error but continue processing other validators
			k.Logger().Error("failed to store daily reward", "validator", validatorAddr, "day", day, "error", err)
		}
	}

	return nil
}

// StoreDailyRewardInternal stores a daily reward for a validator on a specific day.
// This is an internal helper separate from the generated StoreDailyReward.
func (k Keeper) StoreDailyRewardInternal(ctx context.Context, validatorAddr string, day uint64, amount string) error {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.DailyRewardKeyPrefix))

	// Construct composite key: validatorAddress:day
	index := fmt.Sprintf("%s:%d", validatorAddr, day)

	dailyReward := types.DailyReward{
		Id:               index,
		ValidatorAddress: validatorAddr,
		Day:              day,
		Amount:           amount,
		Creator:          k.authority,
	}

	b := k.cdc.MustMarshal(&dailyReward)
	store.Set(types.DailyRewardKey(index), b)

	return nil
}

// GetDailyRewardInternal retrieves the daily reward for a validator on a specific day.
// This is an internal helper separate from the generated GetDailyReward.
func (k Keeper) GetDailyRewardInternal(ctx context.Context, validatorAddr string, day uint64) (math.LegacyDec, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.DailyRewardKeyPrefix))

	index := fmt.Sprintf("%s:%d", validatorAddr, day)
	b := store.Get(types.DailyRewardKey(index))

	if b == nil {
		return math.LegacyZeroDec(), false
	}

	var dr types.DailyReward
	k.cdc.MustUnmarshal(b, &dr)

	amount := ParseDecFromString(dr.Amount)
	return amount, true
}

// GetDailyRewardsForValidator retrieves all daily rewards for a validator within a date range.
// Used for cycle aggregation.
func (k Keeper) GetDailyRewardsForValidator(ctx context.Context, validatorAddr string, startDay, endDay uint64) []types.DailyReward {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.DailyRewardKeyPrefix))

	var rewards []types.DailyReward

	// Iterate through all daily rewards
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var dr types.DailyReward
		k.cdc.MustUnmarshal(iterator.Value(), &dr)

		// Filter by validator address and day range
		if dr.ValidatorAddress == validatorAddr && dr.Day >= startDay && dr.Day <= endDay {
			rewards = append(rewards, dr)
		}
	}

	return rewards
}

// CalculateAndStoreCycleRewards calculates and stores cycle rewards for all validators
// at the end of a cycle. This aggregates daily rewards for the past 30 days.
//
// Logic:
// Every 30 days:
// - For each validator in eligibleValidator store:
//   - Sum dailyReward for the last cycle_days
//   - Store in cycleReward/{cycle}/{validator}
//
// - Then reset proposerCount for next cycle
func (k Keeper) CalculateAndStoreCycleRewards(ctx context.Context, currentDay uint64) error {
	// Note: Params will be populated after proto regeneration
	// For now, use hardcoded defaults
	cycleDays := uint64(30)

	if cycleDays == 0 {
		return fmt.Errorf("cycle_days must be greater than 0")
	}

	// Calculate cycle number
	cycle := k.GetCycleFromDay(currentDay, cycleDays)

	// Calculate the start day of this cycle
	// Note: This is called at the end of a cycle, so we need to aggregate the previous cycle
	cycleStartDay := (cycle - 1) * cycleDays
	cycleEndDay := cycleStartDay + cycleDays - 1

	// If currentDay doesn't align with cycle boundary, adjust
	if currentDay%cycleDays != 0 {
		return fmt.Errorf("CalculateAndStoreCycleRewards should only be called at cycle boundaries")
	}

	// Get all eligible validators
	eligibleValidators := k.GetAllEligibleValidator(ctx)

	// For each eligible validator, sum their daily rewards for the cycle
	for _, ev := range eligibleValidators {
		validatorAddr := ev.ValidatorAddress

		// Get all daily rewards for this validator in this cycle
		dailyRewards := k.GetDailyRewardsForValidator(ctx, validatorAddr, cycleStartDay, cycleEndDay)

		// Sum the rewards
		totalCycleReward := math.LegacyZeroDec()
		for _, dr := range dailyRewards {
			amountDec := ParseDecFromString(dr.Amount)
			totalCycleReward = totalCycleReward.Add(amountDec)
		}

		// Store the cycle reward
		err := k.StoreCycleRewardInternal(ctx, cycle, validatorAddr, totalCycleReward.String())
		if err != nil {
			k.Logger().Error("failed to store cycle reward", "cycle", cycle, "validator", validatorAddr, "error", err)
		}
	}

	// Reset proposer counts for the next cycle
	// This clears proposer counts for the past cycle to avoid re-aggregation
	k.ResetProposerCountsForCycle(ctx, cycleStartDay, cycleEndDay)

	return nil
}

// StoreCycleRewardInternal stores a cycle reward for a validator in a specific cycle.
func (k Keeper) StoreCycleRewardInternal(ctx context.Context, cycle uint64, validatorAddr string, amount string) error {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.CycleRewardKeyPrefix))

	// Construct composite key: cycle:validatorAddress
	index := fmt.Sprintf("%d:%s", cycle, validatorAddr)

	cycleReward := types.CycleReward{
		Id:               index,
		Cycle:            cycle,
		ValidatorAddress: validatorAddr,
		Amount:           amount,
		Creator:          k.authority,
	}

	b := k.cdc.MustMarshal(&cycleReward)
	store.Set(types.CycleRewardKey(index), b)

	return nil
}

// GetCycleRewardInternal retrieves the cycle reward for a validator in a specific cycle.
func (k Keeper) GetCycleRewardInternal(ctx context.Context, cycle uint64, validatorAddr string) (math.LegacyDec, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.CycleRewardKeyPrefix))

	index := fmt.Sprintf("%d:%s", cycle, validatorAddr)
	b := store.Get(types.CycleRewardKey(index))

	if b == nil {
		return math.LegacyZeroDec(), false
	}

	var cr types.CycleReward
	k.cdc.MustUnmarshal(b, &cr)

	amount := ParseDecFromString(cr.Amount)
	return amount, true
}

// GetAllCycleRewardsForCycle retrieves all cycle rewards for a specific cycle.
// Used in queries to return all validators and their rewards for a cycle.
func (k Keeper) GetAllCycleRewardsForCycle(ctx context.Context, cycle uint64) []types.CycleReward {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.CycleRewardKeyPrefix))

	var rewards []types.CycleReward

	// Iterate through all cycle rewards
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var cr types.CycleReward
		k.cdc.MustUnmarshal(iterator.Value(), &cr)

		// Filter by cycle
		if cr.Cycle == cycle {
			rewards = append(rewards, cr)
		}
	}

	return rewards
}

// ResetProposerCountsForCycle removes proposer counts for a specific cycle range.
// This prevents re-aggregation of rewards in subsequent cycles.
func (k Keeper) ResetProposerCountsForCycle(ctx context.Context, startDay, endDay uint64) error {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProposerCountKeyPrefix))

	// Collect keys to delete (can't modify store while iterating)
	keysToDelete := [][]byte{}

	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var pc types.ProposerCount
		k.cdc.MustUnmarshal(iterator.Value(), &pc)

		// Delete if day falls within cycle range
		if pc.Day >= startDay && pc.Day <= endDay {
			keysToDelete = append(keysToDelete, iterator.Key())
		}
	}

	// Delete collected keys
	for _, key := range keysToDelete {
		store.Delete(key)
	}

	return nil
}
