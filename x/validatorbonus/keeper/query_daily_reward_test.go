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

func TestDailyRewardQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.ValidatorbonusKeeper(t)
	msgs := createNDailyReward(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetDailyRewardRequest
		response *types.QueryGetDailyRewardResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetDailyRewardRequest{
				Id: msgs[0].Id,
			},
			response: &types.QueryGetDailyRewardResponse{DailyReward: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetDailyRewardRequest{
				Id: msgs[1].Id,
			},
			response: &types.QueryGetDailyRewardResponse{DailyReward: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetDailyRewardRequest{
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
			response, err := keeper.DailyReward(ctx, tc.request)
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

func TestDailyRewardQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.ValidatorbonusKeeper(t)
	msgs := createNDailyReward(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllDailyRewardRequest {
		return &types.QueryAllDailyRewardRequest{
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
			resp, err := keeper.DailyRewardAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.DailyReward), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.DailyReward),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.DailyRewardAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.DailyReward), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.DailyReward),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.DailyRewardAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.DailyReward),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.DailyRewardAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
