package jupiter

import (
	"bytes"
	"errors"

	"github.com/blockchain-develop/solana-parser/program"

	"github.com/blockchain-develop/solana-parser/log"
	"github.com/blockchain-develop/solana-parser/types"
	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/programs/jupiter"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *jupiter.Instruction, in *types.Instruction, meta *types.Meta) error

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

var (
	Instruction_AnchorSelfCPILog = ag_binary.TypeID([8]byte{228, 69, 165, 46, 81, 203, 154, 29})
	Event_Create                 = [8]byte{0x1b, 0x72, 0xa9, 0x4d, 0xde, 0xeb, 0x63, 0x76}
	Event_Swap                   = [8]byte{0xbd, 0xdb, 0x7f, 0xd3, 0x4e, 0xe6, 0x61, 0xee}
)

func init() {
	program.RegisterParser(jupiter.ProgramID, jupiter.ProgramName, program.Swap, 1, ProgramParser)
	RegisterParser(uint64(jupiter.Instruction_Route.Uint32()), ParseSwap)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	dec := ag_binary.NewBorshDecoder(in.RawInstruction.DataBytes)
	typeID, err := dec.ReadTypeID()
	// 特殊的，处理log
	if typeID == Instruction_AnchorSelfCPILog {
		return ParseSwapLog(nil, in, meta)
	}
	if typeID != jupiter.Instruction_Route {
		return nil
	}
	inst, err := jupiter.DecodeInstruction(in.RawInstruction.AccountValues, in.RawInstruction.DataBytes)
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

// Swap
func ParseSwap(inst *jupiter.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*jupiter.Route)
	route := &types.Route{
		Router: in.RawInstruction.ProgID,
		User:   inst1.GetUserTransferAuthorityAccount().PublicKey,
	}
	stepSize := len(*inst1.RoutePlan)
	if len(in.Children) != stepSize*2 {
		err := errors.New("jupiter route is invalid")
		log.Logger.Error("jupiter", "error", err)
		return err
	}
	for i := 0; i < stepSize; i++ {
		swapIn := in.Children[i*2]
		eventIn := in.Children[i*2+1]
		var swap *types.Swap
		if len(swapIn.Event) == 1 {
			swap = swapIn.Event[0].(*types.Swap)
		}
		routePlan := &types.RoutePlan{}
		var swapEvent *types.SwapEvent
		if len(eventIn.Event) == 1 {
			swapEvent = eventIn.Event[0].(*types.SwapEvent)
		}
		route.RouteSteps = append(route.RouteSteps, &types.RouteStep{
			Swap:      swap,
			RoutePlan: routePlan,
			SwapEvent: swapEvent,
		})
	}
	in.Event = []interface{}{route}
	return nil
}

func ParseSwapLog(inst *jupiter.Instruction, in *types.Instruction, meta *types.Meta) error {
	data := in.RawInstruction.DataBytes
	dec := ag_binary.NewBorshDecoder(data)
	instId, _ := dec.ReadBytes(8)
	eventId, _ := dec.ReadBytes(8)
	if bytes.Compare(instId, Instruction_AnchorSelfCPILog[:]) != 0 || bytes.Compare(eventId, jupiter.SwapEventEventDataDiscriminator[:]) != 0 {
		return nil
	}
	var swapEvent jupiter.SwapEvent
	if err := dec.Decode(&swapEvent); err != nil {
		return err
	}
	mySwapEvent := types.SwapEvent{
		Amm:          swapEvent.Amm,
		InputMint:    swapEvent.InputMint,
		InputAmount:  swapEvent.InputAmount,
		OutputMint:   swapEvent.OutputMint,
		OutputAmount: swapEvent.OutputAmount,
	}
	in.Event = []interface{}{&mySwapEvent}
	return nil
}

// Default
func ParseDefault(inst *jupiter.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

// Fault
func ParseFault(inst *jupiter.Instruction, in *types.Instruction, meta *types.Meta) error {
	panic("not supported")
}
