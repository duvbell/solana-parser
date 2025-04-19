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
	Accounts         []*solana.AccountMeta
	TokenAccounts    map[solana.PublicKey]*TokenAccount
	MintAccounts     map[solana.PublicKey]*MintAccount
	TokenPreBalance  map[solana.PublicKey]decimal.Decimal
	TokenPostBalance map[solana.PublicKey]decimal.Decimal
	SolPreBalance    map[solana.PublicKey]decimal.Decimal
	SolPostBalance   map[solana.PublicKey]decimal.Decimal
	ErrorMessage     []byte
}

type Block struct {
	Hash        solana.Hash
	Time        uint64
	Slot        uint64
	Transaction []*Transaction
}

type Transaction struct {
	Hash         solana.Signature
	Time         uint64
	Slot         uint64
	Instructions []*Instruction
	Meta         *Meta
	Seq          int
}

type Instruction struct {
	Seq               int
	RawInstruction    *solana.GenericInstruction
	ParsedInstruction interface{}
	Event             []interface{}
	Receipt           []interface{}
	Children          []*Instruction
}

func (in *Instruction) FindChildTransferByTo(to solana.PublicKey) *Transfer {
	for _, item := range in.Children {
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

func (in *Instruction) FindChildTransferByFrom(from solana.PublicKey) *Transfer {
	for _, item := range in.Children {
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

func (in *Instruction) FindChildMintTos() []*MintTo {
	mintTos := make([]*MintTo, 0)
	for _, item := range in.Children {
		if len(item.Event) != 1 {
			continue
		}
		switch item.Event[0].(type) {
		case *MintTo:
			mintTos = append(mintTos, item.Event[0].(*MintTo))
		}
	}
	return mintTos
}

func (in *Instruction) FindChildMintToByTo(to solana.PublicKey) *MintTo {
	for _, item := range in.Children {
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

func (in *Instruction) FindChildrenByProgram(id solana.PublicKey) []*Instruction {
	instructions := make([]*Instruction, 0)
	for _, item := range in.Children {
		if item.RawInstruction.ProgID == id {
			instructions = append(instructions, item)
		}
	}
	return instructions
}

func (in *Instruction) FindChildByProgram(id solana.PublicKey) *Instruction {
	for _, item := range in.Children {
		if item.RawInstruction.ProgID == id {
			return item
		}
	}
	return nil
}
