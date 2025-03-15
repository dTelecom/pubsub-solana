package pubsub

import (
	"context"

	"github.com/gagliardetto/solana-go"

	"github.com/dTelecom/pubsub-solana/internal/contract_client"
)

//go:generate ../../bin/mockgen -source $GOFILE -destination=mocks/mocks.go -package mocks

type ContractClient interface {
	// IsSigner проверяет является ли переданный ключ публичным ключом клиента
	IsSigner(key solana.PublicKey) bool
	// SendMessage отправляет сообщение получателю
	SendMessage(ctx context.Context, receiver solana.PublicKey, content []byte) (solana.Signature, error)
	// MarkAsRead помечает сообщение как прочитанное
	MarkAsRead(ctx context.Context, sender solana.PublicKey, timestamp int64) (solana.Signature, error)
	// IncomingMessageSubscribe делает подписку на изменение данных во входящих сообщениях
	IncomingMessageSubscribe(ctx context.Context, sender solana.PublicKey, handler func(context.Context, contract_client.MessageData))
	// OutgoingMessageSubscribe делает подписку на изменение данных в исходящих сообщениях
	OutgoingMessageSubscribe(ctx context.Context, receiver solana.PublicKey, handler func(context.Context, contract_client.MessageData))
}

type MessageIdGenerator interface {
	Generate() string
}

type DataEncoder interface {
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
}
