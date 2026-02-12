package commands

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/urfave/cli"
)

var sendOnionCommand = cli.Command{
	Name:     "sendonion",
	Category: "Peers",
	Usage: "Send an onion message to a peer (legacy) or to a " +
		"destination via pathfinding",
	Description: `
	Send an onion message. Two modes are supported:

	Legacy mode (requires --peer, --pathkey, --onion):
	  lncli sendonion --peer <pubkey> --pathkey <key> --onion <blob>

	Pathfinding mode (requires --destination):
	  lncli sendonion --destination <pubkey> [--tlv <type>=<hex_value>]
	`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "peer",
			Usage: "(legacy) hex-encoded pubkey of the peer " +
				"to send to",
		},
		cli.StringFlag{
			Name:  "pathkey",
			Usage: "(legacy) hex-encoded path key",
		},
		cli.StringFlag{
			Name:  "onion",
			Usage: "(legacy) hex-encoded onion blob",
		},
		cli.StringFlag{
			Name: "destination",
			Usage: "hex-encoded pubkey of the destination " +
				"node; pathfinding will be used to find " +
				"a route",
		},
		cli.StringSliceFlag{
			Name: "tlv",
			Usage: "final hop TLV record as type=hex_value " +
				"(can be repeated), e.g. --tlv " +
				"77017=deadbeef",
		},
	},
	Action: actionDecorator(sendOnion),
}

func sendOnion(ctx *cli.Context) error {
	ctxc := getContext()
	client, cleanUp := getClient(ctx)
	defer cleanUp()

	req := &lnrpc.SendOnionMessageRequest{}

	destination := ctx.String("destination")
	peerStr := ctx.String("peer")

	switch {
	case destination != "" && peerStr != "":
		return fmt.Errorf("cannot set both --destination and --peer")

	case destination != "":
		dest, err := hex.DecodeString(destination)
		if err != nil {
			return fmt.Errorf("invalid destination hex: %w", err)
		}
		req.Destination = dest

		// Parse TLV flags.
		tlvs := ctx.StringSlice("tlv")
		if len(tlvs) > 0 {
			req.FinalHopTlvs = make(map[uint64][]byte)
			for _, entry := range tlvs {
				parts := strings.SplitN(entry, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid tlv "+
						"format %q, expected "+
						"type=hex_value", entry)
				}

				tlvType, err := strconv.ParseUint(
					parts[0], 10, 64,
				)
				if err != nil {
					return fmt.Errorf("invalid tlv "+
						"type %q: %w",
						parts[0], err)
				}

				val, err := hex.DecodeString(parts[1])
				if err != nil {
					return fmt.Errorf("invalid tlv "+
						"hex value %q: %w",
						parts[1], err)
				}

				req.FinalHopTlvs[tlvType] = val
			}
		}

	case peerStr != "":
		peer, err := hex.DecodeString(peerStr)
		if err != nil {
			return fmt.Errorf("invalid peer hex: %w", err)
		}
		req.Peer = peer

		pathKey, err := hex.DecodeString(ctx.String("pathkey"))
		if err != nil {
			return fmt.Errorf("invalid pathkey hex: %w", err)
		}
		req.PathKey = pathKey

		onion, err := hex.DecodeString(ctx.String("onion"))
		if err != nil {
			return fmt.Errorf("invalid onion hex: %w", err)
		}
		req.Onion = onion

	default:
		return fmt.Errorf("must set either --destination or --peer")
	}

	resp, err := client.SendOnionMessage(ctxc, req)
	if err != nil {
		return err
	}

	printRespJSON(resp)

	return nil
}

var subscribeOnionCommand = cli.Command{
	Name:     "subscribeonion",
	Category: "Peers",
	Usage:    "Subscribe to incoming onion messages",
	Action:   actionDecorator(subscribeOnion),
}

func subscribeOnion(ctx *cli.Context) error {
	ctxc := getContext()
	client, cleanUp := getClient(ctx)
	defer cleanUp()

	stream, err := client.SubscribeOnionMessages(
		ctxc, &lnrpc.SubscribeOnionMessagesRequest{},
	)
	if err != nil {
		return err
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}

		printRespJSON(msg)
	}
}
