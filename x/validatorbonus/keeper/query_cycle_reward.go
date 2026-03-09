package keeper

import (
	"context"

	"blockmazechain/x/validatorbonus/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CycleRewardAll(ctx context.Context, req *types.QueryAllCycleRewardRequest) (*types.QueryAllCycleRewardResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var cycleRewards []types.CycleReward

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	cycleRewardStore := prefix.NewStore(store, types.KeyPrefix(types.CycleRewardKeyPrefix))

	pageRes, err := query.Paginate(cycleRewardStore, req.Pagination, func(key []byte, value []byte) error {
		var cycleReward types.CycleReward
		if err := k.cdc.Unmarshal(value, &cycleReward); err != nil {
			return err
		}

		cycleRewards = append(cycleRewards, cycleReward)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllCycleRewardResponse{CycleReward: cycleRewards, Pagination: pageRes}, nil
}

func (k Keeper) CycleReward(ctx context.Context, req *types.QueryGetCycleRewardRequest) (*types.QueryGetCycleRewardResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetCycleReward(
		ctx,
		req.Id,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetCycleRewardResponse{CycleReward: val}, nil
}
