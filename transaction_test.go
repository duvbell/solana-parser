package transaction

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestTransaction_Parse(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetParsedTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("kPJDvuHhph8nMZ5Uwfiq1tFEoANR4QY3CkxFFzVUCULw2k5HEVHwFvw6Q8yJGt3PQowaRDcGKy6WptToZnKkBtL"),
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
	tx := New()
	err = tx.Parse(transaction)
	if err != nil {
		panic(err)
	}
	err = tx.ParseActions(DefaultParse)
	if err != nil {
		panic(err)
	}
	txJson, _ := json.MarshalIndent(tx, "", "    ")
	os.WriteFile(fmt.Sprintf("tx.json"), txJson, 0644)
}

func TestBlock_Parse(t *testing.T) {
	client := rpc.New(rpc.MainNetBeta_RPC)
	rewards := false
	version := uint64(0)
	r, err := client.GetParsedBlockWithOpts(
		context.Background(),
		313040629,
		&rpc.GetBlockOpts{
			Encoding:                       solana.EncodingJSONParsed,
			TransactionDetails:             rpc.TransactionDetailsFull,
			Rewards:                        &rewards,
			Commitment:                     rpc.CommitmentConfirmed,
			MaxSupportedTransactionVersion: &version,
		},
	)
	if err != nil {
		panic(err)
	}
	rJson, _ := json.MarshalIndent(r, "", "    ")
	os.WriteFile(fmt.Sprintf("block.json"), rJson, 0644)

	for i, tx := range r.Transactions {
		myTx := New()
		err = myTx.Parse(&tx)
		if err != nil {
			panic(err)
		}
		myTx.Seq = i + 1
		if len(myTx.Instructions) == 0 {
			// no instruction
			continue
		}
		fmt.Printf("hash: %s\n", tx.Transaction.Signatures[0].String())
		err = myTx.ParseActions(DefaultParse)
		if err != nil {
			panic(err)
		}
	}
}

type Request struct {
	Jsonrpc string        `json:"jsonrpc"`
	Id      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func TestBlock_Raw(t *testing.T) {
	slot := 263991594
	options := make(map[string]interface{})
	options["encoding"] = "json"
	options["maxSupportedTransactionVersion"] = 1
	options["transactionDetails"] = "full"
	options["rewards"] = false

	req := Request{
		Jsonrpc: "2.0",
		Id:      1,
		Method:  "getBlock",
		Params: []interface{}{
			slot,
			options,
		},
	}
	reqJson, _ := json.Marshal(req)
	httpRequest, err := http.NewRequest("POST", "https://api.mainnet-beta.solana.com", bytes.NewBuffer(reqJson))
	if err != nil {
		panic(err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Accepts", "application/json")

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		panic(err)
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != 200 {
		fmt.Printf("response status code: %d", httpResponse.StatusCode)
		panic("response status invalid")
	}
	rsp, _ := ioutil.ReadAll(httpResponse.Body)
	//
	os.WriteFile(fmt.Sprintf("block.json"), rsp, 0644)
}
