package keeper_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "blockmazechain/testutil/keeper"
	"blockmazechain/testutil/nullify"
	"blockmazechain/x/validatorbonus/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestCycleRewardQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.ValidatorbonusKeeper(t)
	msgs := createNCycleReward(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetCycleRewardRequest
		response *types.QueryGetCycleRewardResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetCycleRewardRequest{
				Id: msgs[0].Id,
			},
			response: &types.QueryGetCycleRewardResponse{CycleReward: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetCycleRewardRequest{
				Id: msgs[1].Id,
			},
			response: &types.QueryGetCycleRewardResponse{CycleReward: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetCycleRewardRequest{
				Id: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.CycleReward(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestCycleRewardQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.ValidatorbonusKeeper(t)
	msgs := createNCycleReward(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllCycleRewardRequest {
		return &types.QueryAllCycleRewardRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.CycleRewardAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.CycleReward), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.CycleReward),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.CycleRewardAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.CycleReward), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.CycleReward),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.CycleRewardAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.CycleReward),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.CycleRewardAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
