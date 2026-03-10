package keeper

import (
	"context"
	"time"

	"cosmossdk.io/math"
)

// GetDayFromTime calculates the day index from the provided time and activation time.
// day = (blockTime - activationTime) / secondsPerDay
func (k Keeper) GetDayFromTime(ctx context.Context, blockTime time.Time, activationTime time.Time) uint64 {
	secondsElapsed := blockTime.Unix() - activationTime.Unix()
	if secondsElapsed < 0 {
		return 0
	}

	secondsPerDay := k.GetParams(ctx).SecondsPerDay
	if secondsPerDay == 0 {
		secondsPerDay = 86400
	}

	return uint64(secondsElapsed) / secondsPerDay
}

// GetCycleFromDay calculates the cycle number from the day index.
// cycle = day / cycle_days
func (k Keeper) GetCycleFromDay(day uint64, cycleDays uint64) uint64 {
	if cycleDays == 0 {
		return 0
	}
	return day / cycleDays
}

// CalculateDailyRewardShare calculates the daily reward amount for a validator.
// Formula: validatorShare = proposerBlocks / totalBlocksThatDay
//
//	dailyReward = validatorShare * (total_reward_pool / (cycle_days * total_cycles))
func (k Keeper) CalculateDailyRewardShare(
	validatorProposerBlocks math.Int,
	totalBlocksPerDay math.Int,
	totalRewardPool math.LegacyDec,
	cycleDays uint64,
	totalCycles uint64,
) math.LegacyDec {
	if totalBlocksPerDay.IsZero() || cycleDays == 0 || totalCycles == 0 {
		return math.LegacyZeroDec()
	}

	// Convert to Decimal for precise calculations
	validatorBlocksDec := math.LegacyNewDecFromInt(validatorProposerBlocks)
	totalBlocksDec := math.LegacyNewDecFromInt(totalBlocksPerDay)

	// validatorShare = proposerBlocks / totalBlocksThatDay
	validatorShare := validatorBlocksDec.Quo(totalBlocksDec)

	// Calculate daily reward pool per cycle: total_reward_pool / (cycle_days * total_cycles)
	totalDays := math.NewInt(int64(cycleDays) * int64(totalCycles))
	dailyRewardPool := totalRewardPool.Quo(math.LegacyNewDecFromInt(totalDays))

	// dailyReward = validatorShare * dailyRewardPool
	dailyReward := validatorShare.Mul(dailyRewardPool)

	return dailyReward
}

// ParseDecFromString parses a string into a LegacyDec value.
// Returns zero decimal if parsing fails.
func ParseDecFromString(value string) math.LegacyDec {
	dec, err := math.LegacyNewDecFromStr(value)
	if err != nil {
		return math.LegacyZeroDec()
	}
	return dec
}

// ConstructKey constructs a composite key from multiple parts.
// Used for complex keys like validatorAddress + day or cycle + validatorAddress.
func ConstructKey(parts ...string) string {
	key := ""
	for i, part := range parts {
		if i > 0 {
			key += ":"
		}
		key += part
	}
	return key
}
