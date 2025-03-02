package contract_client

import (
	"context"

	"github.com/gagliardetto/solana-go"
)

// InitMessageAccount initialization of the message account
func (c *SolanaClient) InitMessageAccount(ctx context.Context, sender solana.PublicKey) (solana.Signature, error) {
	messagePubkey, _, _ := solana.FindProgramAddress([][]byte{
		[]byte("message"),
		sender.Bytes(),
		c.signer.PublicKey().Bytes(),
	}, programID)

	args, err := encodeArgs("initialize_message_account", sender)
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
					{PublicKey: solana.SystemProgramID, IsWritable: false, IsSigner: false},
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
