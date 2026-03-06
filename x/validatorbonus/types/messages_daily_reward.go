package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateDailyReward{}

func NewMsgCreateDailyReward(
	creator string,
	index string,
	validatorAddress string,
	day uint64,
	amount string,

) *MsgCreateDailyReward {
	return &MsgCreateDailyReward{
		Creator:          creator,
		Index:            index,
		ValidatorAddress: validatorAddress,
		Day:              day,
		Amount:           amount,
	}
}

func (msg *MsgCreateDailyReward) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateDailyReward{}

func NewMsgUpdateDailyReward(
	creator string,
	index string,
	validatorAddress string,
	day uint64,
	amount string,

) *MsgUpdateDailyReward {
	return &MsgUpdateDailyReward{
		Creator:          creator,
		Index:            index,
		ValidatorAddress: validatorAddress,
		Day:              day,
		Amount:           amount,
	}
}

func (msg *MsgUpdateDailyReward) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteDailyReward{}

func NewMsgDeleteDailyReward(
	creator string,
	index string,

) *MsgDeleteDailyReward {
	return &MsgDeleteDailyReward{
		Creator: creator,
		Index:   index,
	}
}

func (msg *MsgDeleteDailyReward) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
