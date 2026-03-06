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

func (k Keeper) EligibleValidatorAll(ctx context.Context, req *types.QueryAllEligibleValidatorRequest) (*types.QueryAllEligibleValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var eligibleValidators []types.EligibleValidator

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	eligibleValidatorStore := prefix.NewStore(store, types.KeyPrefix(types.EligibleValidatorKeyPrefix))

	pageRes, err := query.Paginate(eligibleValidatorStore, req.Pagination, func(key []byte, value []byte) error {
		var eligibleValidator types.EligibleValidator
		if err := k.cdc.Unmarshal(value, &eligibleValidator); err != nil {
			return err
		}

		eligibleValidators = append(eligibleValidators, eligibleValidator)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllEligibleValidatorResponse{EligibleValidator: eligibleValidators, Pagination: pageRes}, nil
}

func (k Keeper) EligibleValidator(ctx context.Context, req *types.QueryGetEligibleValidatorRequest) (*types.QueryGetEligibleValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetEligibleValidator(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetEligibleValidatorResponse{EligibleValidator: val}, nil
}
