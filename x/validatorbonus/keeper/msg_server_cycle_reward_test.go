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

func TestCycleRewardMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.ValidatorbonusKeeper(t)
	srv := keeper.NewMsgServerImpl(k)
	creator := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateCycleReward{Creator: creator,
			Index: strconv.Itoa(i),
		}
		_, err := srv.CreateCycleReward(ctx, expected)
		require.NoError(t, err)
		rst, found := k.GetCycleReward(ctx,
			expected.Index,
		)
		require.True(t, found)
		require.Equal(t, expected.Creator, rst.Creator)
	}
}

func TestCycleRewardMsgServerUpdate(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgUpdateCycleReward
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateCycleReward{Creator: creator,
				Index: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgUpdateCycleReward{Creator: "B",
				Index: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateCycleReward{Creator: creator,
				Index: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.ValidatorbonusKeeper(t)
			srv := keeper.NewMsgServerImpl(k)
			expected := &types.MsgCreateCycleReward{Creator: creator,
				Index: strconv.Itoa(0),
			}
			_, err := srv.CreateCycleReward(ctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateCycleReward(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetCycleReward(ctx,
					expected.Index,
				)
				require.True(t, found)
				require.Equal(t, expected.Creator, rst.Creator)
			}
		})
	}
}

func TestCycleRewardMsgServerDelete(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgDeleteCycleReward
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteCycleReward{Creator: creator,
				Index: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDeleteCycleReward{Creator: "B",
				Index: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteCycleReward{Creator: creator,
				Index: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.ValidatorbonusKeeper(t)
			srv := keeper.NewMsgServerImpl(k)

			_, err := srv.CreateCycleReward(ctx, &types.MsgCreateCycleReward{Creator: creator,
				Index: strconv.Itoa(0),
			})
			require.NoError(t, err)
			_, err = srv.DeleteCycleReward(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetCycleReward(ctx,
					tc.request.Index,
				)
				require.False(t, found)
			}
		})
	}
}
