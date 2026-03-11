// Copyright 2024 The gwat Authors
// This file is part of the gwat library.
//
// The gwat library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

package core

import (
	"errors"
	"math/big"
	"testing"

	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/rawdb"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/state"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/vm"
	"gitlab.waterfall.network/waterfall/protocol/gwat/crypto"
	"gitlab.waterfall.network/waterfall/protocol/gwat/params"
)

// newEIP7702StateTransition builds a minimal StateTransition for testing
// applyAuthorization / validateAuthorization.  It uses eip3860Config
// (chain ID 1337, Prague active from slot 0).
func newEIP7702StateTransition(t *testing.T, statedb *state.StateDB) *StateTransition {
	t.Helper()
	evm := vm.NewEVM(vm.BlockContext{
		CanTransfer: func(db vm.StateDB, addr common.Address, amount *big.Int) bool { return true },
		Transfer:    func(db vm.StateDB, a, b common.Address, v *big.Int) {},
		BaseFee:     new(big.Int),
	}, vm.TxContext{}, statedb, eip3860Config, vm.Config{})
	msg := NewMockMessage(common.Address{}, &common.Address{}, new(big.Int), 0, nil)
	return NewStateTransition(evm, nil, nil, msg, new(GasPool).AddGas(0))
}

// TestEIP7702DelegatedEOACanSend verifies that an EOA with EIP-7702 delegation
// code is still allowed to send regular transactions (preCheck must not reject
// it with ErrSenderNoEOA).
func TestEIP7702DelegatedEOACanSend(t *testing.T) {
	delegateAddr := common.HexToAddress("0x00000000000000000000000000000000000000bb")

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	senderAddr := crypto.PubkeyToAddress(key.PublicKey)

	db, err := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	if err != nil {
		t.Fatalf("create state: %v", err)
	}
	db.CreateAccount(senderAddr)
	// Set EIP-7702 delegation code on the sender — simulates a previously
	// applied SetCode transaction.
	db.SetCode(senderAddr, types.AddressToDelegation(delegateAddr))

	toAddr := common.HexToAddress("0x1234")
	const gasLimit = 21000
	msg := types.NewMessage(
		senderAddr, &toAddr,
		0,            // nonce matches state nonce (0)
		new(big.Int), // value
		gasLimit,
		new(big.Int), // gasPrice
		new(big.Int), // gasFeeCap
		new(big.Int), // gasTipCap
		nil,          // data
		nil,          // accessList
		false,        // isFake = false → triggers preCheck validation
	)

	evm := vm.NewEVM(vm.BlockContext{
		CanTransfer: func(db vm.StateDB, addr common.Address, amount *big.Int) bool { return true },
		Transfer:    func(db vm.StateDB, a, b common.Address, v *big.Int) {},
		BaseFee:     new(big.Int),
	}, vm.TxContext{}, db, eip3860Config, vm.Config{NoBaseFee: true})

	st := NewStateTransition(evm, nil, nil, msg, new(GasPool).AddGas(gasLimit))
	if err := st.preCheck(); err != nil {
		// This is the bug: delegated EOA incorrectly rejected as "sender not an eoa"
		t.Errorf("preCheck rejected delegated EOA: %v", err)
	}
}

// TestApplyAuthorization covers the EIP-7702 applyAuthorization logic.
func TestApplyAuthorization(t *testing.T) {
	delegateAddr := common.HexToAddress("0x00000000000000000000000000000000000000bb")

	newState := func(t *testing.T) *state.StateDB {
		t.Helper()
		db, err := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
		if err != nil {
			t.Fatalf("create state: %v", err)
		}
		return db
	}

	t.Run("sets delegation code on authority", func(t *testing.T) {
		key, err := crypto.GenerateKey()
		if err != nil {
			t.Fatalf("generate key: %v", err)
		}
		authorityAddr := crypto.PubkeyToAddress(key.PublicKey)

		db := newState(t)
		db.CreateAccount(authorityAddr)
		// nonce is 0 by default

		auth, err := types.SignSetCode(key, types.SetCodeAuthorization{
			ChainID: big.NewInt(1337),
			Address: delegateAddr,
			Nonce:   0,
		})
		if err != nil {
			t.Fatalf("SignSetCode: %v", err)
		}

		st := newEIP7702StateTransition(t, db)
		if err := st.applyAuthorization(&auth); err != nil {
			t.Fatalf("applyAuthorization: %v", err)
		}

		got := db.GetCode(authorityAddr)
		wantTarget, ok := types.ParseDelegation(got)
		if !ok {
			t.Fatalf("expected delegation code, got %x", got)
		}
		if wantTarget != delegateAddr {
			t.Errorf("delegate target: got %s, want %s", wantTarget, delegateAddr)
		}
		if have := db.GetNonce(authorityAddr); have != 1 {
			t.Errorf("nonce after apply: got %d, want 1", have)
		}
	})

	t.Run("refunds new-account gas when authority already exists", func(t *testing.T) {
		key, err := crypto.GenerateKey()
		if err != nil {
			t.Fatalf("generate key: %v", err)
		}
		authorityAddr := crypto.PubkeyToAddress(key.PublicKey)

		db := newState(t)
		db.CreateAccount(authorityAddr) // account exists

		auth, err := types.SignSetCode(key, types.SetCodeAuthorization{
			ChainID: big.NewInt(1337),
			Address: delegateAddr,
			Nonce:   0,
		})
		if err != nil {
			t.Fatalf("SignSetCode: %v", err)
		}

		st := newEIP7702StateTransition(t, db)
		if err := st.applyAuthorization(&auth); err != nil {
			t.Fatalf("applyAuthorization: %v", err)
		}

		wantRefund := params.CallNewAccountGas - params.TxAuthTupleGas
		if got := db.GetRefund(); got != wantRefund {
			t.Errorf("refund: got %d, want %d", got, wantRefund)
		}
	})

	t.Run("no refund when authority is a new account", func(t *testing.T) {
		key, err := crypto.GenerateKey()
		if err != nil {
			t.Fatalf("generate key: %v", err)
		}
		// Do NOT create the account — it does not exist in state.

		auth, err := types.SignSetCode(key, types.SetCodeAuthorization{
			ChainID: big.NewInt(1337),
			Address: delegateAddr,
			Nonce:   0,
		})
		if err != nil {
			t.Fatalf("SignSetCode: %v", err)
		}

		db := newState(t)
		st := newEIP7702StateTransition(t, db)
		if err := st.applyAuthorization(&auth); err != nil {
			t.Fatalf("applyAuthorization: %v", err)
		}

		if got := db.GetRefund(); got != 0 {
			t.Errorf("expected no refund for new account, got %d", got)
		}
	})

	t.Run("zero address clears delegation code", func(t *testing.T) {
		key, err := crypto.GenerateKey()
		if err != nil {
			t.Fatalf("generate key: %v", err)
		}
		authorityAddr := crypto.PubkeyToAddress(key.PublicKey)

		db := newState(t)
		db.CreateAccount(authorityAddr)
		// pre-set a delegation so there is something to clear
		db.SetCode(authorityAddr, types.AddressToDelegation(delegateAddr))

		auth, err := types.SignSetCode(key, types.SetCodeAuthorization{
			ChainID: big.NewInt(1337),
			Address: common.Address{}, // zero → clear
			Nonce:   0,
		})
		if err != nil {
			t.Fatalf("SignSetCode: %v", err)
		}

		st := newEIP7702StateTransition(t, db)
		if err := st.applyAuthorization(&auth); err != nil {
			t.Fatalf("applyAuthorization: %v", err)
		}

		if got := db.GetCode(authorityAddr); len(got) != 0 {
			t.Errorf("expected code cleared, got %x", got)
		}
	})

	t.Run("wrong chain ID is rejected", func(t *testing.T) {
		key, err := crypto.GenerateKey()
		if err != nil {
			t.Fatalf("generate key: %v", err)
		}

		auth, err := types.SignSetCode(key, types.SetCodeAuthorization{
			ChainID: big.NewInt(9999), // does not match eip3860Config (1337)
			Address: delegateAddr,
			Nonce:   0,
		})
		if err != nil {
			t.Fatalf("SignSetCode: %v", err)
		}

		db := newState(t)
		st := newEIP7702StateTransition(t, db)
		if err := st.applyAuthorization(&auth); !errors.Is(err, ErrAuthorizationWrongChainID) {
			t.Errorf("expected ErrAuthorizationWrongChainID, got %v", err)
		}
	})

	t.Run("nonce mismatch is rejected", func(t *testing.T) {
		key, err := crypto.GenerateKey()
		if err != nil {
			t.Fatalf("generate key: %v", err)
		}
		authorityAddr := crypto.PubkeyToAddress(key.PublicKey)

		db := newState(t)
		db.CreateAccount(authorityAddr)
		db.SetNonce(authorityAddr, 5) // state nonce = 5

		auth, err := types.SignSetCode(key, types.SetCodeAuthorization{
			ChainID: big.NewInt(1337),
			Address: delegateAddr,
			Nonce:   3, // mismatch: auth nonce 3 ≠ state nonce 5
		})
		if err != nil {
			t.Fatalf("SignSetCode: %v", err)
		}

		st := newEIP7702StateTransition(t, db)
		if err := st.applyAuthorization(&auth); !errors.Is(err, ErrAuthorizationNonceMismatch) {
			t.Errorf("expected ErrAuthorizationNonceMismatch, got %v", err)
		}
	})

	t.Run("authority with non-delegation code is rejected", func(t *testing.T) {
		key, err := crypto.GenerateKey()
		if err != nil {
			t.Fatalf("generate key: %v", err)
		}
		authorityAddr := crypto.PubkeyToAddress(key.PublicKey)

		db := newState(t)
		db.CreateAccount(authorityAddr)
		db.SetCode(authorityAddr, []byte{0x60, 0x00}) // ordinary bytecode, not a delegation

		auth, err := types.SignSetCode(key, types.SetCodeAuthorization{
			ChainID: big.NewInt(1337),
			Address: delegateAddr,
			Nonce:   0,
		})
		if err != nil {
			t.Fatalf("SignSetCode: %v", err)
		}

		st := newEIP7702StateTransition(t, db)
		if err := st.applyAuthorization(&auth); !errors.Is(err, ErrAuthorizationDestinationHasCode) {
			t.Errorf("expected ErrAuthorizationDestinationHasCode, got %v", err)
		}
	})
}
