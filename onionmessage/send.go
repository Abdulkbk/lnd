package onionmessage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	sphinx "github.com/lightningnetwork/lightning-onion"
	"github.com/lightningnetwork/lnd/actor"
	"github.com/lightningnetwork/lnd/fn/v2"
	graphdb "github.com/lightningnetwork/lnd/graph/db"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/lightningnetwork/lnd/record"
	"github.com/lightningnetwork/lnd/routing/route"
)

// SendConfig contains the dependencies needed to find a path and send an
// onion message to a destination node.
type SendConfig struct {
	// Graph provides access to the channel graph for pathfinding.
	Graph graphdb.NodeTraverser

	// OurPubKey is our node's public key, used as the BFS source.
	OurPubKey route.Vertex

	// Receptionist is used to look up peer actors for sending messages.
	Receptionist *actor.Receptionist

	// MaxHops is the maximum number of hops for the BFS path search.
	MaxHops int
}

// SendToDestination finds a path to the destination node, constructs a
// blinded onion message, and sends it via the first hop's peer actor.
func SendToDestination(ctx context.Context, cfg *SendConfig,
	destination route.Vertex, finalHopTLVs []*lnwire.FinalHopTLV,
	replyPath *sphinx.BlindedPath) error {

	// Find the shortest path to the destination.
	path, err := FindPath(
		cfg.Graph, cfg.OurPubKey, destination, cfg.MaxHops,
	)
	if err != nil {
		return err
	}

	if len(path.Hops) == 0 {
		return fmt.Errorf("path to self is not supported")
	}

	// Build the blinded path and onion message for the discovered route.
	onionMsg, blindingKey, err := buildOnionMessageForPath(
		path, replyPath, finalHopTLVs,
	)
	if err != nil {
		return fmt.Errorf("failed to build onion message: %w", err)
	}

	// Send via the first hop's peer actor.
	firstHop := path.Hops[0]

	return sendToFirstHop(ctx, cfg.Receptionist, firstHop, blindingKey,
		onionMsg)
}

// SendDirectToDestination builds a blinded onion message for the given
// pre-built path (no pathfinding) and sends it to the first hop's peer actor.
// This is used as a fallback when graph pathfinding fails but the destination
// may be a directly connected peer.
func SendDirectToDestination(ctx context.Context, cfg *SendConfig,
	path *OnionMessagePath, finalHopTLVs []*lnwire.FinalHopTLV,
	replyPath *sphinx.BlindedPath) error {

	if len(path.Hops) == 0 {
		return fmt.Errorf("path must have at least one hop")
	}

	onionMsg, blindingKey, err := buildOnionMessageForPath(
		path, replyPath, finalHopTLVs,
	)
	if err != nil {
		return fmt.Errorf("failed to build onion message: %w", err)
	}

	firstHop := path.Hops[0]

	return sendToFirstHop(ctx, cfg.Receptionist, firstHop, blindingKey,
		onionMsg)
}

// buildOnionMessageForPath constructs a blinded onion message for the given
// path. It returns the serialized onion blob and the blinding session public
// key needed for the first hop.
func buildOnionMessageForPath(path *OnionMessagePath,
	replyPath *sphinx.BlindedPath,
	finalHopTLVs []*lnwire.FinalHopTLV) ([]byte, *btcec.PublicKey,
	error) {

	hops := path.Hops

	// Build HopInfo list for sphinx.BuildBlindedPath.
	hopInfos := make([]*sphinx.HopInfo, len(hops))

	for i, hop := range hops {
		pubKey, err := btcec.ParsePubKey(hop[:])
		if err != nil {
			return nil, nil, fmt.Errorf("invalid pubkey at "+
				"hop %d: %w", i, err)
		}

		var routeData *record.BlindedRouteData

		// Final hop gets empty route data.
		if i == len(hops)-1 {
			routeData = &record.BlindedRouteData{}
		} else {
			// Non-final hops get NextNodeID pointing to the next
			// hop.
			nextPub, err := btcec.ParsePubKey(hops[i+1][:])
			if err != nil {
				return nil, nil, fmt.Errorf("invalid next "+
					"pubkey at hop %d: %w", i, err)
			}

			nextNode := fn.NewLeft[*btcec.PublicKey,
				lnwire.ShortChannelID](nextPub)

			routeData = record.NewNonFinalBlindedRouteDataOnionMessage( //nolint:lll
				nextNode, nil, nil,
			)
		}

		encoded, err := record.EncodeBlindedRouteData(routeData)
		if err != nil {
			return nil, nil, fmt.Errorf("encode route data "+
				"hop %d: %w", i, err)
		}

		hopInfos[i] = &sphinx.HopInfo{
			NodePub:   pubKey,
			PlainText: encoded,
		}
	}

	// Build the blinded path with a fresh session key.
	sessionKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("generate session key: %w", err)
	}

	blindedPath, err := sphinx.BuildBlindedPath(sessionKey, hopInfos)
	if err != nil {
		return nil, nil, fmt.Errorf("build blinded path: %w", err)
	}

	// Convert to a sphinx payment path for onion construction.
	sphinxPath, err := route.OnionMessageBlindedPathToSphinxPath(
		blindedPath.Path, replyPath, finalHopTLVs,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("convert to sphinx path: %w", err)
	}

	// Build the onion packet.
	onionSessionKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("generate onion key: %w", err)
	}

	onionPkt, err := sphinx.NewOnionPacket(
		sphinxPath, onionSessionKey, nil,
		sphinx.DeterministicPacketFiller,
		sphinx.WithMaxPayloadSize(sphinx.MaxRoutingPayloadSize),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("create onion packet: %w", err)
	}

	var buf bytes.Buffer
	if err := onionPkt.Encode(&buf); err != nil {
		return nil, nil, fmt.Errorf("encode onion packet: %w", err)
	}

	return buf.Bytes(), blindedPath.SessionKey.PubKey(), nil
}

// sendToFirstHop looks up the peer actor for the given node and sends the
// onion message.
func sendToFirstHop(ctx context.Context, receptionist *actor.Receptionist,
	firstHop route.Vertex, blindingKey *btcec.PublicKey,
	onionBlob []byte) error {

	var pubKeyBytes [33]byte
	copy(pubKeyBytes[:], firstHop[:])

	actorRefOpt := findPeerActor(receptionist, pubKeyBytes)

	if actorRefOpt.IsNone() {
		return ErrPeerActorNotFound
	}

	actorRefOpt.WhenSome(func(actorRef actor.ActorRef[*Request,
		*Response]) {

		onionMsg := lnwire.NewOnionMessage(blindingKey, onionBlob)
		req := &Request{msg: *onionMsg}
		actorRef.Tell(ctx, req)
	})

	return nil
}
