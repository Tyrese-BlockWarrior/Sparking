package sparkingwater

import (
	"io/ioutil"
	"sparkingwater/x/sparkingwater/keeper"
	"sparkingwater/x/sparkingwater/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState, moduleAddress sdk.AccAddress) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)

	wasmCode, err := ioutil.ReadFile("./cw_admin_factory.wasm")
	if err != nil {
		panic(err)
	}

	code_id, _, err := k.PermissionedWasm.Create(
		ctx,
		moduleAddress,
		wasmCode,
		&wasmtypes.AccessConfig{
			Permission: wasmtypes.AccessTypeOnlyAddress,
			Addresses:  []string{moduleAddress.String()},
		})
	if err != nil {
		panic(err)
	}

	addr, _, err := k.PermissionedWasm.Instantiate(
		ctx,
		code_id,
		moduleAddress,
		moduleAddress,
		[]byte("{}"),
		"le contract",
		sdk.Coins{})
	if err != nil {
		panic(err)
	}

	k.Contract = addr
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
