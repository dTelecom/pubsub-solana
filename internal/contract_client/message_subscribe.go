package contract_client

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

func (c *SolanaClient) IncomingMessageSubscribe(ctx context.Context, sender solana.PublicKey, handler func(context.Context, MessageData)) error {
	return c.messageSubscribe(ctx, c.getIncomingMessagePubkey(sender), handler)
}

func (c *SolanaClient) OutgoingMessageSubscribe(ctx context.Context, sender solana.PublicKey, handler func(context.Context, MessageData)) error {
	return c.messageSubscribe(ctx, c.getOutgoingMessagePubkey(sender), handler)
}

func (c *SolanaClient) messageSubscribe(ctx context.Context, messagePubkey solana.PublicKey, handler func(context.Context, MessageData)) error {
	for k, v := range c.logContext {
		ctx = context.WithValue(ctx, k, v)
	}

	sub, err := c.wsClient.AccountSubscribe(messagePubkey, c.commitment)
	if err != nil {
		return fmt.Errorf("Message subscription error: %s", err)
	}

	if c.isEphemeral {
		// Airdrop to trigger lazy reload
		_, _ = c.rpcClient.RequestAirdrop(ctx, messagePubkey, 1, "")
		// if err != nil {
		// 	return fmt.Errorf("Failed to lazy reload message account: %s", err)
		// }
	}

	initialMessageData, err := c.getMessageData(ctx, messagePubkey)
	if err != nil {
		sub.Unsubscribe()

		return fmt.Errorf("Cannot get initial message data: %s", err)
	}

	go func() {
		defer sub.Unsubscribe()

		handler(ctx, initialMessageData)

		for {
			got, err := sub.Recv(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					break
				}
				log.Printf("Receive error: %v", err)
				continue
			}

			data := got.Value.Data.GetBinary()
			data = data[8:] // remove descriminator

			var messageData MessageData
			if err := borsh.Deserialize(&messageData, data); err != nil {
				log.Printf("Cannot deserialize message data: %v", err)
			}

			handler(ctx, messageData)
		}
	}()

	return nil
}
