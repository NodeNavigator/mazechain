package keeper

import (
	"context"

	"blockmazechain/x/validatorbonus/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateProposerCount(goCtx context.Context, msg *types.MsgCreateProposerCount) (*types.MsgCreateProposerCountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetProposerCount(
		ctx,
		msg.Index,
	)
	if isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var proposerCount = types.ProposerCount{
		Creator:          msg.Creator,
		Id:               msg.Index,
		ValidatorAddress: msg.ValidatorAddress,
		Day:              msg.Day,
		Count:            msg.Count,
	}

	k.SetProposerCount(
		ctx,
		proposerCount,
	)
	return &types.MsgCreateProposerCountResponse{}, nil
}

func (k msgServer) UpdateProposerCount(goCtx context.Context, msg *types.MsgUpdateProposerCount) (*types.MsgUpdateProposerCountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetProposerCount(
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

	var proposerCount = types.ProposerCount{
		Creator:          msg.Creator,
		Id:               msg.Index,
		ValidatorAddress: msg.ValidatorAddress,
		Day:              msg.Day,
		Count:            msg.Count,
	}

	k.SetProposerCount(ctx, proposerCount)

	return &types.MsgUpdateProposerCountResponse{}, nil
}

func (k msgServer) DeleteProposerCount(goCtx context.Context, msg *types.MsgDeleteProposerCount) (*types.MsgDeleteProposerCountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetProposerCount(
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

	k.RemoveProposerCount(
		ctx,
		msg.Index,
	)

	return &types.MsgDeleteProposerCountResponse{}, nil
}
