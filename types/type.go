package types

import (
	"github.com/gagliardetto/solana-go"
	"github.com/shopspring/decimal"
)

type TokenAccount struct {
	Owner     *solana.PublicKey
	ProgramId *solana.PublicKey
	Mint      solana.PublicKey
}

type MintAccount struct {
	Mint     solana.PublicKey
	Decimals uint8
}

type Meta struct {
	Accounts      []*solana.AccountMeta
	TokenAccounts map[solana.PublicKey]*TokenAccount
	MintAccounts  map[solana.PublicKey]*MintAccount
	PreBalance    map[solana.PublicKey]decimal.Decimal
	PostBalance   map[solana.PublicKey]decimal.Decimal
	ErrorMessage  []byte
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
	Meta         *Meta
	Seq          int
}

func (tx *Transaction) FindNextTransferByTo(index int, to solana.PublicKey) *Transfer {
	for i := index; i < len(tx.Instructions); i++ {
		item := tx.Instructions[i]
		if len(item.Event) != 1 {
			continue
		}
		switch item.Event[0].(type) {
		case *Transfer:
			transfer := item.Event[0].(*Transfer)
			if transfer.To == to {
				return transfer
			}
		}
	}
	return nil
}

func (tx *Transaction) FindNextTransferByFrom(index int, from solana.PublicKey) *Transfer {
	for i := index; i < len(tx.Instructions); i++ {
		item := tx.Instructions[i]
		if len(item.Event) != 1 {
			continue
		}
		switch item.Event[0].(type) {
		case *Transfer:
			transfer := item.Event[0].(*Transfer)
			if transfer.From == from {
				return transfer
			}
		}
	}
	return nil
}

func (tx *Transaction) FindNextMintTo(index int, to solana.PublicKey) *MintTo {
	for i := index; i < len(tx.Instructions); i++ {
		item := tx.Instructions[i]
		if len(item.Event) != 1 {
			continue
		}
		switch item.Event[0].(type) {
		case *MintTo:
			mintTo := item.Event[0].(*MintTo)
			if mintTo.Account == to {
				return mintTo
			}
		}
	}
	return nil
}

func (tx *Transaction) FindNextInstructionByProgram(index int, id solana.PublicKey) (*Instruction, int) {
	for i := index + 1; i < len(tx.Instructions); i++ {
		item := tx.Instructions[i]
		if item.RawInstruction.ProgID == id {
			return item, i
		}
	}
	return nil, -1
}

type Instruction struct {
	Seq               int
	RawInstruction    *solana.GenericInstruction
	ParsedInstruction interface{}
	Event             []interface{}
	Receipt           []interface{}
}
