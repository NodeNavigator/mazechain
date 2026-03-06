package types

import (
	"testing"

	"blockmazechain/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateCycleReward_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateCycleReward
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateCycleReward{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateCycleReward{
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

func TestMsgUpdateCycleReward_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateCycleReward
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateCycleReward{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateCycleReward{
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

func TestMsgDeleteCycleReward_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteCycleReward
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteCycleReward{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteCycleReward{
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
