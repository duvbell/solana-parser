package pump

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blockchain-develop/solana-parser/log"
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/programs/pumpfun"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

var (
	Instruction_AnchorSelfCPILog = ag_binary.TypeID([8]byte{228, 69, 165, 46, 81, 203, 154, 29})
)

func init() {
	program.RegisterParser(pumpfun.ProgramID, pumpfun.ProgramName, program.Swap, ProgramParser)
	RegisterParser(uint64(pumpfun.Instruction_Initialize.Uint32()), ParseInitialize)
	RegisterParser(uint64(pumpfun.Instruction_Create.Uint32()), ParseCreate)
	RegisterParser(uint64(pumpfun.Instruction_Buy.Uint32()), ParseBuy)
	RegisterParser(uint64(pumpfun.Instruction_Sell.Uint32()), ParseSell)
	RegisterParser(uint64(pumpfun.Instruction_Withdraw.Uint32()), ParseWithdraw)
	RegisterParser(uint64(pumpfun.Instruction_SetParams.Uint32()), ParseDefault)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	dec := ag_binary.NewBorshDecoder(in.Instruction.Data)
	typeID, err := dec.ReadTypeID()
	if typeID == Instruction_AnchorSelfCPILog {
		return nil
	}
	inst, err := pumpfun.DecodeInstruction(in.AccountMetas(meta.Accounts), in.Instruction.Data)
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

// Initialize
func ParseInitialize(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	//log.Logger.Info("ignore parse initialize", "program", pumpfun.ProgramName)
}

// Create
func ParseCreate(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	//log.Logger.Info("ignore parse create", "program", pumpfun.ProgramName)
	inst1 := inst.Impl.(*pumpfun.Create)
	memeMint := &types.MemeCreate{
		Dex:                    in.Instruction.ProgramId,
		Mint:                   inst1.GetMintAccount().PublicKey,
		User:                   inst1.GetUserAccount().PublicKey,
		BondingCurve:           inst1.GetBondingCurveAccount().PublicKey,
		AssociatedBondingCurve: inst1.GetAssociatedBondingCurveAccount().PublicKey,
	}
	mintTos := in.FindChildrenMintTos()
	if len(mintTos) >= 1 {
		memeMint.MintTo = mintTos[0]
	}
	in.Event = []interface{}{memeMint}

	children := in.FindChildrenPrograms(pumpfun.ProgramID)
	if len(children) > 0 {
		myLog := children[0]
		data, err := myLog.Instruction.Data.MarshalJSON()
		if err != nil {
			return
		}
		var tradeEvent pumpfun.TradeEvent
		if err = ag_binary.NewBorshDecoder(data).Decode(&tradeEvent); err != nil {
			return
		}
		memeBuyEvent := types.MemeBuyEvent{}
		in.Receipt = []interface{}{&memeBuyEvent}
	}
}

// Buy
func ParseBuy(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	//log.Logger.Info("ignore parse buy", "program", pumpfun.ProgramName)
	inst1 := inst.Impl.(*pumpfun.Buy)
	memeBuy := &types.MemeBuy{
		Dex:                    in.Instruction.ProgramId,
		Mint:                   inst1.GetMintAccount().PublicKey,
		User:                   inst1.GetUserAccount().PublicKey,
		BondingCurve:           inst1.GetBondingCurveAccount().PublicKey,
		AssociatedBondingCurve: inst1.GetAssociatedBondingCurveAccount().PublicKey,
	}
	transfers := in.FindChildrenTransfers()
	if len(transfers) >= 2 {
		memeBuy.MintTransfer = transfers[0]
		memeBuy.SolTransfer = transfers[1]
	}
	if len(transfers) >= 3 {
		memeBuy.FeeTransfer = transfers[2]
	}
	in.Event = []interface{}{memeBuy}

	children := in.FindChildrenPrograms(pumpfun.ProgramID)
	if len(children) > 0 {
		myLog := children[0]
		data, err := myLog.Instruction.Data.MarshalJSON()
		if err != nil {
			return
		}
		var tradeEvent pumpfun.TradeEvent
		if err = ag_binary.NewBorshDecoder(data).Decode(&tradeEvent); err != nil {
			return
		}
		memeSellEvent := types.MemeSellEvent{}
		in.Receipt = []interface{}{&memeSellEvent}
	}
}

// Sell
func ParseSell(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	//log.Logger.Info("ignore parse sell", "program", pumpfun.ProgramName)
	inst1 := inst.Impl.(*pumpfun.Sell)
	memeSell := &types.MemeSell{
		Dex:                    in.Instruction.ProgramId,
		Mint:                   inst1.GetMintAccount().PublicKey,
		User:                   inst1.GetUserAccount().PublicKey,
		BondingCurve:           inst1.GetBondingCurveAccount().PublicKey,
		AssociatedBondingCurve: inst1.GetAssociatedBondingCurveAccount().PublicKey,
	}
	transfers := in.FindChildrenTransfers()
	if len(transfers) >= 1 {
		memeSell.MintTransfer = transfers[0]
	}
	in.Event = []interface{}{memeSell}

	children := in.FindChildrenPrograms(pumpfun.ProgramID)
	if len(children) > 0 {
		myLog := children[0]
		myEventJson := myLog.Event[0].(json.RawMessage)
		var createEvent types.MemeSellEvent
		json.Unmarshal(myEventJson, &createEvent)
		in.Receipt = []interface{}{&createEvent}
	}
}

// Sell
func ParseWithdraw(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse withdraw", "program", pumpfun.ProgramName)
}

func ParseLog(in *types.Instruction, meta *types.Meta) {
	data := []byte(in.Instruction.Data)
	data = data[8:]
	event, err := types.DecodeEventsFromEmitCPI(data)
	if err != nil {
		return
	}
	fmt.Printf("%s\n", string(event))
	in.Event = []interface{}{json.RawMessage(event)}
}

// Default
func ParseDefault(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
