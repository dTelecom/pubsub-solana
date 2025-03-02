package contract_client

import (
	"github.com/gagliardetto/solana-go"
)

func (c *SolanaClient) IsSigner(key solana.PublicKey) bool {
	return c.signer.PublicKey().Equals(key)
}
