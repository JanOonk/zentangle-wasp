package test

import (
	"fmt"
	"testing"

	"github.com/iotaledger/wasp/contracts/wasm/testwasmlib/go/testwasmlibclient"
	"github.com/iotaledger/wasp/packages/vm/wasmlib/go/wasmclient"
	"github.com/stretchr/testify/require"
)

// hardcoded seed and chain ID, taken from wasp-cli.json
// note that normally the chain has already been set up and
// the contract has already been deployed in some way, so
// these values are usually available from elsewhere
const (
	mySeed    = "6C6tRksZDWeDTCzX4Q7R2hbpyFV86cSGLVxdkFKSB3sv"
	myChainID = "jn52vSuUUYY22T1mV2ny14EADYBu3ofyewLRSsVRnjpz"
)

func setupClient(t *testing.T) *testwasmlibclient.TestWasmLibService {
	require.True(t, wasmclient.SeedIsValid(mySeed))
	require.True(t, wasmclient.ChainIsValid(myChainID))

	// we're testing against wasp-cluster, so defaults will do
	svcClient := wasmclient.DefaultServiceClient()

	// create the service for the testwasmlib smart contract
	svc, err := testwasmlibclient.NewTestWasmLibService(svcClient, myChainID)
	require.NoError(t, err)

	// we'll use the first address in the seed to sign requests
	svc.SignRequests(wasmclient.SeedToKeyPair(mySeed, 0))
	return svc
}

func TestClientEvents(t *testing.T) {
	svc := setupClient(t)

	// get new triggerEvent interface, pass params, and post the request
	f := svc.TriggerEvent()
	f.Name("Lala")
	f.Address(wasmclient.SeedToAddress(mySeed, 0))
	req1 := f.Post()
	require.NoError(t, req1.Error())

	// err := svc.WaitRequest(req1)
	// require.NoError(t, err)

	// get new triggerEvent interface, pass params, and post the request
	f = svc.TriggerEvent()
	f.Name("Trala")
	f.Address(wasmclient.SeedToAddress(mySeed, 1))
	req2 := f.Post()
	require.NoError(t, req2.Error())

	err := svc.WaitRequest(req2)
	require.NoError(t, err)
}

func TestClientRandom(t *testing.T) {
	svc := setupClient(t)

	// generate new random value
	f := svc.Random()
	req := f.Post()
	require.NoError(t, req.Error())

	err := svc.WaitRequest(req)
	require.NoError(t, err)

	// get current random value
	v := svc.GetRandom()
	res := v.Call()
	require.NoError(t, v.Error())
	require.GreaterOrEqual(t, res.Random(), int64(0))
	fmt.Println("Random: ", res.Random())
}

func TestClientArray(t *testing.T) {
	svc := setupClient(t)

	v := svc.ArrayLength()
	v.Name("Bands")
	res := v.Call()
	require.NoError(t, v.Error())
	require.EqualValues(t, 0, res.Length())

	f := svc.ArraySet()
	f.Name("Bands")
	f.Index(0)
	f.Value("Dire Straits")
	req := f.Post()
	require.NoError(t, req.Error())
	err := svc.WaitRequest(req)
	require.NoError(t, err)

	v = svc.ArrayLength()
	v.Name("Bands")
	res = v.Call()
	require.NoError(t, v.Error())
	require.EqualValues(t, 1, res.Length())

	c := svc.ArrayClear()
	c.Name("Bands")
	req = c.Post()
	require.NoError(t, req.Error())
	err = svc.WaitRequest(req)
	require.NoError(t, err)

	v = svc.ArrayLength()
	v.Name("Bands")
	res = v.Call()
	require.NoError(t, v.Error())
	require.EqualValues(t, 0, res.Length())
}
