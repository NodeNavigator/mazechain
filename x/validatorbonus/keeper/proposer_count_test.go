package keeper_test

import (
	"context"
	"strconv"
	"testing"

	keepertest "blockmazechain/testutil/keeper"
	"blockmazechain/testutil/nullify"
	"blockmazechain/x/validatorbonus/keeper"
	"blockmazechain/x/validatorbonus/types"

	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNProposerCount(keeper keeper.Keeper, ctx context.Context, n int) []types.ProposerCount {
	items := make([]types.ProposerCount, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)

		keeper.SetProposerCount(ctx, items[i])
	}
	return items
}

func TestProposerCountGet(t *testing.T) {
	keeper, ctx := keepertest.ValidatorbonusKeeper(t)
	items := createNProposerCount(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetProposerCount(ctx,
			item.Index,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestProposerCountRemove(t *testing.T) {
	keeper, ctx := keepertest.ValidatorbonusKeeper(t)
	items := createNProposerCount(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveProposerCount(ctx,
			item.Index,
		)
		_, found := keeper.GetProposerCount(ctx,
			item.Index,
		)
		require.False(t, found)
	}
}

func TestProposerCountGetAll(t *testing.T) {
	keeper, ctx := keepertest.ValidatorbonusKeeper(t)
	items := createNProposerCount(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllProposerCount(ctx)),
	)
}
