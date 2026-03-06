package keeper

import (
	"blockmazechain/x/validatorbonus/types"
)

var _ types.QueryServer = Keeper{}
