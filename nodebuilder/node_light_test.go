package nodebuilder

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/libp2p/go-libp2p/core/crypto"
	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	nodebuilder "github.com/celestiaorg/celestia-node/nodebuilder/node"
	"github.com/celestiaorg/celestia-node/nodebuilder/p2p"
	"github.com/celestiaorg/celestia-node/nodebuilder/state"
)

func TestNewLightWithP2PKey(t *testing.T) {
	key, _, err := crypto.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err)
	node := TestNode(t, nodebuilder.Light, p2p.WithP2PKey(key))
	assert.True(t, node.Host.ID().MatchesPrivateKey(key))
}

func TestNewLightWithHost(t *testing.T) {
	nw, _ := mocknet.WithNPeers(1)
	node := TestNode(t, nodebuilder.Light, p2p.WithHost(nw.Hosts()[0]))
	assert.Equal(t, nw.Peers()[0], node.Host.ID())
}

func TestLight_WithMutualPeers(t *testing.T) {
	peers := []string{
		"/ip6/100:0:114b:abc5:e13a:c32f:7a9e:f00a/tcp/2121/p2p/12D3KooWSRqDfpLsQxpyUhLC9oXHD2WuZ2y5FWzDri7LT4Dw9fSi",
		"/ip4/192.168.1.10/tcp/2121/p2p/12D3KooWSRqDfpLsQxpyUhLC9oXHD2WuZ2y5FWzDri7LT4Dw9fSi",
	}
	cfg := DefaultConfig(nodebuilder.Light)
	cfg.P2P.MutualPeers = peers
	node := TestNodeWithConfig(t, nodebuilder.Light, cfg)

	require.NotNil(t, node)
	assert.Equal(t, node.Config.P2P.MutualPeers, peers)
}

func TestLight_WithBlockedAddresses(t *testing.T) {
	addrs := []string{
		"/ip4/171.240.143.94/tcp/2121",
		"/ip4/113.172.240.47/tcp/2121",
		"/ip4/113.172.255.29/tcp/2121",
		"/ip4/113.172.192.189/tcp/2121",
		"/ip4/49.12.151.204/tcp/2121",
	}
	cfg := DefaultConfig(nodebuilder.Light)
	cfg.P2P.BlockAddresses = addrs
	node := TestNodeWithConfig(t, nodebuilder.Light, cfg)

	require.NotNil(t, node)
	assert.Equal(t, node.Config.P2P.BlockAddresses, addrs)
}

func TestLight_WithBlockedSubnets(t *testing.T) {
	subnets := []string{
		"/ip4/10.0.0.0/ipcidr/8",
		"/ip4/100.64.0.0/ipcidr/10",
		"/ip4/169.254.0.0/ipcidr/16",
		"/ip4/172.16.0.0/ipcidr/12",
		"/ip4/192.0.0.0/ipcidr/24",
		"/ip4/192.0.2.0/ipcidr/24",
		"/ip4/192.168.0.0/ipcidr/16",
		"/ip4/198.18.0.0/ipcidr/15",
		"/ip4/198.51.100.0/ipcidr/24",
		"/ip4/203.0.113.0/ipcidr/24",
		"/ip4/240.0.0.0/ipcidr/4",
		"/ip6/100::/ipcidr/64",
		"/ip6/2001:2::/ipcidr/48",
		"/ip6/2001:db8::/ipcidr/32",
		"/ip6/fc00::/ipcidr/7",
		"/ip6/fe80::/ipcidr/10",
	}
	cfg := DefaultConfig(nodebuilder.Light)
	cfg.P2P.BlockSubnets = subnets
	node := TestNodeWithConfig(t, nodebuilder.Light, cfg)

	require.NotNil(t, node)
	assert.Equal(t, node.Config.P2P.BlockSubnets, subnets)
}

func TestLight_WithNetwork(t *testing.T) {
	node := TestNode(t, nodebuilder.Light)
	require.NotNil(t, node)
	assert.Equal(t, p2p.Private, node.Network)
}

// TestLight_WithStubbedCoreAccessor ensures that a node started without
// a core connection will return a stubbed StateModule.
func TestLight_WithStubbedCoreAccessor(t *testing.T) {
	node := TestNode(t, nodebuilder.Light)
	_, err := node.StateServ.Balance(context.Background())
	assert.ErrorIs(t, state.ErrNoStateAccess, err)
}
