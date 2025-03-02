package pubsub

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

func (p *PubSub) Start(ctx context.Context, nodes []solana.PublicKey) error {
	for _, node := range nodes {
		if p.contractClient.IsSigner(node) {
			continue
		}

		p.recipients[node] = &recipientType{
			key:          node,
			sending:      false,
			messageQueue: [][]byte{},
		}

		if err := p.contractMagicblockClient.IncomingMessageSubscribe(ctx, node, p.makeIncomingHandler(node)); err != nil {
			return fmt.Errorf("failed to subscribe to incoming messages: %w", err)
		}

		if err := p.contractMagicblockClient.OutgoingMessageSubscribe(ctx, node, p.makeOutgoingHandler(node)); err != nil {
			return fmt.Errorf("failed to subscribe to outgoing messages: %w", err)
		}
	}

	return nil
}
