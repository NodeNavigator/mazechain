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

func createNDailyReward(keeper keeper.Keeper, ctx context.Context, n int) []types.DailyReward {
	items := make([]types.DailyReward, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)

		keeper.SetDailyReward(ctx, items[i])
	}
	return items
}

func TestDailyRewardGet(t *testing.T) {
	keeper, ctx := keepertest.ValidatorbonusKeeper(t)
	items := createNDailyReward(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetDailyReward(ctx,
			item.Index,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestDailyRewardRemove(t *testing.T) {
	keeper, ctx := keepertest.ValidatorbonusKeeper(t)
	items := createNDailyReward(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveDailyReward(ctx,
			item.Index,
		)
		_, found := keeper.GetDailyReward(ctx,
			item.Index,
		)
		require.False(t, found)
	}
}

func TestDailyRewardGetAll(t *testing.T) {
	keeper, ctx := keepertest.ValidatorbonusKeeper(t)
	items := createNDailyReward(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllDailyReward(ctx)),
	)
}
