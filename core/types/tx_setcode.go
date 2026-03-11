// Copyright 2024 The gwat Authors
// This file is part of the gwat library.
//
// The gwat library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The gwat library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the gwat library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"math/big"

	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/common/hexutil"
	"gitlab.waterfall.network/waterfall/protocol/gwat/crypto"
)

// DelegationPrefix is used by code to denote the account is delegating to
// another account.
var DelegationPrefix = []byte{0xef, 0x01, 0x00}

// ParseDelegation tries to parse the address from a delegation slice.
func ParseDelegation(b []byte) (common.Address, bool) {
	if len(b) != 23 || !bytes.HasPrefix(b, DelegationPrefix) {
		return common.Address{}, false
	}
	return common.BytesToAddress(b[len(DelegationPrefix):]), true
}

// AddressToDelegation adds the delegation prefix to the specified address.
func AddressToDelegation(addr common.Address) []byte {
	return append(DelegationPrefix, addr.Bytes()...)
}

// SetCodeAuthorization is an authorization from an account to deploy code at its address.
type SetCodeAuthorization struct {
	ChainID *big.Int
	Address common.Address
	Nonce   uint64
	V       uint8
	R       *big.Int
	S       *big.Int
}

// authorizationJSON is the JSON representation of SetCodeAuthorization.
type authorizationJSON struct {
	ChainID *hexutil.Big   `json:"chainId"`
	Address common.Address `json:"address"`
	Nonce   hexutil.Uint64 `json:"nonce"`
	V       hexutil.Uint64 `json:"yParity"`
	R       *hexutil.Big   `json:"r"`
	S       *hexutil.Big   `json:"s"`
}

// MarshalJSON marshals SetCodeAuthorization as JSON.
func (a SetCodeAuthorization) MarshalJSON() ([]byte, error) {
	enc := authorizationJSON{
		ChainID: (*hexutil.Big)(a.ChainID),
		Address: a.Address,
		Nonce:   hexutil.Uint64(a.Nonce),
		V:       hexutil.Uint64(a.V),
		R:       (*hexutil.Big)(a.R),
		S:       (*hexutil.Big)(a.S),
	}
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals SetCodeAuthorization from JSON.
func (a *SetCodeAuthorization) UnmarshalJSON(input []byte) error {
	var dec authorizationJSON
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.ChainID == nil {
		return errors.New("missing required field 'chainId' in authorization")
	}
	a.ChainID = dec.ChainID.ToInt()
	a.Address = dec.Address
	a.Nonce = uint64(dec.Nonce)
	a.V = uint8(dec.V)
	if dec.R == nil {
		return errors.New("missing required field 'r' in authorization")
	}
	a.R = dec.R.ToInt()
	if dec.S == nil {
		return errors.New("missing required field 's' in authorization")
	}
	a.S = dec.S.ToInt()
	return nil
}

// SignSetCode creates a signed SetCode authorization.
func SignSetCode(prv *ecdsa.PrivateKey, auth SetCodeAuthorization) (SetCodeAuthorization, error) {
	sighash := auth.SigHash()
	sig, err := crypto.Sign(sighash[:], prv)
	if err != nil {
		return SetCodeAuthorization{}, err
	}
	R, S, _ := decodeSignature(sig)
	return SetCodeAuthorization{
		ChainID: new(big.Int).Set(auth.ChainID),
		Address: auth.Address,
		Nonce:   auth.Nonce,
		V:       sig[64],
		R:       R,
		S:       S,
	}, nil
}

// SigHash returns the hash of SetCodeAuthorization for signing.
func (a *SetCodeAuthorization) SigHash() common.Hash {
	return prefixedRlpHash(0x05, []interface{}{
		a.ChainID,
		a.Address,
		a.Nonce,
	})
}

// Authority recovers the authorizing account of an authorization.
func (a *SetCodeAuthorization) Authority() (common.Address, error) {
	sighash := a.SigHash()
	if !crypto.ValidateSignatureValues(a.V, a.R, a.S, true) {
		return common.Address{}, ErrInvalidSig
	}
	r, s := a.R.Bytes(), a.S.Bytes()
	var sig [crypto.SignatureLength]byte
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = a.V
	pub, err := crypto.Ecrecover(sighash[:], sig[:])
	if err != nil {
		return common.Address{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return common.Address{}, errors.New("invalid public key")
	}
	var addr common.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return addr, nil
}

// SetCodeTx implements the EIP-7702 transaction type which temporarily installs
// the code at the signer's address.
type SetCodeTx struct {
	ChainID    *big.Int
	Nonce      uint64
	GasTipCap  *big.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *big.Int // a.k.a. maxFeePerGas
	Gas        uint64
	To         common.Address // EIP-7702 requires a destination; not a contract-creation tx
	Value      *big.Int
	Data       []byte
	AccessList AccessList
	AuthList   []SetCodeAuthorization

	// Signature values
	V *big.Int
	R *big.Int
	S *big.Int
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *SetCodeTx) copy() TxData {
	cpy := &SetCodeTx{
		Nonce:      tx.Nonce,
		To:         tx.To,
		Data:       common.CopyBytes(tx.Data),
		Gas:        tx.Gas,
		AccessList: make(AccessList, len(tx.AccessList)),
		AuthList:   make([]SetCodeAuthorization, len(tx.AuthList)),
		Value:      new(big.Int),
		ChainID:    new(big.Int),
		GasTipCap:  new(big.Int),
		GasFeeCap:  new(big.Int),
		V:          new(big.Int),
		R:          new(big.Int),
		S:          new(big.Int),
	}
	copy(cpy.AccessList, tx.AccessList)
	for i, auth := range tx.AuthList {
		a := SetCodeAuthorization{
			Address: auth.Address,
			Nonce:   auth.Nonce,
			V:       auth.V,
			ChainID: new(big.Int),
			R:       new(big.Int),
			S:       new(big.Int),
		}
		if auth.ChainID != nil {
			a.ChainID.Set(auth.ChainID)
		}
		if auth.R != nil {
			a.R.Set(auth.R)
		}
		if auth.S != nil {
			a.S.Set(auth.S)
		}
		cpy.AuthList[i] = a
	}
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	if tx.ChainID != nil {
		cpy.ChainID.Set(tx.ChainID)
	}
	if tx.GasTipCap != nil {
		cpy.GasTipCap.Set(tx.GasTipCap)
	}
	if tx.GasFeeCap != nil {
		cpy.GasFeeCap.Set(tx.GasFeeCap)
	}
	if tx.V != nil {
		cpy.V.Set(tx.V)
	}
	if tx.R != nil {
		cpy.R.Set(tx.R)
	}
	if tx.S != nil {
		cpy.S.Set(tx.S)
	}
	return cpy
}

// accessors for innerTx.
func (tx *SetCodeTx) txType() byte           { return SetCodeTxType }
func (tx *SetCodeTx) chainID() *big.Int      { return tx.ChainID }
func (tx *SetCodeTx) accessList() AccessList { return tx.AccessList }
func (tx *SetCodeTx) data() []byte           { return tx.Data }
func (tx *SetCodeTx) gas() uint64            { return tx.Gas }
func (tx *SetCodeTx) gasFeeCap() *big.Int    { return tx.GasFeeCap }
func (tx *SetCodeTx) gasTipCap() *big.Int    { return tx.GasTipCap }
func (tx *SetCodeTx) gasPrice() *big.Int     { return tx.GasFeeCap }
func (tx *SetCodeTx) value() *big.Int        { return tx.Value }
func (tx *SetCodeTx) nonce() uint64          { return tx.Nonce }
func (tx *SetCodeTx) to() *common.Address    { tmp := tx.To; return &tmp }

func (tx *SetCodeTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

func (tx *SetCodeTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID = chainID
	tx.V, tx.R, tx.S = v, r, s
}
