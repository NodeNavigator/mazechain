package types

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	// this line is used by starport scaffolding # 1
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateEligibleValidator{},
		&MsgUpdateEligibleValidator{},
		&MsgDeleteEligibleValidator{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateProposerCount{},
		&MsgUpdateProposerCount{},
		&MsgDeleteProposerCount{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateDailyReward{},
		&MsgUpdateDailyReward{},
		&MsgDeleteDailyReward{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateCycleReward{},
		&MsgUpdateCycleReward{},
		&MsgDeleteCycleReward{},
	)
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
