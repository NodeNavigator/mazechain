package keeper_test

import (
	"strconv"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	keepertest "blockmazechain/testutil/keeper"
	"blockmazechain/x/validatorbonus/keeper"
	"blockmazechain/x/validatorbonus/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

// authority is the module authority address used in the test keeper.
var testAuthority = authtypes.NewModuleAddress(govtypes.ModuleName).String()

func TestEligibleValidatorMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.ValidatorbonusKeeper(t)
	srv := keeper.NewMsgServerImpl(k)

	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateEligibleValidator{
			Creator:          testAuthority, // must use module authority
			Index:            strconv.Itoa(i),
			ValidatorAddress: "addr" + strconv.Itoa(i),
			JoinTime:         100,
		}
		_, err := srv.CreateEligibleValidator(ctx, expected)
		require.NoError(t, err)
		rst, found := k.GetEligibleValidator(ctx, expected.Index)
		require.True(t, found)
		require.Equal(t, expected.Creator, rst.Creator)
	}
}

func TestEligibleValidatorMsgServerCreateUnauthorized(t *testing.T) {
	k, ctx := keepertest.ValidatorbonusKeeper(t)
	srv := keeper.NewMsgServerImpl(k)

	// Any non-authority address must be rejected
	_, err := srv.CreateEligibleValidator(ctx, &types.MsgCreateEligibleValidator{
		Creator: "random-wallet-address",
		Index:   "0",
	})
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
}
