package sparkingwater_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "sparkingwater/testutil/keeper"
	"sparkingwater/testutil/nullify"
	"sparkingwater/x/sparkingwater"
	"sparkingwater/x/sparkingwater/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.SparkingwaterKeeper(t)
	sparkingwater.InitGenesis(ctx, *k, genesisState)
	got := sparkingwater.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
