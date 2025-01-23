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

func TestTransaction_Parse(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetParsedTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("5JkHtCDnJRVYrHYg883BdkKnvZe2iPwW7LooTFSu2yjMPinbcRFTXuKZ3A8tEpz8N82Tn9nudxoTSemzLa19nKho"),
		&rpc.GetParsedTransactionOpts{
			Commitment:                     rpc.CommitmentConfirmed,
			MaxSupportedTransactionVersion: &rpc.MaxSupportedTransactionVersion1,
		})
	if err != nil {
		panic(err)
	}
	transaction := &rpc.ParsedTransactionWithMeta{
		Slot:        result.Slot,
		BlockTime:   result.BlockTime,
		Transaction: result.Transaction,
		Meta:        result.Meta,
	}
	tx := NewTransaction()
	err = tx.Parse(transaction)
	if err != nil {
		panic(err)
	}
	err = tx.ParseActions(DefaultParse)
	if err != nil {
		panic(err)
	}
	txJson, _ := json.MarshalIndent(tx, "", "    ")
	os.WriteFile(fmt.Sprintf("tx.json"), txJson, 0644)
}
