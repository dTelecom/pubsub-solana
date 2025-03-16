package contract_client

import (
	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"

	"github.com/dTelecom/pubsub-solana/internal/common"
)

// Smart contract address
var programID = solana.MustPublicKeyFromBase58("5NTgsFtbN9X3XEjepBeYbzRcwnguaE7niz74QhRBMqGU")

var zeroSignature = solana.Signature{}

// SolanaClient - structure for working with Solana
type SolanaClient struct {
	logger        common.Logger
	isEphemeral   bool
	counterPubkey solana.PublicKey
	rpcClient     *rpc.Client
	wsURL         string
	wsClient      *ws.Client
	wsMu          sync.Mutex
	signer        solana.PrivateKey
	commitment    rpc.CommitmentType

	logContext map[string]string
}

// New create new client
func New(logger common.Logger, isEphemeral bool, rpcURL, wsURL string, signer solana.PrivateKey) *SolanaClient {
	counterPubkey, _, _ := solana.FindProgramAddress([][]byte{[]byte("node_counter")}, programID)

	return &SolanaClient{
		logger:        logger,
		isEphemeral:   isEphemeral,
		counterPubkey: counterPubkey,
		rpcClient:     rpc.New(rpcURL),
		wsURL:         wsURL,

		signer:     signer,
		commitment: rpc.CommitmentConfirmed,

		logContext: map[string]string{"rpcURL": rpcURL, "wsURL": wsURL},
	}
}
