package keeper_test

import (
	"strconv"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	keepertest "blockmazechain/testutil/keeper"
	"blockmazechain/x/validatorbonus/keeper"
	"blockmazechain/x/validatorbonus/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestEligibleValidatorMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.ValidatorbonusKeeper(t)
	srv := keeper.NewMsgServerImpl(k)
	creator := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateEligibleValidator{Creator: creator,
			Index: strconv.Itoa(i),
		}
		_, err := srv.CreateEligibleValidator(ctx, expected)
		require.NoError(t, err)
		rst, found := k.GetEligibleValidator(ctx,
			expected.Index,
		)
		require.True(t, found)
		require.Equal(t, expected.Creator, rst.Creator)
	}
}

func TestEligibleValidatorMsgServerUpdate(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgUpdateEligibleValidator
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateEligibleValidator{Creator: creator,
				Index: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgUpdateEligibleValidator{Creator: "B",
				Index: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateEligibleValidator{Creator: creator,
				Index: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.ValidatorbonusKeeper(t)
			srv := keeper.NewMsgServerImpl(k)
			expected := &types.MsgCreateEligibleValidator{Creator: creator,
				Index: strconv.Itoa(0),
			}
			_, err := srv.CreateEligibleValidator(ctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateEligibleValidator(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetEligibleValidator(ctx,
					expected.Index,
				)
				require.True(t, found)
				require.Equal(t, expected.Creator, rst.Creator)
			}
		})
	}
}

func TestEligibleValidatorMsgServerDelete(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgDeleteEligibleValidator
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteEligibleValidator{Creator: creator,
				Index: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDeleteEligibleValidator{Creator: "B",
				Index: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteEligibleValidator{Creator: creator,
				Index: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.ValidatorbonusKeeper(t)
			srv := keeper.NewMsgServerImpl(k)

			_, err := srv.CreateEligibleValidator(ctx, &types.MsgCreateEligibleValidator{Creator: creator,
				Index: strconv.Itoa(0),
			})
			require.NoError(t, err)
			_, err = srv.DeleteEligibleValidator(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetEligibleValidator(ctx,
					tc.request.Index,
				)
				require.False(t, found)
			}
		})
	}
}
