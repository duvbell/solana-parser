package solanaparser

import (
	"encoding/json"
	"errors"
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

type Transaction struct {
	Hash         solana.Signature
	Instructions []*Instruction
	Meta         Meta
	BlockHash    solana.Hash
	Seq          int
}

func NewTransaction() *Transaction {
	return &Transaction{
		Meta: Meta{
			Accounts:    make(map[solana.PublicKey]*solana.AccountMeta),
			TokenOwner:  make(map[solana.PublicKey]solana.PublicKey),
			TokenMint:   make(map[solana.PublicKey]solana.PublicKey),
			PreBalance:  make(map[solana.PublicKey]decimal.Decimal),
			PostBalance: make(map[solana.PublicKey]decimal.Decimal),
		},
	}
}

func (t *Transaction) Parse(tx *rpc.ParsedTransactionWithMeta) error {
	if tx.Meta == nil || tx.Transaction == nil {
		return errors.New("transaction meta or data is missing")
	}
	meta := tx.Meta
	transaction := tx.Transaction
	t.Hash = transaction.Signatures[0]
	if meta.Err != nil {
		// if failed, ignore this transaction
		errJson, _ := json.Marshal(meta.Err)
		t.Meta.ErrorMessage = errJson
		return nil
	}
	message := transaction.Message
	instructions := message.Instructions
	if len(instructions) == 0 {
		return nil
	}
	if instructions[0].ProgramId == solana.VoteProgramID {
		return nil
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
		current := &Instruction{
			Seq:         index + 1,
			Instruction: instruction,
			Children:    nil,
		}
		t.Instructions = append(t.Instructions, current)
	}
	innerInstructions := meta.InnerInstructions
	for _, innerInstruction := range innerInstructions {
		parent := t.Instructions[innerInstruction.Index]
		parent.parseInstructions(innerInstruction.Instructions)
	}
	return nil
}

func (t *Transaction) ParseActions(parsers map[solana.PublicKey]Parser) error {
	for _, instruction := range t.Instructions {
		instruction.instructionActions(parsers, &t.Meta)
	}
	return nil
}
