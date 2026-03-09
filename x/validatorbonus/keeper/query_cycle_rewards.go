package keeper

import (
	"context"

	"blockmazechain/x/validatorbonus/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CycleRewards queries all validator rewards for a specific cycle.
//
// Logic:
// Iterate over cycleReward/{cycle} prefix and return:
// - validators[]: array of validator addresses
// - rewards[]: array of corresponding reward amounts (as JSON strings)
func (k Keeper) CycleRewards(goCtx context.Context, req *types.QueryCycleRewardsRequest) (*types.QueryCycleRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get all cycle rewards for this cycle
	cycleRewards := k.GetAllCycleRewardsForCycle(ctx, req.Cycle)

	// Build response arrays
	validators := make([]string, 0, len(cycleRewards))
	rewards := make([]string, 0, len(cycleRewards))

	for _, cr := range cycleRewards {
		validators = append(validators, cr.ValidatorAddress)
		rewards = append(rewards, cr.Amount)
	}

	return &types.QueryCycleRewardsResponse{
		Validators: validators,
		Rewards:    rewards,
	}, nil
}
