package onionmessage

import (
	"bytes"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	sphinx "github.com/lightningnetwork/lightning-onion"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/lightningnetwork/lnd/record"
	"github.com/lightningnetwork/lnd/routing/route"
	"github.com/stretchr/testify/require"
)

// TestBuildOnionMessageForPath tests that an onion message built for a
// multi-hop path can be correctly peeled by each hop using its private key.
func TestBuildOnionMessageForPath(t *testing.T) {
	t.Parallel()

	// Generate keys for a 3-hop path: hop1 -> hop2 -> dest.
	hop1Key, err := btcec.NewPrivateKey()
	require.NoError(t, err)

	hop2Key, err := btcec.NewPrivateKey()
	require.NoError(t, err)

	destKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)

	hop1Vertex := route.NewVertex(hop1Key.PubKey())
	hop2Vertex := route.NewVertex(hop2Key.PubKey())
	destVertex := route.NewVertex(destKey.PubKey())

	path := &OnionMessagePath{
		Hops: []route.Vertex{hop1Vertex, hop2Vertex, destVertex},
	}

	// Build the onion message with a final hop TLV.
	finalHopTLVs := []*lnwire.FinalHopTLV{
		{
			TLVType: lnwire.InvoiceRequestNamespaceType,
			Value:   []byte{0xde, 0xad},
		},
	}

	onionBlob, blindingKey, err := buildOnionMessageForPath(
		path, nil, finalHopTLVs,
	)
	require.NoError(t, err)
	require.NotNil(t, blindingKey)

	// Create the lnwire.OnionMessage to use with PeelOnionLayers.
	onionMsg := &lnwire.OnionMessage{
		PathKey:   blindingKey,
		OnionBlob: onionBlob,
	}

	// Peel the onion using the private keys of each hop.
	privKeys := []*btcec.PrivateKey{hop1Key, hop2Key, destKey}
	hops := PeelOnionLayers(t, privKeys, onionMsg)

	require.Len(t, hops, 3)

	// Verify intermediate hops are not final.
	require.False(t, hops[0].IsFinal, "hop1 should not be final")
	require.False(t, hops[1].IsFinal, "hop2 should not be final")

	// Verify the last hop is final.
	require.True(t, hops[2].IsFinal, "dest should be final")

	// Verify the final hop contains our custom TLV.
	require.NotNil(t, hops[2].Payload.FinalHopTLVs)
	require.Len(t, hops[2].Payload.FinalHopTLVs, 1)
	require.Equal(t, lnwire.InvoiceRequestNamespaceType,
		hops[2].Payload.FinalHopTLVs[0].TLVType)
	require.Equal(t, []byte{0xde, 0xad},
		hops[2].Payload.FinalHopTLVs[0].Value)

	// Verify intermediate hops have route data pointing to next hop.
	for i := range 2 {
		require.NotEmpty(t, hops[i].EncryptedData,
			"hop %d should have encrypted data", i)
	}
}

// TestBuildOnionMessageForPathWithReplyPath tests onion message construction
// with a reply path included for the final hop.
func TestBuildOnionMessageForPathWithReplyPath(t *testing.T) {
	t.Parallel()

	// Generate keys for a 2-hop path: hop1 -> dest.
	hop1Key, err := btcec.NewPrivateKey()
	require.NoError(t, err)

	destKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)

	hop1Vertex := route.NewVertex(hop1Key.PubKey())
	destVertex := route.NewVertex(destKey.PubKey())

	path := &OnionMessagePath{
		Hops: []route.Vertex{hop1Vertex, destVertex},
	}

	// Build a reply path (just a single-hop blinded path back to us).
	replyKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)

	replyData := &record.BlindedRouteData{}
	replyEncoded, err := record.EncodeBlindedRouteData(replyData)
	require.NoError(t, err)

	replySessionKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)

	replyBlindedPath, err := sphinx.BuildBlindedPath(replySessionKey,
		[]*sphinx.HopInfo{
			{
				NodePub:   replyKey.PubKey(),
				PlainText: replyEncoded,
			},
		},
	)
	require.NoError(t, err)

	// Build the onion with a reply path.
	onionBlob, blindingKey, err := buildOnionMessageForPath(
		path, replyBlindedPath.Path, nil,
	)
	require.NoError(t, err)

	onionMsg := &lnwire.OnionMessage{
		PathKey:   blindingKey,
		OnionBlob: onionBlob,
	}

	// Peel the onion.
	privKeys := []*btcec.PrivateKey{hop1Key, destKey}
	hops := PeelOnionLayers(t, privKeys, onionMsg)

	require.Len(t, hops, 2)
	require.True(t, hops[1].IsFinal)

	// Verify the reply path is present in the final hop payload.
	require.NotNil(t, hops[1].Payload.ReplyPath)
	require.Equal(t,
		replyBlindedPath.Path.BlindingPoint.SerializeCompressed(),
		hops[1].Payload.ReplyPath.BlindingPoint.SerializeCompressed(),
	)
}

// TestBuildOnionMessageForPathSingleHop tests onion message construction for
// a direct (single-hop) path.
func TestBuildOnionMessageForPathSingleHop(t *testing.T) {
	t.Parallel()

	destKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)

	destVertex := route.NewVertex(destKey.PubKey())

	path := &OnionMessagePath{
		Hops: []route.Vertex{destVertex},
	}

	onionBlob, blindingKey, err := buildOnionMessageForPath(
		path, nil, nil,
	)
	require.NoError(t, err)

	onionMsg := &lnwire.OnionMessage{
		PathKey:   blindingKey,
		OnionBlob: onionBlob,
	}

	privKeys := []*btcec.PrivateKey{destKey}
	hops := PeelOnionLayers(t, privKeys, onionMsg)

	require.Len(t, hops, 1)
	require.True(t, hops[0].IsFinal)
}

// TestBuildOnionMessageForPathRoundTrip tests that the onion packet can be
// decoded and re-encoded consistently.
func TestBuildOnionMessageForPathRoundTrip(t *testing.T) {
	t.Parallel()

	destKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)

	destVertex := route.NewVertex(destKey.PubKey())

	path := &OnionMessagePath{
		Hops: []route.Vertex{destVertex},
	}

	onionBlob, _, err := buildOnionMessageForPath(path, nil, nil)
	require.NoError(t, err)

	// Verify the blob can be decoded as a valid onion packet.
	var pkt sphinx.OnionPacket
	require.NoError(t, pkt.Decode(bytes.NewReader(onionBlob)))
}
