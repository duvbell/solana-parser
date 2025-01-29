package raydium_amm

import (
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/raydium_amm"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(raydium_amm.ProgramID, ProgramParser)
	RegisterParser(uint64(raydium_amm.Instruction_Initialize2), ParseInitialize2)
	RegisterParser(uint64(raydium_amm.Instruction_Deposit), ParseDeposit)
	RegisterParser(uint64(raydium_amm.Instruction_SwapBaseIn), ParseSwapBaseIn)
	RegisterParser(uint64(raydium_amm.Instruction_SwapBaseOut), ParseSwapBaseOut)
	RegisterParser(uint64(raydium_amm.Instruction_Withdraw), ParseWithdraw)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) {
	inst, err := raydium_amm.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
	if err != nil {
		return
	}
	id := uint64(inst.TypeID.Uint32())
	parser, ok := Parsers[id]
	if !ok {
		return
	}
	parser(inst, in, meta)
}

func insertAccount(accounts []*solana.AccountMeta, index int) []*solana.AccountMeta {
	s := make([]*solana.AccountMeta, 0)
	s = append(s, accounts[0:index]...)
	s = append(s, &solana.AccountMeta{
		PublicKey:  solana.MustPublicKeyFromBase58("11111111111111111111111111111111"),
		IsWritable: true,
		IsSigner:   false,
	})
	s = append(s, accounts[index:]...)
	return s
}

// Initialize2
func ParseInitialize2(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_amm.Initialize2)
	// the latest three transfer
	index := len(in.Children)
	t1 := in.Children[index-3].Event[0].(*types.Transfer)
	t2 := in.Children[index-2].Event[0].(*types.Transfer)
	t3 := in.Children[index-1].Event[0].(*types.MintTo)
	createPool := &types.CreatePool{
		Pool:      inst1.GetAmmAccount().PublicKey,
		TokenA:    inst1.GetPcMintAccount().PublicKey,
		TokenB:    inst1.GetCoinMintAccount().PublicKey,
		TokenLP:   inst1.GetLpMintAccount().PublicKey,
		AccountA:  inst1.GetPoolPcTokenAccountAccount().PublicKey,
		AccountB:  inst1.GetPoolCoinTokenAccountAccount().PublicKey,
		AccountLP: inst1.GetPoolTempLpAccount().PublicKey,
		User:      inst1.GetAmmAuthorityAccount().PublicKey,
	}
	addLiquidity := &types.AddLiquidity{
		Pool:           inst1.GetAmmAccount().PublicKey,
		User:           inst1.GetAmmAuthorityAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
		TokenLpMint:    t3,
	}
	pool := &types.Pool{
		Hash:     inst1.GetAmmAccount().PublicKey,
		MintA:    inst1.GetCoinMintAccount().PublicKey,
		MintB:    inst1.GetPcMintAccount().PublicKey,
		MintLp:   inst1.GetLpMintAccount().PublicKey,
		VaultA:   inst1.GetPoolCoinTokenAccountAccount().PublicKey,
		VaultB:   inst1.GetPoolPcTokenAccountAccount().PublicKey,
		ReserveA: 0,
		ReserveB: 0,
	}
	in.Event = []interface{}{createPool, addLiquidity}
	in.Receipt = []interface{}{pool}
}

// Deposit
func ParseDeposit(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_amm.Deposit)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	addLiquidity := &types.AddLiquidity{
		Pool:           inst1.GetAmmAccount().PublicKey,
		User:           inst1.GetUserOwnerAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
	}
	panic("not supported")
	in.Event = []interface{}{addLiquidity, addLiquidity}
}

// SwapBaseIn
func ParseSwapBaseIn(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_amm.SwapBaseIn)
	if len(inst1.GetAccounts()) == 17 {
		inst1.SetAccounts(insertAccount(inst1.GetAccounts(), 4))
	}
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	swap := &types.Swap{
		Pool:           inst1.GetAmmAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
		User:           inst1.GetUserSourceOwnerAccount().PublicKey,
	}
	in.Event = []interface{}{swap}
}

// SwapBaseOut
func ParseSwapBaseOut(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_amm.SwapBaseOut)
	if len(inst1.GetAccounts()) == 17 {
		inst1.SetAccounts(insertAccount(inst1.GetAccounts(), 4))
	}
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	swap := &types.Swap{
		Pool:           inst1.GetAmmAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
		User:           inst1.GetUserSourceOwnerAccount().PublicKey,
	}
	in.Event = []interface{}{swap}
}

// Withdraw
func ParseWithdraw(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_amm.Withdraw)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	t3 := in.Children[1].Event[0].(*types.Burn)
	removeLiquidity := &types.RemoveLiquidity{
		Pool:           inst1.GetAmmAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
		TokenLpBurn:    t3,
		User:           inst1.GetUserOwnerAccount().PublicKey,
	}
	in.Event = []interface{}{removeLiquidity}
}

// Default
func ParseDefault(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
