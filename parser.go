package solanaparser

import (
	"encoding/json"
	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/raydiumamm"
	amm_v4 "github.com/gagliardetto/solana-go/programs/raydiumclmm"
	stable_swap "github.com/gagliardetto/solana-go/programs/stabblestableswap"
	whirlpool "github.com/gagliardetto/solana-go/programs/whirlpool"
)

var (
	DefaultParse      = make(map[solana.PublicKey]Parser, 0)
	WhirlPool         = solana.MustPublicKeyFromBase58("whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc")
	RaydiumClmm       = solana.MustPublicKeyFromBase58("CAMMCzo5YL8w4VFF8KVHrK22GGUsp5VTaW7grrKgrWqK")
	RaydiumAMM        = solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")
	StabbleStableSwap = solana.MustPublicKeyFromBase58("swapNyd8XiQwJ6ianp9snpu4brUqFxadzvHebnAXjJZ")
)

func init() {
	RegisterParser(RaydiumAMM, RaydiumAmmParser)
	RegisterParser(RaydiumClmm, RaydiumClmmParser)
	RegisterParser(WhirlPool, WhirlPoolParser)
	RegisterParser(solana.SystemProgramID, SystemParser)
	RegisterParser(solana.TokenProgramID, TokenParser)
	RegisterParser(solana.Token2022ProgramID, Token2022Parser)
	RegisterParser(StabbleStableSwap, StabbleStableSwapParser)
}

func RegisterParser(program solana.PublicKey, p Parser) {
	DefaultParse[program] = p
}

type Parser func(in *Instruction, meta *Meta) ([]interface{}, []interface{})

type SystemTransfer struct {
	Destination solana.PublicKey `json:"destination"`
	Lamports    uint64           `json:"lamports"`
	Source      solana.PublicKey `json:"source"`
}

type SystemInstruction struct {
	T    string `json:"type"`
	Info interface{}
	Raw  json.RawMessage `json:"info"`
}

func (j *SystemInstruction) UnmarshalJSON(data []byte) error {
	type Aux SystemInstruction
	aux := (*Aux)(j)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	var object interface{}
	switch j.T {
	case "transfer":
		object = &SystemTransfer{}
	default:
		return nil
	}
	if err := json.Unmarshal(j.Raw, object); err != nil {
		return err
	}
	j.Info = object
	return nil
}

func SystemParser(in *Instruction, meta *Meta) ([]interface{}, []interface{}) {
	inJson, _ := in.Instruction.Parsed.MarshalJSON()
	var systemInstruction SystemInstruction
	err := json.Unmarshal(inJson, &systemInstruction)
	if err != nil {
		return nil, nil
	}
	switch systemInstruction.Info.(type) {
	case *SystemInstruction:
		systemTransfer := systemInstruction.Info.(*SystemTransfer)
		transfer := &Transfer{
			Mint:   solana.MustPublicKeyFromBase58("11111111111111111111111111111111"),
			Amount: systemTransfer.Lamports,
			From:   systemTransfer.Source,
			To:     systemTransfer.Destination,
		}
		return []interface{}{transfer}, nil
	default:
		return nil, nil
	}
}

type TokenTransfer struct {
	Destination solana.PublicKey `json:"destination"`
	Lamports    uint64           `json:"amount,string"`
	Source      solana.PublicKey `json:"source"`
	Authority   solana.PublicKey `json:"authority"`
	Mint        solana.PublicKey `json:"mint"`
	TokenAmount struct {
		Amount   uint64 `json:"amount,string"`
		Decimals uint64 `json:"decimals"`
	} `json:"tokenAmount"`
}

type TokenMint struct {
	Account   solana.PublicKey `json:"account"`
	Amount    uint64           `json:"amount,string"`
	Authority solana.PublicKey `json:"mintAuthority"`
	Mint      solana.PublicKey `json:"mint"`
}

type TokenBurn struct {
	Account   solana.PublicKey `json:"account"`
	Amount    uint64           `json:"amount,string"`
	Authority solana.PublicKey `json:"authority"`
	Mint      solana.PublicKey `json:"mint"`
}

type TokenInitialize struct {
	Account solana.PublicKey `json:"account"`
	Owner   solana.PublicKey `json:"owner"`
	Mint    solana.PublicKey `json:"mint"`
}

type TokenInstruction struct {
	T    string `json:"type"`
	Info interface{}
	Raw  json.RawMessage `json:"info"`
}

func (j *TokenInstruction) UnmarshalJSON(data []byte) error {
	type Aux TokenInstruction
	aux := (*Aux)(j)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	var object interface{}
	switch j.T {
	case "transfer":
		object = &TokenTransfer{}
	case "transferChecked":
		object = &TokenTransfer{}
	case "mintTo":
		object = &TokenMint{}
	case "burn":
		object = &TokenBurn{}
	case "initializeAccount":
		object = &TokenInitialize{}
	case "initializeAccount3":
		object = &TokenInitialize{}
	default:
		return nil
	}
	if err := json.Unmarshal(j.Raw, object); err != nil {
		return err
	}
	j.Info = object
	return nil
}

func TokenParser(in *Instruction, meta *Meta) ([]interface{}, []interface{}) {
	inJson, _ := in.Instruction.Parsed.MarshalJSON()
	var tokenInstruction TokenInstruction
	json.Unmarshal(inJson, &tokenInstruction)
	switch tokenInstruction.Info.(type) {
	case *TokenTransfer:
		tokenTransfer := tokenInstruction.Info.(*TokenTransfer)
		mint := meta.TokenMint[tokenTransfer.Source]
		from := tokenTransfer.Source
		if k, ok := meta.TokenOwner[tokenTransfer.Source]; ok {
			from = k
		}
		to := tokenTransfer.Destination
		if k, ok := meta.TokenOwner[tokenTransfer.Destination]; ok {
			to = k
		}
		amount := tokenTransfer.Lamports
		if amount == 0 {
			amount = tokenTransfer.TokenAmount.Amount
		}
		transfer := &Transfer{
			Mint:   mint,
			Amount: amount,
			From:   from,
			To:     to,
		}
		return []interface{}{transfer}, nil
	case *TokenMint:
		tokenMint := tokenInstruction.Info.(*TokenMint)
		account := tokenMint.Account
		if k, ok := meta.TokenOwner[tokenMint.Account]; ok {
			account = k
		}
		mintTo := &MintTo{
			Mint:    tokenMint.Mint,
			Amount:  tokenMint.Amount,
			Account: account,
		}
		return []interface{}{mintTo}, nil
	case *TokenBurn:
		tokenBurn := tokenInstruction.Info.(*TokenBurn)
		account := tokenBurn.Account
		if k, ok := meta.TokenOwner[tokenBurn.Account]; ok {
			account = k
		}
		burn := &Burn{
			Mint:    tokenBurn.Mint,
			Amount:  tokenBurn.Amount,
			Account: account,
		}
		return []interface{}{burn}, nil
	case *TokenInitialize:
		tokenInitialize := tokenInstruction.Info.(*TokenInitialize)
		init := &Initialize{
			Mint:    tokenInitialize.Mint,
			Account: tokenInitialize.Account,
			Owner:   tokenInitialize.Owner,
		}
		// update token owner & mint by spl token instructions
		meta.TokenOwner[init.Account] = init.Owner
		meta.TokenMint[init.Account] = init.Mint
		return []interface{}{init}, nil
	default:
		return nil, nil
	}
}

func Token2022Parser(in *Instruction, meta *Meta) ([]interface{}, []interface{}) {
	inJson, _ := in.Instruction.Parsed.MarshalJSON()
	var tokenInstruction TokenInstruction
	json.Unmarshal(inJson, &tokenInstruction)
	switch tokenInstruction.Info.(type) {
	case *TokenTransfer:
		tokenTransfer := tokenInstruction.Info.(*TokenTransfer)
		mint := meta.TokenMint[tokenTransfer.Source]
		from := tokenTransfer.Source
		if k, ok := meta.TokenOwner[tokenTransfer.Source]; ok {
			from = k
		}
		to := tokenTransfer.Destination
		if k, ok := meta.TokenOwner[tokenTransfer.Destination]; ok {
			to = k
		}
		amount := tokenTransfer.Lamports
		if amount == 0 {
			amount = tokenTransfer.TokenAmount.Amount
		}
		transfer := &Transfer{
			Mint:   mint,
			Amount: amount,
			From:   from,
			To:     to,
		}
		return []interface{}{transfer}, nil
	case *TokenMint:
		tokenMint := tokenInstruction.Info.(*TokenMint)
		account := tokenMint.Account
		if k, ok := meta.TokenOwner[tokenMint.Account]; ok {
			account = k
		}
		mintTo := &MintTo{
			Mint:    tokenMint.Mint,
			Amount:  tokenMint.Amount,
			Account: account,
		}
		return []interface{}{mintTo}, nil
	case *TokenBurn:
		tokenBurn := tokenInstruction.Info.(*TokenBurn)
		account := tokenBurn.Account
		if k, ok := meta.TokenOwner[tokenBurn.Account]; ok {
			account = k
		}
		burn := &Burn{
			Mint:    tokenBurn.Mint,
			Amount:  tokenBurn.Amount,
			Account: account,
		}
		return []interface{}{burn}, nil
	case *TokenInitialize:
		tokenInitialize := tokenInstruction.Info.(*TokenInitialize)
		init := &Initialize{
			Mint:    tokenInitialize.Mint,
			Account: tokenInitialize.Account,
			Owner:   tokenInitialize.Owner,
		}
		// update token owner & mint by spl token instructions
		meta.TokenOwner[init.Account] = init.Owner
		meta.TokenMint[init.Account] = init.Mint
		return []interface{}{init}, nil
	default:
		return nil, nil
	}
}

func RaydiumAmmParser(in *Instruction, meta *Meta) ([]interface{}, []interface{}) {
	inst := new(raydium_amm.Instruction)
	instruction := in.Instruction
	err := ag_binary.NewBorshDecoder(instruction.Data).Decode(inst)
	if err != nil {
		return nil, nil
	}
	accounts := make([]*solana.AccountMeta, 0)
	for _, item := range instruction.Accounts {
		account := meta.Accounts[item]
		accounts = append(accounts, account)
	}
	insertAccount := func(accounts []*solana.AccountMeta, index int) []*solana.AccountMeta {
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
	switch inst.TypeID.Uint8() {
	case raydium_amm.Instruction_Initialize2:
		inst1 := inst.Impl.(*raydium_amm.Initialize2)
		inst1.SetAccounts(accounts)
		// the latest three transfer
		index := len(in.Children)
		t1 := in.Children[index-3].Event[0].(*Transfer)
		t2 := in.Children[index-2].Event[0].(*Transfer)
		t3 := in.Children[index-1].Event[0].(*MintTo)
		createPool := &CreatePool{
			Pool:      inst1.GetAmmAccount().PublicKey,
			TokenA:    inst1.GetPcMintAccount().PublicKey,
			TokenB:    inst1.GetCoinMintAccount().PublicKey,
			TokenLP:   inst1.GetLpMintAccount().PublicKey,
			AccountA:  inst1.GetPoolPcTokenAccountAccount().PublicKey,
			AccountB:  inst1.GetPoolCoinTokenAccountAccount().PublicKey,
			AccountLP: inst1.GetPoolTempLpAccount().PublicKey,
			User:      inst1.GetAmmAuthorityAccount().PublicKey,
		}
		addLiquidity := &AddLiquidity{
			Pool:           inst1.GetAmmAccount().PublicKey,
			User:           inst1.GetAmmAuthorityAccount().PublicKey,
			TokenATransfer: t1,
			TokenBTransfer: t2,
			TokenLpMint:    t3,
		}
		pool := &Pool{
			Hash:     inst1.GetAmmAccount().PublicKey,
			MintA:    inst1.GetCoinMintAccount().PublicKey,
			MintB:    inst1.GetPcMintAccount().PublicKey,
			MintLp:   inst1.GetLpMintAccount().PublicKey,
			VaultA:   inst1.GetPoolCoinTokenAccountAccount().PublicKey,
			VaultB:   inst1.GetPoolPcTokenAccountAccount().PublicKey,
			ReserveA: 0,
			ReserveB: 0,
		}
		return []interface{}{createPool, addLiquidity}, []interface{}{pool}
	case raydium_amm.Instruction_Deposit:
		inst1 := inst.Impl.(*raydium_amm.Deposit)
		inst1.SetAccounts(accounts)
		t1 := in.Children[0].Event[0].(*Transfer)
		t2 := in.Children[1].Event[0].(*Transfer)
		addLiquidity := &AddLiquidity{
			Pool:           inst1.GetAmmAccount().PublicKey,
			User:           inst1.GetUserOwnerAccount().PublicKey,
			TokenATransfer: t1,
			TokenBTransfer: t2,
		}
		panic("not supported")
		return []interface{}{addLiquidity}, []interface{}{}
	case raydium_amm.Instruction_SwapBaseIn:
		inst1 := inst.Impl.(*raydium_amm.SwapBaseIn)
		if len(accounts) == 17 {
			accounts = insertAccount(accounts, 4)
		}
		inst1.SetAccounts(accounts)
		t1 := in.Children[0].Event[0].(*Transfer)
		t2 := in.Children[1].Event[0].(*Transfer)
		swap := &Swap{
			Pool:           inst1.GetAmmAccount().PublicKey,
			TokenATransfer: t1,
			TokenBTransfer: t2,
			User:           inst1.GetUserSourceOwnerAccount().PublicKey,
		}
		return []interface{}{swap}, []interface{}{}
	case raydium_amm.Instruction_SwapBaseOut:
		inst1 := inst.Impl.(*raydium_amm.SwapBaseOut)
		if len(accounts) == 17 {
			accounts = insertAccount(accounts, 4)
		}
		inst1.SetAccounts(accounts)
		t1 := in.Children[0].Event[0].(*Transfer)
		t2 := in.Children[1].Event[0].(*Transfer)
		swap := &Swap{
			Pool:           inst1.GetAmmAccount().PublicKey,
			TokenATransfer: t1,
			TokenBTransfer: t2,
			User:           inst1.GetUserSourceOwnerAccount().PublicKey,
		}
		return []interface{}{swap}, []interface{}{}
	case raydium_amm.Instruction_Withdraw:
		inst1 := inst.Impl.(*raydium_amm.Withdraw)
		inst1.SetAccounts(accounts)
		t1 := in.Children[0].Event[0].(*Transfer)
		t2 := in.Children[1].Event[0].(*Transfer)
		t3 := in.Children[1].Event[0].(*Burn)
		removeLiquidity := &RemoveLiquidity{
			Pool:           inst1.GetAmmAccount().PublicKey,
			TokenATransfer: t1,
			TokenBTransfer: t2,
			TokenLpBurn:    t3,
			User:           inst1.GetUserOwnerAccount().PublicKey,
		}
		panic("not supported")
		return []interface{}{removeLiquidity}, []interface{}{}
	default:
		return nil, nil
	}
}

func RaydiumClmmParser(in *Instruction, meta *Meta) ([]interface{}, []interface{}) {
	inst := new(amm_v4.Instruction)
	err := ag_binary.NewBorshDecoder(in.Instruction.Data).Decode(inst)
	if err != nil {
		return nil, nil
	}
	accounts := make([]*solana.AccountMeta, 0)
	for _, item := range in.Instruction.Accounts {
		accounts = append(accounts, &solana.AccountMeta{
			PublicKey:  item,
			IsWritable: false,
			IsSigner:   false,
		})
	}
	switch inst.TypeID {
	case amm_v4.Instruction_CreatePool:
		inst1 := inst.Impl.(*amm_v4.CreatePool)
		inst1.SetAccounts(accounts)
		pool := &Pool{
			Hash:     inst1.GetPoolStateAccount().PublicKey,
			MintA:    inst1.GetTokenMint0Account().PublicKey,
			MintB:    inst1.GetTokenMint1Account().PublicKey,
			MintLp:   inst1.GetTokenVault1Account().PublicKey,
			VaultA:   inst1.GetTokenVault1Account().PublicKey,
			VaultB:   inst1.GetTokenVault1Account().PublicKey,
			ReserveA: 0,
			ReserveB: 0,
		}
		panic("not supported")
		return nil, []interface{}{pool}
	case amm_v4.Instruction_IncreaseLiquidityV2:
		inst1 := inst.Impl.(*amm_v4.IncreaseLiquidityV2)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event[0].(*Transfer)
		t2 := in.Children[1].Event[0].(*Transfer)
		//
		addLiquidity := &AddLiquidity{
			Pool:           inst1.GetPoolStateAccount().PublicKey,
			User:           inst1.Get(0).PublicKey,
			TokenATransfer: t1,
			TokenBTransfer: t2,
		}
		panic("not supported")
		return []interface{}{addLiquidity}, nil
	case amm_v4.Instruction_DecreaseLiquidityV2:
		inst1 := inst.Impl.(*amm_v4.DecreaseLiquidityV2)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event[0].(*Transfer)
		t2 := in.Children[1].Event[0].(*Transfer)
		//
		removeLiquidity := &RemoveLiquidity{
			Pool:           inst1.GetPoolStateAccount().PublicKey,
			User:           inst1.Get(0).PublicKey,
			TokenATransfer: t1,
			TokenBTransfer: t2,
		}
		panic("not supported")
		return []interface{}{removeLiquidity}, nil
	case amm_v4.Instruction_Swap:
		inst1 := inst.Impl.(*amm_v4.Swap)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event[0].(*Transfer)
		t2 := in.Children[1].Event[0].(*Transfer)
		//
		swap := &Swap{
			Pool:           inst1.GetPoolStateAccount().PublicKey,
			User:           inst1.GetPayerAccount().PublicKey,
			TokenATransfer: t1,
			TokenBTransfer: t2,
		}
		return []interface{}{swap}, nil
	case amm_v4.Instruction_SwapV2:
		inst1 := inst.Impl.(*amm_v4.SwapV2)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event[0].(*Transfer)
		t2 := in.Children[1].Event[0].(*Transfer)
		//
		swap := &Swap{
			Pool:           inst1.GetPoolStateAccount().PublicKey,
			User:           inst1.Get(0).PublicKey,
			TokenATransfer: t1,
			TokenBTransfer: t2,
		}
		return []interface{}{swap}, nil
	default:
		return nil, nil
	}
}

func WhirlPoolParser(in *Instruction, meta *Meta) ([]interface{}, []interface{}) {
	inst := new(whirlpool.Instruction)
	err := ag_binary.NewBorshDecoder(in.Instruction.Data).Decode(inst)
	if err != nil {
		return nil, nil
	}
	accounts := make([]*solana.AccountMeta, 0)
	for _, item := range in.Instruction.Accounts {
		accounts = append(accounts, &solana.AccountMeta{
			PublicKey:  item,
			IsWritable: false,
			IsSigner:   false,
		})
	}
	switch inst.TypeID {
	case whirlpool.Instruction_InitializePool:
		inst1 := inst.Impl.(*whirlpool.InitializePool)
		inst1.SetAccounts(accounts)
		//
		swap := &Swap{
			Pool: inst1.GetWhirlpoolAccount().PublicKey,
			User: inst1.GetFunderAccount().PublicKey,
		}
		panic("not supported")
		return []interface{}{swap}, nil
	case whirlpool.Instruction_Swap:
		inst1 := inst.Impl.(*whirlpool.Swap)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event[0].(*Transfer)
		t2 := in.Children[1].Event[0].(*Transfer)
		//
		swap := &Swap{
			Pool:           inst1.GetWhirlpoolAccount().PublicKey,
			User:           inst1.GetTokenAuthorityAccount().PublicKey,
			TokenATransfer: t1,
			TokenBTransfer: t2,
		}
		panic("not supported")
		return []interface{}{swap}, nil
	case whirlpool.Instruction_SwapV2:
		inst1 := inst.Impl.(*whirlpool.SwapV2)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event[0].(*Transfer)
		t2 := in.Children[1].Event[0].(*Transfer)
		//
		swap := &Swap{
			Pool:           inst1.GetWhirlpoolAccount().PublicKey,
			User:           inst1.GetTokenAuthorityAccount().PublicKey,
			TokenATransfer: t1,
			TokenBTransfer: t2,
		}
		return []interface{}{swap}, nil
	case whirlpool.Instruction_IncreaseLiquidity:
		inst1 := inst.Impl.(*whirlpool.IncreaseLiquidity)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event[0].(*Transfer)
		t2 := in.Children[1].Event[0].(*Transfer)
		//
		addLiquidity := &AddLiquidity{
			Pool:           inst1.GetWhirlpoolAccount().PublicKey,
			User:           inst1.GetPositionAuthorityAccount().PublicKey,
			TokenATransfer: t1,
			TokenBTransfer: t2,
		}
		panic("not supported")
		return []interface{}{addLiquidity}, nil
	case whirlpool.Instruction_DecreaseLiquidity:
		inst1 := inst.Impl.(*whirlpool.DecreaseLiquidity)
		inst1.SetAccounts(accounts)
		//
		t1 := in.Children[0].Event[0].(*Transfer)
		t2 := in.Children[1].Event[0].(*Transfer)
		//
		removeLiquidity := &RemoveLiquidity{
			Pool:           inst1.GetWhirlpoolAccount().PublicKey,
			User:           inst1.GetPositionAuthorityAccount().PublicKey,
			TokenATransfer: t1,
			TokenBTransfer: t2,
		}
		panic("not supported")
		return []interface{}{removeLiquidity}, nil
	default:
		return nil, nil
	}
}

func StabbleStableSwapParser(in *Instruction, meta *Meta) ([]interface{}, []interface{}) {
	inst := new(stable_swap.Instruction)
	err := ag_binary.NewBorshDecoder(in.Instruction.Data).Decode(inst)
	if err != nil {
		return nil, nil
	}
	accounts := make([]*solana.AccountMeta, 0)
	for _, item := range in.Instruction.Accounts {
		accounts = append(accounts, &solana.AccountMeta{
			PublicKey:  item,
			IsWritable: false,
			IsSigner:   false,
		})
	}
	switch inst.TypeID {
	case stable_swap.Instruction_Deposit:
		inst1 := inst.Impl.(*stable_swap.Deposit)
		inst1.SetAccounts(accounts)
		t1 := in.Children[0].Event[0].(*Transfer)
		t2 := in.Children[1].Event[0].(*Transfer)
		addLiquidity := &AddLiquidity{
			Pool:           inst1.GetPoolAccount().PublicKey,
			User:           inst1.GetUserAccount().PublicKey,
			TokenATransfer: t1,
			TokenBTransfer: t2,
		}
		panic("not supported")
		return []interface{}{addLiquidity}, nil
	case stable_swap.Instruction_Withdraw:
		inst1 := inst.Impl.(*stable_swap.Withdraw)
		inst1.SetAccounts(accounts)
		t1 := in.Children[0].Event[0].(*Transfer)
		t2 := in.Children[1].Event[0].(*Transfer)
		removeLiquidity := &RemoveLiquidity{
			Pool:           inst1.GetPoolAccount().PublicKey,
			User:           inst1.GetUserAccount().PublicKey,
			TokenATransfer: t1,
			TokenBTransfer: t2,
		}
		panic("not supported")
		return []interface{}{removeLiquidity}, nil
	case stable_swap.Instruction_Swap:
		inst1 := inst.Impl.(*stable_swap.Swap)
		inst1.SetAccounts(accounts)
		// the first one is user deposit
		// the second is vault withdraw
		t1 := in.Children[0].Event[0].(*Transfer)
		t2 := in.Children[1].Children[1].Event[0].(*Transfer)
		swap := &Swap{
			Pool:           inst1.GetPoolAccount().PublicKey,
			User:           inst1.GetUserAccount().PublicKey,
			TokenATransfer: t1,
			TokenBTransfer: t2,
		}
		return []interface{}{swap}, nil
	case stable_swap.Instruction_SwapV2:
		inst1 := inst.Impl.(*stable_swap.SwapV2)
		inst1.SetAccounts(accounts)
		// the first one is user deposit
		// the latest one is credit
		t1 := in.Children[0].Event[0].(*Transfer)
		t2 := in.Children[len(in.Children)-1].Event[0].(*Transfer)
		//
		swap := &Swap{
			Pool:           inst1.GetPoolAccount().PublicKey,
			User:           inst1.GetUserAccount().PublicKey,
			TokenATransfer: t1,
			TokenBTransfer: t2,
		}
		panic("not supported")
		return []interface{}{swap}, nil
	case stable_swap.Instruction_Initialize:
		inst1 := inst.Impl.(*stable_swap.Initialize)
		inst1.SetAccounts(accounts)
		//
		panic("not supported")
		return nil, nil
	default:
		return nil, nil
	}
}
