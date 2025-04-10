package contract_client

import (
	"context"
	"errors"
	"fmt"

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
			if c.isEphemeral {
				// Airdrop to trigger lazy reload
				_, _ = c.rpcClient.RequestAirdrop(ctx, messagePubkey, 1, "")
				if ctx.Err() != nil {
					return
				}
			}

			wsClient, err := c.getWSClient(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return
				}
				c.logger.Errorw("Failed to get websocket client", err)
				continue
			}

			sub, err := wsClient.AccountSubscribe(messagePubkey, c.commitment)
			if err != nil {
				c.logger.Errorw("Message subscription error", err)
				c.clearWSClient(wsClient)
				continue
			}

			initialMessageData, err := c.getMessageData(ctx, messagePubkey)
			if err != nil {
				sub.Unsubscribe()
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return
				}
				c.logger.Errorw("Cannot get initial message data", err)
				continue
			}

			handler(ctx, initialMessageData)

			for {
				got, err := sub.Recv(ctx)
				if err != nil {
					sub.Unsubscribe()
					if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
						return
					}
					c.logger.Errorw("Receive error", err)
					// websocket connection is probably broken
					c.clearWSClient(wsClient)
					continue mainFor
				}

				data := got.Value.Data.GetBinary()
				data = data[8:] // remove descriminator

				var messageData MessageData
				if err := borsh.Deserialize(&messageData, data); err != nil {
					c.logger.Warnw("Cannot deserialize message data", err)
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
