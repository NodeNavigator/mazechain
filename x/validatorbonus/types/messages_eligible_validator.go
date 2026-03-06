package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateEligibleValidator{}

func NewMsgCreateEligibleValidator(
	creator string,
	index string,
	validatorAddress string,
	joinTime int32,

) *MsgCreateEligibleValidator {
	return &MsgCreateEligibleValidator{
		Creator:          creator,
		Index:            index,
		ValidatorAddress: validatorAddress,
		JoinTime:         joinTime,
	}
}

func (msg *MsgCreateEligibleValidator) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateEligibleValidator{}

func NewMsgUpdateEligibleValidator(
	creator string,
	index string,
	validatorAddress string,
	joinTime int32,

) *MsgUpdateEligibleValidator {
	return &MsgUpdateEligibleValidator{
		Creator:          creator,
		Index:            index,
		ValidatorAddress: validatorAddress,
		JoinTime:         joinTime,
	}
}

func (msg *MsgUpdateEligibleValidator) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteEligibleValidator{}

func NewMsgDeleteEligibleValidator(
	creator string,
	index string,

) *MsgDeleteEligibleValidator {
	return &MsgDeleteEligibleValidator{
		Creator: creator,
		Index:   index,
	}
}

func (msg *MsgDeleteEligibleValidator) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
