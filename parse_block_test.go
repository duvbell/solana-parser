package solanaparser

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
	"time"
)

func TestBlock_Scan(t *testing.T) {
	client := rpc.New(rpc.MainNetBeta_RPC)
	rewards := false
	version := uint64(0)
	slot := uint64(315806607)
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second * 2)
		slot += 1
		r, err := client.GetParsedBlockWithOpts(
			context.Background(),
			slot,
			&rpc.GetBlockOpts{
				Encoding:                       solana.EncodingJSONParsed,
				TransactionDetails:             rpc.TransactionDetailsFull,
				Rewards:                        &rewards,
				Commitment:                     rpc.CommitmentConfirmed,
				MaxSupportedTransactionVersion: &version,
			},
		)
		if err != nil {
			fmt.Printf("error getting block %d: %s\n", i, err)
			continue
		}
		fmt.Printf("============================================ block: %d\n", slot)
		ParseBlock(slot, r)
	}
}

func TestBlock_Parse(t *testing.T) {
	client := rpc.New(rpc.MainNetBeta_RPC)
	rewards := false
	version := uint64(0)
	slot := uint64(315806587)
	r, err := client.GetParsedBlockWithOpts(
		context.Background(),
		slot,
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

	ParseBlock(slot, r)
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
