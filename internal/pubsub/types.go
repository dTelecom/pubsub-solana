package pubsub

import (
	"context"

	"github.com/gagliardetto/solana-go"
)

type Event struct {
	ID         string
	FromPeerId string
	Message    []byte
}

type Handler func(context.Context, Event)

type msgType struct {
	ID    string `borsh:"id"`
	Topic string `borsh:"topic"`
	Value []byte `borsh:"value"`
}

type recipientType struct {
	key          solana.PublicKey
	messageQueue chan []byte
}
