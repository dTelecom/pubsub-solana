package pubsub

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

func (p *PubSub) Start(ctx context.Context, nodeKeys []solana.PublicKey) error {
	for _, nodeKey := range nodeKeys {
		if p.contractClient.IsSigner(nodeKey) {
			continue
		}

		recipient := &recipientType{
			key:          nodeKey,
			messageQueue: make(chan []byte, 100),
		}

		p.recipients[nodeKey] = recipient

		if err := p.contractMagicblockClient.IncomingMessageSubscribe(ctx, nodeKey, p.makeIncomingHandler(nodeKey)); err != nil {
			return fmt.Errorf("failed to subscribe to incoming messages: %w", err)
		}

		if err := p.contractMagicblockClient.OutgoingMessageSubscribe(ctx, nodeKey, p.makeOutgoingHandler(ctx, recipient)); err != nil {
			return fmt.Errorf("failed to subscribe to outgoing messages: %w", err)
		}
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case message, ok := <-p.messageQueue:
				if !ok {
					return
				}
				for _, recipient := range p.recipients {
					select {
					case recipient.messageQueue <- message:
						continue
					default:
						fmt.Printf("recipient message queue is full: %v\n", recipient.key)
					}
				}
			}
		}
	}()

	return nil
}
