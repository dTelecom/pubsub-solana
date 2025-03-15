package contract_client

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

func (c *SolanaClient) WaitForProgramReadiness(ctx context.Context) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			info, err := c.rpcClient.GetAccountInfoWithOpts(
				ctx,
				programID,
				&rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentFinalized},
			)
			if err == nil && info != nil && len(info.Bytes()) > 0 {
				log.Println("Program is active and ready for work!")
				return nil
			}
			log.Printf("Waiting program readiness. Cause: %s\n", err)
		}
	}
}

func (c *SolanaClient) WaitForTransactionConfirmation(ctx context.Context, signature solana.Signature) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			status, err := c.rpcClient.GetSignatureStatuses(ctx, false, signature)
			if err != nil {
				return err
			}

			if status.Value[0] != nil && status.Value[0].ConfirmationStatus == rpc.ConfirmationStatusType(c.commitment) {
				log.Println("Transaction is confirmed!")
				return nil
			}

			log.Println("Waiting confirmation of transaction Ожидание подтверждения транзакции...")
		}
	}
}

func (c *SolanaClient) getLatestBlockhash(ctx context.Context) (solana.Hash, error) {
	blockhashResp, err := c.rpcClient.GetLatestBlockhash(ctx, c.commitment)
	if err != nil {
		return solana.Hash{}, fmt.Errorf("❌ Get latest blockhash error: %v", err)
	}
	return blockhashResp.Value.Blockhash, nil
}

func getDiscriminator(functionName string) []byte {
	hash := sha256.Sum256([]byte("global:" + functionName))
	return hash[:8]
}

func encodeArgs(method string, data interface{}) ([]byte, error) {
	args, err := borsh.Serialize(data)
	if err != nil {
		return nil, err
	}
	return append(getDiscriminator(method), args...), nil
}

func (c *SolanaClient) sendTransaction(ctx context.Context, tx *solana.Transaction) (solana.Signature, error) {
	_, err := tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if key == c.signer.PublicKey() {
				return &c.signer
			}
			return nil
		},
	)
	if err != nil {
		return zeroSignature, fmt.Errorf("❌ Transaction sign error: %v", err)
	}

	signagure, err := c.rpcClient.SendTransactionWithOpts(
		ctx,
		tx,
		rpc.TransactionOpts{
			SkipPreflight:       true,
			PreflightCommitment: "",
			MaxRetries:          pointer.ToUint(5),
		})
	if err != nil {
		return zeroSignature, fmt.Errorf("❌ Transaction send error: %v", err)
	}

	fmt.Println("✅ Transaction has been sent!")
	return signagure, nil
}

func (c *SolanaClient) getIncomingMessagePubkey(sender solana.PublicKey) solana.PublicKey {
	// todo: do not ignore error
	messagePubkey, _, _ := solana.FindProgramAddress([][]byte{
		[]byte("message"),
		sender.Bytes(),
		c.signer.PublicKey().Bytes(),
	}, programID)

	return messagePubkey
}

func (c *SolanaClient) getOutgoingMessagePubkey(receiver solana.PublicKey) solana.PublicKey {
	// todo: do not ignore error
	messagePubkey, _, _ := solana.FindProgramAddress([][]byte{
		[]byte("message"),
		c.signer.PublicKey().Bytes(),
		receiver.Bytes(),
	}, programID)

	return messagePubkey
}

func (c *SolanaClient) getMessageData(ctx context.Context, messagePubKey solana.PublicKey) (MessageData, error) {
	var res MessageData

	accountInfo, err := c.rpcClient.GetAccountInfoWithOpts(
		ctx,
		messagePubKey,
		&rpc.GetAccountInfoOpts{
			Commitment: c.commitment,
			DataSlice:  nil,
		})
	if err != nil {
		return res, fmt.Errorf("Cannot get MessageData account: %w", err)
	}

	if err := borsh.Deserialize(&res, accountInfo.Value.Data.GetBinary()[8:]); err != nil {
		return res, fmt.Errorf("Deserialization of MessageData error: %w", err)
	}
	return res, nil
}
