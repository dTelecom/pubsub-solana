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

	recipients map[solana.PublicKey]*recipientType

	subscriptions   map[string][]Handler
	subscriptionsMu sync.RWMutex
}

func New(client, magicblockClient ContractClient, messageIdGenerator MessageIdGenerator, dataEncoder DataEncoder) *PubSub {
	return &PubSub{
		contractClient:           client,
		contractMagicblockClient: magicblockClient,
		messageIdGenerator:       messageIdGenerator,
		dataEncoder:              dataEncoder,
		recipients:               map[solana.PublicKey]*recipientType{},
		subscriptions:            map[string][]Handler{},
	}
}
