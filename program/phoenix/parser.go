package phoenix_v1

import (
	"github.com/blockchain-develop/solana-parser/log"
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go/programs/phoenix_v1"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(phoenix_v1.ProgramID, ProgramParser)
	RegisterParser(uint64(phoenix_v1.Instruction_Swap.Uint32()), ParseSwap)
	RegisterParser(uint64(phoenix_v1.Instruction_SwapWithFreeFunds.Uint32()), ParseSwapWithFreeFunds)
	RegisterParser(uint64(phoenix_v1.Instruction_PlaceLimitOrder.Uint32()), ParsePlaceLimitOrder)
	RegisterParser(uint64(phoenix_v1.Instruction_PlaceLimitOrderWithFreeFunds.Uint32()), ParsePlaceLimitOrderWithFreeFunds)
	RegisterParser(uint64(phoenix_v1.Instruction_ReduceOrder.Uint32()), ParseReduceOrder)
	RegisterParser(uint64(phoenix_v1.Instruction_ReduceOrderWithFreeFunds.Uint32()), ParseReduceOrderWithFreeFunds)
	RegisterParser(uint64(phoenix_v1.Instruction_CancelAllOrders.Uint32()), ParseCancelAllOrders)
	RegisterParser(uint64(phoenix_v1.Instruction_CancelAllOrdersWithFreeFunds.Uint32()), ParseCancelAllOrdersWithFreeFunds)
	RegisterParser(uint64(phoenix_v1.Instruction_CancelUpTo.Uint32()), ParseCancelUpTo)
	RegisterParser(uint64(phoenix_v1.Instruction_CancelUpToWithFreeFunds.Uint32()), ParseCancelUpToWithFreeFunds)
	RegisterParser(uint64(phoenix_v1.Instruction_CancelMultipleOrdersById.Uint32()), ParseCancelMultipleOrdersById)
	RegisterParser(uint64(phoenix_v1.Instruction_CancelMultipleOrdersByIdWithFreeFunds.Uint32()), ParseCancelMultipleOrdersByIdWithFreeFunds)
	RegisterParser(uint64(phoenix_v1.Instruction_WithdrawFunds.Uint32()), ParseWithdrawFunds)
	RegisterParser(uint64(phoenix_v1.Instruction_DepositFunds.Uint32()), ParseDepositFunds)
	RegisterParser(uint64(phoenix_v1.Instruction_RequestSeat.Uint32()), ParseRequestSeat)
	RegisterParser(uint64(phoenix_v1.Instruction_Log.Uint32()), ParseLog)
	RegisterParser(uint64(phoenix_v1.Instruction_PlaceMultiplePostOnlyOrders.Uint32()), ParsePlaceMultiplePostOnlyOrders)
	RegisterParser(uint64(phoenix_v1.Instruction_PlaceMultiplePostOnlyOrdersWithFreeFunds.Uint32()), ParsePlaceMultiplePostOnlyOrdersWithFreeFunds)
	RegisterParser(uint64(phoenix_v1.Instruction_InitializeMarket.Uint32()), ParseInitializeMarket)
	RegisterParser(uint64(phoenix_v1.Instruction_ClaimAuthority.Uint32()), ParseClaimAuthority)
	RegisterParser(uint64(phoenix_v1.Instruction_NameSuccessor.Uint32()), ParseNameSuccessor)
	RegisterParser(uint64(phoenix_v1.Instruction_ChangeMarketStatus.Uint32()), ParseChangeMarketStatus)
	RegisterParser(uint64(phoenix_v1.Instruction_ChangeSeatStatus.Uint32()), ParseChangeSeatStatus)
	RegisterParser(uint64(phoenix_v1.Instruction_RequestSeatAuthorized.Uint32()), ParseRequestSeatAuthorized)
	RegisterParser(uint64(phoenix_v1.Instruction_EvictSeat.Uint32()), ParseEvictSeat)
	RegisterParser(uint64(phoenix_v1.Instruction_ForceCancelOrders.Uint32()), ParseForceCancelOrders)
	RegisterParser(uint64(phoenix_v1.Instruction_CollectFees.Uint32()), ParseCollectFees)
	RegisterParser(uint64(phoenix_v1.Instruction_ChangeFeeRecipient.Uint32()), ParseChangeFeeRecipient)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) {
	inst, err := phoenix_v1.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
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

func ParseSwap(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse swap", "program", phoenix_v1.ProgramName)
}
func ParseSwapWithFreeFunds(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse swap", "program", phoenix_v1.ProgramName)
}
func ParsePlaceLimitOrder(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParsePlaceLimitOrderWithFreeFunds(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseReduceOrder(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseReduceOrderWithFreeFunds(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseCancelAllOrders(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseCancelAllOrdersWithFreeFunds(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseCancelUpTo(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseCancelUpToWithFreeFunds(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseCancelMultipleOrdersById(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseCancelMultipleOrdersByIdWithFreeFunds(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseWithdrawFunds(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseDepositFunds(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseRequestSeat(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseLog(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParsePlaceMultiplePostOnlyOrders(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParsePlaceMultiplePostOnlyOrdersWithFreeFunds(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseInitializeMarket(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseClaimAuthority(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseNameSuccessor(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseChangeMarketStatus(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseChangeSeatStatus(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseRequestSeatAuthorized(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseEvictSeat(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseForceCancelOrders(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseCollectFees(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseChangeFeeRecipient(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// Default
func ParseDefault(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
