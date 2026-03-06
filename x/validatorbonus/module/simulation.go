package validatorbonus

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"blockmazechain/testutil/sample"
	validatorbonussimulation "blockmazechain/x/validatorbonus/simulation"
	"blockmazechain/x/validatorbonus/types"
)

// avoid unused import issue
var (
	_ = validatorbonussimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgCreateEligibleValidator = "op_weight_msg_eligible_validator"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateEligibleValidator int = 100

	opWeightMsgUpdateEligibleValidator = "op_weight_msg_eligible_validator"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateEligibleValidator int = 100

	opWeightMsgDeleteEligibleValidator = "op_weight_msg_eligible_validator"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteEligibleValidator int = 100

	opWeightMsgCreateProposerCount = "op_weight_msg_proposer_count"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateProposerCount int = 100

	opWeightMsgUpdateProposerCount = "op_weight_msg_proposer_count"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateProposerCount int = 100

	opWeightMsgDeleteProposerCount = "op_weight_msg_proposer_count"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteProposerCount int = 100

	opWeightMsgCreateDailyReward = "op_weight_msg_daily_reward"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateDailyReward int = 100

	opWeightMsgUpdateDailyReward = "op_weight_msg_daily_reward"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateDailyReward int = 100

	opWeightMsgDeleteDailyReward = "op_weight_msg_daily_reward"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteDailyReward int = 100

	opWeightMsgCreateCycleReward = "op_weight_msg_cycle_reward"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateCycleReward int = 100

	opWeightMsgUpdateCycleReward = "op_weight_msg_cycle_reward"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateCycleReward int = 100

	opWeightMsgDeleteCycleReward = "op_weight_msg_cycle_reward"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteCycleReward int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	validatorbonusGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		EligibleValidatorList: []types.EligibleValidator{
			{
				Creator: sample.AccAddress(),
				Index:   "0",
			},
			{
				Creator: sample.AccAddress(),
				Index:   "1",
			},
		},
		ProposerCountList: []types.ProposerCount{
			{
				Creator: sample.AccAddress(),
				Index:   "0",
			},
			{
				Creator: sample.AccAddress(),
				Index:   "1",
			},
		},
		DailyRewardList: []types.DailyReward{
			{
				Creator: sample.AccAddress(),
				Index:   "0",
			},
			{
				Creator: sample.AccAddress(),
				Index:   "1",
			},
		},
		CycleRewardList: []types.CycleReward{
			{
				Creator: sample.AccAddress(),
				Index:   "0",
			},
			{
				Creator: sample.AccAddress(),
				Index:   "1",
			},
		},
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&validatorbonusGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateEligibleValidator int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateEligibleValidator, &weightMsgCreateEligibleValidator, nil,
		func(_ *rand.Rand) {
			weightMsgCreateEligibleValidator = defaultWeightMsgCreateEligibleValidator
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateEligibleValidator,
		validatorbonussimulation.SimulateMsgCreateEligibleValidator(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateEligibleValidator int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateEligibleValidator, &weightMsgUpdateEligibleValidator, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateEligibleValidator = defaultWeightMsgUpdateEligibleValidator
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateEligibleValidator,
		validatorbonussimulation.SimulateMsgUpdateEligibleValidator(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteEligibleValidator int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteEligibleValidator, &weightMsgDeleteEligibleValidator, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteEligibleValidator = defaultWeightMsgDeleteEligibleValidator
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteEligibleValidator,
		validatorbonussimulation.SimulateMsgDeleteEligibleValidator(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateProposerCount int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateProposerCount, &weightMsgCreateProposerCount, nil,
		func(_ *rand.Rand) {
			weightMsgCreateProposerCount = defaultWeightMsgCreateProposerCount
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateProposerCount,
		validatorbonussimulation.SimulateMsgCreateProposerCount(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateProposerCount int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateProposerCount, &weightMsgUpdateProposerCount, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateProposerCount = defaultWeightMsgUpdateProposerCount
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateProposerCount,
		validatorbonussimulation.SimulateMsgUpdateProposerCount(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteProposerCount int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteProposerCount, &weightMsgDeleteProposerCount, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteProposerCount = defaultWeightMsgDeleteProposerCount
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteProposerCount,
		validatorbonussimulation.SimulateMsgDeleteProposerCount(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateDailyReward int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateDailyReward, &weightMsgCreateDailyReward, nil,
		func(_ *rand.Rand) {
			weightMsgCreateDailyReward = defaultWeightMsgCreateDailyReward
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateDailyReward,
		validatorbonussimulation.SimulateMsgCreateDailyReward(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateDailyReward int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateDailyReward, &weightMsgUpdateDailyReward, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateDailyReward = defaultWeightMsgUpdateDailyReward
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateDailyReward,
		validatorbonussimulation.SimulateMsgUpdateDailyReward(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteDailyReward int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteDailyReward, &weightMsgDeleteDailyReward, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteDailyReward = defaultWeightMsgDeleteDailyReward
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteDailyReward,
		validatorbonussimulation.SimulateMsgDeleteDailyReward(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateCycleReward int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateCycleReward, &weightMsgCreateCycleReward, nil,
		func(_ *rand.Rand) {
			weightMsgCreateCycleReward = defaultWeightMsgCreateCycleReward
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateCycleReward,
		validatorbonussimulation.SimulateMsgCreateCycleReward(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateCycleReward int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateCycleReward, &weightMsgUpdateCycleReward, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateCycleReward = defaultWeightMsgUpdateCycleReward
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateCycleReward,
		validatorbonussimulation.SimulateMsgUpdateCycleReward(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteCycleReward int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteCycleReward, &weightMsgDeleteCycleReward, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteCycleReward = defaultWeightMsgDeleteCycleReward
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteCycleReward,
		validatorbonussimulation.SimulateMsgDeleteCycleReward(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateEligibleValidator,
			defaultWeightMsgCreateEligibleValidator,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				validatorbonussimulation.SimulateMsgCreateEligibleValidator(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateEligibleValidator,
			defaultWeightMsgUpdateEligibleValidator,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				validatorbonussimulation.SimulateMsgUpdateEligibleValidator(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeleteEligibleValidator,
			defaultWeightMsgDeleteEligibleValidator,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				validatorbonussimulation.SimulateMsgDeleteEligibleValidator(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateProposerCount,
			defaultWeightMsgCreateProposerCount,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				validatorbonussimulation.SimulateMsgCreateProposerCount(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateProposerCount,
			defaultWeightMsgUpdateProposerCount,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				validatorbonussimulation.SimulateMsgUpdateProposerCount(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeleteProposerCount,
			defaultWeightMsgDeleteProposerCount,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				validatorbonussimulation.SimulateMsgDeleteProposerCount(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateDailyReward,
			defaultWeightMsgCreateDailyReward,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				validatorbonussimulation.SimulateMsgCreateDailyReward(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateDailyReward,
			defaultWeightMsgUpdateDailyReward,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				validatorbonussimulation.SimulateMsgUpdateDailyReward(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeleteDailyReward,
			defaultWeightMsgDeleteDailyReward,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				validatorbonussimulation.SimulateMsgDeleteDailyReward(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateCycleReward,
			defaultWeightMsgCreateCycleReward,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				validatorbonussimulation.SimulateMsgCreateCycleReward(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateCycleReward,
			defaultWeightMsgUpdateCycleReward,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				validatorbonussimulation.SimulateMsgUpdateCycleReward(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeleteCycleReward,
			defaultWeightMsgDeleteCycleReward,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				validatorbonussimulation.SimulateMsgDeleteCycleReward(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
