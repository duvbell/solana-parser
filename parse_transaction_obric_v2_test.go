package solanaparser

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"os"
	"testing"
)

func TestTransaction_ObricV2_ParseSwap(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("4PpWCbNeUwcL1ntTj23K9M1uNozVjeN1jbUS8HREPHRDTNxy41xf3EPt3CVYekynnQkWYUgTKfYppb7KEZy4NT2q"),
		&rpc.GetTransactionOpts{
			Commitment:                     rpc.CommitmentConfirmed,
			MaxSupportedTransactionVersion: &rpc.MaxSupportedTransactionVersion1,
		})
	if err != nil {
		panic(err)
	}
	transaction, _ := result.Transaction.GetTransaction()
	transactionParsed := &rpc.TransactionParsed{
		Transaction: transaction,
		Meta:        result.Meta,
	}
	txRawJson, _ := json.MarshalIndent(transactionParsed, "", "    ")
	os.WriteFile(fmt.Sprintf("tx_raw.json"), txRawJson, 0644)
	tx := ParseTransaction(0, transaction, result.Meta)
	if tx == nil {
		panic("invalid transaction")
	}
	txJson, _ := json.MarshalIndent(tx, "", "    ")
	os.WriteFile(fmt.Sprintf("tx.json"), txJson, 0644)
}
