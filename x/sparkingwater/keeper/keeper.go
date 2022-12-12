package keeper

import (
	"fmt"

	"sparkingwater/x/sparkingwater/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey
		memKey   storetypes.StoreKey

		paramstore paramtypes.Subspace

		// the reason this is a pointer is because when you
		// make one of these, you get a pointer back. it is
		// possible there are subtleties here where if this
		// was not a pointer some state would get updated in a
		// copy and not the original keeper. i do not want to
		// play with fire.
		PermissionedWasm *wasmkeeper.PermissionedKeeper
		// The address of the contract that was instantiated at
		// genisis.
		Contract sdk.AccAddress
	}
)

type SparkingAuthPolicy struct{}

func (p SparkingAuthPolicy) CanCreateCode(wasmkeeper.ChainAccessConfigs, sdk.AccAddress, wasmtypes.AccessConfig) bool {
	return true
}

func (p SparkingAuthPolicy) CanInstantiateContract(wasmtypes.AccessConfig, sdk.AccAddress) bool {
	return true
}

func (p SparkingAuthPolicy) CanModifyContract(sdk.AccAddress, sdk.AccAddress) bool {
	return false
}

func (p SparkingAuthPolicy) CanModifyCodeAccessConfig(sdk.AccAddress, sdk.AccAddress, bool) bool {
	return false
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,

	wasmKeeper wasmkeeper.Keeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	permissionedWasm := wasmkeeper.NewPermissionedKeeper(wasmKeeper, SparkingAuthPolicy{})

	return &Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		memKey:           memKey,
		paramstore:       ps,
		PermissionedWasm: permissionedWasm,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
