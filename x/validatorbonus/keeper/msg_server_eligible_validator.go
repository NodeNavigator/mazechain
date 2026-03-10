package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"blockmazechain/x/validatorbonus/types"
)

// CreateEligibleValidator adds a validator to the eligible list.
// Only the module authority (governance or upgrade handler) can call this.
func (k msgServer) CreateEligibleValidator(goCtx context.Context, msg *types.MsgCreateEligibleValidator) (*types.MsgCreateEligibleValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetEligibleValidator(ctx, msg.Index)
	if isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "validator with this index already registered")
	}

	var eligibleValidator = types.EligibleValidator{
		Creator:          msg.Creator,
		Id:               msg.Index,
		ValidatorAddress: msg.ValidatorAddress,
		JoinTime:         msg.JoinTime,
	}

	k.SetEligibleValidator(ctx, eligibleValidator)
	return &types.MsgCreateEligibleValidatorResponse{}, nil
}
