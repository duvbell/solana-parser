package transaction

import (
	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/raydiumamm"
	"github.com/shopspring/decimal"
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
