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

func createNCycleReward(keeper keeper.Keeper, ctx context.Context, n int) []types.CycleReward {
	items := make([]types.CycleReward, n)
	for i := range items {
		items[i].Id = strconv.Itoa(i)

		keeper.SetCycleReward(ctx, items[i])
	}
	return items
}

func TestCycleRewardGet(t *testing.T) {
	keeper, ctx := keepertest.ValidatorbonusKeeper(t)
	items := createNCycleReward(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetCycleReward(ctx,
			item.Id,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestCycleRewardRemove(t *testing.T) {
	keeper, ctx := keepertest.ValidatorbonusKeeper(t)
	items := createNCycleReward(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveCycleReward(ctx,
			item.Id,
		)
		_, found := keeper.GetCycleReward(ctx,
			item.Id,
		)
		require.False(t, found)
	}
}

func TestCycleRewardGetAll(t *testing.T) {
	keeper, ctx := keepertest.ValidatorbonusKeeper(t)
	items := createNCycleReward(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllCycleReward(ctx)),
	)
}
