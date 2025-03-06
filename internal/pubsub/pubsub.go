package pubsub

import (
	"sync"

	"github.com/gagliardetto/solana-go"
)

type PubSub struct {
	contractClient           ContractClient
	contractMagicblockClient ContractClient
	messageIdGenerator       MessageIdGenerator
	dataEncoder              DataEncoder

	messageQueue chan []byte
	recipients   map[solana.PublicKey]*recipientType

	subscriptions   map[string][]Handler
	subscriptionsMu sync.RWMutex
}

func New(client, magicblockClient ContractClient, messageIdGenerator MessageIdGenerator, dataEncoder DataEncoder) *PubSub {
	return &PubSub{
		contractClient:           client,
		contractMagicblockClient: magicblockClient,
		messageIdGenerator:       messageIdGenerator,
		dataEncoder:              dataEncoder,
		messageQueue:             make(chan []byte, 100),
		recipients:               map[solana.PublicKey]*recipientType{},
		subscriptions:            map[string][]Handler{},
	}
}
