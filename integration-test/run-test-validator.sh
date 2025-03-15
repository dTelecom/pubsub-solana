DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

solana-test-validator \
  -r \
  --account mAGicPQYBMvcYveUZA5F5UNNwyHvfYh5xkLS2Fr1mev \
  $DIR/accounts/validator-authority.json \
  --account LUzidNSiPNjYNkxZcUm5hYHwnWPwsUfh2US1cpWwaBm \
  $DIR/accounts/luzid-authority.json \
  --limit-ledger-size \
  1000000 \
  --bpf-program \
  DELeGGvXpWV2fqJUhqcF5ZSYMS4JTLjteaAMARRSaeSh \
  $DIR/elfs/dlp.so \
  --bpf-program 5NTgsFtbN9X3XEjepBeYbzRcwnguaE7niz74QhRBMqGU \
  $DIR/elfs/messenger.so
