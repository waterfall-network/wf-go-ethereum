package core

import (
	"errors"
	"testing"

	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	"gitlab.waterfall.network/waterfall/protocol/gwat/params"
)

// ── mock ─────────────────────────────────────────────────────────────────────

// mockFixChain implements fixValSyncChain for testing.
type mockFixChain struct {
	slotInfo        *types.SlotInfo
	config          *params.ChainConfig
	genesis         *types.Block
	txBlockHashes   map[common.Hash]common.Hash
	valSyncData     map[common.Hash]*types.ValidatorSync
	restoreErr      error
	restoreCalls    []common.Hash
	getTxBlockCalls []common.Hash
	setSyncCalls    []*types.ValidatorSync
}

func (m *mockFixChain) GetSlotInfo() *types.SlotInfo { return m.slotInfo }
func (m *mockFixChain) Config() *params.ChainConfig  { return m.config }
func (m *mockFixChain) Genesis() *types.Block        { return m.genesis }
func (m *mockFixChain) GetTxBlockHash(txHash common.Hash) common.Hash {
	m.getTxBlockCalls = append(m.getTxBlockCalls, txHash)
	return m.txBlockHashes[txHash]
}
func (m *mockFixChain) RestoreTxLookupEntries(blHash common.Hash) error {
	m.restoreCalls = append(m.restoreCalls, blHash)
	return m.restoreErr
}
func (m *mockFixChain) GetValidatorSyncData(initTxHash common.Hash) *types.ValidatorSync {
	return m.valSyncData[initTxHash]
}
func (m *mockFixChain) SetValidatorSyncData(vs *types.ValidatorSync) {
	m.setSyncCalls = append(m.setSyncCalls, vs)
	m.valSyncData[vs.InitTxHash] = vs
}

// newMockFixChain creates a mock with the given forkSlot and SlotsPerEpoch=32.
// The genesis block is non-mainnet.
func newMockFixChain(forkSlot uint64) *mockFixChain {
	return &mockFixChain{
		slotInfo:      &types.SlotInfo{SlotsPerEpoch: 32},
		config:        &params.ChainConfig{ForkSlotValSyncProc: forkSlot},
		genesis:       types.NewBlockWithHeader(&types.Header{}),
		txBlockHashes: make(map[common.Hash]common.Hash),
		valSyncData:   make(map[common.Hash]*types.ValidatorSync),
	}
}

// testFixData returns a minimal fixData map with a single op for unit tests.
// forkSlot=320 → forkEpoch=10 (320/32), procEpoch=14.
func testFixData() (map[common.Hash]*FixOp, *FixOp) {
	initTxHash := common.HexToHash("0xaaaa000000000000000000000000000000000000000000000000000000000000")
	initTxBlock := common.HexToHash("0xbbbb000000000000000000000000000000000000000000000000000000000000")
	op := &FixOp{
		OpType:      types.Activate,
		Index:       99,
		Creator:     common.HexToAddress("0xcccc"),
		InitTxHash:  initTxHash,
		InitTxBlock: initTxBlock,
	}
	return map[common.Hash]*FixOp{initTxHash: op}, op
}

// ── FixOp.CreateValidatorSync ─────────────────────────────────────────────────

func TestFixOp_CreateValidatorSync(t *testing.T) {
	initTxHash := common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")
	creator := common.HexToAddress("0xaaaa")
	op := &FixOp{
		OpType:     types.Activate,
		Index:      42,
		Creator:    creator,
		InitTxHash: initTxHash,
	}
	const procEpoch = uint64(100)

	vs := op.CreateValidatorSync(procEpoch)

	if vs.InitTxHash != initTxHash {
		t.Errorf("InitTxHash: got %v, want %v", vs.InitTxHash, initTxHash)
	}
	if vs.OpType != types.Activate {
		t.Errorf("OpType: got %v, want Activate", vs.OpType)
	}
	if vs.ProcEpoch != procEpoch {
		t.Errorf("ProcEpoch: got %d, want %d", vs.ProcEpoch, procEpoch)
	}
	if vs.Index != 42 {
		t.Errorf("Index: got %d, want 42", vs.Index)
	}
	if vs.Creator != creator {
		t.Errorf("Creator: got %v, want %v", vs.Creator, creator)
	}
	if vs.TxHash != nil {
		t.Errorf("TxHash: got %v, want nil", vs.TxHash)
	}
	if vs.Amount != nil {
		t.Errorf("Amount: got %v, want nil", vs.Amount)
	}
	if vs.Balance != nil {
		t.Errorf("Balance: got %v, want nil", vs.Balance)
	}
}

// ── mainnetValSyncFixData integrity ──────────────────────────────────────────

func TestMainnetValSyncFixData_Count(t *testing.T) {
	const want = 37
	if got := len(mainnetValSyncFixData); got != want {
		t.Errorf("len(mainnetValSyncFixData) = %d, want %d", got, want)
	}
}

func TestMainnetValSyncFixData_KeyEqualsInitTxHash(t *testing.T) {
	for key, op := range mainnetValSyncFixData {
		if key != op.InitTxHash {
			t.Errorf("index %d: map key %v != InitTxHash %v", op.Index, key, op.InitTxHash)
		}
	}
}

func TestMainnetValSyncFixData_KnownEntries(t *testing.T) {
	tests := []struct {
		initTxHash  string
		opType      types.ValidatorSyncOp
		index       uint64
		creator     string
		initTxBlock string
	}{
		{
			// first pre-fork entry
			initTxHash:  "0x4f2c8e7b7b9eb70fa1941236519714e5a670f4618efbcfa7a1a2327fc4285bed",
			opType:      types.Deactivate,
			index:       2125,
			creator:     "0xd374819f2f66d2828b06ef7777e8ab97e6a2ebca",
			initTxBlock: "0x5bb0f5335b15d5bd68bc7d56cbd6598ac3a2a49ad051616d16f16017af13bf9a",
		},
		{
			// first post-fork entry (has TxHash)
			initTxHash:  "0x1c893460442fdc5a78ce2de79d11a774b1c3345bfabcff7f427c2106aa8f0289",
			opType:      types.Activate,
			index:       46793,
			creator:     "0x134828c57593d6f913dcea89e5fb2dc233d65e63",
			initTxBlock: "0xae0f21c95aaae244bf609d929927a0e51f1457e4a9b37fc04a147f5257df0e23",
		},
		{
			// last post-fork entry
			initTxHash:  "0xbdbf54628dbaf2816a98af968214040ac466f63eb41e56b0f0fbf5c0f29ed106",
			opType:      types.Activate,
			index:       48374,
			creator:     "0xd67d91840744ed2998e14ab54ab3e675e6ed09fc",
			initTxBlock: "0xf3b90ed5118a2a07182531a2b9105ed6675b2138e029cc95122c38245387f0e2",
		},
	}
	for _, tc := range tests {
		key := common.HexToHash(tc.initTxHash)
		op, ok := mainnetValSyncFixData[key]
		if !ok {
			t.Errorf("entry %s not found in mainnetValSyncFixData", tc.initTxHash)
			continue
		}
		if op.OpType != tc.opType {
			t.Errorf("index %d OpType: got %v, want %v", tc.index, op.OpType, tc.opType)
		}
		if op.Index != tc.index {
			t.Errorf("initTxHash %s: Index got %d, want %d", tc.initTxHash, op.Index, tc.index)
		}
		if op.Creator != common.HexToAddress(tc.creator) {
			t.Errorf("index %d Creator: got %v, want %v", tc.index, op.Creator, tc.creator)
		}
		if op.InitTxBlock != common.HexToHash(tc.initTxBlock) {
			t.Errorf("index %d InitTxBlock: got %v, want %v", tc.index, op.InitTxBlock, tc.initTxBlock)
		}
	}
}

// ── FixValidatorSyncOps (public wrapper) ───────────────────────────────────

func TestGetFixValidatorSyncOps_NonMainnetReturnsEmpty(t *testing.T) {
	// forkSlot=320 → forkEpoch=10, procEpoch=14; epoch 10 is in window,
	// but non-mainnet genesis means fixData is nil → result must be empty.
	bc := newMockFixChain(320)
	result := FixValidatorSyncOps(bc, 10, nil)
	if len(result) != 0 {
		t.Errorf("expected empty result for non-mainnet, got %d ops", len(result))
	}
}

// ── getFixValidatorSyncOpsFromData (internal) ─────────────────────────────────

func TestGetFixValidatorSyncOpsFromData_EpochBounds(t *testing.T) {
	// forkSlot=320 → forkEpoch=10, procEpoch=14; window is [10, 14).
	fixData, op := testFixData()
	tests := []struct {
		currEpoch uint64
		wantOps   int
	}{
		{9, 0},  // before forkEpoch
		{10, 1}, // == forkEpoch (window start)
		{13, 1}, // last epoch inside window
		{14, 0}, // == procEpoch (window end, exclusive)
		{15, 0}, // after procEpoch
	}
	for _, tc := range tests {
		bc := newMockFixChain(320)
		bc.txBlockHashes[op.InitTxHash] = op.InitTxBlock // TxLookup already present

		result := getFixValidatorSyncOpsFromData(bc, tc.currEpoch, fixData)
		if len(result) != tc.wantOps {
			t.Errorf("currEpoch=%d: got %d ops, want %d", tc.currEpoch, len(result), tc.wantOps)
		}
	}
}

func TestGetFixValidatorSyncOpsFromData_NilFixData(t *testing.T) {
	bc := newMockFixChain(320)
	result := getFixValidatorSyncOpsFromData(bc, 10, nil)
	if len(result) != 0 {
		t.Errorf("expected empty result for nil fixData, got %d ops", len(result))
	}
}

func TestGetFixValidatorSyncOpsFromData_TxLookupMissing_Restored(t *testing.T) {
	// GetTxBlockHash returns zero (wrong) → RestoreTxLookupEntries must be called.
	fixData, op := testFixData()
	bc := newMockFixChain(320)
	// txBlockHashes is empty → GetTxBlockHash returns common.Hash{}

	result := getFixValidatorSyncOpsFromData(bc, 10, fixData)

	if len(result) != 1 {
		t.Fatalf("expected 1 op, got %d", len(result))
	}
	if len(bc.restoreCalls) != 1 {
		t.Fatalf("expected 1 RestoreTxLookupEntries call, got %d", len(bc.restoreCalls))
	}
	if bc.restoreCalls[0] != op.InitTxBlock {
		t.Errorf("RestoreTxLookupEntries called with %v, want %v", bc.restoreCalls[0], op.InitTxBlock)
	}
}

func TestGetFixValidatorSyncOpsFromData_TxLookupPresent_NotRestored(t *testing.T) {
	// GetTxBlockHash returns the expected block → no restore needed.
	fixData, op := testFixData()
	bc := newMockFixChain(320)
	bc.txBlockHashes[op.InitTxHash] = op.InitTxBlock

	result := getFixValidatorSyncOpsFromData(bc, 10, fixData)

	if len(result) != 1 {
		t.Fatalf("expected 1 op, got %d", len(result))
	}
	if len(bc.restoreCalls) != 0 {
		t.Errorf("expected no RestoreTxLookupEntries calls, got %d", len(bc.restoreCalls))
	}
}

func TestGetFixValidatorSyncOpsFromData_TxLookupRestoreFails_OmitsOp(t *testing.T) {
	// If RestoreTxLookupEntries returns an error, the op must be skipped.
	fixData, _ := testFixData()
	bc := newMockFixChain(320)
	bc.restoreErr = errors.New("disk error")

	result := getFixValidatorSyncOpsFromData(bc, 10, fixData)

	if len(result) != 0 {
		t.Errorf("expected op to be omitted on restore error, got %d ops", len(result))
	}
}

func TestGetFixValidatorSyncOpsFromData_ZeroInitTxBlock_SkipsTxLookupCheck(t *testing.T) {
	// InitTxBlock == zero hash → TxLookup check must be skipped entirely.
	initTxHash := common.HexToHash("0xdddd000000000000000000000000000000000000000000000000000000000000")
	op := &FixOp{
		OpType:      types.Activate,
		Index:       1,
		Creator:     common.HexToAddress("0x1111"),
		InitTxHash:  initTxHash,
		InitTxBlock: common.Hash{}, // zero → skip
	}
	fixData := map[common.Hash]*FixOp{initTxHash: op}
	bc := newMockFixChain(320)

	result := getFixValidatorSyncOpsFromData(bc, 10, fixData)

	if len(result) != 1 {
		t.Fatalf("expected 1 op, got %d", len(result))
	}
	if len(bc.getTxBlockCalls) != 0 {
		t.Errorf("expected no GetTxBlockHash calls for zero InitTxBlock, got %d", len(bc.getTxBlockCalls))
	}
	if len(bc.restoreCalls) != 0 {
		t.Errorf("expected no RestoreTxLookupEntries calls, got %d", len(bc.restoreCalls))
	}
}

func TestGetFixValidatorSyncOpsFromData_ProcEpochInResult(t *testing.T) {
	// forkSlot=320 → forkEpoch=10, procEpoch=14; returned ops must have ProcEpoch=14.
	fixData, op := testFixData()
	bc := newMockFixChain(320)
	bc.txBlockHashes[op.InitTxHash] = op.InitTxBlock

	result := getFixValidatorSyncOpsFromData(bc, 10, fixData)

	if len(result) != 1 {
		t.Fatalf("expected 1 op, got %d", len(result))
	}
	const wantProcEpoch = uint64(14) // forkEpoch(10) + 4
	if result[0].ProcEpoch != wantProcEpoch {
		t.Errorf("ProcEpoch: got %d, want %d", result[0].ProcEpoch, wantProcEpoch)
	}
}

func TestGetFixValidatorSyncOpsFromData_ZeroTxHash_Reset(t *testing.T) {
	// Stored op has TxHash == zero hash (written by fixMainnet0).
	// getFixValidatorSyncOpsFromData must reset TxHash→nil and ProcEpoch=procEpoch and save.
	fixData, op := testFixData()
	bc := newMockFixChain(320)
	bc.txBlockHashes[op.InitTxHash] = op.InitTxBlock

	zeroHash := common.Hash{}
	bc.valSyncData[op.InitTxHash] = &types.ValidatorSync{
		InitTxHash: op.InitTxHash,
		OpType:     op.OpType,
		Index:      op.Index,
		Creator:    op.Creator,
		TxHash:     &zeroHash,
		ProcEpoch:  0,
	}

	result := getFixValidatorSyncOpsFromData(bc, 10, fixData)

	if len(result) != 1 {
		t.Fatalf("expected 1 op, got %d", len(result))
	}
	if len(bc.setSyncCalls) != 1 {
		t.Fatalf("expected 1 SetValidatorSyncData call, got %d", len(bc.setSyncCalls))
	}
	saved := bc.setSyncCalls[0]
	if saved.TxHash != nil {
		t.Errorf("saved TxHash: got %v, want nil", saved.TxHash)
	}
	const wantProcEpoch = uint64(14)
	if saved.ProcEpoch != wantProcEpoch {
		t.Errorf("saved ProcEpoch: got %d, want %d", saved.ProcEpoch, wantProcEpoch)
	}
}

func TestGetFixValidatorSyncOpsFromData_NonZeroTxHash_NotReset(t *testing.T) {
	// Stored op has a real (non-zero) TxHash → must NOT be overwritten.
	fixData, op := testFixData()
	bc := newMockFixChain(320)
	bc.txBlockHashes[op.InitTxHash] = op.InitTxBlock

	realHash := common.HexToHash("0x1234000000000000000000000000000000000000000000000000000000000000")
	bc.valSyncData[op.InitTxHash] = &types.ValidatorSync{
		InitTxHash: op.InitTxHash,
		OpType:     op.OpType,
		Index:      op.Index,
		Creator:    op.Creator,
		TxHash:     &realHash,
		ProcEpoch:  5,
	}

	result := getFixValidatorSyncOpsFromData(bc, 10, fixData)

	if len(result) != 1 {
		t.Fatalf("expected 1 op, got %d", len(result))
	}
	if len(bc.setSyncCalls) != 0 {
		t.Errorf("expected no SetValidatorSyncData calls, got %d", len(bc.setSyncCalls))
	}
}

func TestGetFixValidatorSyncOpsFromData_NilStoredOp_NotReset(t *testing.T) {
	// No stored op → GetValidatorSyncData returns nil → no save.
	fixData, op := testFixData()
	bc := newMockFixChain(320)
	bc.txBlockHashes[op.InitTxHash] = op.InitTxBlock
	// valSyncData is empty

	result := getFixValidatorSyncOpsFromData(bc, 10, fixData)

	if len(result) != 1 {
		t.Fatalf("expected 1 op, got %d", len(result))
	}
	if len(bc.setSyncCalls) != 0 {
		t.Errorf("expected no SetValidatorSyncData calls, got %d", len(bc.setSyncCalls))
	}
}

// ── mergeValSyncOps ───────────────────────────────────────────────────────────

func TestMergeValSyncOps_EmptyFixes_ReturnsBase(t *testing.T) {
	hash := common.HexToHash("0x0101010101010101010101010101010101010101010101010101010101010101")
	base := []*types.ValidatorSync{{InitTxHash: hash}}

	result := mergeValSyncOps(base, nil)

	if len(result) != 1 || result[0].InitTxHash != hash {
		t.Errorf("expected base slice unchanged, got %v", result)
	}
}

func TestMergeValSyncOps_NoDuplicates_AppendsFixes(t *testing.T) {
	h1 := common.HexToHash("0x0101010101010101010101010101010101010101010101010101010101010101")
	h2 := common.HexToHash("0x0202020202020202020202020202020202020202020202020202020202020202")
	base := []*types.ValidatorSync{{InitTxHash: h1}}
	fixes := []*types.ValidatorSync{{InitTxHash: h2}}

	result := mergeValSyncOps(base, fixes)

	if len(result) != 2 {
		t.Fatalf("expected 2 ops, got %d", len(result))
	}
}

func TestMergeValSyncOps_DuplicateInitTxHash_FixWins(t *testing.T) {
	hash := common.HexToHash("0x0303030303030303030303030303030303030303030303030303030303030303")
	baseOp := &types.ValidatorSync{InitTxHash: hash, ProcEpoch: 5}
	fixOp := &types.ValidatorSync{InitTxHash: hash, ProcEpoch: 14}
	base := []*types.ValidatorSync{baseOp}
	fixes := []*types.ValidatorSync{fixOp}

	result := mergeValSyncOps(base, fixes)

	if len(result) != 1 {
		t.Fatalf("expected 1 op after dedup, got %d", len(result))
	}
	if result[0].ProcEpoch != 14 {
		t.Errorf("fix op should win: ProcEpoch got %d, want 14", result[0].ProcEpoch)
	}
}

func TestMergeValSyncOps_MultipleDuplicates_AllDropped(t *testing.T) {
	h1 := common.HexToHash("0x0404040404040404040404040404040404040404040404040404040404040404")
	h2 := common.HexToHash("0x0505050505050505050505050505050505050505050505050505050505050505")
	h3 := common.HexToHash("0x0606060606060606060606060606060606060606060606060606060606060606")
	base := []*types.ValidatorSync{
		{InitTxHash: h1, ProcEpoch: 1},
		{InitTxHash: h2, ProcEpoch: 2},
		{InitTxHash: h3, ProcEpoch: 3},
	}
	fixes := []*types.ValidatorSync{
		{InitTxHash: h1, ProcEpoch: 14},
		{InitTxHash: h3, ProcEpoch: 14},
	}

	result := mergeValSyncOps(base, fixes)

	// h2 from base + h1,h3 from fixes = 3 total
	if len(result) != 3 {
		t.Fatalf("expected 3 ops, got %d", len(result))
	}
	// h2 base entry must be preserved
	found := false
	for _, op := range result {
		if op.InitTxHash == h2 && op.ProcEpoch == 2 {
			found = true
		}
	}
	if !found {
		t.Errorf("base op with h2/ProcEpoch=2 should be preserved")
	}
}

func TestMergeValSyncOps_NilBase_OnlyFixes(t *testing.T) {
	hash := common.HexToHash("0x0707070707070707070707070707070707070707070707070707070707070707")
	fixes := []*types.ValidatorSync{{InitTxHash: hash, ProcEpoch: 14}}

	result := mergeValSyncOps(nil, fixes)

	if len(result) != 1 || result[0].InitTxHash != hash {
		t.Errorf("expected single fix op, got %v", result)
	}
}

// ── FixValidatorSyncOps deduplication ─────────────────────────────────────────

func TestFixValidatorSyncOps_NonMainnet_PassesThroughBase(t *testing.T) {
	// Non-mainnet: no fix ops → base returned unchanged.
	bc := newMockFixChain(320)
	hash := common.HexToHash("0x0808080808080808080808080808080808080808080808080808080808080808")
	base := []*types.ValidatorSync{{InitTxHash: hash}}

	result := FixValidatorSyncOps(bc, 10, base)

	// fixData is nil (non-mainnet) → fixValSyncOps is empty → mergeValSyncOps returns base
	if len(result) != 1 || result[0].InitTxHash != hash {
		t.Errorf("expected base unchanged for non-mainnet, got %v", result)
	}
}
