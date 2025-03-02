package pubsub

import (
	"context"
	"sync"

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
	sending      bool
	messageQueue [][]byte
	mu           sync.Mutex
}
