package contract_client

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

// SendMessage put message to message account
func (c *SolanaClient) SendMessage(ctx context.Context, receiver solana.PublicKey, content []byte) (solana.Signature, error) {
	c.logger.Debugw("➡️ Sending message to %v\n", receiver)

	if len(content) > 887 {
		return zeroSignature, fmt.Errorf("message is too long (max. 887 bytes): %v", len(content))
	}

	messagePubkey, _, _ := solana.FindProgramAddress([][]byte{
		[]byte("message"),
		c.signer.PublicKey().Bytes(),
		receiver.Bytes(),
	}, programID)

	args, err := encodeArgs("send_message", struct {
		Receiver solana.PublicKey
		Content  []byte
	}{
		Receiver: receiver,
		Content:  content,
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
					{PublicKey: receiver, IsWritable: false, IsSigner: false},
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
