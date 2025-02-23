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

type Parser func(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(phoenix_v1.ProgramID, phoenix_v1.ProgramName, program.OrderBook, 1, ProgramParser)
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

func ProgramParser(transaction *types.Transaction, index int) error {
	in := transaction.Instructions[index]
	inst, err := phoenix_v1.DecodeInstruction(in.Raw.AccountValues, in.Raw.DataBytes)
	if err != nil {
		return err
	}
	id := uint64(inst.TypeID.Uint32())
	parser, ok := Parsers[id]
	if !ok {
		return errors.New("parser not found")
	}
	return parser(inst, transaction, index)
}

func ParseSwap(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	inst1 := inst.Impl.(*phoenix_v1.Swap)
	in := transaction.Instructions[index]
	swap := &types.Swap{
		Dex:  in.Raw.ProgID,
		Pool: inst1.GetMarketAccount().PublicKey,
		User: inst1.GetTraderAccount().PublicKey,
	}
	if transfer := transaction.FindNextTransferByTo(index, inst1.GetBaseVaultAccount().PublicKey); transfer != nil {
		swap.InputTransfer = transfer
	}
	if transfer := transaction.FindNextTransferByTo(index, inst1.GetQuoteVaultAccount().PublicKey); transfer != nil {
		swap.InputTransfer = transfer
	}
	if transfer := transaction.FindNextTransferByFrom(index, inst1.GetBaseVaultAccount().PublicKey); transfer != nil {
		swap.OutputTransfer = transfer
	}
	if transfer := transaction.FindNextTransferByFrom(index, inst1.GetQuoteVaultAccount().PublicKey); transfer != nil {
		swap.OutputTransfer = transfer
	}
	in.Event = []interface{}{swap}
	return nil
}
func ParseSwapWithFreeFunds(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	log.Logger.Info("ignore parse swap", "program", phoenix_v1.ProgramName)
	return nil
}
func ParsePlaceLimitOrder(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	log.Logger.Info("ignore parse place limit order", "program", phoenix_v1.ProgramName)
	return nil
}
func ParsePlaceLimitOrderWithFreeFunds(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	log.Logger.Info("ignore parse place limit order with free funds", "program", phoenix_v1.ProgramName)
	return nil
}
func ParseReduceOrder(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseReduceOrderWithFreeFunds(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseCancelAllOrders(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseCancelAllOrdersWithFreeFunds(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseCancelUpTo(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseCancelUpToWithFreeFunds(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseCancelMultipleOrdersById(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseCancelMultipleOrdersByIdWithFreeFunds(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseWithdrawFunds(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseDepositFunds(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseRequestSeat(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseLog(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParsePlaceMultiplePostOnlyOrders(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParsePlaceMultiplePostOnlyOrdersWithFreeFunds(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseInitializeMarket(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseClaimAuthority(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseNameSuccessor(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseChangeMarketStatus(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseChangeSeatStatus(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseRequestSeatAuthorized(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseEvictSeat(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseForceCancelOrders(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseCollectFees(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}
func ParseChangeFeeRecipient(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}

// Default
func ParseDefault(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	return nil
}

// Fault
func ParseFault(inst *phoenix_v1.Instruction, transaction *types.Transaction, index int) error {
	panic("not supported")
}
