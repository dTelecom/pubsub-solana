package contract_client

import (
	"context"

	"github.com/gagliardetto/solana-go"
)

var delegationProgramID = solana.MustPublicKeyFromBase58("DELeGGvXpWV2fqJUhqcF5ZSYMS4JTLjteaAMARRSaeSh")

// DelegateMessageAccount delegates message account for using in magicblock
func (c *SolanaClient) DelegateMessageAccount(ctx context.Context, sender solana.PublicKey) (solana.Signature, error) {
	messagePubkey, _, _ := solana.FindProgramAddress([][]byte{
		[]byte("message"),
		sender.Bytes(),
		c.signer.PublicKey().Bytes(),
	}, programID)

	bufferMessagePubkey, _, _ := solana.FindProgramAddress([][]byte{
		[]byte("buffer"),
		messagePubkey.Bytes(),
	}, programID)

	delegationRecordMessagePubkey, _, _ := solana.FindProgramAddress([][]byte{
		[]byte("delegation"),
		messagePubkey.Bytes(),
	}, delegationProgramID)

	delegationMetadataMessagePubkey, _, _ := solana.FindProgramAddress([][]byte{
		[]byte("delegation-metadata"),
		messagePubkey.Bytes(),
	}, delegationProgramID)

	args, err := encodeArgs("delegate_message_account", sender)
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
					{PublicKey: c.signer.PublicKey(), IsWritable: true, IsSigner: true}, // receiver
					{PublicKey: bufferMessagePubkey, IsWritable: true},                  // buffer_message
					{PublicKey: delegationRecordMessagePubkey, IsWritable: true},        // delegation_record_message
					{PublicKey: delegationMetadataMessagePubkey, IsWritable: true},      // delegation_metadata_message,
					{PublicKey: messagePubkey, IsWritable: true},                        // message
					{PublicKey: programID},              // owner_program
					{PublicKey: delegationProgramID},    // delegation_program
					{PublicKey: solana.SystemProgramID}, // system_program
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
