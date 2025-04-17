package raydium_cp

import (
	"errors"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/raydium_cp"
	"github.com/solana-parser/program"
	"github.com/solana-parser/types"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) error

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(raydium_cp.ProgramID, raydium_cp.ProgramName, program.Swap, 1, ProgramParser)
	RegisterParser(uint64(raydium_cp.Instruction_CreateAmmConfig.Uint32()), ParseCreateAmmConfig)
	RegisterParser(uint64(raydium_cp.Instruction_UpdateAmmConfig.Uint32()), ParseUpdateAmmConfig)
	RegisterParser(uint64(raydium_cp.Instruction_UpdatePoolStatus.Uint32()), ParseUpdatePoolStatus)
	RegisterParser(uint64(raydium_cp.Instruction_CollectProtocolFee.Uint32()), ParseCollectProtocolFee)
	RegisterParser(uint64(raydium_cp.Instruction_CollectFundFee.Uint32()), ParseCollectFundFee)
	RegisterParser(uint64(raydium_cp.Instruction_Initialize.Uint32()), ParseInitialize)
	RegisterParser(uint64(raydium_cp.Instruction_Deposit.Uint32()), ParseDeposit)
	RegisterParser(uint64(raydium_cp.Instruction_Withdraw.Uint32()), ParseWithdraw)
	RegisterParser(uint64(raydium_cp.Instruction_SwapBaseInput.Uint32()), ParseSwapBaseInput)
	RegisterParser(uint64(raydium_cp.Instruction_SwapBaseOutput.Uint32()), ParseSwapBaseOutput)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	inst, err := raydium_cp.DecodeInstruction(in.RawInstruction.AccountValues, in.RawInstruction.DataBytes)
	if err != nil {
		return err
	}
	id := uint64(inst.TypeID.Uint32())
	parser, ok := Parsers[id]
	if !ok {
		return errors.New("parser not found")
	}
	return parser(inst, in, meta)
}

func ParseCreateAmmConfig(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}
func ParseUpdateAmmConfig(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}
func ParseUpdatePoolStatus(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}
func ParseCollectProtocolFee(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}
func ParseCollectFundFee(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}
func ParseInitialize(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) error {
	// log.Logger.Info("ignore parse initialize", "program", raydium_cp.ProgramName)
	inst1 := inst.Impl.(*raydium_cp.Initialize)
	createPool := &types.CreatePool{
		Dex:     in.RawInstruction.ProgID,
		Pool:    inst1.GetPoolStateAccount().PublicKey,
		User:    inst1.GetCreatorAccount().PublicKey,
		TokenA:  inst1.GetToken0MintAccount().PublicKey,
		TokenB:  inst1.GetToken1MintAccount().PublicKey,
		TokenLP: inst1.GetLpMintAccount().PublicKey,
		VaultA:  inst1.GetToken0VaultAccount().PublicKey,
		VaultB:  inst1.GetToken1VaultAccount().PublicKey,
		VaultLP: solana.PublicKey{},
	}
	addLiquidity := &types.AddLiquidity{
		Dex:  in.RawInstruction.ProgID,
		Pool: inst1.GetPoolStateAccount().PublicKey,
		User: inst1.GetCreatorAccount().PublicKey,
	}
	addLiquidity.TokenATransfer = in.FindNextTransferByTo(inst1.GetToken0VaultAccount().PublicKey)
	addLiquidity.TokenBTransfer = in.FindNextTransferByTo(inst1.GetToken1VaultAccount().PublicKey)
	in.Event = []interface{}{createPool, addLiquidity}
	return nil
}
func ParseDeposit(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*raydium_cp.Deposit)
	addLiquidity := &types.AddLiquidity{
		Dex:  in.RawInstruction.ProgID,
		Pool: inst1.GetPoolStateAccount().PublicKey,
		User: inst1.GetOwnerAccount().PublicKey,
	}
	addLiquidity.TokenATransfer = in.FindNextTransferByTo(inst1.GetToken0VaultAccount().PublicKey)
	addLiquidity.TokenBTransfer = in.FindNextTransferByTo(inst1.GetToken1VaultAccount().PublicKey)
	in.Event = []interface{}{addLiquidity}
	return nil
}
func ParseWithdraw(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*raydium_cp.Withdraw)
	removeLiquidity := &types.RemoveLiquidity{
		Dex:  in.RawInstruction.ProgID,
		Pool: inst1.GetPoolStateAccount().PublicKey,
		User: inst1.GetOwnerAccount().PublicKey,
	}
	removeLiquidity.TokenATransfer = in.FindNextTransferByFrom(inst1.GetToken0VaultAccount().PublicKey)
	removeLiquidity.TokenBTransfer = in.FindNextTransferByFrom(inst1.GetToken1VaultAccount().PublicKey)
	in.Event = []interface{}{removeLiquidity}
	return nil
}
func ParseSwapBaseInput(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*raydium_cp.SwapBaseInput)
	swap := &types.Swap{
		Dex:  in.RawInstruction.ProgID,
		Pool: inst1.GetPoolStateAccount().PublicKey,
		User: inst1.GetPayerAccount().PublicKey,
	}
	swap.InputTransfer = in.FindNextTransferByTo(inst1.GetInputVaultAccount().PublicKey)
	swap.OutputTransfer = in.FindNextTransferByTo(inst1.GetOutputVaultAccount().PublicKey)
	in.Event = []interface{}{swap}
	return nil
}
func ParseSwapBaseOutput(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*raydium_cp.SwapBaseOutput)
	swap := &types.Swap{
		Dex:  in.RawInstruction.ProgID,
		Pool: inst1.GetPoolStateAccount().PublicKey,
		User: inst1.GetPayerAccount().PublicKey,
	}
	swap.InputTransfer = in.FindNextTransferByTo(inst1.GetInputVaultAccount().PublicKey)
	swap.OutputTransfer = in.FindNextTransferByTo(inst1.GetOutputVaultAccount().PublicKey)
	in.Event = []interface{}{swap}
	return nil
}

// Default
func ParseDefault(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

// Fault
func ParseFault(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) error {
	panic("not supported")
}
