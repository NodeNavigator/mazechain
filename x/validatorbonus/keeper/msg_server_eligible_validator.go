package keeper

import (
	"context"

	"blockmazechain/x/validatorbonus/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateEligibleValidator(goCtx context.Context, msg *types.MsgCreateEligibleValidator) (*types.MsgCreateEligibleValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetEligibleValidator(
		ctx,
		msg.Index,
	)
	if isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var eligibleValidator = types.EligibleValidator{
		Creator:          msg.Creator,
		Index:            msg.Index,
		ValidatorAddress: msg.ValidatorAddress,
		JoinTime:         msg.JoinTime,
	}

	k.SetEligibleValidator(
		ctx,
		eligibleValidator,
	)
	return &types.MsgCreateEligibleValidatorResponse{}, nil
}

func (k msgServer) UpdateEligibleValidator(goCtx context.Context, msg *types.MsgUpdateEligibleValidator) (*types.MsgUpdateEligibleValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetEligibleValidator(
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

	var eligibleValidator = types.EligibleValidator{
		Creator:          msg.Creator,
		Index:            msg.Index,
		ValidatorAddress: msg.ValidatorAddress,
		JoinTime:         msg.JoinTime,
	}

	k.SetEligibleValidator(ctx, eligibleValidator)

	return &types.MsgUpdateEligibleValidatorResponse{}, nil
}

func (k msgServer) DeleteEligibleValidator(goCtx context.Context, msg *types.MsgDeleteEligibleValidator) (*types.MsgDeleteEligibleValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetEligibleValidator(
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

	k.RemoveEligibleValidator(
		ctx,
		msg.Index,
	)

	return &types.MsgDeleteEligibleValidatorResponse{}, nil
}
