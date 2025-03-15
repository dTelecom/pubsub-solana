//go:build integration

// integration-test/run-test-validator.sh

package contract_client_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gagliardetto/solana-go"

	"github.com/dTelecom/pubsub-solana/internal/contract_client"
)

func TestClient_Happy(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	t.Cleanup(func() {
		cancel()
	})

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	// Create sender wallet
	// signer := solana.NewWallet()
	signer, err := solana.PrivateKeyFromSolanaKeygenFile(filepath.Join(homeDir, ".config/solana/id.json"))
	if err != nil {
		t.Fatalf("Error loading wallet: %v", err)
	}

	c := contract_client.New(false, "http://127.0.0.1:8899", "ws://127.0.0.1:8900", signer)

	t.Log("📌 Sender wallet:", signer.PublicKey().String())

	// 🔹 Waiting for smart contract deployment
	err = c.WaitForProgramReadiness(ctx)
	if err != nil {
		t.Fatalf("Failed to wait for program deployment: %v", err)
	}

	// 🔹 Message initialization
	t.Run("InitMessageAccount", func(t *testing.T) {
		signature, err := c.InitMessageAccount(ctx, signer.PublicKey())
		require.NoError(t, err, "Message initialization error")

		t.Log("Message initialized")

		// 🔹 Waiting for message to initialize
		err = c.WaitForTransactionConfirmation(ctx, signature)
		require.NoError(t, err, "Failed to wait for transaction confirmation")
	})

	const firstMessageContent = "Hello, Solana!"
	// 🔹 Sending message
	t.Run("SendMessage", func(t *testing.T) {
		signature, err := c.SendMessage(ctx, signer.PublicKey(), []byte(firstMessageContent))
		require.NoError(t, err, "Message sending error")

		t.Log("Message sent")

		// 🔹 Waiting for message to reach recipient
		err = c.WaitForTransactionConfirmation(ctx, signature)
		require.NoError(t, err, "Failed to wait for transaction confirmation")

		messageData, err := c.GetIncomingMessageData(ctx, signer.PublicKey())
		require.NoError(t, err, "Error getting data from Message")

		require.False(t, messageData.Read, "Message should be unread")
		require.Equal(t, []byte(firstMessageContent), messageData.Content, "Unexpected message text")
		require.Greater(t, messageData.TimeStamp, int64(0), "Message timestamp must be greater than zero")
	})

	// 🔹 Marking message as read
	t.Run("MarkAsRead", func(t *testing.T) {
		messageData, err := c.GetIncomingMessageData(ctx, signer.PublicKey())
		require.NoError(t, err, "Error getting data from Message")

		signature, err := c.MarkAsRead(ctx, signer.PublicKey(), messageData.TimeStamp)
		require.NoError(t, err, "Error marking message")

		t.Log("Message read")

		// 🔹 Waiting for transaction to complete
		err = c.WaitForTransactionConfirmation(ctx, signature)
		require.NoError(t, err, "Failed to wait for transaction confirmation")

		messageData, err = c.GetIncomingMessageData(ctx, signer.PublicKey())
		require.NoError(t, err, "Error getting data from Message")

		require.True(t, messageData.Read, "Message should have been marked as read")
	})

	t.Run("MessageSubscribe", func(t *testing.T) {
		ch := make(chan contract_client.MessageData, 1)
		c.IncomingMessageSubscribe(ctx, signer.PublicKey(), func(_ context.Context, msg contract_client.MessageData) {
			ch <- msg
		})

		select {
		case messageData := <-ch:
			require.Equal(t, []byte(firstMessageContent), messageData.Content, "Unexpected message text")
			require.Truef(t, messageData.Read, "Message should be read")
		case <-time.After(10 * time.Second):
			require.FailNow(t, "timeout")
		}

		const secondMessageContent = "Hello, again!!!"
		signature, err := c.SendMessage(ctx, signer.PublicKey(), []byte(secondMessageContent))
		require.NoError(t, err, "Message sending error")

		// 🔹 Waiting for message to reach recipient
		err = c.WaitForTransactionConfirmation(ctx, signature)
		require.NoError(t, err, "Failed to wait for transaction confirmation")

		select {
		case messageData := <-ch:
			require.Equal(t, []byte(secondMessageContent), messageData.Content, "Unexpected message text")
			require.Falsef(t, messageData.Read, "Message should be unread")
		case <-time.After(10 * time.Second):
			require.FailNow(t, "timeout")
		}
	})
}
