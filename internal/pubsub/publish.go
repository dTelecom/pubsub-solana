package pubsub

import (
	"errors"
	"fmt"

	"github.com/near/borsh-go"
)

func (p *PubSub) Publish(topic string, value []byte) (string, error) {
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

	select {
	case p.messageQueue <- encoded:
		return msg.ID, nil
	default:
		return "", errors.New("global message queue is full")
	}
}
