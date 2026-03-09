package validatorbonus

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "blockmazechain/api/blockmazechain/validatorbonus"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "EligibleValidatorAll",
					Use:       "list-eligible-validator",
					Short:     "List all eligible-validator",
				},
				{
					RpcMethod:      "EligibleValidator",
					Use:            "show-eligible-validator [id]",
					Short:          "Shows a eligible-validator",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod: "ProposerCountAll",
					Use:       "list-proposer-count",
					Short:     "List all proposer-count",
				},
				{
					RpcMethod:      "ProposerCount",
					Use:            "show-proposer-count [id]",
					Short:          "Shows a proposer-count",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod: "DailyRewardAll",
					Use:       "list-daily-reward",
					Short:     "List all daily-reward",
				},
				{
					RpcMethod:      "DailyReward",
					Use:            "show-daily-reward [id]",
					Short:          "Shows a daily-reward",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod: "CycleRewardAll",
					Use:       "list-cycle-reward",
					Short:     "List all cycle-reward",
				},
				{
					RpcMethod:      "CycleReward",
					Use:            "show-cycle-reward [id]",
					Short:          "Shows a cycle-reward",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod:      "ValidatorCycleReward",
					Use:            "validator-cycle-reward [validator-address] [cycle] [reward] [is-validator]",
					Short:          "Query validator-cycle-reward",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "validatorAddress"}, {ProtoField: "cycle"}, {ProtoField: "reward"}, {ProtoField: "isValidator"}},
				},

				{
					RpcMethod:      "CycleRewards",
					Use:            "cycle-rewards [cycle]",
					Short:          "Query paginated cycle rewards for a cycle",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "cycle"}},
				},

				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateEligibleValidator",
					Use:            "create-eligible-validator [index] [validatorAddress] [joinTime]",
					Short:          "Create a new eligible-validator",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "validatorAddress"}, {ProtoField: "joinTime"}},
				},
				{
					RpcMethod:      "UpdateEligibleValidator",
					Use:            "update-eligible-validator [index] [validatorAddress] [joinTime]",
					Short:          "Update eligible-validator",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "validatorAddress"}, {ProtoField: "joinTime"}},
				},
				{
					RpcMethod:      "DeleteEligibleValidator",
					Use:            "delete-eligible-validator [index]",
					Short:          "Delete eligible-validator",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				{
					RpcMethod:      "CreateProposerCount",
					Use:            "create-proposer-count [index] [validatorAddress] [day] [count]",
					Short:          "Create a new proposer-count",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "validatorAddress"}, {ProtoField: "day"}, {ProtoField: "count"}},
				},
				{
					RpcMethod:      "UpdateProposerCount",
					Use:            "update-proposer-count [index] [validatorAddress] [day] [count]",
					Short:          "Update proposer-count",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "validatorAddress"}, {ProtoField: "day"}, {ProtoField: "count"}},
				},
				{
					RpcMethod:      "DeleteProposerCount",
					Use:            "delete-proposer-count [index]",
					Short:          "Delete proposer-count",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				{
					RpcMethod:      "CreateDailyReward",
					Use:            "create-daily-reward [index] [validatorAddress] [day] [amount]",
					Short:          "Create a new daily-reward",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "validatorAddress"}, {ProtoField: "day"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "UpdateDailyReward",
					Use:            "update-daily-reward [index] [validatorAddress] [day] [amount]",
					Short:          "Update daily-reward",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "validatorAddress"}, {ProtoField: "day"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "DeleteDailyReward",
					Use:            "delete-daily-reward [index]",
					Short:          "Delete daily-reward",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				{
					RpcMethod:      "CreateCycleReward",
					Use:            "create-cycle-reward [index] [cycle] [validatorAddress] [amount]",
					Short:          "Create a new cycle-reward",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "cycle"}, {ProtoField: "validatorAddress"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "UpdateCycleReward",
					Use:            "update-cycle-reward [index] [cycle] [validatorAddress] [amount]",
					Short:          "Update cycle-reward",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "cycle"}, {ProtoField: "validatorAddress"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "DeleteCycleReward",
					Use:            "delete-cycle-reward [index]",
					Short:          "Delete cycle-reward",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
