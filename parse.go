package solanaparser

import (
	"encoding/json"
	"github.com/blockchain-develop/solana-parser/log"
	"github.com/blockchain-develop/solana-parser/program"
	_ "github.com/blockchain-develop/solana-parser/program/lifinity"
	_ "github.com/blockchain-develop/solana-parser/program/meteora_dlmm"
	_ "github.com/blockchain-develop/solana-parser/program/meteora_pools"
	_ "github.com/blockchain-develop/solana-parser/program/meteora_vault"
	_ "github.com/blockchain-develop/solana-parser/program/phoenix"
	_ "github.com/blockchain-develop/solana-parser/program/pump"
	_ "github.com/blockchain-develop/solana-parser/program/raydium_amm"
	_ "github.com/blockchain-develop/solana-parser/program/raydium_clmm"
	_ "github.com/blockchain-develop/solana-parser/program/raydium_cp"
	_ "github.com/blockchain-develop/solana-parser/program/spl_token"
	_ "github.com/blockchain-develop/solana-parser/program/spl_token_2022"
	_ "github.com/blockchain-develop/solana-parser/program/stable_swap"
	_ "github.com/blockchain-develop/solana-parser/program/stable_vault"
	_ "github.com/blockchain-develop/solana-parser/program/system"
	_ "github.com/blockchain-develop/solana-parser/program/whirlpool"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/shopspring/decimal"
)

func ParseBlock(slot uint64, b *rpc.GetParsedBlockResult) *types.Block {
	log.Logger.Info("parse block", "slot", slot, "hash", b.Blockhash.String())
	block := &types.Block{}
	block.Slot = slot
	block.Time = uint64(*b.BlockTime)
	block.Hash = b.Blockhash
	myTxs := make([]*types.Transaction, 0)
	for i, tx := range b.Transactions {
		myTx := ParseTransaction(i+1, &tx)
		if myTx == nil {
			continue
		}
		if len(myTx.Instructions) == 0 {
			// no instruction
			continue
		}
		myTxs = append(myTxs, myTx)
	}
	block.Transaction = myTxs
	return block
}

func ParseTransaction(seq int, tx *rpc.ParsedTransactionWithMeta) *types.Transaction {
	log.Logger.Info("parse transaction", "seq", seq, "tx", tx.Transaction.Signatures[0].String())
	if tx.Meta == nil || tx.Transaction == nil {
		log.Logger.Error("parse transaction: meta or transaction is missing")
		return nil
	}
	t := &types.Transaction{
		Meta: types.Meta{
			Accounts:    make(map[solana.PublicKey]*solana.AccountMeta),
			TokenOwner:  make(map[solana.PublicKey]solana.PublicKey),
			TokenMint:   make(map[solana.PublicKey]solana.PublicKey),
			PreBalance:  make(map[solana.PublicKey]decimal.Decimal),
			PostBalance: make(map[solana.PublicKey]decimal.Decimal),
		},
		Seq: seq,
	}
	meta := tx.Meta
	transaction := tx.Transaction
	t.Hash = transaction.Signatures[0]
	if meta.Err != nil {
		// if failed, ignore this transaction
		errJson, _ := json.Marshal(meta.Err)
		t.Meta.ErrorMessage = errJson
		return t
	}
	message := transaction.Message
	instructions := message.Instructions
	if len(instructions) == 0 {
		return t
	}
	if instructions[0].ProgramId == solana.VoteProgramID {
		return t
	}
	// account infos
	for _, item := range message.AccountKeys {
		t.Meta.Accounts[item.PublicKey] = &solana.AccountMeta{
			PublicKey:  item.PublicKey,
			IsWritable: item.Writable,
			IsSigner:   item.Signer,
		}
	}
	for _, item := range meta.PostTokenBalances {
		account := message.AccountKeys[item.AccountIndex]
		t.Meta.TokenOwner[account.PublicKey] = *item.Owner
		t.Meta.TokenMint[account.PublicKey] = item.Mint
		t.Meta.PostBalance[account.PublicKey], _ = decimal.NewFromString(item.UiTokenAmount.Amount)
	}
	for _, item := range meta.PreTokenBalances {
		account := message.AccountKeys[item.AccountIndex]
		t.Meta.PreBalance[account.PublicKey], _ = decimal.NewFromString(item.UiTokenAmount.Amount)
	}
	for index, instruction := range instructions {
		instruction.StackHeight = 1
		myInstruction := &types.Instruction{
			Seq:         index + 1,
			Instruction: instruction,
			Event:       nil,
			Receipt:     nil,
			Children:    nil,
		}
		t.Instructions = append(t.Instructions, myInstruction)
	}
	innerInstructions := meta.InnerInstructions
	for _, innerInstruction := range innerInstructions {
		parent := t.Instructions[innerInstruction.Index]
		build(parent, innerInstruction.Instructions)
	}
	for _, instruction := range t.Instructions {
		parse(instruction, &t.Meta)
	}
	return t
}

func split(subIns []*rpc.ParsedInstruction) []int {
	currentHeight := subIns[0].StackHeight
	indexes := make([]int, 0)
	for index, item := range subIns {
		if item.StackHeight == currentHeight {
			indexes = append(indexes, index)
		}
	}
	return indexes
}

func build(parent *types.Instruction, subIns []*rpc.ParsedInstruction) {
	if len(subIns) == 0 {
		return
	}
	// ins split by stack height
	indexes := split(subIns)
	indexes = append(indexes, len(subIns))
	for i := 0; i < len(indexes)-1; i++ {
		index1 := indexes[i]
		index2 := indexes[i+1]
		current := &types.Instruction{
			Seq:         i + 1,
			Instruction: subIns[index1],
			Children:    nil,
		}
		parent.Children = append(parent.Children, current)
		build(current, subIns[index1+1:index2])
	}
}

func parse(in *types.Instruction, meta *types.Meta) {
	for _, child := range in.Children {
		parse(child, meta)
	}
	program.Parse(in, meta)
}
