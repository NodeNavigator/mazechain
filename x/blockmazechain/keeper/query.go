package keeper

import (
	"blockmazechain/x/blockmazechain/types"
)

var _ types.QueryServer = Keeper{}
