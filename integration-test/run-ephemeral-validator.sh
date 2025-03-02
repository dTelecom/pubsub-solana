docker run -e ACCOUNTS_REMOTE=https://api.devnet.solana.com \
           -e VALIDATOR_MILLIS_PER_SLOT=50 \
           -e ACCOUNTS_LIFECYCLE=ephemeral \
           -p 8899:8899 -p 8900:8900 -p 10000:10000 \
           magicblocklabs/validator

