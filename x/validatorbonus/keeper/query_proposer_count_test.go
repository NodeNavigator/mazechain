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

func TestProposerCountQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.ValidatorbonusKeeper(t)
	msgs := createNProposerCount(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetProposerCountRequest
		response *types.QueryGetProposerCountResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetProposerCountRequest{
				Id: msgs[0].Id,
			},
			response: &types.QueryGetProposerCountResponse{ProposerCount: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetProposerCountRequest{
				Id: msgs[1].Id,
			},
			response: &types.QueryGetProposerCountResponse{ProposerCount: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetProposerCountRequest{
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
			response, err := keeper.ProposerCount(ctx, tc.request)
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

func TestProposerCountQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.ValidatorbonusKeeper(t)
	msgs := createNProposerCount(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllProposerCountRequest {
		return &types.QueryAllProposerCountRequest{
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
			resp, err := keeper.ProposerCountAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ProposerCount), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ProposerCount),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.ProposerCountAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ProposerCount), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ProposerCount),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.ProposerCountAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.ProposerCount),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.ProposerCountAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
