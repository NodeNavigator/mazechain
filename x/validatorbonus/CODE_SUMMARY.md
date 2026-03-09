# Validator Bonus Module - Code Summary

## Files Created/Modified

### 1. **helpers.go** - Utility Functions
Location: `x/validatorbonus/keeper/helpers.go`

Core utility functions:
- `GetDayFromTime()` - Convert timestamp to day index
- `GetCycleFromDay()` - Convert day to cycle number
- `CalculateDailyRewardShare()` - Calculate validator's daily reward share
- `ParseDecFromString()` - Safe decimal string parsing
- `ConstructKey()` - Build composite store keys

**Usage**: Mathematical calculations and key construction for storage operations.

---

### 2. **begin_block.go** - Block Proposer Tracking
Location: `x/validatorbonus/keeper/begin_block.go`

Functions:
- `BeginBlocker()` - Main entry point called each block
- `IncrementProposerCount()` - Increment daily proposal count
- `GetProposerCountForDay()` - Query proposal count for validator-day
- `GetAllProposerCountsForDay()` - Get all proposals for a day
- `GetTotalBlocksForDay()` - Calculate total blocks in a day

**Operation**: On each block, checks if proposer is eligible and increments their proposal counter for the current day.

---

### 3. **reward.go** - Reward Calculation & Aggregation
Location: `x/validatorbonus/keeper/reward.go`

Functions:
- `CalculateAndStoreDailyRewards()` - End-of-day reward calculation
- `StoreDailyRewardInternal()` - Store daily reward
- `GetDailyRewardInternal()` - Query daily reward
- `GetDailyRewardsForValidator()` - Get rewards for date range
- `CalculateAndStoreCycleRewards()` - End-of-cycle aggregation
- `StoreCycleRewardInternal()` - Store cycle reward
- `GetCycleRewardInternal()` - Query cycle reward
- `GetAllCycleRewardsForCycle()` - Get all rewards for a cycle
- `ResetProposerCountsForCycle()` - Clean up after aggregation

**Operation**: 
- Daily: Calculate and store rewards based on proposal shares
- Cycle: Aggregate daily rewards and reset counters

---

### 4. **end_block.go** - Day/Cycle Boundary Detection
Location: `x/validatorbonus/keeper/end_block.go`

Functions:
- `EndBlocker()` - Main entry point, detects boundaries
- `GetGenesisTime()` / `SetGenesisTime()` - Manage genesis time
- `GetLastProcessedDay()` / `SetLastProcessedDay()` - Track processed days

**Operation**: On each block end, checks if day or cycle boundary crossed and triggers calculations.

---

### 5. **query_validator_cycle_reward.go** - Single Validator Query
Location: `x/validatorbonus/keeper/query_validator_cycle_reward.go`

Function:
- `ValidatorCycleReward()` - Query implementation

**Input**: validatorAddress, cycle
**Output**: reward amount, isValidator flag

**Logic**: Check eligibility, lookup reward, return result

---

### 6. **query_cycle_rewards.go** - Cycle Rewards Query
Location: `x/validatorbonus/keeper/query_cycle_rewards.go`

Function:
- `CycleRewards()` - Query implementation

**Input**: cycle
**Output**: validators[], rewards[] (as JSON)

**Logic**: Get all rewards for cycle, return arrays

---

### 7. **module.go** - Module Integration (MODIFIED)
Location: `x/validatorbonus/module/module.go`

Changes:
- Updated `BeginBlock()` to call `k.BeginBlocker(ctx)`
- Updated `EndBlock()` to call `k.EndBlocker(ctx)`

**Impact**: Wires the reward logic into the block lifecycle.

---

### 8. **params.go** - Parameter Definitions (MODIFIED)
Location: `x/validatorbonus/types/params.go`

Functions:
- `NewParams()` - Create params with values
- `DefaultParams()` - Create default parameters
- `ParamSetPairs()` - Define parameter validators
- `Validate()` - Validate parameter values

**Parameters**:
- `TotalRewardPool` (string) - Total rewards to distribute
- `CycleDays` (uint64) - Days per cycle (default: 30)
- `TotalCycles` (uint64) - Number of cycles (default: 15)

---

### 9. **params.proto** - Proto Definition (MODIFIED)
Location: `proto/blockmazechain/validatorbonus/params.proto`

Added fields:
```protobuf
string total_reward_pool = 1;
uint64 cycle_days = 2;
uint64 total_cycles = 3;
```

---

## Data Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                         BLOCK LIFECYCLE                         │
└─────────────────────────────────────────────────────────────────┘

START OF BLOCK
    ↓
┌─────────────────────────────────────────────────────────────────┐
│ BeginBlocker (begin_block.go)                                   │
│ 1. Get block proposer                                           │
│ 2. Check if eligible validator                                 │
│ 3. Calculate day index                                          │
│ 4. Increment proposerCount/{validator}:{day}                    │
└─────────────────────────────────────────────────────────────────┘
    ↓
[Execute transactions]
    ↓
END OF BLOCK
    ↓
┌─────────────────────────────────────────────────────────────────┐
│ EndBlocker (end_block.go)                                       │
│ 1. Get current day                                              │
│ 2. Compare with lastProcessedDay                                │
│                                                                 │
│ IF day changed:                                                 │
│ ├─ Call CalculateAndStoreDailyRewards()                        │
│ │  └─ For each validator with proposals today:                 │
│ │     ├─ Calculate share: blocks / totalBlocks                 │
│ │     ├─ Calculate reward: share × dailyPool                   │
│ │     └─ Store dailyReward/{validator}:{day}                   │
│ │                                                              │
│ └─ IF cycle changed:                                           │
│    ├─ Call CalculateAndStoreCycleRewards()                     │
│    │  └─ For each eligible validator:                          │
│    │     ├─ Sum dailyReward for cycle                          │
│    │     ├─ Store cycleReward/{cycle}:{validator}              │
│    │     └─ ResetProposerCountsForCycle()                      │
│    └─ Update lastProcessedDay                                  │
└─────────────────────────────────────────────────────────────────┘
    ↓
NEXT BLOCK

┌─────────────────────────────────────────────────────────────────┐
│                      QUERY ENDPOINTS                            │
└─────────────────────────────────────────────────────────────────┘

QueryValidatorCycleReward:
  Input: {validatorAddress, cycle}
  ├─ Check eligibleValidator store
  └─ Lookup cycleReward/{cycle}:{validator}

QueryCycleRewards:
  Input: {cycle}
  └─ Iterate cycleReward with cycle prefix
     ├─ Build validators[]
     └─ Build rewards[]
```

## Storage Layout

```
KV Store Organization:

eligibleValidator/
  {validatorAddress} → EligibleValidator{address, joinTime}

proposerCount/
  {validatorAddress}:{day} → ProposerCount{address, day, count}
  [Reset after each cycle]

dailyReward/
  {validatorAddress}:{day} → DailyReward{address, day, amount}

cycleReward/
  {cycle}:{validatorAddress} → CycleReward{cycle, address, amount}

Module State:
  genesis_time → int64 (Unix timestamp)
  last_processed_day → uint64
```

## Key Design Decisions

### 1. **Composite Keys**
Used format: `key1:key2` for multi-field keys
- Enables prefix iteration
- Maintains human-readable format
- Efficient for lookups

### 2. **Lazy Evaluation**
- Calculations happen at day/cycle boundaries, not every block
- Reduces BeginBlocker overhead
- Concentrates computation in EndBlocker

### 3. **Hardcoded Defaults**
- Currently using hardcoded cycle_days=30, totalCycles=15
- Will use params after proto regeneration
- Allows testing before full deployment

### 4. **Separate Internal Methods**
- `*Internal()` suffixed methods for internal operations
- Avoid conflicts with generated keeper methods
- Clear separation of concerns

### 5. **No Automatic Reward Claiming**
- Module only calculates and stores rewards
- Does not automatically transfer tokens
- Allows flexible reward distribution mechanism

## Performance Characteristics

| Operation | Complexity | Frequency |
|-----------|-----------|-----------|
| BeginBlocker | O(1) | Every block |
| EndBlocker (no boundary) | O(1) | Every block |
| End-of-Day | O(V) | Once per day |
| End-of-Cycle | O(V) | Once per cycle |
| Query validator reward | O(1) | On demand |
| Query cycle rewards | O(V) | On demand |

Where V = number of validators

## Testing Considerations

1. **Unit Tests**: Test individual functions with mocked stores
2. **Integration Tests**: Test full flow over multiple days
3. **Edge Cases**: 
   - Single validator
   - Zero blocks
   - Cycle boundaries
   - Genesis block (day 0)
4. **Performance**: Test with 100+ validators

## Future Enhancements

1. **Day-First Key Indexing**
   - Current: `validator:day`
   - Proposed: `day:validator`
   - Enables faster day-based iteration

2. **Automatic Reward Distribution**
   - Claim mechanism
   - Batched transfers
   - Fee collection

3. **Dynamic Thresholds**
   - Adjustable eligibility window
   - Variable reward pools
   - Governance-controlled parameters

4. **Historical Queries**
   - Query rewards for multiple cycles
   - Leaderboards
   - Validator statistics

## Deployment Steps

1. Update proto files (already done)
2. Run `buf generate` to create Go types
3. Fix params.go references (already done)
4. Build: `go build ./...`
5. Test: `go test ./...`
6. Set genesis parameters and eligible validators
7. Start blockchain node

See `SETUP.md` for detailed deployment instructions.
