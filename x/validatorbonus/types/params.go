package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyTotalRewardPool = []byte("TotalRewardPool")
	KeyCycleDays       = []byte("CycleDays")
	KeyTotalCycles     = []byte("TotalCycles")
	KeySecondsPerDay   = []byte("SecondsPerDay")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(totalRewardPool string, cycleDays, totalCycles, secondsPerDay uint64) Params {
	return Params{
		TotalRewardPool: totalRewardPool,
		CycleDays:       cycleDays,
		TotalCycles:     totalCycles,
		SecondsPerDay:   secondsPerDay,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		"1000000000", // Default 1 billion in smallest unit
		30,           // Default 30 days per cycle
		15,           // Default 15 cycles
		// 86400,        // Default 86400 seconds per day (24h)
		10, // Default 86400 seconds per day (24h)
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyTotalRewardPool, &p.TotalRewardPool, validateTotalRewardPool),
		paramtypes.NewParamSetPair(KeyCycleDays, &p.CycleDays, validateCycleDays),
		paramtypes.NewParamSetPair(KeyTotalCycles, &p.TotalCycles, validateTotalCycles),
		paramtypes.NewParamSetPair(KeySecondsPerDay, &p.SecondsPerDay, validateSecondsPerDay),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateTotalRewardPool(p.TotalRewardPool); err != nil {
		return err
	}
	if err := validateCycleDays(p.CycleDays); err != nil {
		return err
	}
	if err := validateTotalCycles(p.TotalCycles); err != nil {
		return err
	}
	if err := validateSecondsPerDay(p.SecondsPerDay); err != nil {
		return err
	}
	return nil
}

func validateTotalRewardPool(v interface{}) error {
	totalRewardPool, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}
	if totalRewardPool == "" {
		return fmt.Errorf("total_reward_pool cannot be empty")
	}
	return nil
}

func validateCycleDays(v interface{}) error {
	cycleDays, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}
	if cycleDays == 0 {
		return fmt.Errorf("cycle_days must be greater than 0")
	}
	return nil
}

func validateTotalCycles(v interface{}) error {
	totalCycles, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}
	if totalCycles == 0 {
		return fmt.Errorf("total_cycles must be greater than 0")
	}
	return nil
}

func validateSecondsPerDay(v interface{}) error {
	secondsPerDay, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}
	if secondsPerDay == 0 {
		return fmt.Errorf("seconds_per_day must be greater than 0")
	}
	return nil
}
