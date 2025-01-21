package solanaparser

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type Block struct {
	Hash        solana.Hash
	Time        uint64
	Slot        uint64
	Transaction []*Transaction
}

func New() *Block {
	return &Block{}
}

func (block *Block) ParseWithParsers(slot uint64, b *rpc.GetParsedBlockResult, parsers map[solana.PublicKey]Parser) {
	block.Slot = slot
	block.Time = uint64(*b.BlockTime)
	block.Hash = b.Blockhash
	myTxs := make([]*Transaction, 0)
	for i, tx := range b.Transactions {
		myTx := NewTransaction()
		err := myTx.Parse(&tx)
		if err != nil {
			continue
		}
		myTx.Seq = i + 1
		if len(myTx.Instructions) == 0 {
			// no instruction
			continue
		}
		err = myTx.ParseActions(parsers)
		if err != nil {
			// do nothing
		}
		myTxs = append(myTxs, myTx)
	}
	block.Transaction = myTxs
}

func (block *Block) Parse(slot uint64, b *rpc.GetParsedBlockResult) {
	block.Slot = slot
	block.Time = uint64(*b.BlockTime)
	block.Hash = b.Blockhash
	myTxs := make([]*Transaction, 0)
	for i, tx := range b.Transactions {
		myTx := NewTransaction()
		err := myTx.Parse(&tx)
		if err != nil {
			continue
		}
		myTx.Seq = i + 1
		if len(myTx.Instructions) == 0 {
			// no instruction
			continue
		}
		err = myTx.ParseActions(DefaultParse)
		if err != nil {
			// do nothing
		}
		myTxs = append(myTxs, myTx)
	}
	block.Transaction = myTxs
}
