package contract_client

import (
	"context"

	"github.com/gagliardetto/solana-go"
)

type MessageData struct {
	TimeStamp int64  `borsh:"timestamp"`
	Read      bool   `borsh:"read"`
	Content   []byte `borsh:"content"`
}

func (c *SolanaClient) GetIncomingMessageData(ctx context.Context, sender solana.PublicKey) (MessageData, error) {
	return c.getMessageData(ctx, c.getIncomingMessagePubkey(sender))
}
