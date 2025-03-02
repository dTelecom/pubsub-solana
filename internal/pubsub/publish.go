package pubsub

import (
	"context"
	"fmt"

	"github.com/near/borsh-go"
)

func (p *PubSub) Publish(ctx context.Context, topic string, value []byte) (string, error) {
	msg := msgType{
		ID:    p.messageIdGenerator.Generate(),
		Topic: topic,
		Value: value,
	}

	serialized, err := borsh.Serialize(msg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal publish message: %w", err)
	}

	encoded, err := p.dataEncoder.Encode(serialized)
	if err != nil {
		return "", fmt.Errorf("failed to encode publish message: %w", err)
	}

	for _, recipient := range p.recipients {
		recipient.mu.Lock()
		if !recipient.sending {
			p.sendMessage(ctx, recipient, encoded)
		} else {
			recipient.messageQueue = append(recipient.messageQueue, encoded)
		}
		recipient.mu.Unlock()
	}

	return msg.ID, nil
}
