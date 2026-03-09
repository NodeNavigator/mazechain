package keeper

import (
	"context"

	"blockmazechain/x/validatorbonus/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateCycleReward(goCtx context.Context, msg *types.MsgCreateCycleReward) (*types.MsgCreateCycleRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetCycleReward(
		ctx,
		msg.Index,
	)
	if isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var cycleReward = types.CycleReward{
		Creator:          msg.Creator,
		Id:               msg.Index,
		Cycle:            msg.Cycle,
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           msg.Amount,
	}

	k.SetCycleReward(
		ctx,
		cycleReward,
	)
	return &types.MsgCreateCycleRewardResponse{}, nil
}

func (k msgServer) UpdateCycleReward(goCtx context.Context, msg *types.MsgUpdateCycleReward) (*types.MsgUpdateCycleRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetCycleReward(
		ctx,
		msg.Index,
	)
	if !isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var cycleReward = types.CycleReward{
		Creator:          msg.Creator,
		Id:               msg.Index,
		Cycle:            msg.Cycle,
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           msg.Amount,
	}

	k.SetCycleReward(ctx, cycleReward)

	return &types.MsgUpdateCycleRewardResponse{}, nil
}

func (k msgServer) DeleteCycleReward(goCtx context.Context, msg *types.MsgDeleteCycleReward) (*types.MsgDeleteCycleRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetCycleReward(
		ctx,
		msg.Index,
	)
	if !isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveCycleReward(
		ctx,
		msg.Index,
	)

	return &types.MsgDeleteCycleRewardResponse{}, nil
}
