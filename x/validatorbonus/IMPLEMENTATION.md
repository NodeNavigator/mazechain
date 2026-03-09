# Validator Bonus Module - Complete Implementation Guide

## Overview

This document provides a complete implementation of the `x/validatorbonus` Cosmos SDK module for distributing rewards to validators based on block proposer counts.

## System Design

### Program Duration
- **Total Duration**: 15 cycles
- **Cycle Length**: 30 days per cycle
- **Total Days**: 450 days (15 × 30)

### Reward Distribution Flow

```
Block Proposal (BeginBlocker)
    ↓
Track proposer in proposerCount/{validatorAddress}/{day}
    ↓
End of Day (EndBlocker - Day Change)
    ↓
Calculate daily rewards for each validator
Store in dailyReward/{validator}/{day}
    ↓
End of Cycle (EndBlocker - Cycle Change)
    ↓
Aggregate daily rewards for the cycle
Store in cycleReward/{cycle}/{validator}
Reset proposerCount for next cycle
```

## Data Structures

### Stores (KV Stores)

1. **eligibleValidator/{validatorAddress}**
   - Stores validators that joined within 30 days of chain genesis
   - Key: validatorAddress
   - Value: EligibleValidator{validatorAddress, joinTime}

2. **proposerCount/{validatorAddress}:{day}**
   - Tracks block proposals per validator per day
   - Key: validatorAddress:day
   - Value: ProposerCount{validatorAddress, day, count}

3. **dailyReward/{validatorAddress}:{day}**
   - Stores calculated daily rewards
   - Key: validatorAddress:day
   - Value: DailyReward{validatorAddress, day, amount}

4. **cycleReward/{cycle}:{validatorAddress}**
   - Stores aggregated cycle rewards
   - Key: cycle:validatorAddress
   - Value: CycleReward{cycle, validatorAddress, amount}

### Parameters

```go
type Params struct {
    TotalRewardPool string // Total reward pool to distribute (string for decimal precision)
    CycleDays       uint64 // Days per cycle (default: 30)
    TotalCycles     uint64 // Total number of cycles (default: 15)
}
```

## Implementation Files

### 1. `helpers.go` - Utility Functions

**Key Functions:**
- `GetDayFromTime(blockTime, genesisTime)` → uint64
  - Converts time to day index: (blockTime - genesisTime) / 86400
  
- `GetCycleFromDay(day, cycleDays)` → uint64
  - Calculates cycle from day: day / cycleDays

- `CalculateDailyRewardShare(validatorBlocks, totalBlocks, totalRewardPool, cycleDays, totalCycles)` → Decimal
  - Calculates daily reward using the formula:
    ```
    validatorShare = validatorBlocks / totalBlocks
    dailyRewardPool = totalRewardPool / (cycleDays × totalCycles)
    dailyReward = validatorShare × dailyRewardPool
    ```

- `ParseDecFromString(value)` → Decimal
  - Safely parses string to Decimal, returning zero on error

### 2. `begin_block.go` - Block Proposer Tracking

**BeginBlocker Function:**
1. Retrieves the block proposer from the context
2. Converts consensus address to validator operator address
3. Checks if the validator is in the eligibleValidator store
4. If eligible: increments `proposerCount/{validatorAddress}/{day}`

**Key Methods:**
- `BeginBlocker(ctx)` - Main entry point
- `IncrementProposerCount(ctx, validatorAddr, day)` - Safely increments count
- `GetProposerCountForDay(ctx, validatorAddr, day)` → uint64
- `GetAllProposerCountsForDay(ctx, day)` → []ProposerCount
- `GetTotalBlocksForDay(ctx, day)` → uint64

### 3. `reward.go` - Reward Calculation & Aggregation

**Daily Reward Calculation:**
```go
CalculateAndStoreDailyRewards(ctx, day)
  ├─ Get total blocks for the day
  ├─ For each proposer in proposerCount for that day:
  │   ├─ Calculate their share: proposerBlocks / totalBlocks
  │   ├─ Calculate reward: share × (totalRewardPool / totalDays)
  │   └─ Store in dailyReward/{validator}/{day}
  └─ Called when day boundary is crossed
```

**Cycle Reward Aggregation:**
```go
CalculateAndStoreCycleRewards(ctx, cycleEndDay+1)
  ├─ For each eligible validator:
  │   ├─ Sum all dailyReward for days in this cycle
  │   └─ Store in cycleReward/{cycle}/{validator}
  ├─ Reset proposerCount for the cycle days
  └─ Called when cycle boundary is crossed
```

**Key Methods:**
- `CalculateAndStoreDailyRewards(ctx, day)` - End-of-day reward calculation
- `StoreDailyReward(ctx, validatorAddr, day, amount)` - Store daily reward
- `GetDailyReward(ctx, validatorAddr, day)` → (Decimal, bool)
- `GetDailyRewardsForValidator(ctx, validatorAddr, startDay, endDay)` → []DailyReward
- `CalculateAndStoreCycleRewards(ctx, cycleEndDay+1)` - End-of-cycle aggregation
- `StoreCycleReward(ctx, cycle, validatorAddr, amount)` - Store cycle reward
- `GetCycleReward(ctx, cycle, validatorAddr)` → (Decimal, bool)
- `GetAllCycleRewardsForCycle(ctx, cycle)` → []CycleReward
- `ResetProposerCountsForCycle(ctx, startDay, endDay)` - Clean up after aggregation

### 4. `end_block.go` - End-of-Day/Cycle Logic

**EndBlocker Function:**
Monitors day and cycle boundaries to trigger:
1. **End-of-Day**: Calculate daily rewards when transitioning to a new day
2. **End-of-Cycle**: Aggregate rewards and reset counters when transitioning to a new cycle

**Day Tracking:**
- Stores `genesis_time` to calculate current day index
- Tracks `last_processed_day` to detect day transitions

**Key Methods:**
- `EndBlocker(ctx)` - Main entry point
- `GetGenesisTime(ctx)` → int64
- `SetGenesisTime(ctx, genesisTime)` - Called during genesis
- `GetLastProcessedDay(ctx)` → uint64
- `SetLastProcessedDay(ctx, day)` - Updated each block

### 5. `query_validator_cycle_reward.go` - Query Implementation

**Query: validator_cycle_reward**
```protobuf
Request:
  - validatorAddress (string)
  - cycle (uint64)

Response:
  - reward (string) - Decimal amount or "0"
  - isValidator (bool) - Whether validator is eligible
```

**Logic:**
1. Check if validator exists in eligibleValidator store
2. If not found: return `{reward: "0", isValidator: false}`
3. If found: lookup cycleReward/{cycle}/{validator}
4. If no reward found: return `{reward: "0", isValidator: true}`
5. Otherwise: return the found reward with `isValidator: true`

### 6. `query_cycle_rewards.go` - Query Implementation

**Query: cycle_rewards**
```protobuf
Request:
  - cycle (uint64)

Response:
  - validators (string) - JSON array of validator addresses
  - rewards (string) - JSON array of reward amounts
```

**Logic:**
1. Get all cycleReward entries for the specified cycle
2. Build parallel arrays of validators and rewards
3. Marshal arrays to JSON and return as strings

## Eligibility Rules

**ELIGIBILITY REQUIREMENT:**
Only validators that joined within 30 days of chain genesis are eligible for rewards.

**How It Works:**
1. During genesis or early chain operation, eligible validators are registered:
   ```go
   SetEligibleValidator(ctx, EligibleValidator{
       ValidatorAddress: "cosmosvaloper1...",
       JoinTime: blockHeight or blockTime,
   })
   ```

2. The BeginBlocker checks `eligibleValidator` store before incrementing proposer count

3. Only eligible validators have their proposals tracked and rewarded

## Reward Calculation Formula

For a validator on a specific day:

```
dailyReward = (proposerBlocks / totalBlocksPerDay) × (totalRewardPool / totalRewardDays)

where:
  proposerBlocks = number of blocks this validator proposed that day
  totalBlocksPerDay = sum of all proposer blocks that day
  totalRewardPool = module parameter (e.g., "1000000000")
  totalRewardDays = cycleDays × totalCycles (450 days for default settings)
```

**Example:**
- Total reward pool: 1,000,000 tokens
- Cycle days: 30
- Total cycles: 15
- Total reward days: 450

If validator proposed 5 out of 10,800 blocks (30 seconds/block = 10,800 blocks/day):
```
daily_reward_share = 5 / 10800
daily_pool = 1,000,000 / 450 = 2,222.22 tokens per day
daily_reward = (5 / 10800) × 2,222.22 = 1.03 tokens
```

## Integration with Module

### Module Registration

The BeginBlocker and EndBlocker are wired in `module/module.go`:

```go
func (am AppModule) BeginBlock(ctx context.Context) error {
    return am.keeper.BeginBlocker(ctx)
}

func (am AppModule) EndBlocker(ctx context.Context) error {
    return am.keeper.EndBlocker(ctx)
}
```

### Genesis Initialization

During genesis setup, you should:
1. Set module parameters
2. Register eligible validators
3. Initialize genesis time

```go
// In InitGenesis (from module/genesis.go or handlers)
params := types.NewParams(
    "1000000000",  // total_reward_pool
    30,            // cycle_days
    15,            // total_cycles
)
k.SetParams(ctx, params)

// Register eligible validators
k.SetEligibleValidator(ctx, types.EligibleValidator{
    ValidatorAddress: "cosmosvaloper1...",
    JoinTime: genesisTime,
})

// Initialize genesis time in keeper
k.SetGenesisTime(ctx, genesisTime)
k.SetLastProcessedDay(ctx, 0)
```

## Production Considerations

### Performance Optimizations

1. **Avoid Heavy Iteration in BeginBlocker:**
   - ✅ O(1) operations: single validator lookup and count increment
   - ❌ Don't iterate all validators on each block

2. **Defer Complex Calculations to EndBlocker:**
   - Daily reward calculation happens once per day transition
   - Cycle aggregation happens once per cycle boundary
   - Not on every block

3. **Minimize KV Writes:**
   - ProposerCount only written when incremented
   - DailyReward written once per day per validator
   - CycleReward written once per cycle per validator

4. **Use Prefix Stores Efficiently:**
   - All queries use prefix iteration with early termination
   - Composite keys (e.g., `addr:day`) allow prefix scanning

### Testing Recommendations

1. **Unit Tests:**
   - Test day/cycle calculations with various timestamps
   - Test reward calculations with edge cases (0 blocks, single validator)
   - Test storage and retrieval operations

2. **Integration Tests:**
   - Simulate multiple days and cycles
   - Verify proposer tracking and aggregation
   - Test query responses

3. **Edge Cases:**
   - Genesis block (day 0)
   - Cycle boundaries
   - Multiple validators with different proposal rates
   - Large numbers of validators

## Migration Notes

If upgrading from an earlier version:
1. Ensure params are properly initialized
2. Set genesis time if not already stored
3. Migration handlers can reconstruct proposer counts if needed

## Governance

Module parameters can be updated via governance:
- `total_reward_pool`: Update the total rewards to distribute
- `cycle_days`: Adjust cycle length (requires careful migration)
- `total_cycles`: Extend or reduce the program duration
