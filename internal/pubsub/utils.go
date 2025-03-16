package pubsub

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"

	"github.com/dTelecom/pubsub-solana/internal/contract_client"
)

const baseDelay = 100 * time.Millisecond
const maxDelay = 30 + time.Second

func (p *PubSub) makeIncomingHandler(sender solana.PublicKey) func(ctx context.Context, data contract_client.MessageData) {
	return func(ctx context.Context, data contract_client.MessageData) {
		if data.Read {
			return
		}

		defer func() {
			if _, err := p.contractMagicblockClient.MarkAsRead(ctx, sender, data.TimeStamp); err != nil {
				p.logger.Errorw("Failed to mark message as read", err)
			}
		}()

		if len(data.Content) > 0 {
			decoded, err := p.dataEncoder.Decode(data.Content)
			if err != nil {
				p.logger.Errorw("Failed to decode incoming message", err)
				return
			}

			var msg msgType
			if err := borsh.Deserialize(&msg, decoded); err != nil {
				p.logger.Errorw("Failed to deserialize incoming message", err)
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

func (p *PubSub) makeOutgoingHandler(ctx context.Context, recipient *recipientType) func(context.Context, contract_client.MessageData) {
	var (
		mu          sync.Mutex
		cond        = sync.NewCond(&mu)
		hasBeenRead bool
	)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case message, ok := <-recipient.messageQueue:
				if !ok {
					return
				}

				mu.Lock()
				for !hasBeenRead {
					cond.Wait()
				}
				mu.Unlock()

				for attempt := 0; ; attempt++ {
					_, err := p.contractMagicblockClient.SendMessage(ctx, recipient.key, message)
					if err == nil {
						break
					}
					if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
						return
					}
					p.logger.Errorw("Failed to send message", err)
					delay := baseDelay * (1 << (attempt - 1))
					if delay > maxDelay {
						delay = maxDelay
					}
					select {
					case <-ctx.Done():
						return
					case <-time.After(delay):
					}
				}

				mu.Lock()
				hasBeenRead = false
				mu.Unlock()
			}
		}
	}()

	return func(ctx context.Context, data contract_client.MessageData) {
		if !data.Read {
			return
		}

		mu.Lock()
		hasBeenRead = true
		cond.Broadcast()
		mu.Unlock()
	}
}
