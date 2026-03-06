package types

import (
	"testing"

	"blockmazechain/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateEligibleValidator_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateEligibleValidator
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateEligibleValidator{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateEligibleValidator{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgUpdateEligibleValidator_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateEligibleValidator
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateEligibleValidator{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateEligibleValidator{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgDeleteEligibleValidator_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteEligibleValidator
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteEligibleValidator{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteEligibleValidator{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
