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
	"context"
	"fmt"

	"gitlab.waterfall.network/waterfall/protocol/gwat/core"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/rawdb"
	gwatcoretypes "gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	valStore "gitlab.waterfall.network/waterfall/protocol/gwat/validator/storage"
	"gitlab.waterfall.network/waterfall/protocol/gwat/validator/validatorsync"
	"gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/iface"
	wftypes "gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/types"
	wfera "gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/types/era"
	wfcommon "gitlab.waterfall.network/waterfall/protocol/wf-types/common"
)

// blockChainWrapper wraps *core.BlockChain and implements both iface.BlockChain
// and iface.ValidatorChain. All method signatures match exactly; only the
// concrete hash/address types differ and are converted at the boundary.
type blockChainWrapper struct{ inner *core.BlockChain }

// ---------- iface.BlockChain — Config ----------

func (w *blockChainWrapper) Config() *iface.ChainConfig {
	return wfChainConfig(w.inner.Config())
}

// ---------- iface.BlockChain — Block queries ----------

func (w *blockChainWrapper) GetLastFinalizedBlock() iface.Block {
	return wrapBlock(w.inner.GetLastFinalizedBlock())
}

func (w *blockChainWrapper) GetBlockByHash(hash wfcommon.Hash) iface.Block {
	return wrapBlock(w.inner.GetBlockByHash(gwatHash(hash)))
}

func (w *blockChainWrapper) GetBlocksByHashes(hashes wfcommon.HashArray) iface.BlockMap {
	return wrapBlockMap(w.inner.GetBlocksByHashes(gwatHashArray(hashes)))
}

func (w *blockChainWrapper) GetBlock(ctx context.Context, hash wfcommon.Hash) iface.Block {
	return wrapBlock(w.inner.GetBlock(ctx, gwatHash(hash)))
}

func (w *blockChainWrapper) GetLastFinalizedNumber() uint64 {
	return w.inner.GetLastFinalizedNumber()
}

func (w *blockChainWrapper) GetLastFinalizedHeader() *wftypes.HeaderInfo {
	return wfHeaderInfo(w.inner.GetLastFinalizedHeader())
}

func (w *blockChainWrapper) GetHeaderByHash(hash wfcommon.Hash) *wftypes.HeaderInfo {
	return wfHeaderInfo(w.inner.GetHeaderByHash(gwatHash(hash)))
}

func (w *blockChainWrapper) GetHeaderByNumber(number uint64) *wftypes.HeaderInfo {
	return wfHeaderInfo(w.inner.GetHeaderByNumber(number))
}

func (w *blockChainWrapper) GetHeadersByHashes(hashes wfcommon.HashArray) iface.HeaderMap {
	return wfHeaderMap(w.inner.GetHeadersByHashes(gwatHashArray(hashes)))
}

func (w *blockChainWrapper) GetBlockDag(hash wfcommon.Hash) *wftypes.BlockDAG {
	return wfBlockDAG(w.inner.GetBlockDag(gwatHash(hash)))
}

func (w *blockChainWrapper) GetTips() wftypes.Tips {
	return wfTips(w.inner.GetTips())
}

func (w *blockChainWrapper) GetOptimisticSpines(gtSlot uint64) ([]wfcommon.HashArray, error) {
	raw, err := w.inner.GetOptimisticSpines(gtSlot)
	if err != nil || raw == nil {
		return nil, err
	}
	out := make([]wfcommon.HashArray, len(raw))
	for i, ha := range raw {
		out[i] = wfHashArray(ha)
	}
	return out, nil
}

func (w *blockChainWrapper) GetFinalizedNumberByHash(hash wfcommon.Hash) *uint64 {
	return rawdb.ReadFinalizedNumberByHash(w.inner.Database(), gwatHash(hash))
}

func (w *blockChainWrapper) GetFinalizedHashByNumber(number uint64) wfcommon.Hash {
	return wfHash(w.inner.ReadFinalizedHashByNumber(number))
}

func (w *blockChainWrapper) GetBlockFinalizedNumber(hash wfcommon.Hash) *uint64 {
	return w.inner.GetBlockFinalizedNumber(gwatHash(hash))
}

// ---------- iface.BlockChain — Finalization mutations ----------

// FinalizeBlock adapts iface.BlockChain.FinalizeBlock to gwat's two-step API:
// UpdateFinalizingState(block, stateBlock) + WriteFinalizedBlock(finNr, block, isHead).
func (w *blockChainWrapper) FinalizeBlock(blockHash wfcommon.Hash, finNr uint64, stateBlockHash wfcommon.Hash, isHead bool) error {
	block := w.inner.GetBlockByHash(gwatHash(blockHash))
	if block == nil {
		return fmt.Errorf("pluginimpl: FinalizeBlock: block not found: %s", blockHash)
	}
	stateBlock := w.inner.GetBlockByHash(gwatHash(stateBlockHash))
	if stateBlock == nil {
		return fmt.Errorf("pluginimpl: FinalizeBlock: state block not found: %s", stateBlockHash)
	}
	// SetNumber must be called before UpdateFinalizingState so that NewEVMBlockContext
	// can safely dereference header.Number (e.g. at the ForkSlotValSyncProc fork point).
	// This mirrors gwat's internal finalizer (dag/finalizer/finalizer.go:156-157).
	block.SetNumber(&finNr)
	if err := w.inner.UpdateFinalizingState(block, stateBlock); err != nil {
		return err
	}
	return w.inner.WriteFinalizedBlock(finNr, block, isHead)
}

func (w *blockChainWrapper) FinalizeTips(finHashes wfcommon.HashArray, lastFinHash wfcommon.Hash, lastFinHeight uint64) {
	w.inner.FinalizeTips(gwatHashArray(finHashes), gwatHash(lastFinHash), lastFinHeight)
}

func (w *blockChainWrapper) RollbackFinalization(spineHash wfcommon.Hash, lfNr uint64) error {
	return w.inner.RollbackFinalization(gwatHash(spineHash), lfNr)
}

func (w *blockChainWrapper) SaveBlockDag(blockDag *wftypes.BlockDAG) {
	w.inner.SaveBlockDag(gwatBlockDAG(blockDag))
}

func (w *blockChainWrapper) SetRollbackActive()   { w.inner.SetRollbackActive() }
func (w *blockChainWrapper) ResetRollbackActive() { w.inner.ResetRollbackActive() }

// ---------- iface.BlockChain — DAG traversal ----------

func (w *blockChainWrapper) CollectAncestorsAftCpByTips(
	parents wfcommon.HashArray, cpHash wfcommon.Hash,
) (bool, iface.HeaderMap, wfcommon.HashArray, wftypes.Tips) {
	ok, hm, unloaded, tips := w.inner.CollectAncestorsAftCpByTips(gwatHashArray(parents), gwatHash(cpHash))
	return ok, wfHeaderMap(hm), wfHashArray(unloaded), wfTips(tips)
}

func (w *blockChainWrapper) CollectAncestorsAftCpByParents(
	parents wfcommon.HashArray, cpHash wfcommon.Hash,
) (bool, iface.HeaderMap, wfcommon.HashArray, error) {
	ok, hm, unloaded, err := w.inner.CollectAncestorsAftCpByParents(gwatHashArray(parents), gwatHash(cpHash))
	return ok, wfHeaderMap(hm), wfHashArray(unloaded), err
}

// ---------- iface.BlockChain — Slot / era / sync state ----------

func (w *blockChainWrapper) GetSlotInfo() *wftypes.SlotInfo {
	return wfSlotInfo(w.inner.GetSlotInfo())
}

func (w *blockChainWrapper) SetSlotInfo(si *wftypes.SlotInfo) error {
	return w.inner.SetSlotInfo(gwatSlotInfo(si))
}

func (w *blockChainWrapper) GetEraInfo() *wfera.EraInfo {
	return wfEraInfo(w.inner.GetEraInfo())
}

func (w *blockChainWrapper) HandleEra(cp *wftypes.Checkpoint) error {
	return w.inner.HandleEra(gwatCheckpoint(cp))
}

func (w *blockChainWrapper) GetLastCoordinatedCheckpoint() *wftypes.Checkpoint {
	return wfCheckpoint(w.inner.GetLastCoordinatedCheckpoint())
}

func (w *blockChainWrapper) SetLastCoordinatedCheckpoint(cp *wftypes.Checkpoint) {
	w.inner.SetLastCoordinatedCheckpoint(gwatCheckpoint(cp))
}

func (w *blockChainWrapper) IsSynced() bool            { return w.inner.IsSynced() }
func (w *blockChainWrapper) SetIsSynced(synced bool)   { w.inner.SetIsSynced(synced) }
func (w *blockChainWrapper) ResetSyncCheckpointCache() { w.inner.ResetSyncCheckpointCache() }

func (w *blockChainWrapper) SetSyncCheckpointCache(cp *wftypes.Checkpoint) {
	w.inner.SetSyncCheckpointCache(gwatCheckpoint(cp))
}

// ---------- iface.BlockChain — Validator sync ----------

func (w *blockChainWrapper) AppendNotProcessedValidatorSyncData(valSyncData []*wftypes.ValidatorSync) {
	w.inner.AppendNotProcessedValidatorSyncData(gwatValidatorSyncSlice(valSyncData))
}

// CleanInvalidNotProcessedValidatorSync delegates to gwat using the same
// validation function that gwat's own dag package uses.
func (w *blockChainWrapper) CleanInvalidNotProcessedValidatorSync() {
	w.inner.CleanInvalidNotProcessedValidatorSync(validatorsync.ValidateCreateTxValidatorSyncOp)
}

func (w *blockChainWrapper) FixValidatorSyncOps(finEpoch uint64, valSyncData []*wftypes.ValidatorSync) []*wftypes.ValidatorSync {
	result := core.FixValidatorSyncOps(w.inner, finEpoch, gwatValidatorSyncSlice(valSyncData))
	return wfValidatorSyncSlice(result)
}

// ---------- iface.BlockChain — Mutex ----------

func (w *blockChainWrapper) DagMuLock()   { w.inner.DagMuLock() }
func (w *blockChainWrapper) DagMuUnlock() { w.inner.DagMuUnlock() }

// ---------- iface.BlockChain — Validator storage ----------

func (w *blockChainWrapper) ValidatorStorage() iface.CreatorsBySlotStorage {
	return &creatorsBySlotWrapper{storage: w.inner.ValidatorStorage(), bc: w.inner}
}

// creatorsBySlotWrapper adapts gwat's valStore.Storage.GetCreatorsBySlot
// (which takes a blockchain context and variadic slot filter) to the narrow
// iface.CreatorsBySlotStorage interface (slot as a plain argument).
type creatorsBySlotWrapper struct {
	storage valStore.Storage
	bc      *core.BlockChain
}

func (c *creatorsBySlotWrapper) GetCreatorsBySlot(slot uint64) ([]wfcommon.Address, error) {
	addrs, err := c.storage.GetCreatorsBySlot(c.bc, slot)
	if err != nil || addrs == nil {
		return nil, err
	}
	out := make([]wfcommon.Address, len(addrs))
	for i, a := range addrs {
		out[i] = wfAddr(a)
	}
	return out, nil
}

// ---------- iface.ValidatorChain — Slot / era ----------

func (w *blockChainWrapper) EpochToEra(epoch uint64) *wfera.Era {
	return wfEra(w.inner.EpochToEra(epoch))
}

func (w *blockChainWrapper) ReadEra(eraNumber uint64) *wfera.Era {
	return wfEra(rawdb.ReadEra(w.inner.Database(), eraNumber))
}

// ---------- iface.ValidatorChain — State ----------

func (w *blockChainWrapper) StateAt(root wfcommon.Hash) (iface.StateDB, error) {
	return wrapStateDB(w.inner.StateAt(gwatHash(root)))
}

// ---------- iface.ValidatorChain — Validator sync data ----------

func (w *blockChainWrapper) GetValidatorSyncData(initTxHash wfcommon.Hash) *wftypes.ValidatorSync {
	return wfValidatorSync(w.inner.GetValidatorSyncData(gwatHash(initTxHash)))
}

func (w *blockChainWrapper) GetNotProcessedValidatorSyncData() map[wfcommon.Hash]*wftypes.ValidatorSync {
	return wfValidatorSyncMap(w.inner.GetNotProcessedValidatorSyncData())
}

func (w *blockChainWrapper) SetValidatorSyncData(vs *wftypes.ValidatorSync) {
	w.inner.SetValidatorSyncData(gwatValidatorSync(vs))
}

// ---------- iface.ValidatorChain — Transactions ----------

func (w *blockChainWrapper) GetTransaction(txHash wfcommon.Hash) (*wftypes.Transaction, wfcommon.Hash, uint64) {
	tx, blHash, index := w.inner.GetTransaction(gwatHash(txHash))
	if tx == nil {
		return nil, wfcommon.Hash{}, 0
	}
	raw, _ := tx.MarshalBinary()
	return &wftypes.Transaction{Raw: raw, Data: tx.Data()}, wfHash(blHash), index
}

func (w *blockChainWrapper) GetTransactionReceipt(txHash wfcommon.Hash) (*wftypes.Receipt, wfcommon.Hash, uint64) {
	rc, blHash, index := w.inner.GetTransactionReceipt(gwatHash(txHash))
	if rc == nil {
		// The receipt of a tx of the block being processed is not written yet,
		// while its TxLookup entry already resolves the block hash. Callers rely
		// on blHash to detect a preceding validator op within the same block,
		// so it must be propagated even when the receipt is not available.
		return nil, wfHash(blHash), index
	}
	return &wftypes.Receipt{Status: rc.Status}, wfHash(blHash), index
}

// GetTransactionSender recovers the sender of the transaction identified by
// txHash using gwat's latest signer derived from the chain config.
func (w *blockChainWrapper) GetTransactionSender(txHash wfcommon.Hash) (wfcommon.Address, error) {
	tx, _, _ := w.inner.GetTransaction(gwatHash(txHash))
	if tx == nil {
		return wfcommon.Address{}, fmt.Errorf("pluginimpl: GetTransactionSender: tx not found: %s", txHash)
	}
	signer := gwatcoretypes.LatestSigner(w.inner.Config())
	from, err := gwatcoretypes.Sender(signer, tx)
	if err != nil {
		return wfcommon.Address{}, err
	}
	return wfAddr(from), nil
}

// ---------- iface.ValidatorChain — Chain fixes ----------

func (w *blockChainWrapper) FixValidatorSyncOpProcessing(
	ctx iface.ProcessorCtx,
	op iface.ValidatorSyncOpFix,
	txHash wfcommon.Hash,
	from, to wfcommon.Address,
) (bool, []byte, error) {
	return w.inner.FixValidatorSyncOpProcessing(ctx, op, gwatHash(txHash), gwatAddr(from), gwatAddr(to))
}
