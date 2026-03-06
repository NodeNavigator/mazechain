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

func TestDailyRewardMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.ValidatorbonusKeeper(t)
	srv := keeper.NewMsgServerImpl(k)
	creator := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateDailyReward{Creator: creator,
			Index: strconv.Itoa(i),
		}
		_, err := srv.CreateDailyReward(ctx, expected)
		require.NoError(t, err)
		rst, found := k.GetDailyReward(ctx,
			expected.Index,
		)
		require.True(t, found)
		require.Equal(t, expected.Creator, rst.Creator)
	}
}

func TestDailyRewardMsgServerUpdate(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgUpdateDailyReward
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateDailyReward{Creator: creator,
				Index: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgUpdateDailyReward{Creator: "B",
				Index: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateDailyReward{Creator: creator,
				Index: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.ValidatorbonusKeeper(t)
			srv := keeper.NewMsgServerImpl(k)
			expected := &types.MsgCreateDailyReward{Creator: creator,
				Index: strconv.Itoa(0),
			}
			_, err := srv.CreateDailyReward(ctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateDailyReward(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetDailyReward(ctx,
					expected.Index,
				)
				require.True(t, found)
				require.Equal(t, expected.Creator, rst.Creator)
			}
		})
	}
}

func TestDailyRewardMsgServerDelete(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgDeleteDailyReward
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteDailyReward{Creator: creator,
				Index: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDeleteDailyReward{Creator: "B",
				Index: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteDailyReward{Creator: creator,
				Index: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.ValidatorbonusKeeper(t)
			srv := keeper.NewMsgServerImpl(k)

			_, err := srv.CreateDailyReward(ctx, &types.MsgCreateDailyReward{Creator: creator,
				Index: strconv.Itoa(0),
			})
			require.NoError(t, err)
			_, err = srv.DeleteDailyReward(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetDailyReward(ctx,
					tc.request.Index,
				)
				require.False(t, found)
			}
		})
	}
}
