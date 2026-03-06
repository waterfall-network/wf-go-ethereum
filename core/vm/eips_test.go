// Copyright 2024 The gwat Authors
// This file is part of the gwat library.
//
// The gwat library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

package vm

import (
	"math"
	"math/big"
	"testing"

	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/rawdb"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/state"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	"gitlab.waterfall.network/waterfall/protocol/gwat/params"
)

// cancunConfig is a chain config with ForkSlotValSyncProc=0 so IsCancun is active from slot 0.
var cancunConfig = &params.ChainConfig{
	ChainID:                big.NewInt(1337),
	SecondsPerSlot:         4,
	SlotsPerEpoch:          32,
	EpochsPerEra:           8,
	TransitionPeriod:       2,
	ValidatorsPerSlot:      6,
	EffectiveBalance:       big.NewInt(3200),
	ValidatorOpExpireSlots: 14400,
	ForkSlotSubNet1:        math.MaxUint64,
	ForkSlotDelegate:       0,
	ForkSlotPrefixFin:      0,
	ForkSlotShanghai:       0,
	ForkSlotValOpTracking:  0,
	ForkSlotReduceBaseFee:  0,
	ForkSlotValSyncProc:    0,
	ForkSlotUpValsPerSlot:  0,
	UpValidatorsPerSlot:    25,
}

// TestEIP7702DelegationResolution verifies EIP-7702 delegation behaviour:
//   - CALL on an EOA with a delegation designator executes the delegate's code
//     in the context of the EOA (storage writes go to the EOA).
//   - Delegation to an account with no code is a no-op.
//   - Regular (non-delegation) code is executed as-is, without following.
func TestEIP7702DelegationResolution(t *testing.T) {
	var (
		eoaAddr      = common.HexToAddress("0x00000000000000000000000000000000000000aa")
		delegateAddr = common.HexToAddress("0x00000000000000000000000000000000000000bb")
		storageSlot  = common.Hash{}
		wantStored   = common.Hash{31: 0x01} // big-endian uint256(1)
	)

	// Delegate bytecode: PUSH1 1, PUSH1 0, SSTORE
	// Stores the value 1 at storage slot 0 of whichever account is being executed.
	delegateCode := []byte{byte(PUSH1), 0x01, byte(PUSH1), 0x00, byte(SSTORE)}

	newTestEVM := func(statedb *state.StateDB) *EVM {
		blockCtx := BlockContext{
			CanTransfer: func(db StateDB, addr common.Address, amount *big.Int) bool {
				return db.GetBalance(addr).Cmp(amount) >= 0
			},
			Transfer: func(db StateDB, from, to common.Address, amount *big.Int) {
				db.SubBalance(from, amount)
				db.AddBalance(to, amount)
			},
			BaseFee: new(big.Int),
		}
		return NewEVM(blockCtx, TxContext{}, statedb, cancunConfig, Config{})
	}

	t.Run("call executes delegate code in EOA context", func(t *testing.T) {
		statedb, err := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
		if err != nil {
			t.Fatalf("failed to create state: %v", err)
		}
		statedb.CreateAccount(delegateAddr)
		statedb.SetCode(delegateAddr, delegateCode)

		statedb.CreateAccount(eoaAddr)
		statedb.SetCode(eoaAddr, types.AddressToDelegation(delegateAddr))

		// EIP-2929: SSTORE requires the executing address to be in the access list.
		statedb.AddAddressToAccessList(eoaAddr)

		evm := newTestEVM(statedb)
		_, _, err = evm.Call(AccountRef(common.Address{}), eoaAddr, nil, 1_000_000, new(big.Int))
		if err != nil {
			t.Fatalf("unexpected call error: %v", err)
		}
		// Delegate code ran in EOA's context: slot 0 of eoaAddr must be 1.
		if got := statedb.GetState(eoaAddr, storageSlot); got != wantStored {
			t.Errorf("eoaAddr storage: got %s, want %s", got, wantStored)
		}
		// Delegate's own storage must be untouched.
		if got := statedb.GetState(delegateAddr, storageSlot); got != (common.Hash{}) {
			t.Errorf("delegateAddr storage should be empty, got %s", got)
		}
	})

	t.Run("delegation to empty account — no-op", func(t *testing.T) {
		statedb, err := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
		if err != nil {
			t.Fatalf("failed to create state: %v", err)
		}
		// eoaAddr delegates to an account that has no code.
		statedb.CreateAccount(eoaAddr)
		statedb.SetCode(eoaAddr, types.AddressToDelegation(delegateAddr))

		evm := newTestEVM(statedb)
		ret, _, err := evm.Call(AccountRef(common.Address{}), eoaAddr, nil, 1_000_000, new(big.Int))
		if err != nil {
			t.Fatalf("unexpected call error: %v", err)
		}
		if len(ret) != 0 {
			t.Errorf("expected empty return, got %x", ret)
		}
		if got := statedb.GetState(eoaAddr, storageSlot); got != (common.Hash{}) {
			t.Errorf("eoaAddr storage should be empty, got %s", got)
		}
	})

	t.Run("regular code not treated as delegation", func(t *testing.T) {
		statedb, err := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
		if err != nil {
			t.Fatalf("failed to create state: %v", err)
		}
		// eoaAddr has ordinary code (not a delegation designator).
		statedb.CreateAccount(eoaAddr)
		statedb.SetCode(eoaAddr, delegateCode)

		// EIP-2929: SSTORE requires the executing address to be in the access list.
		statedb.AddAddressToAccessList(eoaAddr)

		evm := newTestEVM(statedb)
		_, _, err = evm.Call(AccountRef(common.Address{}), eoaAddr, nil, 1_000_000, new(big.Int))
		if err != nil {
			t.Fatalf("unexpected call error: %v", err)
		}
		// Own code ran directly: slot 0 of eoaAddr must be 1.
		if got := statedb.GetState(eoaAddr, storageSlot); got != wantStored {
			t.Errorf("eoaAddr storage: got %s, want %s", got, wantStored)
		}
	})
}

// TestEIP6780Selfdestruct verifies EIP-6780 behaviour:
//   - For an existing contract, SELFDESTRUCT only transfers balance; code is preserved.
//   - For a contract created in the same transaction, SELFDESTRUCT fully destroys it.
func TestEIP6780Selfdestruct(t *testing.T) {
	var (
		contractAddr   = common.HexToAddress("0x00000000000000000000000000000000000000aa")
		beneficiary    = common.HexToAddress("0x00000000000000000000000000000000000000bb")
		initialBalance = big.NewInt(100)
	)

	// Bytecode: PUSH20 <beneficiary> SELFDESTRUCT
	selfdestructCode := append([]byte{byte(PUSH20)}, beneficiary.Bytes()...)
	selfdestructCode = append(selfdestructCode, byte(SELFDESTRUCT))

	newTestEVM := func(statedb *state.StateDB) *EVM {
		blockCtx := BlockContext{
			CanTransfer: func(db StateDB, addr common.Address, amount *big.Int) bool {
				return db.GetBalance(addr).Cmp(amount) >= 0
			},
			Transfer: func(db StateDB, from, to common.Address, amount *big.Int) {
				db.SubBalance(from, amount)
				db.AddBalance(to, amount)
			},
			BaseFee: new(big.Int),
		}
		return NewEVM(blockCtx, TxContext{}, statedb, cancunConfig, Config{})
	}

	t.Run("existing contract — balance transferred, code preserved", func(t *testing.T) {
		statedb, err := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
		if err != nil {
			t.Fatalf("failed to create state: %v", err)
		}
		statedb.CreateAccount(contractAddr)
		statedb.SetCode(contractAddr, selfdestructCode)
		statedb.SetBalance(contractAddr, initialBalance)
		// Do NOT call CreateContract — newContract remains false.

		evm := newTestEVM(statedb)
		evm.Call(AccountRef(common.Address{}), contractAddr, nil, 1_000_000, new(big.Int))

		if statedb.HasSuicided(contractAddr) {
			t.Error("existing contract must not be marked as suicided")
		}
		if got := statedb.GetCode(contractAddr); len(got) == 0 {
			t.Error("existing contract code must be preserved after SELFDESTRUCT")
		}
		if got := statedb.GetBalance(beneficiary); got.Cmp(initialBalance) != 0 {
			t.Errorf("beneficiary balance: got %s, want %s", got, initialBalance)
		}
	})

	t.Run("new contract (same tx) — fully destroyed", func(t *testing.T) {
		statedb, err := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
		if err != nil {
			t.Fatalf("failed to create state: %v", err)
		}
		statedb.CreateAccount(contractAddr)
		statedb.SetCode(contractAddr, selfdestructCode)
		statedb.SetBalance(contractAddr, initialBalance)
		statedb.CreateContract(contractAddr) // marks newContract = true

		evm := newTestEVM(statedb)
		evm.Call(AccountRef(common.Address{}), contractAddr, nil, 1_000_000, new(big.Int))

		if !statedb.HasSuicided(contractAddr) {
			t.Error("contract created in same tx must be marked as suicided after SELFDESTRUCT")
		}
		if got := statedb.GetBalance(beneficiary); got.Cmp(initialBalance) != 0 {
			t.Errorf("beneficiary balance: got %s, want %s", got, initialBalance)
		}
	})

	t.Run("existing contract selfdestructs to self — no-op", func(t *testing.T) {
		statedb, err := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
		if err != nil {
			t.Fatalf("failed to create state: %v", err)
		}
		// Code: PUSH20 <self> SELFDESTRUCT
		selfCode := append([]byte{byte(PUSH20)}, contractAddr.Bytes()...)
		selfCode = append(selfCode, byte(SELFDESTRUCT))

		statedb.CreateAccount(contractAddr)
		statedb.SetCode(contractAddr, selfCode)
		statedb.SetBalance(contractAddr, initialBalance)

		evm := newTestEVM(statedb)
		evm.Call(AccountRef(common.Address{}), contractAddr, nil, 1_000_000, new(big.Int))

		if statedb.HasSuicided(contractAddr) {
			t.Error("existing contract selfdestructing to self must not be marked as suicided")
		}
		// Balance must remain unchanged (self-to-self transfer skipped).
		if got := statedb.GetBalance(contractAddr); got.Cmp(initialBalance) != 0 {
			t.Errorf("contract balance: got %s, want %s", got, initialBalance)
		}
	})
}
