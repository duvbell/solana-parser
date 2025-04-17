package pump

import (
	"bytes"
	"errors"

	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/programs/pumpfun"
	"github.com/solana-parser/log"
	"github.com/solana-parser/program"
	"github.com/solana-parser/types"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) error

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

var (
	Instruction_AnchorSelfCPILog = ag_binary.TypeID([8]byte{228, 69, 165, 46, 81, 203, 154, 29})
	Event_Create                 = [8]byte{0x1b, 0x72, 0xa9, 0x4d, 0xde, 0xeb, 0x63, 0x76}
	Event_Swap                   = [8]byte{0xbd, 0xdb, 0x7f, 0xd3, 0x4e, 0xe6, 0x61, 0xee}
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

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	dec := ag_binary.NewBorshDecoder(in.RawInstruction.DataBytes)
	typeID, err := dec.ReadTypeID()
	if typeID == Instruction_AnchorSelfCPILog {
		return nil
	}
	inst, err := pumpfun.DecodeInstruction(in.RawInstruction.AccountValues, in.RawInstruction.DataBytes)
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

// Initialize
func ParseInitialize(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) error {
	//log.Logger.Info("ignore parse initialize", "program", pumpfun.ProgramName)
	return nil
}

// Create
func ParseCreate(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) error {
	//log.Logger.Info("ignore parse create", "program", pumpfun.ProgramName)
	inst1 := inst.Impl.(*pumpfun.Create)
	in.ParsedInstruction = inst1
	memeMint := &types.MemeCreate{
		Dex:                    in.RawInstruction.ProgID,
		Mint:                   inst1.GetMintAccount().PublicKey,
		User:                   inst1.GetUserAccount().PublicKey,
		BondingCurve:           inst1.GetBondingCurveAccount().PublicKey,
		AssociatedBondingCurve: inst1.GetAssociatedBondingCurveAccount().PublicKey,
	}
	memeMint.MintTo = in.FindNextMintTo(inst1.GetAssociatedBondingCurveAccount().PublicKey)
	in.Event = []interface{}{memeMint}

	myLog := in.FindNextProgram(pumpfun.ProgramID)
	if myLog == nil {
		return nil
	}
	data := myLog.RawInstruction.DataBytes
	dec := ag_binary.NewBorshDecoder(data)
	instId, _ := dec.ReadBytes(8)
	eventId, _ := dec.ReadBytes(8)
	if bytes.Compare(instId, Instruction_AnchorSelfCPILog[:]) != 0 || bytes.Compare(eventId, Event_Create[:]) != 0 {
		return nil
	}
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

	return nil
}

// Buy
func ParseBuy(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) error {
	//log.Logger.Info("ignore parse buy", "program", pumpfun.ProgramName)
	inst1 := inst.Impl.(*pumpfun.Buy)
	in.ParsedInstruction = inst1
	memeBuy := &types.MemeBuy{
		Dex:                    in.RawInstruction.ProgID,
		Mint:                   inst1.GetMintAccount().PublicKey,
		User:                   inst1.GetUserAccount().PublicKey,
		BondingCurve:           inst1.GetBondingCurveAccount().PublicKey,
		AssociatedBondingCurve: inst1.GetAssociatedBondingCurveAccount().PublicKey,
	}
	memeBuy.SolTransfer = in.FindNextTransferByTo(inst1.GetBondingCurveAccount().PublicKey)
	memeBuy.FeeTransfer = in.FindNextTransferByTo(inst1.GetFeeRecipientAccount().PublicKey)
	memeBuy.MintTransfer = in.FindNextTransferByFrom(inst1.GetAssociatedBondingCurveAccount().PublicKey)
	in.Event = []interface{}{memeBuy}

	myLog := in.FindNextProgram(pumpfun.ProgramID)
	if myLog == nil {
		return nil
	}
	data := myLog.RawInstruction.DataBytes
	dec := ag_binary.NewBorshDecoder(data)
	instId, _ := dec.ReadBytes(8)
	eventId, _ := dec.ReadBytes(8)
	if bytes.Compare(instId, Instruction_AnchorSelfCPILog[:]) != 0 || bytes.Compare(eventId, Event_Swap[:]) != 0 {
		return nil
	}
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
	return nil
}

// Sell
func ParseSell(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) error {
	//log.Logger.Info("ignore parse sell", "program", pumpfun.ProgramName)
	inst1 := inst.Impl.(*pumpfun.Sell)
	in.ParsedInstruction = inst1
	memeSell := &types.MemeSell{
		Dex:                    in.RawInstruction.ProgID,
		Mint:                   inst1.GetMintAccount().PublicKey,
		User:                   inst1.GetUserAccount().PublicKey,
		BondingCurve:           inst1.GetBondingCurveAccount().PublicKey,
		AssociatedBondingCurve: inst1.GetAssociatedBondingCurveAccount().PublicKey,
	}
	memeSell.MintTransfer = in.FindNextTransferByTo(inst1.GetAssociatedBondingCurveAccount().PublicKey)
	in.Event = []interface{}{memeSell}

	myLog := in.FindNextProgram(pumpfun.ProgramID)
	if myLog == nil {
		return nil
	}
	data := myLog.RawInstruction.DataBytes
	dec := ag_binary.NewBorshDecoder(data)
	instId, _ := dec.ReadBytes(8)
	eventId, _ := dec.ReadBytes(8)
	if bytes.Compare(instId, Instruction_AnchorSelfCPILog[:]) != 0 || bytes.Compare(eventId, Event_Swap[:]) != 0 {
		return nil
	}
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
	return nil
}

// Sell
func ParseWithdraw(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) error {
	log.Logger.Info("ignore parse withdraw", "program", pumpfun.ProgramName)
	return nil
}

// Default
func ParseDefault(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

// Fault
func ParseFault(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
