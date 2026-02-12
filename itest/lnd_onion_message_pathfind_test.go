package itest

import (
	"testing"
	"time"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lntest"
	"github.com/lightningnetwork/lnd/lntest/node"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/stretchr/testify/require"
)

// testOnionMessagePathfinding tests sending onion messages using automatic
// pathfinding through the graph.
func testOnionMessagePathfinding(ht *lntest.HarnessTest) {
	// Spin up three nodes for the test network.
	alice := ht.NewNodeWithCoins("Alice", nil)
	bob := ht.NewNodeWithCoins("Bob", nil)
	carol := ht.NewNode("Carol", nil)

	// Connect nodes so they can share gossip and forward messages.
	ht.ConnectNodesPerm(alice, bob)
	ht.ConnectNodesPerm(bob, carol)

	// Open channels so nodes appear in the graph with edges.
	chanPoint := ht.OpenChannel(
		alice, bob, lntest.OpenChannelParams{Amt: 100000},
	)
	ht.AssertChannelInGraph(alice, chanPoint)

	chanPoint2 := ht.OpenChannel(
		bob, carol, lntest.OpenChannelParams{Amt: 100000},
	)
	ht.AssertChannelInGraph(alice, chanPoint2)

	testCases := []struct {
		name string
		test func(ht *lntest.HarnessTest, alice, bob,
			carol *node.HarnessNode)
	}{
		{
			name: "multi-hop pathfinding",
			test: testMultiHopPathfinding,
		},
		{
			name: "direct peer fallback",
			test: testDirectPeerFallback,
		},
	}

	for _, tc := range testCases {
		success := ht.Run(tc.name, func(t *testing.T) {
			tc.test(ht, alice, bob, carol)
		})
		if !success {
			break
		}
	}
}

// testMultiHopPathfinding tests that Alice can send an onion message to Carol
// via pathfinding (Alice -> Bob -> Carol) using the destination field.
func testMultiHopPathfinding(ht *lntest.HarnessTest, alice, bob,
	carol *node.HarnessNode) {

	// Subscribe to onion messages on Carol.
	msgClient, cancel := carol.RPC.SubscribeOnionMessages()
	defer cancel()

	messages := make(chan *lnrpc.OnionMessageUpdate)
	go func() {
		for {
			msg, err := msgClient.Recv()
			if err != nil {
				return
			}
			select {
			case messages <- msg:
			case <-ht.Context().Done():
				return
			}
		}
	}()

	// Send an onion message from Alice to Carol using pathfinding.
	// Alice doesn't need to construct the onion manually — the system
	// will find a path through Bob and build the message.
	tlvType := uint64(lnwire.InvoiceRequestNamespaceType)
	alice.RPC.SendOnionMessage(&lnrpc.SendOnionMessageRequest{
		Destination:  carol.PubKey[:],
		FinalHopTlvs: map[uint64][]byte{tlvType: {1, 2, 3}},
	})

	// Wait for Carol to receive the message.
	select {
	case msg := <-messages:
		// Verify the custom record arrived.
		require.Equal(
			ht, []byte{1, 2, 3},
			msg.CustomRecords[tlvType],
		)

	case <-time.After(lntest.DefaultTimeout):
		ht.Fatalf("carol did not receive pathfound onion message")
	}
}

// testDirectPeerFallback tests that when pathfinding fails (e.g., destination
// not in graph), the system falls back to direct peer send if the destination
// is a connected peer.
func testDirectPeerFallback(ht *lntest.HarnessTest, alice, _,
	_ *node.HarnessNode) {

	// Create a new node Dave that is connected to Alice but has no
	// channels (so won't be in the graph).
	dave := ht.NewNode("Dave", nil)
	ht.ConnectNodesPerm(alice, dave)

	// Subscribe to onion messages on Dave.
	msgClient, cancel := dave.RPC.SubscribeOnionMessages()
	defer cancel()

	messages := make(chan *lnrpc.OnionMessageUpdate)
	go func() {
		for {
			msg, err := msgClient.Recv()
			if err != nil {
				return
			}
			select {
			case messages <- msg:
			case <-ht.Context().Done():
				return
			}
		}
	}()

	// Send from Alice to Dave using pathfinding mode.
	// Dave is not in the graph, so pathfinding will fail and the system
	// should fall back to direct send.
	tlvType := uint64(lnwire.InvoiceRequestNamespaceType)
	alice.RPC.SendOnionMessage(&lnrpc.SendOnionMessageRequest{
		Destination:  dave.PubKey[:],
		FinalHopTlvs: map[uint64][]byte{tlvType: {4, 5, 6}},
	})

	// Wait for Dave to receive the message.
	select {
	case msg := <-messages:
		require.Equal(
			ht, []byte{4, 5, 6},
			msg.CustomRecords[tlvType],
		)

	case <-time.After(lntest.DefaultTimeout):
		ht.Fatalf("dave did not receive fallback onion message")
	}
}
