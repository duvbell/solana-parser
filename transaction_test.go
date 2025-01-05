package transaction

import (
	"context"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"testing"
)

func TestTransaction_Parse(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetParsedTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("ofFBWuFZYmT1egbjaRMG3NMwJ3TE4YE4ssDjVUsTGM9iNPECkjVfSeyvnSySmKjvsYV6gXPSrwEvvmboicQVho7"),
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
	tx := New()
	err = tx.Parse(transaction)
	if err != nil {
		panic(err)
	}
	err = tx.ParseActions(nil)
	if err != nil {
		panic(err)
	}
}
