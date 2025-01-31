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

func TestTransaction_MeteoraDLMM_Swap(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetParsedTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("2nJG8CNyj9HHpMx5X1WfQ1PAtnTMevg4XMWyKA8GANTFqB8xYT8ZLUqP9k5y7GrkaFjvtRiiyiXhN8HQszjjRVhb"),
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

func TestTransaction_MeteoraDLMM_AddLiquidityByStrategy(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetParsedTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("3ZarkkkDBgKSJBC42uXZT2ZgSGvZ4NzZcPTRsd9tMesVoPXhJUd9mc677MEckGDbQYRwibZJNb8gfjTEuR4Jw2j8"),
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

func TestTransaction_MeteoraDLMM_AddLiquidityByStrategy2(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetParsedTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("44Qd6qaTmK2bGhYw5nTukwKe3y3Uj5sRQcRoPQtxfksXahaSsDm2NULcxbagLmFnZGJWNyR2s8ta8EtyWfAb48NM"),
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

func TestTransaction_MeteoraDLMM_RemoveLiquidityByRange(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetParsedTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("4PS6JiPFv2peTwaGq6AtxR6hj4TeaVNrjL5J3cqRDw3Sc4r7PGLBANdxfqxr6nFA41R6PzM6wW7ijWcE6EieajnG"),
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

func TestTransaction_MeteoraDLMM_SwapExactOut(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetParsedTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("ktF53zG1Ke5a38QRwbhfDpoaMcC7d64LiypJxjqbnxC7QRUj9dB1Nem9xN92mfPrJVbpgJe3ndWzbRLG71a3gvg"),
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

func TestTransaction_MeteoraDLMM_ParseAddLiquidityByWeight(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetParsedTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("4aY1Qshhftiwy4o8UisyWL8cZ7d5YMZwc2iUdntCEKyF4s76myprRvigL2FJir7sMHwuSFfh88SiUTkEGzkjgfiA"),
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
