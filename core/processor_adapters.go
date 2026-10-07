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

package core

import (
	"math/big"

	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/state"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/vm"
	"gitlab.waterfall.network/waterfall/protocol/gwat/token"
	tokenOp "gitlab.waterfall.network/waterfall/protocol/gwat/token/operation"
	"gitlab.waterfall.network/waterfall/protocol/gwat/validator"
	"gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/iface"
	wfcommon "gitlab.waterfall.network/waterfall/protocol/wf-types/common"
)

// ── Token processor adapter ───────────────────────────────────────────────────

// tokenProcessorAdapter wraps gwat's *token.Processor so it satisfies
// iface.TokenProcessor. It decodes the raw op bytes before delegating.
type tokenProcessorAdapter struct {
	inner *token.Processor
}

// NewTokenProcessorAdapter returns an iface.TokenProcessor backed by gwat's
// concrete token processor.
func NewTokenProcessorAdapter(p *token.Processor) iface.TokenProcessor {
	return &tokenProcessorAdapter{inner: p}
}

func (a *tokenProcessorAdapter) IsToken(addr [20]byte) bool {
	return a.inner.IsToken(common.Address(addr))
}

func (a *tokenProcessorAdapter) Call(caller iface.Ref, tokenAddr [20]byte, value *big.Int, op []byte) ([]byte, error) {
	decoded, err := tokenOp.DecodeBytes(op)
	if err != nil {
		return nil, err
	}
	return a.inner.Call(tokenCallerRef{caller}, common.Address(tokenAddr), value, decoded)
}

// tokenCallerRef adapts iface.Ref to gwat's token.Ref (Address() common.Address).
type tokenCallerRef struct{ r iface.Ref }

func (t tokenCallerRef) Address() common.Address { return common.Address(t.r.Address()) }

// ── Validator processor adapter ───────────────────────────────────────────────

// validatorProcessorAdapter wraps gwat's *validator.Processor so it satisfies
// iface.ValidatorProcessor.
type validatorProcessorAdapter struct {
	inner *validator.Processor
}

// NewValidatorProcessorAdapter returns an iface.ValidatorProcessor backed by
// gwat's concrete validator processor.
func NewValidatorProcessorAdapter(p *validator.Processor) iface.ValidatorProcessor {
	return &validatorProcessorAdapter{inner: p}
}

func (a *validatorProcessorAdapter) IsValidatorOp(to *[20]byte) bool {
	if to == nil {
		return false
	}
	addr := common.Address(*to)
	return a.inner.IsValidatorOp(&addr)
}

func (a *validatorProcessorAdapter) Call(caller iface.Ref, to [20]byte, value *big.Int, data []byte, txHash [32]byte) ([]byte, error) {
	return a.inner.Call(
		validatorCallerRef{caller},
		common.Address(to),
		value,
		&validatorMsg{data: data, txHash: common.Hash(txHash)},
	)
}

func (a *validatorProcessorAdapter) GetValidatorsStateAddress() [20]byte {
	return [20]byte(a.inner.GetValidatorsStateAddress())
}

func (a *validatorProcessorAdapter) ProcessRewards(coinbase [20]byte, reward *big.Int) error {
	return a.inner.ProcessRewards(coinbase, reward)
}

// validatorCallerRef adapts iface.Ref to gwat's validator.Ref (Address() common.Address).
type validatorCallerRef struct{ r iface.Ref }

func (v validatorCallerRef) Address() common.Address { return common.Address(v.r.Address()) }

// validatorMsg satisfies the private message interface in gwat/validator.
type validatorMsg struct {
	data   []byte
	txHash common.Hash
}

func (m *validatorMsg) Data() []byte        { return m.data }
func (m *validatorMsg) TxHash() common.Hash { return m.txHash }

// ── StateDB adapter (used by processor factories) ────────────────────────────

// stateDBIfaceAdapter wraps gwat's *state.StateDB to implement iface.StateDB.
// It is used when an external processor factory (e.g. from wf-engine) needs to
// receive an iface.StateDB instead of gwat's concrete type.
type stateDBIfaceAdapter struct {
	inner *state.StateDB
}

func newStateDBIfaceAdapter(s *state.StateDB) iface.StateDB {
	return &stateDBIfaceAdapter{inner: s}
}

func (a *stateDBIfaceAdapter) GetBalance(addr [20]byte) *big.Int {
	return a.inner.GetBalance(common.Address(addr))
}
func (a *stateDBIfaceAdapter) SetBalance(addr [20]byte, amount *big.Int) {
	a.inner.SetBalance(common.Address(addr), amount)
}
func (a *stateDBIfaceAdapter) AddBalance(addr [20]byte, amount *big.Int) {
	a.inner.AddBalance(common.Address(addr), amount)
}
func (a *stateDBIfaceAdapter) SubBalance(addr [20]byte, amount *big.Int) {
	a.inner.SubBalance(common.Address(addr), amount)
}
func (a *stateDBIfaceAdapter) IsValidatorAddress(addr [20]byte) bool {
	return a.inner.IsValidatorAddress(common.Address(addr))
}
func (a *stateDBIfaceAdapter) GetState(addr [20]byte, key [32]byte) [32]byte {
	return [32]byte(a.inner.GetState(common.Address(addr), common.Hash(key)))
}
func (a *stateDBIfaceAdapter) SetState(addr [20]byte, key [32]byte, value [32]byte) {
	a.inner.SetState(common.Address(addr), common.Hash(key), common.Hash(value))
}
func (a *stateDBIfaceAdapter) GetCode(addr [20]byte) []byte {
	return a.inner.GetCode(common.Address(addr))
}
func (a *stateDBIfaceAdapter) SetCode(addr [20]byte, code []byte) {
	a.inner.SetCode(common.Address(addr), code)
}
func (a *stateDBIfaceAdapter) GetNonce(addr [20]byte) uint64 {
	return a.inner.GetNonce(common.Address(addr))
}
func (a *stateDBIfaceAdapter) SetNonce(addr [20]byte, nonce uint64) {
	a.inner.SetNonce(common.Address(addr), nonce)
}
func (a *stateDBIfaceAdapter) Exist(addr [20]byte) bool {
	return a.inner.Exist(common.Address(addr))
}
func (a *stateDBIfaceAdapter) CreateAccount(addr [20]byte) {
	a.inner.CreateAccount(common.Address(addr))
}
func (a *stateDBIfaceAdapter) Suicide(addr [20]byte) bool {
	return a.inner.Suicide(common.Address(addr))
}
func (a *stateDBIfaceAdapter) AddLog(entry *iface.LogEntry) {
	log := &types.Log{
		Address: common.Address(entry.Address),
		Data:    entry.Data,
	}
	log.Topics = make([]common.Hash, len(entry.Topics))
	for i, t := range entry.Topics {
		log.Topics[i] = common.Hash(t)
	}
	a.inner.AddLog(log)
}
func (a *stateDBIfaceAdapter) Snapshot() int {
	return a.inner.Snapshot()
}
func (a *stateDBIfaceAdapter) RevertToSnapshot(id int) {
	a.inner.RevertToSnapshot(id)
}
func (a *stateDBIfaceAdapter) ForEachStorage(addr [20]byte, cb func(key, value [32]byte) bool) error {
	return a.inner.ForEachStorage(common.Address(addr), func(k, v common.Hash) bool {
		return cb([32]byte(k), [32]byte(v))
	})
}

// ── BlockContext converter ────────────────────────────────────────────────────

// vmBlockContextToIface converts gwat's vm.BlockContext to iface.BlockContext
// for use by external processor factories.
func vmBlockContextToIface(ctx vm.BlockContext) iface.BlockContext {
	var random *[32]byte
	if ctx.Random != nil {
		r := [32]byte(*ctx.Random)
		random = &r
	}
	return iface.BlockContext{
		Coinbase:    [20]byte(ctx.Coinbase),
		GasLimit:    ctx.GasLimit,
		BlockHeight: ctx.BlockHeight,
		BlockNumber: ctx.BlockNumber,
		Time:        ctx.Time,
		Difficulty:  ctx.Difficulty,
		BaseFee:     ctx.BaseFee,
		Random:      random,
		Slot:        ctx.Slot,
		Era:         ctx.Era,
		BlockHash:   wfcommon.Hash(ctx.BlockHash),
	}
}
