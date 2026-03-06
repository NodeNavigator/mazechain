package validatorbonus_test

import (
	"testing"

	keepertest "blockmazechain/testutil/keeper"
	"blockmazechain/testutil/nullify"
	validatorbonus "blockmazechain/x/validatorbonus/module"
	"blockmazechain/x/validatorbonus/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		EligibleValidatorList: []types.EligibleValidator{
			{
				Index: "0",
			},
			{
				Index: "1",
			},
		},
		ProposerCountList: []types.ProposerCount{
			{
				Index: "0",
			},
			{
				Index: "1",
			},
		},
		DailyRewardList: []types.DailyReward{
			{
				Index: "0",
			},
			{
				Index: "1",
			},
		},
		CycleRewardList: []types.CycleReward{
			{
				Index: "0",
			},
			{
				Index: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.ValidatorbonusKeeper(t)
	validatorbonus.InitGenesis(ctx, k, genesisState)
	got := validatorbonus.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.EligibleValidatorList, got.EligibleValidatorList)
	require.ElementsMatch(t, genesisState.ProposerCountList, got.ProposerCountList)
	require.ElementsMatch(t, genesisState.DailyRewardList, got.DailyRewardList)
	require.ElementsMatch(t, genesisState.CycleRewardList, got.CycleRewardList)
	// this line is used by starport scaffolding # genesis/test/assert
}
