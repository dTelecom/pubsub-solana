package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DataDog/zstd"
	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"

	"github.com/dTelecom/pubsub-solana/internal/contract_client"
	"github.com/dTelecom/pubsub-solana/internal/pubsub"
)

const (
	solanaRPC     = "https://api.devnet.solana.com"
	solanaWS      = "wss://api.devnet.solana.com/"
	magicblockRPC = "https://devnet.magicblock.app"
	magicblockWS  = "wss://devnet.magicblock.app/"

	// solanaRPC     = "http://localhost:8899"
	// solanaWS      = "ws://localhost:8900/"
	// magicblockRPC = "http://localhost:8899"
	// magicblockWS  = "ws://localhost:8900/"
)

type uuidGenerator struct{}

func (g *uuidGenerator) Generate() string {
	return uuid.New().String()
}

type zstdEncoder struct{}

func (*zstdEncoder) Encode(src []byte) ([]byte, error) {
	return zstd.Compress(nil, src)
}

func (*zstdEncoder) Decode(src []byte) ([]byte, error) {
	return zstd.Decompress(nil, src)
}

func main() {
	initMessageAccountFlag := flag.Bool("init-message-account", false, "InitMessageAccount")
	delegateMessageAccountFlag := flag.Bool("delegate-message-account", false, "DelegateMessageAccount")

	subscribeFlag := flag.Bool("subscribe", false, "Subscribe")
	publishFlag := flag.Bool("publish", false, "Publish")

	walletFlag := flag.String("wallet", "", "Wallet")
	senderKeyFlag := flag.String("sender-key", "", "Sender key")
	nodeKeyFlag := flag.String("node-key", "", "Node key")
	topicFlag := flag.String("topic", "", "Topic")

	flag.Parse()

	if *walletFlag == "" {
		log.Fatal("Wallet is necessary")
	}

	if (*subscribeFlag || *publishFlag) && *topicFlag == "" {
		log.Fatal("Topic is necessary")
	}

	signer, err := solana.PrivateKeyFromSolanaKeygenFile(*walletFlag)
	if err != nil {
		log.Fatalf("Ошибка загрузки кошелька: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	contractClient := contract_client.New(false, solanaRPC, solanaWS, signer)

	if *initMessageAccountFlag || *delegateMessageAccountFlag {
		switch {
		case *initMessageAccountFlag:
			initMessageAccount(ctx, contractClient, solana.MustPublicKeyFromBase58(*senderKeyFlag))
		case *delegateMessageAccountFlag:
			delegateMessageAccount(ctx, contractClient, solana.MustPublicKeyFromBase58(*senderKeyFlag))
		}
	} else {
		contractMagicblockClient := contract_client.New(true, magicblockRPC, magicblockWS, signer)
		ps := pubsub.New(contractClient, contractMagicblockClient, &uuidGenerator{}, &zstdEncoder{})
		if err := ps.Start(ctx, []solana.PublicKey{solana.MustPublicKeyFromBase58(*nodeKeyFlag)}); err != nil {
			log.Fatal(err)
		}

		switch {
		case *subscribeFlag:
			subscribe(ctx, ps, *topicFlag)
		case *publishFlag:
			publish(ctx, ps, *topicFlag)
		}
	}
}

func initMessageAccount(ctx context.Context, contractClient *contract_client.SolanaClient, sender solana.PublicKey) {
	fmt.Println("Initializing message account...")
	signature, err := contractClient.InitMessageAccount(ctx, sender)
	if err != nil {
		log.Fatal(err)
	}

	if err := contractClient.WaitForTransactionConfirmation(ctx, signature); err != nil {
		log.Fatal(err)
	}
}

func delegateMessageAccount(ctx context.Context, contractClient *contract_client.SolanaClient, sender solana.PublicKey) {
	fmt.Println("Delegating message account...")
	signature, err := contractClient.DelegateMessageAccount(ctx, sender)
	if err != nil {
		log.Fatal(err)
	}

	if err := contractClient.WaitForTransactionConfirmation(ctx, signature); err != nil {
		log.Fatal(err)
	}
}

func subscribe(ctx context.Context, ps *pubsub.PubSub, topic string) {
	ps.Subscribe(topic, func(ctx context.Context, event pubsub.Event) {
		fmt.Printf("%v (ID: %v, FromPeerId: %v, wsURL: %v)\n", string(event.Message), event.ID, event.FromPeerId, ctx.Value("wsURL"))
	})

	fmt.Printf("Subscribed to topic: %v\n", topic)
	<-ctx.Done()
}

func publish(ctx context.Context, ps *pubsub.PubSub, topic string) {
	fmt.Printf("Publishing to topic: %v. Write something and press Enter.\n", topic)

	inputCh := make(chan string)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			inputCh <- scanner.Text()
		}
	}()

	for {
		select {
		case line := <-inputCh:
			id, err := ps.Publish(topic, []byte(line))
			if err != nil {
				return
			}
			fmt.Printf("message=%v, id=%v\n", line, id)
		case <-ctx.Done():
			return
		}
	}
}
