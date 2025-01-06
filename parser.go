package transaction

import (
	"encoding/json"
	"fmt"
	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/raydiumamm"
	whirlpool "github.com/gagliardetto/solana-go/programs/whirlpool"
	"github.com/shopspring/decimal"
	amm_v4 "github.com/gagliardetto/solana-go/programs/raydiumclmm"
)

var (
	DefaultParse = make(map[solana.PublicKey]Parser, 0)
	WhirlPool    = solana.MustPublicKeyFromBase58("whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc")
	RaydiumClmm  = solana.MustPublicKeyFromBase58("CAMMCzo5YL8w4VFF8KVHrK22GGUsp5VTaW7grrKgrWqK")
	RaydiumAMM   = solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")
)

func init() {
	RegisterParser(RaydiumAMM, RaydiumAmmParser)
}

func RegisterParser(program solana.PublicKey, p Parser) {
	DefaultParse[program] = p
}

type Parser func(in *Instruction, meta *Meta) (interface{}, interface{})

func RaydiumAmmParser(in *Instruction, meta *Meta) (interface{}, interface{}) {
	inst := new(raydium_amm.Instruction)
	instruction := in.Instruction
	err := ag_binary.NewBorshDecoder(instruction.Data).Decode(inst)
	if err != nil {
		return nil, nil
	}
	accounts := make([]*solana.AccountMeta, 0)
	for _, item := range meta.Accounts {
		accounts = append(accounts, &solana.AccountMeta{
			PublicKey:  item.PublicKey,
			IsWritable: item.Writable,
			IsSigner:   item.Signer,
		})
	}
	switch inst.TypeID.Uint8() {
	case raydium_amm.Instruction_Initialize2:
		inst1 := inst.Impl.(*raydium_amm.Initialize2)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event.(*Transfer)
		t2 := in.Children[1].Event.(*Transfer)
		//
		trade := &Trade{
			Pool:         inst1.GetAmmAccount().PublicKey.String(),
			Type:         CreatePool,
			TokenAAmount: decimal.NewFromInt(int64(t1.Amount)),
			TokenBAmount: decimal.NewFromInt(int64(t2.Amount)),
			User:         inst1.GetUserWalletAccount().PublicKey.String(),
		}
		pool := &Pool{
			Hash:     inst1.GetAmmAccount().PublicKey.String(),
			MintA:    inst1.GetCoinMintAccount().PublicKey.String(),
			MintB:    inst1.GetPcMintAccount().PublicKey.String(),
			MintLp:   inst1.GetLpMintAccount().PublicKey.String(),
			VaultA:   inst1.GetPoolCoinTokenAccountAccount().PublicKey.String(),
			VaultB:   inst1.GetPoolPcTokenAccountAccount().PublicKey.String(),
			VaultLp:  "",
			ReserveA: 0,
			ReserveB: 0,
		}
		return trade, pool
	case raydium_amm.Instruction_Deposit:
		inst1 := inst.Impl.(*raydium_amm.Deposit)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event.(*Transfer)
		t2 := in.Children[1].Event.(*Transfer)
		//
		trade := &Trade{
			Pool:         inst1.GetAmmAccount().PublicKey.String(),
			Type:         AddLiquidity,
			TokenAAmount: decimal.NewFromInt(int64(t1.Amount)),
			TokenBAmount: decimal.NewFromInt(int64(t2.Amount)),
			User:         inst1.GetUserOwnerAccount().PublicKey.String(),
		}
		return trade, nil
	case raydium_amm.Instruction_SwapBaseIn:
		inst1 := inst.Impl.(*raydium_amm.SwapBaseIn)
		inst1.SetAccounts(accounts)
		t1 := in.Children[0].Event.(*Transfer)
		t2 := in.Children[1].Event.(*Transfer)
		//
		trade := &Trade{
			Pool:         inst1.GetAmmAccount().PublicKey.String(),
			Type:         Swap,
			TokenAAmount: decimal.NewFromInt(int64(t1.Amount)),
			TokenBAmount: decimal.NewFromInt(int64(t2.Amount)),
			User:         inst1.GetUserSourceOwnerAccount().PublicKey.String(),
		}
		return trade, nil
	case raydium_amm.Instruction_SwapBaseOut:
		inst1 := inst.Impl.(*raydium_amm.SwapBaseOut)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event.(*Transfer)
		t2 := in.Children[1].Event.(*Transfer)
		//
		trade := &Trade{
			Pool:         inst1.GetAmmAccount().PublicKey.String(),
			Type:         Swap,
			TokenAAmount: decimal.NewFromInt(int64(t1.Amount)),
			TokenBAmount: decimal.NewFromInt(int64(t2.Amount)),
			User:         inst1.GetUserSourceOwnerAccount().PublicKey.String(),
		}
		return trade, nil
	case raydium_amm.Instruction_Withdraw:
		inst1 := inst.Impl.(*raydium_amm.Withdraw)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event.(*Transfer)
		t2 := in.Children[1].Event.(*Transfer)
		//
		trade := &Trade{
			Pool:         inst1.GetAmmAccount().PublicKey.String(),
			Type:         RemoveLiquidity,
			TokenAAmount: decimal.NewFromInt(int64(t1.Amount)),
			TokenBAmount: decimal.NewFromInt(int64(t2.Amount)),
			User:         inst1.GetUserOwnerAccount().PublicKey.String(),
		}
		return trade, nil
	default:
		return nil, nil
	}
}

func SystemParser(in *Instruction, meta *Meta) (interface{}, interface{}) {
	type instruction struct {
		Info struct {
			Destination solana.PublicKey `json:"destination"`
			Lamports    uint64           `json:"lamports"`
			Source      solana.PublicKey `json:"source"`
		} `json:"info"`
		T string `json:"type"`
	}
	inJson, _ := in.Instruction.Parsed.MarshalJSON()
	var myInstruction instruction
	json.Unmarshal(inJson, &myInstruction)
	switch myInstruction.T {
	case "transfer":
		transfer := &Transfer{
			Mint:   "11111111111111111111111111111111",
			Amount: myInstruction.Info.Lamports,
			From:   myInstruction.Info.Source.String(),
			To:     myInstruction.Info.Destination.String(),
		}
		return transfer, nil
	default:
		return nil, nil
	}
}

func TokenParser(in *Instruction, meta *Meta) (interface{}, interface{}) {
	type instruction struct {
		Info struct {
			Destination solana.PublicKey `json:"destination"`
			Lamports    uint64           `json:"amount,string"`
			Source      solana.PublicKey `json:"source"`
			Authority   solana.PublicKey `json:"authority"`
			Mint        solana.PublicKey `json:"mint"`
			TokenAmount struct {
				Amount   uint64 `json:"amount,string"`
				Decimals uint64 `json:"decimals"`
			} `json:"tokenAmount"`
		} `json:"info"`
		T string `json:"type"`
	}
	inJson, _ := in.Instruction.Parsed.MarshalJSON()
	var myInstruction instruction
	json.Unmarshal(inJson, &in)
	switch myInstruction.T {
	case "transfer":
		mint := meta.TokenMint[myInstruction.Info.Source]
		from := myInstruction.Info.Source
		if k, ok := meta.TokenOwner[myInstruction.Info.Source]; ok {
			from = k
		}
		to := myInstruction.Info.Destination
		if k, ok := meta.TokenOwner[myInstruction.Info.Destination]; ok {
			to = k
		}
		transfer := &Transfer{
			Mint:   mint.String(),
			Amount: myInstruction.Info.Lamports,
			From:   from.String(),
			To:     to.String(),
		}
		return transfer, nil
	case "transferChecked":
		mint := meta.TokenMint[myInstruction.Info.Source]
		from := myInstruction.Info.Source
		if k, ok := meta.TokenOwner[myInstruction.Info.Source]; ok {
			from = k
		}
		to := myInstruction.Info.Destination
		if k, ok := meta.TokenOwner[myInstruction.Info.Destination]; ok {
			to = k
		}
		transfer := &Transfer{
			Mint:   mint.String(),
			Amount: myInstruction.Info.TokenAmount.Amount,
			From:   from.String(),
			To:     to.String(),
		}
		return transfer, nil
	default:
		return nil, nil
	}
}

func RaydiumClmmParser(in *Instruction, meta *Meta) (interface{}, interface{}) {
	inst := new(amm_v4.Instruction)
	err := ag_binary.NewBorshDecoder(in.Instruction.Data).Decode(inst)
	if err != nil {
		return nil, nil
	}
	accounts := make([]*solana.AccountMeta, 0)
	for _, item := range in.Instruction.Accounts {
		accounts = append(accounts, &solana.AccountMeta{
			PublicKey:  item,
			IsWritable: false,
			IsSigner:   false,
		})
	}
	switch inst.TypeID {
	case amm_v4.Instruction_CreatePool:
		inst1 := inst.Impl.(*amm_v4.CreatePool)
		inst1.SetAccounts(accounts)
		pool := &Pool{
			Hash:     inst1.GetPoolStateAccount().PublicKey.String(),
			MintA:    inst1.GetTokenMint0Account().PublicKey.String(),
			MintB:    inst1.GetTokenMint1Account().PublicKey.String(),
			MintLp:   inst1.GetTokenVault1Account().PublicKey.String(),
			VaultA:   inst1.GetTokenVault1Account().PublicKey.String(),
			VaultB:   inst1.GetTokenVault1Account().PublicKey.String(),
			VaultLp:  "",
			ReserveA: 0,
			ReserveB: 0,
		}
		return nil, pool
	case amm_v4.Instruction_IncreaseLiquidityV2:
		inst1 := inst.Impl.(*amm_v4.IncreaseLiquidityV2)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event.(*Transfer)
		t2 := in.Children[1].Event.(*Transfer)
		//
		trade := &Trade{
			Pool:         inst1.GetPoolStateAccount().PublicKey.String(),
			Type:         AddLiquidity,
			TokenAAmount: decimal.NewFromInt(int64(t1.Amount)),
			TokenBAmount: decimal.NewFromInt(int64(t2.Amount)),
			User:         inst1.Get(0).PublicKey.String(),
		}
		return trade, nil
	case amm_v4.Instruction_DecreaseLiquidityV2:
		inst1 := inst.Impl.(*amm_v4.DecreaseLiquidityV2)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event.(*Transfer)
		t2 := in.Children[1].Event.(*Transfer)
		//
		trade := &Trade{
			Pool:         inst1.GetPoolStateAccount().PublicKey.String(),
			Type:         RemoveLiquidity,
			TokenAAmount: decimal.NewFromInt(int64(t1.Amount)),
			TokenBAmount: decimal.NewFromInt(int64(t2.Amount)),
			User:         inst1.Get(0).PublicKey.String(),
		}
		return trade, nil
	case amm_v4.Instruction_Swap:
		inst1 := inst.Impl.(*amm_v4.Swap)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event.(*Transfer)
		t2 := in.Children[1].Event.(*Transfer)
		//
		trade := &Trade{
			Pool:         inst1.GetPoolStateAccount().PublicKey.String(),
			Type:         Swap,
			TokenAAmount: decimal.NewFromInt(int64(t1.Amount)),
			TokenBAmount: decimal.NewFromInt(int64(t2.Amount)),
			User:         inst1.GetPayerAccount().PublicKey.String(),
		}
		return trade, nil
	case amm_v4.Instruction_SwapV2:
		inst1 := inst.Impl.(*amm_v4.SwapV2)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event.(*Transfer)
		t2 := in.Children[1].Event.(*Transfer)
		//
		trade := &Trade{
			Pool:         inst1.GetPoolStateAccount().PublicKey.String(),
			Type:         Swap,
			TokenAAmount: decimal.NewFromInt(int64(t1.Amount)),
			TokenBAmount: decimal.NewFromInt(int64(t2.Amount)),
			User:         inst1.Get(0).PublicKey.String(),
		}
		return trade, nil
	default:
		return nil, nil
	}
}

func WhirlPoolParser(in *Instruction, meta *Meta) (interface{}, interface{}) {
	inst := new(whirlpool.Instruction)
	err := ag_binary.NewBorshDecoder(in.Instruction.Data).Decode(inst)
	if err != nil {
		return nil, nil
	}
	accounts := make([]*solana.AccountMeta, 0)
	for _, item := range in.Instruction.Accounts {
		accounts = append(accounts, &solana.AccountMeta{
			PublicKey:  item,
			IsWritable: false,
			IsSigner:   false,
		})
	}
	switch inst.TypeID {
	case whirlpool.Instruction_InitializePool:
		inst1 := inst.Impl.(*whirlpool.InitializePool)
		inst1.SetAccounts(accounts)
		//
		trade := &Trade{
			Pool:         inst1.GetWhirlpoolAccount().PublicKey.String(),
			Type:         Swap,
			TokenAAmount: decimal.NewFromInt(0),
			TokenBAmount: decimal.NewFromInt(0),
			User:         inst1.GetFunderAccount().PublicKey.String(),
		}
		return trade, nil
	case whirlpool.Instruction_Swap:
		inst1 := inst.Impl.(*whirlpool.Swap)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event.(*Transfer)
		t2 := in.Children[1].Event.(*Transfer)
		//
		trade := &Trade{
			Pool:         inst1.GetWhirlpoolAccount().PublicKey.String(),
			Type:         Swap,
			TokenAAmount: decimal.NewFromInt(int64(t1.Amount)),
			TokenBAmount: decimal.NewFromInt(int64(t2.Amount)),
			User:         inst1.GetTokenAuthorityAccount().PublicKey.String(),
		}
		return trade, nil
	case whirlpool.Instruction_IncreaseLiquidity:
		inst1 := inst.Impl.(*whirlpool.IncreaseLiquidity)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event.(*Transfer)
		t2 := in.Children[1].Event.(*Transfer)
		//
		trade := &Trade{
			Pool:         inst1.GetWhirlpoolAccount().PublicKey.String(),
			Type:         AddLiquidity,
			TokenAAmount: decimal.NewFromInt(int64(t1.Amount)),
			TokenBAmount: decimal.NewFromInt(int64(t2.Amount)),
			User:         inst1.GetPositionAuthorityAccount().PublicKey.String(),
		}
		return trade, nil
	case whirlpool.Instruction_DecreaseLiquidity:
		inst1 := inst.Impl.(*whirlpool.DecreaseLiquidity)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event.(*Transfer)
		t2 := in.Children[1].Event.(*Transfer)
		//
		trade := &Trade{
			Pool:         inst1.GetWhirlpoolAccount().PublicKey.String(),
			Type:         RemoveLiquidity,
			TokenAAmount: decimal.NewFromInt(int64(t1.Amount)),
			TokenBAmount: decimal.NewFromInt(int64(t2.Amount)),
			User:         inst1.GetPositionAuthorityAccount().PublicKey.String(),
		}
		return trade, nil
	default:
		return nil, nil
	}
}
