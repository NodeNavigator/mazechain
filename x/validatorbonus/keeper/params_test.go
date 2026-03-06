package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "blockmazechain/testutil/keeper"
	"blockmazechain/x/validatorbonus/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.ValidatorbonusKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
