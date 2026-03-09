package keeper

import (
	"context"
	"encoding/hex"
	"strings"

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

	// Treat ValidatorAddress as a wallet address (account bech32 or 0x) and convert to operator (valoper) address.
	operatorAddr, err := walletToOperatorAddress(req.ValidatorAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid wallet address: %v", err)
	}

	// Check if validator is eligible (by validator operator address)
	_, found := k.GetEligibleValidatorByAddress(ctx, operatorAddr)
	if !found {
		// Validator is not eligible
		return &types.QueryValidatorCycleRewardResponse{
			IsEligible: false,
			Reward:     "0",
		}, nil
	}

	// Get the cycle reward for this validator
	rewardAmount, found := k.GetCycleRewardInternal(ctx, req.Cycle, operatorAddr)

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

// walletToOperatorAddress converts a user-facing wallet address into a validator operator (valoper) bech32 string.
// Supported inputs:
//   - validator operator address (cosmosvaloper...) -> returned as is
//   - account address (cosmos1...) -> converted via sdk.ValAddress(bytes)
//   - 0x-prefixed EVM address (20-byte hex) -> treated as account bytes then converted as above
func walletToOperatorAddress(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", status.Error(codes.InvalidArgument, "empty address")
	}

	// Already a valoper address – pass through
	if strings.HasPrefix(trimmed, "cosmosvaloper") {
		return trimmed, nil
	}

	var acc sdk.AccAddress

	switch {
	case strings.HasPrefix(trimmed, "0x") || strings.HasPrefix(trimmed, "0X"):
		// EVM-style address: 0x + 40 hex chars
		raw := strings.TrimPrefix(strings.TrimPrefix(trimmed, "0x"), "0X")
		bz, err := hex.DecodeString(raw)
		if err != nil {
			return "", err
		}
		acc = sdk.AccAddress(bz)
	default:
		// Assume bech32 account address (cosmos1...)
		var err error
		acc, err = sdk.AccAddressFromBech32(trimmed)
		if err != nil {
			return "", err
		}
	}

	val := sdk.ValAddress(acc)
	return val.String(), nil
}
