package program

import (
	"github.com/blockchain-develop/solana-parser/log"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go"
)

var (
	Parsers     = make(map[solana.PublicKey]Parser)
	Name2Id     = make(map[string]solana.PublicKey)
	Id2Name     = make(map[solana.PublicKey]string)
	Id2Type     = make(map[solana.PublicKey]string)
	Id2Priority = make(map[solana.PublicKey]int)
)

const (
	Token      = "Token"
	Swap       = "Swap"
	StableSwap = "StableSwap"
	OrderBook  = "OrderBook"
)

type Parser func(transaction *types.Transaction, index int) error

func RegisterParser(program solana.PublicKey, name string, t string, priority int, p Parser) {
	Parsers[program] = p
	Id2Priority[program] = priority
	Name2Id[name] = program
	Id2Name[program] = name
	Id2Type[program] = t
}

func RemoveParser(program solana.PublicKey) {
	Parsers[program] = nil
}

func FilterInstruction(in *solana.CompiledInstruction, meta *types.Meta) *types.Instruction {
	programId := meta.Accounts[in.ProgramIDIndex].PublicKey
	p, ok := Parsers[programId]
	if !ok || p == nil {
		return nil
	}
	accountMetas := make([]*solana.AccountMeta, 0)
	for _, accountIndex := range in.Accounts {
		accountMetas = append(accountMetas, meta.Accounts[accountIndex])
	}
	return &types.Instruction{
		Raw: &solana.GenericInstruction{
			AccountValues: accountMetas,
			ProgID:        programId,
			DataBytes:     in.Data,
		},
	}
}

func Parse(transaction *types.Transaction) {
	priority := 0
	for i := 0; i < len(transaction.Instructions); i++ {
		programId := transaction.Instructions[i].Raw.ProgID
		if Id2Priority[programId] == priority {
			parse(transaction, i)
		}
	}
	priority += 1
	for i := 0; i < len(transaction.Instructions); i++ {
		programId := transaction.Instructions[i].Raw.ProgID
		if Id2Priority[programId] == priority {
			parse(transaction, i)
		}
	}
}

func parse(transaction *types.Transaction, index int) {
	in := transaction.Instructions[index]
	parser, ok := Parsers[in.Raw.ProgID]
	if !ok || parser == nil {
		log.Logger.Error("no parser", "program id", in.Raw.ProgID)
		return
	}
	err := parser(transaction, index)
	if err != nil {
		log.Logger.Error("parse error", "program", in.Raw.ProgID, "err", err)
	}
}
