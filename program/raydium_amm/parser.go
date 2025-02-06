package raydium_amm

import (
	"errors"
	"github.com/blockchain-develop/solana-parser/log"
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

/*
coin : token 0
pc : token 1
*/

func init() {
	program.RegisterParser(raydium_amm.ProgramID, raydium_amm.ProgramName, ProgramParser)
	RegisterParser(uint64(raydium_amm.Instruction_Initialize), ParseInitialize)
	RegisterParser(uint64(raydium_amm.Instruction_Initialize2), ParseInitialize2)
	RegisterParser(uint64(raydium_amm.Instruction_MonitorStep), ParseMonitorStep)
	RegisterParser(uint64(raydium_amm.Instruction_Deposit), ParseDeposit)
	RegisterParser(uint64(raydium_amm.Instruction_Withdraw), ParseWithdraw)
	RegisterParser(uint64(raydium_amm.Instruction_MigrateToOpenBook), ParseMigrateToOpenBook)
	RegisterParser(uint64(raydium_amm.Instruction_SetParams), ParseSetParams)
	RegisterParser(uint64(raydium_amm.Instruction_WithdrawPnl), ParseWithdrawPnl)
	RegisterParser(uint64(raydium_amm.Instruction_WithdrawSrm), ParseWithdrawSrm)
	RegisterParser(uint64(raydium_amm.Instruction_SwapBaseIn), ParseSwapBaseIn)
	RegisterParser(uint64(raydium_amm.Instruction_PreInitialize), ParsePreInitialize)
	RegisterParser(uint64(raydium_amm.Instruction_SwapBaseOut), ParseSwapBaseOut)
	RegisterParser(uint64(raydium_amm.Instruction_SimulateInfo), ParseSimulateInfo)
	RegisterParser(uint64(raydium_amm.Instruction_AdminCancelOrders), ParseAdminCancelOrders)
	RegisterParser(uint64(raydium_amm.Instruction_CreateConfigAccount), ParseCreateConfigAccount)
	RegisterParser(uint64(raydium_amm.Instruction_UpdateConfigAccount), ParseUpdateConfigAccount)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	inst, err := raydium_amm.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
	if err != nil {
		return err
	}
	id := uint64(inst.TypeID.Uint32())
	parser, ok := Parsers[id]
	if !ok {
		return errors.New("parser not found")
	}
	parser(inst, in, meta)
	return nil
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

func ParseInitialize(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse initialize", "program", raydium_amm.ProgramName)
}
func ParseInitialize2(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_amm.Initialize2)
	// the latest three transfer
	createPool := &types.CreatePool{
		Dex:     in.Instruction.ProgramId,
		Pool:    inst1.GetAmmAccount().PublicKey,
		User:    inst1.GetUserWalletAccount().PublicKey,
		TokenA:  inst1.GetCoinMintAccount().PublicKey,
		TokenB:  inst1.GetPcMintAccount().PublicKey,
		TokenLP: inst1.GetLpMintAccount().PublicKey,
		VaultA:  inst1.GetPoolCoinTokenAccountAccount().PublicKey,
		VaultB:  inst1.GetPoolPcTokenAccountAccount().PublicKey,
		VaultLP: inst1.GetPoolTempLpAccount().PublicKey,
	}
	addLiquidity := &types.AddLiquidity{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetAmmAccount().PublicKey,
		User: inst1.GetUserWalletAccount().PublicKey,
	}
	transfers := in.FindChildrenTransfers()
	for _, transfer := range transfers {
		if transfer.To == inst1.GetPoolCoinTokenAccountAccount().PublicKey && transfer.Mint == inst1.GetCoinMintAccount().PublicKey {
			addLiquidity.TokenATransfer = transfer
		}
		if transfer.To == inst1.GetPoolPcTokenAccountAccount().PublicKey && transfer.Mint == inst1.GetPcMintAccount().PublicKey {
			addLiquidity.TokenBTransfer = transfer
		}
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
func ParseMonitorStep(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseDeposit(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_amm.Deposit)
	addLiquidity := &types.AddLiquidity{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetAmmAccount().PublicKey,
		User: inst1.GetUserOwnerAccount().PublicKey,
	}
	transfers := in.FindChildrenTransfers()
	for _, transfer := range transfers {
		if transfer.To == inst1.GetPoolCoinTokenAccountAccount().PublicKey {
			addLiquidity.TokenATransfer = transfer
		}
		if transfer.To == inst1.GetPoolPcTokenAccountAccount().PublicKey {
			addLiquidity.TokenBTransfer = transfer
		}
	}
	in.Event = []interface{}{addLiquidity, addLiquidity}
}
func ParseWithdraw(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_amm.Withdraw)
	removeLiquidity := &types.RemoveLiquidity{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetAmmAccount().PublicKey,
		User: inst1.GetUserOwnerAccount().PublicKey,
	}
	transfers := in.FindChildrenTransfers()
	for _, transfer := range transfers {
		if transfer.From == inst1.GetPoolCoinTokenAccountAccount().PublicKey {
			removeLiquidity.TokenATransfer = transfer
		}
		if transfer.From == inst1.GetPoolPcTokenAccountAccount().PublicKey {
			removeLiquidity.TokenBTransfer = transfer
		}
	}
	in.Event = []interface{}{removeLiquidity}
}
func ParseMigrateToOpenBook(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSetParams(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseWithdrawPnl(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseWithdrawSrm(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSwapBaseIn(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_amm.SwapBaseIn)
	if len(inst1.GetAccounts()) == 17 {
		inst1.SetAccounts(insertAccount(inst1.GetAccounts(), 4))
	}
	swap := &types.Swap{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetAmmAccount().PublicKey,
		User: inst1.GetUserSourceOwnerAccount().PublicKey,
	}
	if *inst1.AmountIn > 0 {
		swap.InputTransfer = in.Children[0].Event[0].(*types.Transfer)
		swap.OutputTransfer = in.Children[1].Event[0].(*types.Transfer)
	}
	in.Event = []interface{}{swap}
}
func ParsePreInitialize(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse pre-initialize", "program", raydium_amm.ProgramName)
}
func ParseSwapBaseOut(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_amm.SwapBaseOut)
	if len(inst1.GetAccounts()) == 17 {
		inst1.SetAccounts(insertAccount(inst1.GetAccounts(), 4))
	}
	swap := &types.Swap{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetAmmAccount().PublicKey,
		User: inst1.GetUserSourceOwnerAccount().PublicKey,
	}

	if *inst1.AmountOut > 0 {
		swap.InputTransfer = in.Children[0].Event[0].(*types.Transfer)
		swap.OutputTransfer = in.Children[1].Event[0].(*types.Transfer)
	}

	in.Event = []interface{}{swap}
}
func ParseSimulateInfo(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseAdminCancelOrders(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCreateConfigAccount(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseUpdateConfigAccount(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
}

// Default
func ParseDefault(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
