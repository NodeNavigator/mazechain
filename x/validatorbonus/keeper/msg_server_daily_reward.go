package keeper

import (
	"context"

	"blockmazechain/x/validatorbonus/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateDailyReward(goCtx context.Context, msg *types.MsgCreateDailyReward) (*types.MsgCreateDailyRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetDailyReward(
		ctx,
		msg.Index,
	)
	if isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var dailyReward = types.DailyReward{
		Creator:          msg.Creator,
		Index:            msg.Index,
		ValidatorAddress: msg.ValidatorAddress,
		Day:              msg.Day,
		Amount:           msg.Amount,
	}

	k.SetDailyReward(
		ctx,
		dailyReward,
	)
	return &types.MsgCreateDailyRewardResponse{}, nil
}

func (k msgServer) UpdateDailyReward(goCtx context.Context, msg *types.MsgUpdateDailyReward) (*types.MsgUpdateDailyRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetDailyReward(
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

	var dailyReward = types.DailyReward{
		Creator:          msg.Creator,
		Index:            msg.Index,
		ValidatorAddress: msg.ValidatorAddress,
		Day:              msg.Day,
		Amount:           msg.Amount,
	}

	k.SetDailyReward(ctx, dailyReward)

	return &types.MsgUpdateDailyRewardResponse{}, nil
}

func (k msgServer) DeleteDailyReward(goCtx context.Context, msg *types.MsgDeleteDailyReward) (*types.MsgDeleteDailyRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetDailyReward(
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

	k.RemoveDailyReward(
		ctx,
		msg.Index,
	)

	return &types.MsgDeleteDailyRewardResponse{}, nil
}
