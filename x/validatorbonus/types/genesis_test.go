package types_test

import (
	"testing"

	"blockmazechain/x/validatorbonus/types"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{

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
				// this line is used by starport scaffolding # types/genesis/validField
				Params: types.DefaultParams(),
			},
			valid: true,
		},
		{
			desc: "duplicated eligibleValidator",
			genState: &types.GenesisState{
				EligibleValidatorList: []types.EligibleValidator{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated proposerCount",
			genState: &types.GenesisState{
				ProposerCountList: []types.ProposerCount{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated dailyReward",
			genState: &types.GenesisState{
				DailyRewardList: []types.DailyReward{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated cycleReward",
			genState: &types.GenesisState{
				CycleRewardList: []types.CycleReward{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
