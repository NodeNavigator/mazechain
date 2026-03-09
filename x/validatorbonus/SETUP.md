# Setup and Build Instructions

## Overview

This guide explains how to finalize the validator bonus module implementation, which includes regenerating protocol buffer files and ensuring all dependencies are properly configured.

## Step 1: Regenerate Protocol Buffer Files

The module proto files have been updated with new fields but need to be regenerated to update the Go types. Follow these steps:

### Update Proto Files (Already Done)

The following proto files have been updated:
- `proto/blockmazechain/validatorbonus/params.proto` - Added Params fields
- `proto/blockmazechain/validatorbonus/query.proto` - May need response field updates (optional)

### Regenerate Go Types

From the project root directory:

```bash
# Generate Go code from proto files
buf generate

# Or if using docker:
docker run --rm -v $(pwd):/workspace -w /workspace bufbuild/buf generate
```

This command will:
1. Regenerate `x/validatorbonus/types/params.pb.go` with the new fields
2. Update `x/validatorbonus/types/query.pb.go` if needed
3. Regenerate all related grpc and gw files

### Build and Test

```bash
# Build the module
go build ./x/validatorbonus/...

# Run tests
go test ./x/validatorbonus/...

# Build entire blockchain
make build
```

## Step 2: Initialize Genesis and Parameters

During chain initialization, set up the module parameters and eligible validators.

### Example Genesis Setup

Add to `genesis.json`:

```json
{
  "app_state": {
    "validatorbonus": {
      "params": {
        "total_reward_pool": "1000000000",  // Total tokens to distribute
        "cycle_days": 30,                    // Days per cycle
        "total_cycles": 15                   // Total number of cycles
      },
      "eligible_validators": [
        {
          "validator_address": "cosmosvaloper1...",
          "join_time": 0  // Block height or time when eligible
        }
        // ... more validators
      ]
    }
  }
}
```

Or use a setup function during chain initialization:

```go
// In your genesis or setup code
import "blockmazechain/x/validatorbonus/keeper"
import "blockmazechain/x/validatorbonus/types"

func setupValidatorBonus(ctx sdk.Context, k keeper.Keeper) {
    // Set parameters
    params := types.NewParams(
        "1000000000",  // total_reward_pool
        30,            // cycle_days
        15,            // total_cycles
    )
    k.SetParams(ctx, params)
    
    // Set genesis time
    genesisTime := ctx.BlockTime()
    k.SetGenesisTime(ctx, genesisTime.Unix())
    k.SetLastProcessedDay(ctx, 0)
    
    // Register eligible validators
    for _, validatorAddr := range eligibleValidators {
        k.SetEligibleValidator(ctx, types.EligibleValidator{
            ValidatorAddress: validatorAddr,
            JoinTime: 0,
        })
    }
}
```

## Step 3: Module Features

### BeginBlocker

- Called at the start of each block
- Tracks block proposers and increments their count for the current day
- Only tracks eligible validators
- O(1) operation - very efficient

### EndBlocker

- Called at the end of each block
- Detects day transitions and triggers daily reward calculations
- Detects cycle transitions and triggers cycle aggregation
- Performs calculations only when boundaries are crossed

### Queries

After proto regeneration, the following queries will be available:

#### validator_cycle_reward

Get the reward for a specific validator in a specific cycle:

```bash
blockmazechaind query validatorbonus validator_cycle_reward \
  --validator-address cosmosvaloper1... \
  --cycle 1
```

#### cycle_rewards

Get all validator rewards for a specific cycle:

```bash
blockmazechaind query validatorbonus cycle_rewards \
  --cycle 1
```

## Step 4: Optional Enhancements

### 1. Storage Optimization

The current implementation iterates all proposer counts when calculating daily rewards. For better performance with many validators, consider indexing by day:

```go
// Composite key structure: day:validator
// Allows prefix iteration on day
index := fmt.Sprintf("%d:%s", day, validatorAddr)
```

### 2. Governance Integration

Make parameters updatable via governance:

```bash
# Create a governance proposal to update parameters
blockmazechaind tx gov submit-proposal param-change \
  <proposal.json>
```

### 3. Monitoring and Logging

The module logs important events:

```go
k.Logger().Info("End of day reached", "day", previousDay)
k.Logger().Info("End of cycle reached", "cycle", previousCycle)
k.Logger().Error("failed to calculate daily rewards", ...)
```

Monitor these logs to ensure correct operation.

## Step 5: Testing

### Unit Tests

Create test files in `x/validatorbonus/keeper/`:

```go
// Example test
func TestBeginBlocker(t *testing.T) {
    keeper, ctx := setupKeeper(t)
    
    // Set genesis time
    genesisTime := ctx.BlockTime()
    keeper.SetGenesisTime(ctx, genesisTime.Unix())
    
    // Add eligible validator
    keeper.SetEligibleValidator(ctx, types.EligibleValidator{
        ValidatorAddress: "cosmosvaloper1...",
        JoinTime: 0,
    })
    
    // Simulate block
    err := keeper.BeginBlocker(ctx)
    require.NoError(t, err)
    
    // Verify proposer count incremented
    count := keeper.GetProposerCountForDay(ctx, "cosmosvaloper1...", 0)
    require.Equal(t, uint64(1), count)
}
```

### Integration Tests

Test the full flow over multiple days and cycles:

```go
func TestMultiDayRewards(t *testing.T) {
    keeper, ctx := setupKeeper(t)
    
    // Setup...
    
    // Simulate multiple days
    for day := 0; day < 30; day++ {
        // Process blocks for this day
        // ...
        
        // Check rewards calculated correctly
    }
}
```

## Step 6: Deployment

### Pre-Deployment Checklist

- [ ] Proto files regenerated successfully
- [ ] Code compiles without errors: `go build ./...`
- [ ] Tests pass: `go test ./...`
- [ ] Module registered in app.go
- [ ] BeginBlocker and EndBlocker wired in module
- [ ] Genesis parameters set correctly
- [ ] Eligible validators registered
- [ ] Genesis time initialized

### Production Considerations

1. **Data Backup**: Before any major update, backup the state
2. **Testing**: Test on testnet before mainnet
3. **Monitoring**: Monitor logs and KV store size
4. **Performance**: With many validators, consider optimization from Step 4.1
5. **Governance**: Update parameters via governance if needed

## Troubleshooting

### Issue: Proto Files Not Regenerating

**Solution**: Check `buf.yaml` configuration and ensure all proto dependencies are available:

```bash
buf mod update
buf generate
```

### Issue: Module Parameters Not Set

**Solution**: Ensure `SetParams` is called during genesis:

```go
// In InitGenesis
k.SetParams(ctx, genState.Params)
```

### Issue: No Rewards Calculated

**Solution**: Check:
1. Genesis time is set: `k.GetGenesisTime(ctx)`
2. Eligible validators registered: `k.GetAllEligibleValidator(ctx)`
3. Block time is advancing: `ctx.BlockTime()`
4. Proposers are being tracked: Check proposerCount store

### Issue: High Memory/KV Store Usage

**Solution**: 
1. Reset proposerCount after each cycle (automatic in current implementation)
2. Consider archive/cleanup of old daily/cycle rewards
3. Implement periodic snapshots

## Reference Documentation

- [Cosmos SDK Module Development](https://docs.cosmos.network/v0.53/building-modules/intro.html)
- [Protocol Buffers Guide](https://developers.google.com/protocol-buffers)
- [Buf Documentation](https://buf.build/docs/)
