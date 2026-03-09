package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		EligibleValidatorList: []EligibleValidator{},
		ProposerCountList:     []ProposerCount{},
		DailyRewardList:       []DailyReward{},
		CycleRewardList:       []CycleReward{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in eligibleValidator
	eligibleValidatorIndexMap := make(map[string]struct{})

	for _, elem := range gs.EligibleValidatorList {
		index := string(EligibleValidatorKey(elem.Id))
		if _, ok := eligibleValidatorIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for eligibleValidator")
		}
		eligibleValidatorIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in proposerCount
	proposerCountIndexMap := make(map[string]struct{})

	for _, elem := range gs.ProposerCountList {
		index := string(ProposerCountKey(elem.Id))
		if _, ok := proposerCountIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for proposerCount")
		}
		proposerCountIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in dailyReward
	dailyRewardIndexMap := make(map[string]struct{})

	for _, elem := range gs.DailyRewardList {
		index := string(DailyRewardKey(elem.Id))
		if _, ok := dailyRewardIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for dailyReward")
		}
		dailyRewardIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in cycleReward
	cycleRewardIndexMap := make(map[string]struct{})

	for _, elem := range gs.CycleRewardList {
		index := string(CycleRewardKey(elem.Id))
		if _, ok := cycleRewardIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for cycleReward")
		}
		cycleRewardIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
