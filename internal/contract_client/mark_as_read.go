package contract_client

import (
	"context"

	"github.com/gagliardetto/solana-go"
)

// MarkAsRead marking the message as read
func (c *SolanaClient) MarkAsRead(ctx context.Context, sender solana.PublicKey, timestamp int64) (solana.Signature, error) {
	messagePubkey, _, _ := solana.FindProgramAddress([][]byte{
		[]byte("message"),
		sender.Bytes(),
		c.signer.PublicKey().Bytes(),
	}, programID)

	args, err := encodeArgs("mark_as_read", struct {
		Sender    solana.PublicKey
		Timestamp int64
	}{
		Sender:    sender,
		Timestamp: timestamp,
	})
	if err != nil {
		return zeroSignature, err
	}

	blockhash, err := c.getLatestBlockhash(ctx)
	if err != nil {
		return zeroSignature, err
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			solana.NewInstruction(
				programID,
				solana.AccountMetaSlice{
					{PublicKey: messagePubkey, IsWritable: true, IsSigner: false},
					{PublicKey: c.signer.PublicKey(), IsWritable: true, IsSigner: true},
				},
				args,
			),
		},
		blockhash,
		solana.TransactionPayer(c.signer.PublicKey()),
	)
	if err != nil {
		return zeroSignature, err
	}

	return c.sendTransaction(ctx, tx)
}
