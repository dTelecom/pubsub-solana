package pubsub

import (
	"context"

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

		p.contractMagicblockClient.IncomingMessageSubscribe(ctx, nodeKey, p.makeIncomingHandler(nodeKey))
		p.contractMagicblockClient.OutgoingMessageSubscribe(ctx, nodeKey, p.makeOutgoingHandler(ctx, recipient))
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
						p.logger.Errorw("Message queue is full", nil, "recipient", recipient.key)
					}
				}
			}
		}
	}()

	return nil
}
