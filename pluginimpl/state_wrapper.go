// Copyright 2026 Digital Clever Solution LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pluginimpl

import (
	"math/big"

	"github.com/LFDT-Iguazu/iguazu-types/blockdag/iface"
	gwatcommon "gitlab.waterfall.network/waterfall/protocol/gwat/common"
	gwatstate "gitlab.waterfall.network/waterfall/protocol/gwat/core/state"
	gwatcoretypes "gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
)

// stateDBWrapper wraps *gwatstate.StateDB to implement iface.StateDB.
// All addr/hash parameters are [20]byte / [32]byte aliases in both namespaces,
// so conversions are simple named-type casts.
type stateDBWrapper struct{ inner *gwatstate.StateDB }

func wrapStateDB(s *gwatstate.StateDB, err error) (iface.StateDB, error) {
	if err != nil || s == nil {
		return nil, err
	}
	return &stateDBWrapper{inner: s}, nil
}

func (w *stateDBWrapper) GetBalance(addr [20]byte) *big.Int {
	return w.inner.GetBalance(gwatcommon.Address(addr))
}

func (w *stateDBWrapper) SetBalance(addr [20]byte, amount *big.Int) {
	w.inner.SetBalance(gwatcommon.Address(addr), amount)
}

func (w *stateDBWrapper) AddBalance(addr [20]byte, amount *big.Int) {
	w.inner.AddBalance(gwatcommon.Address(addr), amount)
}

func (w *stateDBWrapper) SubBalance(addr [20]byte, amount *big.Int) {
	w.inner.SubBalance(gwatcommon.Address(addr), amount)
}

func (w *stateDBWrapper) IsValidatorAddress(addr [20]byte) bool {
	return w.inner.IsValidatorAddress(gwatcommon.Address(addr))
}

func (w *stateDBWrapper) GetState(addr [20]byte, key [32]byte) [32]byte {
	return [32]byte(w.inner.GetState(gwatcommon.Address(addr), gwatcommon.Hash(key)))
}

func (w *stateDBWrapper) SetState(addr [20]byte, key [32]byte, value [32]byte) {
	w.inner.SetState(gwatcommon.Address(addr), gwatcommon.Hash(key), gwatcommon.Hash(value))
}

func (w *stateDBWrapper) GetCode(addr [20]byte) []byte {
	return w.inner.GetCode(gwatcommon.Address(addr))
}

func (w *stateDBWrapper) SetCode(addr [20]byte, code []byte) {
	w.inner.SetCode(gwatcommon.Address(addr), code)
}

func (w *stateDBWrapper) GetNonce(addr [20]byte) uint64 {
	return w.inner.GetNonce(gwatcommon.Address(addr))
}

func (w *stateDBWrapper) SetNonce(addr [20]byte, nonce uint64) {
	w.inner.SetNonce(gwatcommon.Address(addr), nonce)
}

func (w *stateDBWrapper) Exist(addr [20]byte) bool {
	return w.inner.Exist(gwatcommon.Address(addr))
}

func (w *stateDBWrapper) CreateAccount(addr [20]byte) {
	w.inner.CreateAccount(gwatcommon.Address(addr))
}

func (w *stateDBWrapper) Suicide(addr [20]byte) bool {
	return w.inner.Suicide(gwatcommon.Address(addr))
}

func (w *stateDBWrapper) AddLog(entry *iface.LogEntry) {
	topics := make([]gwatcommon.Hash, len(entry.Topics))
	for i, t := range entry.Topics {
		topics[i] = gwatcommon.Hash(t)
	}
	w.inner.AddLog(&gwatcoretypes.Log{
		Address: gwatcommon.Address(entry.Address),
		Topics:  topics,
		Data:    entry.Data,
	})
}

func (w *stateDBWrapper) Snapshot() int { return w.inner.Snapshot() }

func (w *stateDBWrapper) RevertToSnapshot(id int) { w.inner.RevertToSnapshot(id) }

func (w *stateDBWrapper) ForEachStorage(addr [20]byte, cb func(key, value [32]byte) bool) error {
	return w.inner.ForEachStorage(gwatcommon.Address(addr), func(k, v gwatcommon.Hash) bool {
		return cb([32]byte(k), [32]byte(v))
	})
}
