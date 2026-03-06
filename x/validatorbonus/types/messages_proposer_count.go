package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateProposerCount{}

func NewMsgCreateProposerCount(
	creator string,
	index string,
	validatorAddress string,
	day uint64,
	count uint64,

) *MsgCreateProposerCount {
	return &MsgCreateProposerCount{
		Creator:          creator,
		Index:            index,
		ValidatorAddress: validatorAddress,
		Day:              day,
		Count:            count,
	}
}

func (msg *MsgCreateProposerCount) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateProposerCount{}

func NewMsgUpdateProposerCount(
	creator string,
	index string,
	validatorAddress string,
	day uint64,
	count uint64,

) *MsgUpdateProposerCount {
	return &MsgUpdateProposerCount{
		Creator:          creator,
		Index:            index,
		ValidatorAddress: validatorAddress,
		Day:              day,
		Count:            count,
	}
}

func (msg *MsgUpdateProposerCount) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteProposerCount{}

func NewMsgDeleteProposerCount(
	creator string,
	index string,

) *MsgDeleteProposerCount {
	return &MsgDeleteProposerCount{
		Creator: creator,
		Index:   index,
	}
}

func (msg *MsgDeleteProposerCount) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
