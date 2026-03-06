package types

import (
	"testing"

	"blockmazechain/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateDailyReward_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateDailyReward
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateDailyReward{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateDailyReward{
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

func TestMsgUpdateDailyReward_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateDailyReward
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateDailyReward{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateDailyReward{
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

func TestMsgDeleteDailyReward_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteDailyReward
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteDailyReward{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteDailyReward{
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
