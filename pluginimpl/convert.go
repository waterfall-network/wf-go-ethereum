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

// Package pluginimpl wires gwat's concrete types to the wf-types/iface interfaces
// so that gwat can be loaded as a plugin by wf-engine.
package pluginimpl

import (
	"math/big"

	gwatcommon "gitlab.waterfall.network/waterfall/protocol/gwat/common"
	gwatcoretypes "gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	gwatparams "gitlab.waterfall.network/waterfall/protocol/gwat/params"
	gwaterra "gitlab.waterfall.network/waterfall/protocol/gwat/validator/era"
	"gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/iface"
	wftypes "gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/types"
	wfera "gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/types/era"
	wfcommon "gitlab.waterfall.network/waterfall/protocol/wf-types/common"
)

// ---------- scalar / primitive conversions ----------

func wfHash(h gwatcommon.Hash) wfcommon.Hash         { return wfcommon.Hash(h) }
func gwatHash(h wfcommon.Hash) gwatcommon.Hash       { return gwatcommon.Hash(h) }
func wfAddr(a gwatcommon.Address) wfcommon.Address   { return wfcommon.Address(a) }
func gwatAddr(a wfcommon.Address) gwatcommon.Address { return gwatcommon.Address(a) }

func wfHashArray(src gwatcommon.HashArray) wfcommon.HashArray {
	if src == nil {
		return nil
	}
	dst := make(wfcommon.HashArray, len(src))
	for i, h := range src {
		dst[i] = wfHash(h)
	}
	return dst
}

func gwatHashArray(src wfcommon.HashArray) gwatcommon.HashArray {
	if src == nil {
		return nil
	}
	dst := make(gwatcommon.HashArray, len(src))
	for i, h := range src {
		dst[i] = gwatHash(h)
	}
	return dst
}

// ---------- HeaderInfo ----------

func wfHeaderInfo(h *gwatcoretypes.Header) *wftypes.HeaderInfo {
	if h == nil {
		return nil
	}
	return &wftypes.HeaderInfo{
		Hash:         wfHash(h.Hash()),
		ParentHashes: wfHashArray(h.ParentHashes),
		Slot:         h.Slot,
		Era:          h.Era,
		Height:       h.Height,
		Nr:           h.Nr(),
		GasLimit:     h.GasLimit,
		Time:         h.Time,
		Coinbase:     wfAddr(h.Coinbase),
		CpHash:       wfHash(h.CpHash),
		Root:         wfHash(h.Root),
	}
}

func wfHeaderMap(src gwatcoretypes.HeaderMap) iface.HeaderMap {
	if src == nil {
		return nil
	}
	dst := make(iface.HeaderMap, len(src))
	for k, v := range src {
		dst[wfHash(k)] = wfHeaderInfo(v)
	}
	return dst
}

// ---------- BlockDAG ----------

func wfBlockDAG(b *gwatcoretypes.BlockDAG) *wftypes.BlockDAG {
	if b == nil {
		return nil
	}
	return &wftypes.BlockDAG{
		Hash:                   wfHash(b.Hash),
		Height:                 b.Height,
		Slot:                   b.Slot,
		CpHash:                 wfHash(b.CpHash),
		CpHeight:               b.CpHeight,
		OrderedAncestorsHashes: wfHashArray(b.OrderedAncestorsHashes),
	}
}

func gwatBlockDAG(b *wftypes.BlockDAG) *gwatcoretypes.BlockDAG {
	if b == nil {
		return nil
	}
	return &gwatcoretypes.BlockDAG{
		Hash:                   gwatHash(b.Hash),
		Height:                 b.Height,
		Slot:                   b.Slot,
		CpHash:                 gwatHash(b.CpHash),
		CpHeight:               b.CpHeight,
		OrderedAncestorsHashes: gwatHashArray(b.OrderedAncestorsHashes),
	}
}

// ---------- Tips ----------

func wfTips(src gwatcoretypes.Tips) wftypes.Tips {
	if src == nil {
		return nil
	}
	dst := make(wftypes.Tips, len(src))
	for k, v := range src {
		dst[wfHash(k)] = wfBlockDAG(v)
	}
	return dst
}

// ---------- SlotInfo ----------

func wfSlotInfo(si *gwatcoretypes.SlotInfo) *wftypes.SlotInfo {
	if si == nil {
		return nil
	}
	return &wftypes.SlotInfo{
		GenesisTime:    si.GenesisTime,
		SecondsPerSlot: si.SecondsPerSlot,
		SlotsPerEpoch:  si.SlotsPerEpoch,
	}
}

func gwatSlotInfo(si *wftypes.SlotInfo) *gwatcoretypes.SlotInfo {
	if si == nil {
		return nil
	}
	return &gwatcoretypes.SlotInfo{
		GenesisTime:    si.GenesisTime,
		SecondsPerSlot: si.SecondsPerSlot,
		SlotsPerEpoch:  si.SlotsPerEpoch,
	}
}

// ---------- Checkpoint ----------

func wfCheckpoint(cp *gwatcoretypes.Checkpoint) *wftypes.Checkpoint {
	if cp == nil {
		return nil
	}
	return &wftypes.Checkpoint{
		Epoch:    cp.Epoch,
		FinEpoch: cp.FinEpoch,
		Root:     wfHash(cp.Root),
		Spine:    wfHash(cp.Spine),
	}
}

func gwatCheckpoint(cp *wftypes.Checkpoint) *gwatcoretypes.Checkpoint {
	if cp == nil {
		return nil
	}
	return &gwatcoretypes.Checkpoint{
		Epoch:    cp.Epoch,
		FinEpoch: cp.FinEpoch,
		Root:     gwatHash(cp.Root),
		Spine:    gwatHash(cp.Spine),
	}
}

// ---------- ValidatorSync ----------

func wfValidatorSync(vs *gwatcoretypes.ValidatorSync) *wftypes.ValidatorSync {
	if vs == nil {
		return nil
	}
	out := &wftypes.ValidatorSync{
		InitTxHash:      wfHash(vs.InitTxHash),
		OpType:          wftypes.ValidatorSyncOp(vs.OpType),
		ProcEpoch:       vs.ProcEpoch,
		Index:           vs.Index,
		Creator:         wfAddr(vs.Creator),
		ActivationEpoch: vs.ActivationEpoch,
		ExitEpoch:       vs.ExitEpoch,
	}
	if vs.Amount != nil {
		out.Amount = new(big.Int).Set(vs.Amount)
	}
	if vs.TxHash != nil {
		h := wfHash(*vs.TxHash)
		out.TxHash = &h
	}
	if vs.Balance != nil {
		out.Balance = new(big.Int).Set(vs.Balance)
	}
	return out
}

func gwatValidatorSync(vs *wftypes.ValidatorSync) *gwatcoretypes.ValidatorSync {
	if vs == nil {
		return nil
	}
	out := &gwatcoretypes.ValidatorSync{
		InitTxHash:      gwatHash(vs.InitTxHash),
		OpType:          gwatcoretypes.ValidatorSyncOp(vs.OpType),
		ProcEpoch:       vs.ProcEpoch,
		Index:           vs.Index,
		Creator:         gwatAddr(vs.Creator),
		ActivationEpoch: vs.ActivationEpoch,
		ExitEpoch:       vs.ExitEpoch,
	}
	if vs.Amount != nil {
		out.Amount = new(big.Int).Set(vs.Amount)
	}
	if vs.TxHash != nil {
		h := gwatHash(*vs.TxHash)
		out.TxHash = &h
	}
	if vs.Balance != nil {
		out.Balance = new(big.Int).Set(vs.Balance)
	}
	return out
}

func wfValidatorSyncSlice(src []*gwatcoretypes.ValidatorSync) []*wftypes.ValidatorSync {
	if src == nil {
		return nil
	}
	dst := make([]*wftypes.ValidatorSync, len(src))
	for i, v := range src {
		dst[i] = wfValidatorSync(v)
	}
	return dst
}

func gwatValidatorSyncSlice(src []*wftypes.ValidatorSync) []*gwatcoretypes.ValidatorSync {
	if src == nil {
		return nil
	}
	dst := make([]*gwatcoretypes.ValidatorSync, len(src))
	for i, v := range src {
		dst[i] = gwatValidatorSync(v)
	}
	return dst
}

func wfValidatorSyncMap(src map[gwatcommon.Hash]*gwatcoretypes.ValidatorSync) map[wfcommon.Hash]*wftypes.ValidatorSync {
	if src == nil {
		return nil
	}
	dst := make(map[wfcommon.Hash]*wftypes.ValidatorSync, len(src))
	for k, v := range src {
		dst[wfHash(k)] = wfValidatorSync(v)
	}
	return dst
}

// ---------- Era ----------

func wfEra(e *gwaterra.Era) *wfera.Era {
	if e == nil {
		return nil
	}
	return &wfera.Era{
		Number:    e.Number,
		From:      e.From,
		To:        e.To,
		Root:      wfHash(e.Root),
		BlockHash: wfHash(e.BlockHash),
	}
}

func wfEraInfo(ei *gwaterra.EraInfo) *wfera.EraInfo {
	if ei == nil {
		return nil
	}
	return wfera.NewEraInfo(wfEra(ei.GetEra()))
}

// ---------- ChainConfig ----------

func wfChainConfig(c *gwatparams.ChainConfig) *iface.ChainConfig {
	if c == nil {
		return nil
	}
	var chainID uint64
	if c.ChainID != nil {
		chainID = c.ChainID.Uint64()
	}
	var validatorsStateAddr wfcommon.Address
	if c.ValidatorsStateAddress != nil {
		validatorsStateAddr = wfAddr(*c.ValidatorsStateAddress)
	}
	var acceptCpRoot map[[32]byte][]uint64
	if c.AcceptCpRootOnFinEpoch != nil {
		acceptCpRoot = make(map[[32]byte][]uint64, len(c.AcceptCpRootOnFinEpoch))
		for k, v := range c.AcceptCpRootOnFinEpoch {
			acceptCpRoot[[32]byte(k)] = v
		}
	}
	return &iface.ChainConfig{
		ChainID:                   chainID,
		EpochsPerEra:              c.EpochsPerEra,
		SlotsPerEpoch:             c.SlotsPerEpoch,
		ValidatorsPerSlot:         c.ValidatorsPerSlot,
		TransitionPeriod:          c.TransitionPeriod,
		StartEpochsPerEra:         c.StartEpochsPerEra,
		AcceptCpRootOnFinEpoch:    acceptCpRoot,
		ValidatorsStateAddress:    validatorsStateAddr,
		WaterfallDummyAddress:     wfAddr(c.WaterfallDummyAddress),
		AllocationContractAddress: wfAddr(c.AllocationContractAddress),
		EffectiveBalance:          c.EffectiveBalance,
		ValidatorOpExpireSlots:    c.ValidatorOpExpireSlots,
		ForkSlotDelegate:          c.ForkSlotDelegate,
		ForkSlotValSyncProc:       c.ForkSlotValSyncProc,
		ForkSlotValOpTracking:     c.ForkSlotValOpTracking,
	}
}
