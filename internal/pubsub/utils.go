package pubsub

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"

	"github.com/dTelecom/pubsub-solana/internal/contract_client"
)

func (p *PubSub) makeIncomingHandler(sender solana.PublicKey) func(ctx context.Context, data contract_client.MessageData) {
	return func(ctx context.Context, data contract_client.MessageData) {
		if data.Read {
			return
		}

		defer func() {
			if _, err := p.contractMagicblockClient.MarkAsRead(ctx, sender, data.TimeStamp); err != nil {
				fmt.Printf("failed to mark message as read: %v\n", err)
			}
		}()

		if len(data.Content) > 0 {
			decoded, err := p.dataEncoder.Decode(data.Content)
			if err != nil {
				fmt.Printf("failed to decode incoming message: %v\n", err)
				return
			}

			var msg msgType
			if err := borsh.Deserialize(&msg, decoded); err != nil {
				fmt.Printf("failed to deserialize incoming message: %v\n", err)
				return
			}

			p.subscriptionsMu.RLock()
			handlers, found := p.subscriptions[msg.Topic]
			p.subscriptionsMu.RUnlock()

			if found && len(handlers) > 0 {
				event := Event{
					ID:         msg.ID,
					FromPeerId: sender.String(),
					Message:    msg.Value,
				}

				for _, h := range handlers {
					h(ctx, event)
				}
			}
		}
	}
}

func (p *PubSub) makeOutgoingHandler(receiver solana.PublicKey) func(context.Context, contract_client.MessageData) {
	return func(ctx context.Context, data contract_client.MessageData) {
		if !data.Read {
			return
		}
		recipient, ok := p.recipients[receiver]
		if !ok {
			fmt.Printf("Unknown recipient: %v", receiver.String())
			return
		}

		recipient.mu.Lock()
		defer recipient.mu.Unlock()

		recipient.sending = false

		if len(recipient.messageQueue) > 0 {
			next := recipient.messageQueue[0]
			recipient.messageQueue = recipient.messageQueue[1:]
			p.sendMessage(ctx, recipient, next)
		}
	}
}

func (p *PubSub) sendMessage(ctx context.Context, recipient *recipientType, content []byte) {
	_, err := p.contractMagicblockClient.SendMessage(ctx, recipient.key, content)
	if err != nil {
		fmt.Printf("failed to send message to %s: %v\n", recipient.key.String(), err)
		recipient.messageQueue = append([][]byte{content}, recipient.messageQueue...)
		return
	}
	recipient.sending = true
}
