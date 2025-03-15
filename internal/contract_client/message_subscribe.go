package contract_client

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/near/borsh-go"
)

func (c *SolanaClient) IncomingMessageSubscribe(ctx context.Context, sender solana.PublicKey, handler func(context.Context, MessageData)) {
	c.messageSubscribe(ctx, c.getIncomingMessagePubkey(sender), handler)
}

func (c *SolanaClient) OutgoingMessageSubscribe(ctx context.Context, sender solana.PublicKey, handler func(context.Context, MessageData)) {
	c.messageSubscribe(ctx, c.getOutgoingMessagePubkey(sender), handler)
}

func (c *SolanaClient) messageSubscribe(ctx context.Context, messagePubkey solana.PublicKey, handler func(context.Context, MessageData)) {
	for k, v := range c.logContext {
		ctx = context.WithValue(ctx, k, v)
	}

	go func() {
	mainFor:
		for {
			wsClient, err := c.getWSClient(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				log.Printf("failed to get websocket client: %s\n", err)
				continue
			}

			sub, err := wsClient.AccountSubscribe(messagePubkey, c.commitment)
			if err != nil {
				log.Printf("message subscription error: %s\n", err)
				c.clearWSClient(wsClient)
				continue
			}

			if c.isEphemeral {
				// Airdrop to trigger lazy reload
				_, _ = c.rpcClient.RequestAirdrop(ctx, messagePubkey, 1, "")
				if err != nil {
					sub.Unsubscribe()
					if errors.Is(err, context.Canceled) {
						return
					}
					log.Printf("Failed to lazy reload message account: %s", err)
					continue
				}
			}

			initialMessageData, err := c.getMessageData(ctx, messagePubkey)
			if err != nil {
				// sub.Unsubscribe()
				if errors.Is(err, context.Canceled) {
					return
				}
				// log.Printf("Cannot get initial message data: %s\n", err)
				// continue
			}

			handler(ctx, initialMessageData)

			for {
				got, err := sub.Recv(ctx)
				if err != nil {
					sub.Unsubscribe()
					if errors.Is(err, context.Canceled) {
						return
					}
					log.Printf("Receive error: %v", err)
					// websocket connection is probably broken
					c.clearWSClient(wsClient)
					continue mainFor
				}

				data := got.Value.Data.GetBinary()
				data = data[8:] // remove descriminator

				var messageData MessageData
				if err := borsh.Deserialize(&messageData, data); err != nil {
					log.Printf("Cannot deserialize message data: %v", err)
				} else {
					handler(ctx, messageData)
				}
			}
		}
	}()
}

func (c *SolanaClient) getWSClient(ctx context.Context) (*ws.Client, error) {
	c.wsMu.Lock()
	defer c.wsMu.Unlock()

	if c.wsClient == nil {
		wsClient, err := ws.Connect(ctx, c.wsURL)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to websocket endpoint: %w", err)
		}
		c.wsClient = wsClient
	}

	return c.wsClient, nil
}

func (c *SolanaClient) clearWSClient(oldClient *ws.Client) {
	c.wsMu.Lock()
	defer c.wsMu.Unlock()

	if c.wsClient == oldClient {
		c.wsClient = nil
	}
}
