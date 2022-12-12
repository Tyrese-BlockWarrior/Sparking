# 🌩️

this repository contains a POC for instantiating a CosmWasm contract
at genesis for a Cosmos SDK based blockchain.

at a high level, this is done by:

1. scaffolding a new chain with [ignite](https://ignite.com)
2. adding the wasmd module (and thus CosmWasm support)
3. creating a new module `x/sparkingwater` which depends on wasmd and
   instantiates a contract at genesis

this targets SDK version 0.46.

## Creating the chain

First, install the [Ignite](https://ignite.com) CLI:

```
curl https://get.ignite.com/cli! | bash
```

(you may need to run this as `sudo`)

Then, make a new chain:

```
ignite scaffold chain sparkingwater --address-prefix sw
```

## Wasmd support

get `wasmd`:

```
go get github.com/CosmWasm/wasmd@v0.29.2
```

If you try and use `wasmd` at this point as a dependency you'll get an
error like this:

```
sparkling-water/app imports
	github.com/CosmWasm/wasmd/x/wasm imports
	github.com/CosmWasm/wasmd/x/wasm/client/rest imports
	github.com/cosmos/cosmos-sdk/types/rest: module github.com/cosmos/cosmos-sdk@latest found (v0.46.6), but does not contain package github.com/cosmos/cosmos-sdk/types/rest
```

This happens because the CosmWasm folks don't seem to be a big fan of
the `0.46` SDK release.

To remedy this, add this to the bottom of your `go.mod` file:

```
replace (
	// https://github.com/CosmWasm/wasmd/pull/938
	github.com/CosmWasm/wasmd => github.com/notional-labs/wasmd v0.29.0-sdk46
)
```

Then, follow the [integration
instructions](https://docs.cosmwasm.com/docs/1.0/integration/).

## Instantiation at genesis

in wasmd there is such a thing as a `PermissionedKeeper` which is a
keeper with a limited set of operations avaliable. when the
sparkingwater keeper is created, we create a new permissioned keeper
which allows to store and instantiate, but not modify contracts.

to do so, we create a new auth policy:

```go
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
```

then, create and store the permissioned keeper during creation of the
sparkingwater keeper:

```go
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
```

that keeper being present, we can then store and instantiate a wasm
contract at genisis like so:

```go
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
```

the commit history of this repository is intentionally quite clean and
follows the steps here. for full code examples, refer to the
corresponding commit.
