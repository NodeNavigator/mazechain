package blockmazechain_test

import (
	"testing"

	keepertest "blockmazechain/testutil/keeper"
	"blockmazechain/testutil/nullify"
	blockmazechain "blockmazechain/x/blockmazechain/module"
	"blockmazechain/x/blockmazechain/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.BlockmazechainKeeper(t)
	blockmazechain.InitGenesis(ctx, k, genesisState)
	got := blockmazechain.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
