package keeper

import (
	"context"
	"fmt"

	"blockmazechain/x/validatorbonus/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CycleRewards queries all validator rewards for a specific cycle.
//
// Logic:
// Iterate over cycleReward/{cycle} prefix and return:
// - validators[]: array of validator addresses
// - rewards[]: array of corresponding reward amounts (as JSON strings)
func (k Keeper) CycleRewards(goCtx context.Context, req *types.QueryCycleRewardsRequest) (*types.QueryCycleRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	cycleStore := prefix.NewStore(store, types.KeyPrefix(types.CycleRewardKeyPrefix))

	var validators []string
	var rewards []string

	// Use prefix indexing by cycle to avoid scanning the entire store
	prefixKey := []byte(fmt.Sprintf("%08d:", req.Cycle))
	cyclePrefixStore := prefix.NewStore(cycleStore, prefixKey)

	pageRes, err := query.Paginate(cyclePrefixStore, req.Pagination, func(key []byte, value []byte) error {
		var cr types.CycleReward
		if err := k.cdc.Unmarshal(value, &cr); err != nil {
			return err
		}

		validators = append(validators, cr.ValidatorAddress)
		rewards = append(rewards, cr.Amount)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCycleRewardsResponse{
		Validators: validators,
		Rewards:    rewards,
		Pagination: pageRes,
	}, nil
}
