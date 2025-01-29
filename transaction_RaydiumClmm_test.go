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

func TestTransaction_RaydiumClmm_SwapV2(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetParsedTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("254kJ7VTZJ4tuU7YGrpD1quntbAxQuSR2kmUUgSmEKE7fNv3PgdT3jDZUkcUSe8vr4Fb5wum5RxQaUo1ZdVsW75T"),
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
	txRawJson, _ := json.MarshalIndent(transaction, "", "    ")
	os.WriteFile(fmt.Sprintf("tx_raw.json"), txRawJson, 0644)
	tx, err := ParseTransaction(0, transaction)
	if err != nil {
		panic(err)
	}
	txJson, _ := json.MarshalIndent(tx, "", "    ")
	os.WriteFile(fmt.Sprintf("tx.json"), txJson, 0644)
}

func TestTransaction_RaydiumClmm_Swap(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetParsedTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("oatTsPD1GHEs9gWXdRbmaZf6EJK44qY7awqptR2WdrjH3qxBKHJYG11Yv9RX3bhciG73CYhGrFGoFmboRbxDo18"),
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
	txRawJson, _ := json.MarshalIndent(transaction, "", "    ")
	os.WriteFile(fmt.Sprintf("tx_raw.json"), txRawJson, 0644)
	tx, err := ParseTransaction(0, transaction)
	if err != nil {
		panic(err)
	}
	txJson, _ := json.MarshalIndent(tx, "", "    ")
	os.WriteFile(fmt.Sprintf("tx.json"), txJson, 0644)
}

func TestTransaction_RaydiumClmm_IncreaseLiquidityV2(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetParsedTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("4HpE61GY92gc5pgfsRsUd8tf2fMuq1nvFJXgkjTgJp2sWyQBoYsnubJGYeyqLz25nmHy8LyFtVu8XNNy3ZqCTap6"),
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
	txRawJson, _ := json.MarshalIndent(transaction, "", "    ")
	os.WriteFile(fmt.Sprintf("tx_raw.json"), txRawJson, 0644)
	tx, err := ParseTransaction(0, transaction)
	if err != nil {
		panic(err)
	}
	txJson, _ := json.MarshalIndent(tx, "", "    ")
	os.WriteFile(fmt.Sprintf("tx.json"), txJson, 0644)
}
