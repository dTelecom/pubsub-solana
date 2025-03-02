solana config set --url localhost
go run ./cmd/main.go --wallet ~/.config/solana/id.json --init-node-counter
go run ./cmd/main.go --wallet ~/.config/solana/id.json --register-node --node-key $(<./node-1-pubkey)
go run ./cmd/main.go --wallet ~/.config/solana/id.json --register-node --node-key $(<./node-2-pubkey)
solana airdrop 5 -k ./node-1-wallet.json
go run ./cmd/main.go --wallet ./node-1-wallet.json --init-message-account --sender-key $(<./node-2-pubkey)
solana airdrop 5 -k ./node-2-wallet.json
go run ./cmd/main.go --wallet ./node-2-wallet.json --init-message-account --sender-key $(<./node-1-pubkey)
