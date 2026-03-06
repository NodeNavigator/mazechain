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

func TestProposerCountMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.ValidatorbonusKeeper(t)
	srv := keeper.NewMsgServerImpl(k)
	creator := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateProposerCount{Creator: creator,
			Index: strconv.Itoa(i),
		}
		_, err := srv.CreateProposerCount(ctx, expected)
		require.NoError(t, err)
		rst, found := k.GetProposerCount(ctx,
			expected.Index,
		)
		require.True(t, found)
		require.Equal(t, expected.Creator, rst.Creator)
	}
}

func TestProposerCountMsgServerUpdate(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgUpdateProposerCount
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateProposerCount{Creator: creator,
				Index: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgUpdateProposerCount{Creator: "B",
				Index: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateProposerCount{Creator: creator,
				Index: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.ValidatorbonusKeeper(t)
			srv := keeper.NewMsgServerImpl(k)
			expected := &types.MsgCreateProposerCount{Creator: creator,
				Index: strconv.Itoa(0),
			}
			_, err := srv.CreateProposerCount(ctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateProposerCount(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetProposerCount(ctx,
					expected.Index,
				)
				require.True(t, found)
				require.Equal(t, expected.Creator, rst.Creator)
			}
		})
	}
}

func TestProposerCountMsgServerDelete(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgDeleteProposerCount
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteProposerCount{Creator: creator,
				Index: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDeleteProposerCount{Creator: "B",
				Index: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteProposerCount{Creator: creator,
				Index: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.ValidatorbonusKeeper(t)
			srv := keeper.NewMsgServerImpl(k)

			_, err := srv.CreateProposerCount(ctx, &types.MsgCreateProposerCount{Creator: creator,
				Index: strconv.Itoa(0),
			})
			require.NoError(t, err)
			_, err = srv.DeleteProposerCount(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetProposerCount(ctx,
					tc.request.Index,
				)
				require.False(t, found)
			}
		})
	}
}
