package validatorbonus

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"blockmazechain/x/validatorbonus/keeper"
	"blockmazechain/x/validatorbonus/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the eligibleValidator
	for _, elem := range genState.EligibleValidatorList {
		k.SetEligibleValidator(ctx, elem)
	}
	// Set all the proposerCount
	for _, elem := range genState.ProposerCountList {
		k.SetProposerCount(ctx, elem)
	}
	// Set all the dailyReward
	for _, elem := range genState.DailyRewardList {
		k.SetDailyReward(ctx, elem)
	}
	// Set all the cycleReward
	for _, elem := range genState.CycleRewardList {
		k.SetCycleReward(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.EligibleValidatorList = k.GetAllEligibleValidator(ctx)
	genesis.ProposerCountList = k.GetAllProposerCount(ctx)
	genesis.DailyRewardList = k.GetAllDailyReward(ctx)
	genesis.CycleRewardList = k.GetAllCycleReward(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
