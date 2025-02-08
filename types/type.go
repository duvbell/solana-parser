package types

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/shopspring/decimal"
)

type Meta struct {
	Accounts     map[solana.PublicKey]*solana.AccountMeta
	TokenOwner   map[solana.PublicKey]solana.PublicKey
	TokenMint    map[solana.PublicKey]solana.PublicKey
	PreBalance   map[solana.PublicKey]decimal.Decimal
	PostBalance  map[solana.PublicKey]decimal.Decimal
	ErrorMessage []byte
}

type Block struct {
	Hash        solana.Hash
	Time        uint64
	Slot        uint64
	Transaction []*Transaction
}

type Transaction struct {
	Hash         solana.Signature
	Instructions []*Instruction
	Meta         Meta
	Seq          int
}

type Instruction struct {
	Seq         int
	Instruction *rpc.ParsedInstruction
	Event       []interface{}
	Receipt     []interface{}
	Children    []*Instruction
}

func (in *Instruction) AccountMetas() []*solana.AccountMeta {
	accounts := make([]*solana.AccountMeta, 0)
	for _, item := range in.Instruction.Accounts {
		accounts = append(accounts, &solana.AccountMeta{
			PublicKey:  item,
			IsWritable: false,
			IsSigner:   false,
		})
	}
	return accounts
}

func (in *Instruction) FindChildrenTransfers() []*Transfer {
	transfers := make([]*Transfer, 0)
	for _, item := range in.Children {
		if len(item.Event) != 1 {
			continue
		}
		switch item.Event[0].(type) {
		case *Transfer:
			transfers = append(transfers, item.Event[0].(*Transfer))
		}
	}
	return transfers
}

func (in *Instruction) FindChildrenPrograms(id solana.PublicKey) []*Instruction {
	instructions := make([]*Instruction, 0)
	for _, item := range in.Children {
		if item.Instruction.ProgramId == id {
			instructions = append(instructions, item)
		}
	}
	return instructions
}
