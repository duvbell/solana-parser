package solanaparser

import (
	"encoding/json"

	"github.com/blockchain-develop/solana-parser/log"
	"github.com/blockchain-develop/solana-parser/program"
	_ "github.com/blockchain-develop/solana-parser/program/lifinity"
	_ "github.com/blockchain-develop/solana-parser/program/meteora_dlmm"
	_ "github.com/blockchain-develop/solana-parser/program/meteora_pools"
	_ "github.com/blockchain-develop/solana-parser/program/phoenix"
	_ "github.com/blockchain-develop/solana-parser/program/pump"
	_ "github.com/blockchain-develop/solana-parser/program/raydium_amm"
	_ "github.com/blockchain-develop/solana-parser/program/raydium_clmm"
	_ "github.com/blockchain-develop/solana-parser/program/raydium_cp"
	_ "github.com/blockchain-develop/solana-parser/program/spl_token"
	_ "github.com/blockchain-develop/solana-parser/program/spl_token_2022"
	_ "github.com/blockchain-develop/solana-parser/program/stable_swap"
	_ "github.com/blockchain-develop/solana-parser/program/system"
	_ "github.com/blockchain-develop/solana-parser/program/whirlpool"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/shopspring/decimal"
)

func ParseBlock(slot uint64, b *rpc.GetBlockResult) *types.Block {
	log.Logger.Info("parse block", "slot", slot)
	block := &types.Block{}
	block.Slot = slot
	if b == nil {
		log.Logger.Info("empty block", "slot", slot)
		return block
	}
	block.Time = uint64(*b.BlockTime)
	block.Hash = b.Blockhash
	myTxs := make([]*types.Transaction, 0)
	for i, _ := range b.Transactions {
		tx := b.Transactions[i].MustGetTransaction()
		meta := b.Transactions[i].Meta
		myTx := ParseTransaction(i+1, tx, meta)
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

func ParseTransaction(seq int, tx *solana.Transaction, meta *rpc.TransactionMeta) *types.Transaction {
	//log.Logger.Info("parse transaction", "seq", seq, "tx", tx.Transaction.Signatures[0].String())
	if meta == nil || tx == nil {
		log.Logger.Error("parse transaction: meta or transaction is missing")
		return nil
	}
	t := &types.Transaction{
		Meta: &types.Meta{
			Accounts:      make([]*solana.AccountMeta, 0),
			TokenAccounts: make(map[solana.PublicKey]*types.TokenAccount),
			MintAccounts:  make(map[solana.PublicKey]*types.MintAccount),
			PreBalance:    make(map[solana.PublicKey]decimal.Decimal),
			PostBalance:   make(map[solana.PublicKey]decimal.Decimal),
		},
		Seq: seq,
	}
	t.Hash = tx.Signatures[0]
	if meta.Err != nil {
		// if failed, ignore this transaction
		errJson, _ := json.Marshal(meta.Err)
		t.Meta.ErrorMessage = errJson
		return t
	}
	message := tx.Message
	if len(message.Instructions) == 0 {
		return t
	}
	// account infos
	readonlySignedAccountsCount := message.Header.NumReadonlySignedAccounts
	readonlyUnsignedAccountsCount := message.Header.NumReadonlyUnsignedAccounts
	requiredSignaturesAccountCount := message.Header.NumRequiredSignatures
	total := len(message.AccountKeys)
	for idx, item := range message.AccountKeys {
		isWritable := (uint8(idx) < requiredSignaturesAccountCount-readonlySignedAccountsCount) || (uint8(idx) >= requiredSignaturesAccountCount && uint8(idx) < uint8(total)-readonlyUnsignedAccountsCount)
		t.Meta.Accounts = append(t.Meta.Accounts, &solana.AccountMeta{
			PublicKey:  item,
			IsWritable: isWritable,
			IsSigner:   uint8(idx) < requiredSignaturesAccountCount,
		})
	}
	for _, item := range meta.LoadedAddresses.Writable {
		t.Meta.Accounts = append(t.Meta.Accounts, &solana.AccountMeta{
			PublicKey:  item,
			IsWritable: true,
			IsSigner:   false,
		})
	}
	for _, item := range meta.LoadedAddresses.ReadOnly {
		t.Meta.Accounts = append(t.Meta.Accounts, &solana.AccountMeta{
			PublicKey:  item,
			IsWritable: false,
			IsSigner:   false,
		})
	}
	for _, item := range meta.PostTokenBalances {
		account := t.Meta.Accounts[item.AccountIndex]
		t.Meta.TokenAccounts[account.PublicKey] = &types.TokenAccount{
			Owner:     item.Owner,
			ProgramId: item.ProgramId,
			Mint:      item.Mint,
		}
		t.Meta.MintAccounts[item.Mint] = &types.MintAccount{
			Mint:     item.Mint,
			Decimals: item.UiTokenAmount.Decimals,
		}
		t.Meta.PostBalance[account.PublicKey], _ = decimal.NewFromString(item.UiTokenAmount.Amount)
	}
	for _, item := range meta.PreTokenBalances {
		account := t.Meta.Accounts[item.AccountIndex]
		_, ok := t.Meta.TokenAccounts[account.PublicKey]
		if !ok {
			t.Meta.TokenAccounts[account.PublicKey] = &types.TokenAccount{
				Owner:     item.Owner,
				ProgramId: item.ProgramId,
				Mint:      item.Mint,
			}
			t.Meta.MintAccounts[item.Mint] = &types.MintAccount{
				Mint:     item.Mint,
				Decimals: item.UiTokenAmount.Decimals,
			}
		}
		t.Meta.PreBalance[account.PublicKey], _ = decimal.NewFromString(item.UiTokenAmount.Amount)
	}
	// add sol
	t.Meta.MintAccounts[solana.PublicKey{}] = &types.MintAccount{
		Mint:     solana.PublicKey{},
		Decimals: 9,
	}
	// todo, get the sol balance
	if t.Meta.Accounts[message.Instructions[0].ProgramIDIndex].PublicKey == solana.VoteProgramID {
		return t
	}
	log.Logger.Trace("parse transaction", "seq", seq, "tx", tx.Signatures[0].String())
	//
	instructions := make([]*types.Instruction, 0)
	index := 0
	for i, _ := range message.Instructions {
		in := program.FilterInstruction(&message.Instructions[i], t.Meta)
		if in != nil {
			in.Seq = index
			index += 1
			instructions = append(instructions, in)
		}
		find := func(index int) int {
			for j := 0; j < len(meta.InnerInstructions); j++ {
				if meta.InnerInstructions[j].Index == uint16(index) {
					return j
				}
			}
			return -1
		}
		inner := find(i)
		if inner == -1 {
			continue
		}
		for j, _ := range meta.InnerInstructions[inner].Instructions {
			innerIn := program.FilterInstruction(&meta.InnerInstructions[inner].Instructions[j], t.Meta)
			if innerIn != nil {
				innerIn.Seq = index
				index += 1
				instructions = append(instructions, innerIn)
			}
		}
	}
	t.Instructions = instructions
	program.Parse(t)
	return t
}
