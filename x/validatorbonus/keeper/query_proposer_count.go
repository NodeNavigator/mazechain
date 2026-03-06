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

func (k Keeper) ProposerCountAll(ctx context.Context, req *types.QueryAllProposerCountRequest) (*types.QueryAllProposerCountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var proposerCounts []types.ProposerCount

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	proposerCountStore := prefix.NewStore(store, types.KeyPrefix(types.ProposerCountKeyPrefix))

	pageRes, err := query.Paginate(proposerCountStore, req.Pagination, func(key []byte, value []byte) error {
		var proposerCount types.ProposerCount
		if err := k.cdc.Unmarshal(value, &proposerCount); err != nil {
			return err
		}

		proposerCounts = append(proposerCounts, proposerCount)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllProposerCountResponse{ProposerCount: proposerCounts, Pagination: pageRes}, nil
}

func (k Keeper) ProposerCount(ctx context.Context, req *types.QueryGetProposerCountRequest) (*types.QueryGetProposerCountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetProposerCount(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetProposerCountResponse{ProposerCount: val}, nil
}
