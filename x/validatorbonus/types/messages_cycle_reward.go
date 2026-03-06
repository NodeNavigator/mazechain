package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateCycleReward{}

func NewMsgCreateCycleReward(
	creator string,
	index string,
	cycle uint64,
	validatorAddress string,
	amount string,

) *MsgCreateCycleReward {
	return &MsgCreateCycleReward{
		Creator:          creator,
		Index:            index,
		Cycle:            cycle,
		ValidatorAddress: validatorAddress,
		Amount:           amount,
	}
}

func (msg *MsgCreateCycleReward) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateCycleReward{}

func NewMsgUpdateCycleReward(
	creator string,
	index string,
	cycle uint64,
	validatorAddress string,
	amount string,

) *MsgUpdateCycleReward {
	return &MsgUpdateCycleReward{
		Creator:          creator,
		Index:            index,
		Cycle:            cycle,
		ValidatorAddress: validatorAddress,
		Amount:           amount,
	}
}

func (msg *MsgUpdateCycleReward) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteCycleReward{}

func NewMsgDeleteCycleReward(
	creator string,
	index string,

) *MsgDeleteCycleReward {
	return &MsgDeleteCycleReward{
		Creator: creator,
		Index:   index,
	}
}

func (msg *MsgDeleteCycleReward) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
