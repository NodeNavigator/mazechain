package keeper

import (
	"context"

	"blockmazechain/x/validatorbonus/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ValidatorCycleReward queries the reward for a specific validator in a specific cycle.
//
// Logic:
// 1. Check if validator exists in eligibleValidator store
// 2. If not:
//   - return isValidator = false, reward = "0"
//
// 3. If exists:
//   - lookup cycleReward/{cycle}/{validator}
//   - return isValidator = true, reward
//
// NOTE: This is a stub implementation. After regenerating proto files with proper
// response fields (Reward and IsValidator), this will be fully functional.
func (k Keeper) ValidatorCycleReward(goCtx context.Context, req *types.QueryValidatorCycleRewardRequest) (*types.QueryValidatorCycleRewardResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if validator is eligible
	_, found := k.GetEligibleValidator(ctx, req.ValidatorAddress)
	if !found {
		// Validator is not eligible
		return &types.QueryValidatorCycleRewardResponse{
			IsEligible: false,
			Reward:     "0",
		}, nil
	}

	// Get the cycle reward for this validator
	rewardAmount, found := k.GetCycleRewardInternal(ctx, req.Cycle, req.ValidatorAddress)

	if !found {
		// No reward found for this validator in this cycle
		return &types.QueryValidatorCycleRewardResponse{
			IsEligible: true,
			Reward:     "0",
		}, nil
	}

	// Return response with reward amount
	return &types.QueryValidatorCycleRewardResponse{
		IsEligible: true,
		Reward:     rewardAmount.String(),
	}, nil
}
