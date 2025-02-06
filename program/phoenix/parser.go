package phoenix_v1

import (
	"errors"
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
	program.RegisterParser(phoenix_v1.ProgramID, phoenix_v1.ProgramName, ProgramParser)
	RegisterParser(uint64(phoenix_v1.Instruction_Swap), ParseSwap)
	RegisterParser(uint64(phoenix_v1.Instruction_SwapWithFreeFunds), ParseSwapWithFreeFunds)
	RegisterParser(uint64(phoenix_v1.Instruction_PlaceLimitOrder), ParsePlaceLimitOrder)
	RegisterParser(uint64(phoenix_v1.Instruction_PlaceLimitOrderWithFreeFunds), ParsePlaceLimitOrderWithFreeFunds)
	RegisterParser(uint64(phoenix_v1.Instruction_ReduceOrder), ParseReduceOrder)
	RegisterParser(uint64(phoenix_v1.Instruction_ReduceOrderWithFreeFunds), ParseReduceOrderWithFreeFunds)
	RegisterParser(uint64(phoenix_v1.Instruction_CancelAllOrders), ParseCancelAllOrders)
	RegisterParser(uint64(phoenix_v1.Instruction_CancelAllOrdersWithFreeFunds), ParseCancelAllOrdersWithFreeFunds)
	RegisterParser(uint64(phoenix_v1.Instruction_CancelUpTo), ParseCancelUpTo)
	RegisterParser(uint64(phoenix_v1.Instruction_CancelUpToWithFreeFunds), ParseCancelUpToWithFreeFunds)
	RegisterParser(uint64(phoenix_v1.Instruction_CancelMultipleOrdersById), ParseCancelMultipleOrdersById)
	RegisterParser(uint64(phoenix_v1.Instruction_CancelMultipleOrdersByIdWithFreeFunds), ParseCancelMultipleOrdersByIdWithFreeFunds)
	RegisterParser(uint64(phoenix_v1.Instruction_WithdrawFunds), ParseWithdrawFunds)
	RegisterParser(uint64(phoenix_v1.Instruction_DepositFunds), ParseDepositFunds)
	RegisterParser(uint64(phoenix_v1.Instruction_RequestSeat), ParseRequestSeat)
	RegisterParser(uint64(phoenix_v1.Instruction_Log), ParseLog)
	RegisterParser(uint64(phoenix_v1.Instruction_PlaceMultiplePostOnlyOrders), ParsePlaceMultiplePostOnlyOrders)
	RegisterParser(uint64(phoenix_v1.Instruction_PlaceMultiplePostOnlyOrdersWithFreeFunds), ParsePlaceMultiplePostOnlyOrdersWithFreeFunds)
	RegisterParser(uint64(phoenix_v1.Instruction_InitializeMarket), ParseInitializeMarket)
	RegisterParser(uint64(phoenix_v1.Instruction_ClaimAuthority), ParseClaimAuthority)
	RegisterParser(uint64(phoenix_v1.Instruction_NameSuccessor), ParseNameSuccessor)
	RegisterParser(uint64(phoenix_v1.Instruction_ChangeMarketStatus), ParseChangeMarketStatus)
	RegisterParser(uint64(phoenix_v1.Instruction_ChangeSeatStatus), ParseChangeSeatStatus)
	RegisterParser(uint64(phoenix_v1.Instruction_RequestSeatAuthorized), ParseRequestSeatAuthorized)
	RegisterParser(uint64(phoenix_v1.Instruction_EvictSeat), ParseEvictSeat)
	RegisterParser(uint64(phoenix_v1.Instruction_ForceCancelOrders), ParseForceCancelOrders)
	RegisterParser(uint64(phoenix_v1.Instruction_CollectFees), ParseCollectFees)
	RegisterParser(uint64(phoenix_v1.Instruction_ChangeFeeRecipient), ParseChangeFeeRecipient)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	inst, err := phoenix_v1.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
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

func ParseSwap(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*phoenix_v1.Swap)
	swap := &types.Swap{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetMarketAccount().PublicKey,
		User: inst1.GetTraderAccount().PublicKey,
	}
	transfers := in.FindChildrenTransfers()
	if len(transfers) >= 2 {
		swap.InputTransfer = transfers[1]
		swap.OutputTransfer = transfers[0]
	}
	in.Event = []interface{}{swap}
}
func ParseSwapWithFreeFunds(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse swap", "program", phoenix_v1.ProgramName)
}
func ParsePlaceLimitOrder(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse place limit order", "program", phoenix_v1.ProgramName)
}
func ParsePlaceLimitOrderWithFreeFunds(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse place limit order with free funds", "program", phoenix_v1.ProgramName)
}
func ParseReduceOrder(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseReduceOrderWithFreeFunds(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCancelAllOrders(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCancelAllOrdersWithFreeFunds(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCancelUpTo(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCancelUpToWithFreeFunds(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCancelMultipleOrdersById(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCancelMultipleOrdersByIdWithFreeFunds(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseWithdrawFunds(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseDepositFunds(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseRequestSeat(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseLog(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParsePlaceMultiplePostOnlyOrders(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParsePlaceMultiplePostOnlyOrdersWithFreeFunds(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseInitializeMarket(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseClaimAuthority(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseNameSuccessor(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseChangeMarketStatus(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseChangeSeatStatus(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseRequestSeatAuthorized(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseEvictSeat(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseForceCancelOrders(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCollectFees(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseChangeFeeRecipient(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
}

// Default
func ParseDefault(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
