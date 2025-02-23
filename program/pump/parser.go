package pump

import (
	"errors"
	"github.com/blockchain-develop/solana-parser/log"
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/programs/pumpfun"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *pumpfun.Instruction, transaction *types.Transaction, index int) error

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

var (
	Instruction_AnchorSelfCPILog = ag_binary.TypeID([8]byte{228, 69, 165, 46, 81, 203, 154, 29})
)

func init() {
	program.RegisterParser(pumpfun.ProgramID, pumpfun.ProgramName, program.Swap, 1, ProgramParser)
	RegisterParser(uint64(pumpfun.Instruction_Initialize.Uint32()), ParseInitialize)
	RegisterParser(uint64(pumpfun.Instruction_Create.Uint32()), ParseCreate)
	RegisterParser(uint64(pumpfun.Instruction_Buy.Uint32()), ParseBuy)
	RegisterParser(uint64(pumpfun.Instruction_Sell.Uint32()), ParseSell)
	RegisterParser(uint64(pumpfun.Instruction_Withdraw.Uint32()), ParseWithdraw)
	RegisterParser(uint64(pumpfun.Instruction_SetParams.Uint32()), ParseDefault)
}

func ProgramParser(transaction *types.Transaction, index int) error {
	in := transaction.Instructions[index]
	dec := ag_binary.NewBorshDecoder(in.Raw.DataBytes)
	typeID, err := dec.ReadTypeID()
	if typeID == Instruction_AnchorSelfCPILog {
		return nil
	}
	inst, err := pumpfun.DecodeInstruction(in.Raw.AccountValues, in.Raw.DataBytes)
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

// Initialize
func ParseInitialize(inst *pumpfun.Instruction, transaction *types.Transaction, index int) error {
	//log.Logger.Info("ignore parse initialize", "program", pumpfun.ProgramName)
	return nil
}

// Create
func ParseCreate(inst *pumpfun.Instruction, transaction *types.Transaction, index int) error {
	//log.Logger.Info("ignore parse create", "program", pumpfun.ProgramName)
	inst1 := inst.Impl.(*pumpfun.Create)
	in := transaction.Instructions[index]
	memeMint := &types.MemeCreate{
		Dex:                    in.Raw.ProgID,
		Mint:                   inst1.GetMintAccount().PublicKey,
		User:                   inst1.GetUserAccount().PublicKey,
		BondingCurve:           inst1.GetBondingCurveAccount().PublicKey,
		AssociatedBondingCurve: inst1.GetAssociatedBondingCurveAccount().PublicKey,
	}
	memeMint.MintTo = transaction.FindNextMintTo(index, inst1.GetAssociatedBondingCurveAccount().PublicKey)
	in.Event = []interface{}{memeMint}

	myLog := transaction.FindNextInstructionByProgram(index, pumpfun.ProgramID)
	if myLog != nil {
		data := myLog.Raw.DataBytes
		dec := ag_binary.NewBorshDecoder(data)
		dec.ReadBytes(8)
		dec.ReadBytes(8)
		var createEvent pumpfun.CreateEvent
		if err := dec.Decode(&createEvent); err != nil {
			return err
		}
		memeCreateEvent := types.MemeCreateEvent{
			Name:         createEvent.Name,
			Symbol:       createEvent.Symbol,
			Uri:          createEvent.Uri,
			Mint:         createEvent.Mint,
			BondingCurve: createEvent.BondingCurve,
			User:         createEvent.User,
		}
		in.Receipt = []interface{}{&memeCreateEvent}
	}
	return nil
}

// Buy
func ParseBuy(inst *pumpfun.Instruction, transaction *types.Transaction, index int) error {
	//log.Logger.Info("ignore parse buy", "program", pumpfun.ProgramName)
	inst1 := inst.Impl.(*pumpfun.Buy)
	in := transaction.Instructions[index]
	memeBuy := &types.MemeBuy{
		Dex:                    in.Raw.ProgID,
		Mint:                   inst1.GetMintAccount().PublicKey,
		User:                   inst1.GetUserAccount().PublicKey,
		BondingCurve:           inst1.GetBondingCurveAccount().PublicKey,
		AssociatedBondingCurve: inst1.GetAssociatedBondingCurveAccount().PublicKey,
	}
	memeBuy.SolTransfer = transaction.FindNextTransferByTo(index, inst1.GetBondingCurveAccount().PublicKey)
	memeBuy.FeeTransfer = transaction.FindNextTransferByTo(index, inst1.GetFeeRecipientAccount().PublicKey)
	memeBuy.MintTransfer = transaction.FindNextTransferByFrom(index, inst1.GetAssociatedBondingCurveAccount().PublicKey)
	in.Event = []interface{}{memeBuy}

	myLog := transaction.FindNextInstructionByProgram(index, pumpfun.ProgramID)
	if myLog != nil {
		data := myLog.Raw.DataBytes
		dec := ag_binary.NewBorshDecoder(data)
		dec.ReadBytes(8)
		dec.ReadBytes(8)
		var tradeEvent pumpfun.TradeEvent
		if err := dec.Decode(&tradeEvent); err != nil {
			return err
		}
		memeBuyEvent := types.MemeBuyEvent{
			Mint:                 tradeEvent.Mint,
			SolAmount:            tradeEvent.SolAmount,
			TokenAmount:          tradeEvent.TokenAmount,
			IsBuy:                tradeEvent.IsBuy,
			User:                 tradeEvent.User,
			Timestamp:            tradeEvent.Timestamp,
			VirtualSolReserves:   tradeEvent.VirtualSolReserves,
			VirtualTokenReserves: tradeEvent.VirtualTokenReserves,
		}
		in.Receipt = []interface{}{&memeBuyEvent}
	}
	return nil
}

// Sell
func ParseSell(inst *pumpfun.Instruction, transaction *types.Transaction, index int) error {
	//log.Logger.Info("ignore parse sell", "program", pumpfun.ProgramName)
	inst1 := inst.Impl.(*pumpfun.Sell)
	in := transaction.Instructions[index]
	memeSell := &types.MemeSell{
		Dex:                    in.Raw.ProgID,
		Mint:                   inst1.GetMintAccount().PublicKey,
		User:                   inst1.GetUserAccount().PublicKey,
		BondingCurve:           inst1.GetBondingCurveAccount().PublicKey,
		AssociatedBondingCurve: inst1.GetAssociatedBondingCurveAccount().PublicKey,
	}
	memeSell.MintTransfer = transaction.FindNextTransferByTo(index, inst1.GetAssociatedBondingCurveAccount().PublicKey)
	in.Event = []interface{}{memeSell}

	myLog := transaction.FindNextInstructionByProgram(index, pumpfun.ProgramID)
	if myLog != nil {
		data := myLog.Raw.DataBytes
		dec := ag_binary.NewBorshDecoder(data)
		dec.ReadBytes(8)
		dec.ReadBytes(8)

		var tradeEvent pumpfun.TradeEvent
		if err := dec.Decode(&tradeEvent); err != nil {
			return err
		}
		memeSellEvent := types.MemeSellEvent{
			Mint:                 tradeEvent.Mint,
			SolAmount:            tradeEvent.SolAmount,
			TokenAmount:          tradeEvent.TokenAmount,
			IsBuy:                tradeEvent.IsBuy,
			User:                 tradeEvent.User,
			Timestamp:            tradeEvent.Timestamp,
			VirtualSolReserves:   tradeEvent.VirtualSolReserves,
			VirtualTokenReserves: tradeEvent.VirtualTokenReserves,
		}
		in.Receipt = []interface{}{&memeSellEvent}
	}
	return nil
}

// Sell
func ParseWithdraw(inst *pumpfun.Instruction, transaction *types.Transaction, index int) error {
	log.Logger.Info("ignore parse withdraw", "program", pumpfun.ProgramName)
	return nil
}

// Default
func ParseDefault(inst *pumpfun.Instruction, transaction *types.Transaction, index int) error {
	return nil
}

// Fault
func ParseFault(inst *pumpfun.Instruction, transaction *types.Transaction, index int) {
	panic("not supported")
}
