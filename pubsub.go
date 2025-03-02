package pubsub_solana

import (
	"context"

	"github.com/DataDog/zstd"
	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"

	"github.com/dTelecom/pubsub-solana/internal/contract_client"
	pubsub_internal "github.com/dTelecom/pubsub-solana/internal/pubsub"
)

type uuidGenerator struct{}

func (*uuidGenerator) Generate() string {
	return uuid.New().String()
}

type zstdEncoder struct{}

func (*zstdEncoder) Encode(src []byte) ([]byte, error) {
	return zstd.Compress(nil, src)
}

func (*zstdEncoder) Decode(src []byte) ([]byte, error) {
	return zstd.Decompress(nil, src)
}

type Event struct {
	ID         string
	FromPeerId string
	Message    []byte
}

type Handler func(context.Context, Event)

type PubSub struct {
	*pubsub_internal.PubSub
	peerId string
}

func New(solanaRPC, solanaWS, ephemeralRPC, ephemeralWS, privateKey string) *PubSub {
	signer := solana.MustPrivateKeyFromBase58(privateKey)

	solanaClient := contract_client.New(false, solanaRPC, solanaWS, signer)
	ephemeralClient := contract_client.New(true, ephemeralRPC, ephemeralWS, signer)

	return &PubSub{
		pubsub_internal.New(solanaClient, ephemeralClient, &uuidGenerator{}, &zstdEncoder{}),
		signer.PublicKey().String(),
	}
}

func (p *PubSub) Subscribe(topic string, handler Handler) {
	p.PubSub.Subscribe(topic, func(ctx context.Context, event pubsub_internal.Event) {
		handler(ctx, Event(event))
	})
}

func (p *PubSub) GetPeerId() string {
	return p.peerId
}
