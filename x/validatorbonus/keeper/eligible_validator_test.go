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

func createNEligibleValidator(keeper keeper.Keeper, ctx context.Context, n int) []types.EligibleValidator {
	items := make([]types.EligibleValidator, n)
	for i := range items {
		items[i].Id = strconv.Itoa(i)
		items[i].ValidatorAddress = "addr" + strconv.Itoa(i) // Fix: provide a non-empty address for indexing
		keeper.SetEligibleValidator(ctx, items[i])
	}
	return items
}

func TestEligibleValidatorGet(t *testing.T) {
	keeper, ctx := keepertest.ValidatorbonusKeeper(t)
	items := createNEligibleValidator(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetEligibleValidator(ctx,
			item.Id,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestEligibleValidatorRemove(t *testing.T) {
	keeper, ctx := keepertest.ValidatorbonusKeeper(t)
	items := createNEligibleValidator(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveEligibleValidator(ctx,
			item.Id,
		)
		_, found := keeper.GetEligibleValidator(ctx,
			item.Id,
		)
		require.False(t, found)
	}
}

func TestEligibleValidatorGetAll(t *testing.T) {
	keeper, ctx := keepertest.ValidatorbonusKeeper(t)
	items := createNEligibleValidator(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllEligibleValidator(ctx)),
	)
}
