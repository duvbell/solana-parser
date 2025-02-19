package solanaparser

import (
	"context"
	"encoding/json"
	"fmt"
	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/pumpfun"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"
	"os"
	"testing"
)

func TestTransaction_PumpFun_ParseCreate(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetParsedTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("2gxXCfpeU4DLPDRjqp5crVvPfCxxuQijp6CB4TgMsUAujLbL9hdRi9enZosdK3bCi7eGeEfsX2vgb8W59Mr6gjqw"),
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
	tx := ParseTransaction(0, transaction)
	if tx == nil {
		panic("invalid transaction")
	}
	txJson, _ := json.MarshalIndent(tx, "", "    ")
	os.WriteFile(fmt.Sprintf("tx.json"), txJson, 0644)
}

func TestTransaction_PumpFun_ParseBuy(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetParsedTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("2gxXCfpeU4DLPDRjqp5crVvPfCxxuQijp6CB4TgMsUAujLbL9hdRi9enZosdK3bCi7eGeEfsX2vgb8W59Mr6gjqw"),
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
	tx := ParseTransaction(0, transaction)
	if tx == nil {
		panic("invalid transaction")
	}
	txJson, _ := json.MarshalIndent(tx, "", "    ")
	os.WriteFile(fmt.Sprintf("tx.json"), txJson, 0644)
}

func TestTransaction_PumpFun_ParseSell(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetParsedTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("45QxHzAfEqnTSuHkQARHbHwJFve37e1fBqsaGGbuYEhrPmpFZdQm4a1c4yrdu3rU7QR9XWfBLw1Vkr4PmpuxVmSZ"),
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
	tx := ParseTransaction(0, transaction)
	if tx == nil {
		panic("invalid transaction")
	}
	txJson, _ := json.MarshalIndent(tx, "", "    ")
	os.WriteFile(fmt.Sprintf("tx.json"), txJson, 0644)
}

func TestTransaction_PumpFun_Event(t *testing.T) {
	data, _ := base58.Decode("FmHeqAfJHjYV3ASURwexkULoe3H2ujsmgUtyL6XXNH9HaMXZVLFipht3eC1mNpTACUNiStJn9ShQg7f4NHG2XMj3fXGcedrVRsM3AoE3jXpg3wE7XJhPUXQnyYyhj7iryk267UvfidwHQXgXRdrVP5YLXUdq9kbD6zadBk5XwKctV5JX4RbdxfAXX1fF6JHmuwWtC6fnGPLFUEAByXdye5Wa2jf3jJiwMHs4ALwbLhWZJkXXPt8tBftaKevkT6oFVw3PKuumPZ2HyMb3YMzfyfW8yUx")
	var tradeEvent pumpfun.CreateEvent

	dec := ag_binary.NewBorshDecoder(data)
	dec.ReadTypeID()
	dec.ReadBytes(8)

	if err := dec.Decode(&tradeEvent); err != nil {
		panic(err)
	}
	eventJson, _ := json.MarshalIndent(tradeEvent, "", "    ")
	fmt.Printf("evnet: %s\n", string(eventJson))
}
