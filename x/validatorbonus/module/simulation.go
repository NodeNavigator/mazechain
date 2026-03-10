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
	_ = validatorbonussimulation.SimulateMsgCreateEligibleValidator
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgCreateEligibleValidator = "op_weight_msg_eligible_validator"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateEligibleValidator int = 100
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
				Id:      "0",
			},
			{
				Creator: sample.AccAddress(),
				Id:      "1",
			},
		},
		ProposerCountList: []types.ProposerCount{
			{
				Creator: sample.AccAddress(),
				Id:      "0",
			},
			{
				Creator: sample.AccAddress(),
				Id:      "1",
			},
		},
		DailyRewardList: []types.DailyReward{
			{
				Creator: sample.AccAddress(),
				Id:      "0",
			},
			{
				Creator: sample.AccAddress(),
				Id:      "1",
			},
		},
		CycleRewardList: []types.CycleReward{
			{
				Creator: sample.AccAddress(),
				Id:      "0",
			},
			{
				Creator: sample.AccAddress(),
				Id:      "1",
			},
		},
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
	}
}
