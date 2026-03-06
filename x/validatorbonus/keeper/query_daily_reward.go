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

func (k Keeper) DailyRewardAll(ctx context.Context, req *types.QueryAllDailyRewardRequest) (*types.QueryAllDailyRewardResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var dailyRewards []types.DailyReward

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	dailyRewardStore := prefix.NewStore(store, types.KeyPrefix(types.DailyRewardKeyPrefix))

	pageRes, err := query.Paginate(dailyRewardStore, req.Pagination, func(key []byte, value []byte) error {
		var dailyReward types.DailyReward
		if err := k.cdc.Unmarshal(value, &dailyReward); err != nil {
			return err
		}

		dailyRewards = append(dailyRewards, dailyReward)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllDailyRewardResponse{DailyReward: dailyRewards, Pagination: pageRes}, nil
}

func (k Keeper) DailyReward(ctx context.Context, req *types.QueryGetDailyRewardRequest) (*types.QueryGetDailyRewardResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetDailyReward(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetDailyRewardResponse{DailyReward: val}, nil
}
