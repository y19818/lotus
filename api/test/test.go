package test

import (
	"context"
	"testing"

	"github.com/multiformats/go-multiaddr"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/miner"
	"github.com/filecoin-project/lotus/node"
)

type TestNode struct {
	api.FullNode
	// ListenAddr is the address on which an API server is listening, if an
	// API server is created for this Node
	ListenAddr multiaddr.Multiaddr
}

type TestStorageNode struct {
	api.StorageMiner
	// ListenAddr is the address on which an API server is listening, if an
	// API server is created for this Node
	ListenAddr multiaddr.Multiaddr

	MineOne func(context.Context, miner.MineReq) error
}

var PresealGenesis = -1

const GenesisPreseals = 2

type StorageMiner struct {
	Full    int
	Preseal int
}

// APIBuilder is a function which is invoked in test suite to provide
// test nodes and networks
//
// storage array defines storage nodes, numbers in the array specify full node
// index the storage node 'belongs' to
type APIBuilder func(t *testing.T, nFull int, storage []StorageMiner, opts ...node.Option) ([]TestNode, []TestStorageNode)
type testSuite struct {
	makeNodes APIBuilder
}

// TestApis is the entry point to API test suite
func TestApis(t *testing.T, b APIBuilder) {
	ts := testSuite{
		makeNodes: b,
	}

	t.Run("version", ts.testVersion)
	t.Run("id", ts.testID)
	t.Run("testConnectTwo", ts.testConnectTwo)
	t.Run("testMining", ts.testMining)
	t.Run("testMiningReal", ts.testMiningReal)
}

var OneMiner = []StorageMiner{{Full: 0, Preseal: PresealGenesis}}

func (ts *testSuite) testVersion(t *testing.T) {
	build.RunningNodeType = build.NodeFull

	ctx := context.Background()
	apis, _ := ts.makeNodes(t, 1, OneMiner)
	api := apis[0]

	v, err := api.Version(ctx)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, v.Version, build.BuildVersion)
}

func (ts *testSuite) testID(t *testing.T) {
	ctx := context.Background()
	apis, _ := ts.makeNodes(t, 1, OneMiner)
	api := apis[0]

	id, err := api.ID(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Regexp(t, "^12", id.Pretty())
}

func (ts *testSuite) testConnectTwo(t *testing.T) {
	ctx := context.Background()
	apis, _ := ts.makeNodes(t, 2, OneMiner)

	p, err := apis[0].NetPeers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(p) != 0 {
		t.Error("Node 0 has a peer")
	}

	p, err = apis[1].NetPeers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(p) != 0 {
		t.Error("Node 1 has a peer")
	}

	addrs, err := apis[1].NetAddrsListen(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if err := apis[0].NetConnect(ctx, addrs); err != nil {
		t.Fatal(err)
	}

	p, err = apis[0].NetPeers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(p) != 1 {
		t.Error("Node 0 doesn't have 1 peer")
	}

	p, err = apis[1].NetPeers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(p) != 1 {
		t.Error("Node 0 doesn't have 1 peer")
	}
}
