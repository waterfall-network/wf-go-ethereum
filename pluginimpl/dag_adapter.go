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
	wfTypes "github.com/LFDT-Iguazu/iguazu-types/blockdag/types"
	wfCommon "github.com/LFDT-Iguazu/iguazu-types/common"
	gwatCommon "gitlab.waterfall.network/waterfall/protocol/gwat/common"
	gwatTypes "gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	"gitlab.waterfall.network/waterfall/protocol/gwat/internal/ethapi"
)

// wfDagAdapter implements ethapi.DagServicer by delegating to an iface.Dag.
// All type conversions are direct field-by-field casts — both sides use the
// same underlying types ([32]byte for Hash, [20]byte for Address, uint8 for
// SyncMode, etc.) so no allocation or serialisation is needed.
type wfDagAdapter struct {
	inner iface.Dag
}

var _ ethapi.DagServicer = (*wfDagAdapter)(nil)

func (a *wfDagAdapter) HandleFinalize(data *gwatTypes.FinalizationParams) *gwatTypes.FinalizationResult {
	return toGwatFinResult(a.inner.HandleFinalize(toWFFinParams(data)))
}

func (a *wfDagAdapter) HandleCoordinatedState() *gwatTypes.FinalizationResult {
	return toGwatFinResult(a.inner.HandleCoordinatedState())
}

func (a *wfDagAdapter) HandleGetCandidates(slot uint64) *gwatTypes.CandidatesResult {
	return toGwatCandidatesResult(a.inner.HandleGetCandidates(slot))
}

func (a *wfDagAdapter) HandleGetOptimisticSpines(fromSpine gwatCommon.Hash) *gwatTypes.OptimisticSpinesResult {
	return toGwatOptimisticSpinesResult(a.inner.HandleGetOptimisticSpines(wfCommon.Hash(fromSpine)))
}

func (a *wfDagAdapter) HandleValidateSpines(spines gwatCommon.HashArray) (bool, error) {
	return a.inner.HandleValidateSpines(toWFHashArray(spines))
}

func (a *wfDagAdapter) HandleValidateFinalization(spines gwatCommon.HashArray) (bool, error) {
	return a.inner.HandleValidateFinalization(toWFHashArray(spines))
}

func (a *wfDagAdapter) HandleSyncSpines(spines gwatCommon.HashArray) (bool, error) {
	return a.inner.HandleSyncSpines(toWFHashArray(spines))
}

func (a *wfDagAdapter) HandleSyncSlotInfo(si gwatTypes.SlotInfo) (bool, error) {
	return a.inner.HandleSyncSlotInfo(wfTypes.SlotInfo{
		GenesisTime:    si.GenesisTime,
		SecondsPerSlot: si.SecondsPerSlot,
		SlotsPerEpoch:  si.SlotsPerEpoch,
	})
}

// ---- gwat → iguazu-types conversions ----

func toWFHash(h gwatCommon.Hash) wfCommon.Hash { return wfCommon.Hash(h) }

func toWFHashArray(arr gwatCommon.HashArray) wfCommon.HashArray {
	if arr == nil {
		return nil
	}
	out := make(wfCommon.HashArray, len(arr))
	for i, h := range arr {
		out[i] = wfCommon.Hash(h)
	}
	return out
}

func toWFAddress(a gwatCommon.Address) wfCommon.Address { return wfCommon.Address(a) }

func toWFCheckpoint(cp *gwatTypes.Checkpoint) *wfTypes.Checkpoint {
	if cp == nil {
		return nil
	}
	return &wfTypes.Checkpoint{
		Epoch:    cp.Epoch,
		FinEpoch: cp.FinEpoch,
		Root:     toWFHash(cp.Root),
		Spine:    toWFHash(cp.Spine),
	}
}

func toWFValidatorSync(vs *gwatTypes.ValidatorSync) *wfTypes.ValidatorSync {
	if vs == nil {
		return nil
	}
	out := &wfTypes.ValidatorSync{
		InitTxHash:      wfCommon.Hash(vs.InitTxHash),
		OpType:          wfTypes.ValidatorSyncOp(vs.OpType),
		ProcEpoch:       vs.ProcEpoch,
		Index:           vs.Index,
		Creator:         toWFAddress(vs.Creator),
		ActivationEpoch: vs.ActivationEpoch,
		ExitEpoch:       vs.ExitEpoch,
	}
	if vs.Amount != nil {
		out.Amount = new(big.Int).Set(vs.Amount)
	}
	if vs.TxHash != nil {
		h := wfCommon.Hash(*vs.TxHash)
		out.TxHash = &h
	}
	if vs.Balance != nil {
		out.Balance = new(big.Int).Set(vs.Balance)
	}
	return out
}

func toWFFinParams(fp *gwatTypes.FinalizationParams) *wfTypes.FinalizationParams {
	if fp == nil {
		return nil
	}
	out := &wfTypes.FinalizationParams{
		Spines:   toWFHashArray(fp.Spines),
		SyncMode: wfTypes.SyncMode(fp.SyncMode),
	}
	if fp.BaseSpine != nil {
		h := wfCommon.Hash(*fp.BaseSpine)
		out.BaseSpine = &h
	}
	out.Checkpoint = toWFCheckpoint(fp.Checkpoint)
	if fp.ValSyncData != nil {
		out.ValSyncData = make([]*wfTypes.ValidatorSync, len(fp.ValSyncData))
		for i, vs := range fp.ValSyncData {
			out.ValSyncData[i] = toWFValidatorSync(vs)
		}
	}
	return out
}

// ---- iguazu-types → gwat conversions ----

func toGwatHash(h wfCommon.Hash) gwatCommon.Hash { return gwatCommon.Hash(h) }

func toGwatHashArray(arr wfCommon.HashArray) gwatCommon.HashArray {
	if arr == nil {
		return nil
	}
	out := make(gwatCommon.HashArray, len(arr))
	for i, h := range arr {
		out[i] = gwatCommon.Hash(h)
	}
	return out
}

func toGwatFinResult(r *wfTypes.FinalizationResult) *gwatTypes.FinalizationResult {
	if r == nil {
		return nil
	}
	out := &gwatTypes.FinalizationResult{
		Error:   r.Error,
		CpEpoch: r.CpEpoch,
	}
	if r.LFSpine != nil {
		h := toGwatHash(*r.LFSpine)
		out.LFSpine = &h
	}
	if r.CpRoot != nil {
		h := toGwatHash(*r.CpRoot)
		out.CpRoot = &h
	}
	return out
}

func toGwatCandidatesResult(r *wfTypes.CandidatesResult) *gwatTypes.CandidatesResult {
	if r == nil {
		return nil
	}
	return &gwatTypes.CandidatesResult{
		Error:      r.Error,
		Candidates: toGwatHashArray(r.Candidates),
	}
}

func toGwatOptimisticSpinesResult(r *wfTypes.OptimisticSpinesResult) *gwatTypes.OptimisticSpinesResult {
	if r == nil {
		return nil
	}
	out := &gwatTypes.OptimisticSpinesResult{
		Error: r.Error,
	}
	if r.Data != nil {
		out.Data = make([]gwatCommon.HashArray, len(r.Data))
		for i, arr := range r.Data {
			out.Data[i] = toGwatHashArray(arr)
		}
	}
	return out
}
