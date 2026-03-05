// Copyright 2024 The gwat Authors
// This file is part of the gwat library.
//
// The gwat library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

package types

import (
	"encoding/json"
	"math/big"
	"testing"

	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/crypto"
)

// TestSignSetCodeAuthorityRoundTrip verifies that signing a SetCodeAuthorization
// and recovering the authority produces the original signer's address.
func TestSignSetCodeAuthorityRoundTrip(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	wantAddr := crypto.PubkeyToAddress(key.PublicKey)

	auth := SetCodeAuthorization{
		ChainID: big.NewInt(1337),
		Address: common.HexToAddress("0x00000000000000000000000000000000000000bb"),
		Nonce:   7,
	}

	signed, err := SignSetCode(key, auth)
	if err != nil {
		t.Fatalf("SignSetCode: %v", err)
	}

	got, err := signed.Authority()
	if err != nil {
		t.Fatalf("Authority: %v", err)
	}
	if got != wantAddr {
		t.Errorf("Authority: got %s, want %s", got, wantAddr)
	}
}

// TestSignSetCodeChainIDWildcard verifies that a zero chain ID authorisation
// can still be signed and recovered correctly.
func TestSignSetCodeChainIDWildcard(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	wantAddr := crypto.PubkeyToAddress(key.PublicKey)

	auth := SetCodeAuthorization{
		ChainID: big.NewInt(0), // wildcard
		Address: common.HexToAddress("0x00000000000000000000000000000000000000cc"),
		Nonce:   0,
	}

	signed, err := SignSetCode(key, auth)
	if err != nil {
		t.Fatalf("SignSetCode: %v", err)
	}

	got, err := signed.Authority()
	if err != nil {
		t.Fatalf("Authority: %v", err)
	}
	if got != wantAddr {
		t.Errorf("Authority: got %s, want %s", got, wantAddr)
	}
}

// TestSetCodeAuthorizationJSONRoundTrip verifies JSON marshal/unmarshal
// of SetCodeAuthorization preserves all fields.
func TestSetCodeAuthorizationJSONRoundTrip(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	original := SetCodeAuthorization{
		ChainID: big.NewInt(1337),
		Address: common.HexToAddress("0x00000000000000000000000000000000000000dd"),
		Nonce:   42,
	}
	signed, err := SignSetCode(key, original)
	if err != nil {
		t.Fatalf("SignSetCode: %v", err)
	}

	data, err := signed.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}

	var got SetCodeAuthorization
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("UnmarshalJSON: %v", err)
	}

	if got.ChainID.Cmp(signed.ChainID) != 0 {
		t.Errorf("ChainID: got %s, want %s", got.ChainID, signed.ChainID)
	}
	if got.Address != signed.Address {
		t.Errorf("Address: got %s, want %s", got.Address, signed.Address)
	}
	if got.Nonce != signed.Nonce {
		t.Errorf("Nonce: got %d, want %d", got.Nonce, signed.Nonce)
	}
	if got.V != signed.V {
		t.Errorf("V: got %d, want %d", got.V, signed.V)
	}
	if got.R.Cmp(signed.R) != 0 {
		t.Errorf("R: got %s, want %s", got.R, signed.R)
	}
	if got.S.Cmp(signed.S) != 0 {
		t.Errorf("S: got %s, want %s", got.S, signed.S)
	}

	// The recovered authority must still be the original signer.
	wantAddr := crypto.PubkeyToAddress(key.PublicKey)
	gotAddr, err := got.Authority()
	if err != nil {
		t.Fatalf("Authority after unmarshal: %v", err)
	}
	if gotAddr != wantAddr {
		t.Errorf("Authority after unmarshal: got %s, want %s", gotAddr, wantAddr)
	}
}
