package core

import (
	"hash"
	"math/big"
	"testing"

	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/rawdb"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	"golang.org/x/crypto/sha3"
)

// txTestHasher is a minimal TrieHasher for constructing test blocks without
// importing the trie package.
type txTestHasher struct{ h hash.Hash }

func newTxTestHasher() *txTestHasher        { return &txTestHasher{h: sha3.NewLegacyKeccak256()} }
func (th *txTestHasher) Reset()             { th.h.Reset() }
func (th *txTestHasher) Update(k, v []byte) { th.h.Write(k); th.h.Write(v) }
func (th *txTestHasher) Hash() common.Hash  { return common.BytesToHash(th.h.Sum(nil)) }

// newTestTx returns an unsigned legacy transaction for unit tests.
func newTestTx(nonce uint64) *types.Transaction {
	return types.NewTransaction(nonce, common.Address{1}, big.NewInt(0), 21000, big.NewInt(1), nil)
}

// ── GetTxBlockHash ────────────────────────────────────────────────────────────

func TestGetTxBlockHash_NotFound(t *testing.T) {
	bc := initValSyncOpChain(t)
	if got := bc.GetTxBlockHash(common.HexToHash("0x1234")); got != (common.Hash{}) {
		t.Errorf("expected zero hash, got %v", got)
	}
}

func TestGetTxBlockHash_FromDB(t *testing.T) {
	bc := initValSyncOpChain(t)
	txHash := common.HexToHash("0xaaaa")
	blockHash := common.HexToHash("0xbbbb")
	rawdb.WriteTxLookupEntry(bc.db, txHash, blockHash)

	if got := bc.GetTxBlockHash(txHash); got != blockHash {
		t.Errorf("got %v, want %v", got, blockHash)
	}
}

func TestGetTxBlockHash_FromCache(t *testing.T) {
	bc := initValSyncOpChain(t)
	txHash := common.HexToHash("0xcccc")
	blockHash := common.HexToHash("0xdddd")

	// WriteTxLookupEntry (method) writes to both DB and cache.
	bc.WriteTxLookupEntry(0, txHash, blockHash, types.ReceiptStatusSuccessful)
	// Remove from DB so only the cache can satisfy the lookup.
	rawdb.DeleteTxLookupEntry(bc.db, txHash)

	if got := bc.GetTxBlockHash(txHash); got != blockHash {
		t.Errorf("got %v, want %v", got, blockHash)
	}
}

// ── RestoreTxLookupEntries ────────────────────────────────────────────────────

func TestRestoreTxLookupEntries_BlockNotFound(t *testing.T) {
	bc := initValSyncOpChain(t)
	if err := bc.RestoreTxLookupEntries(common.HexToHash("0xdeadbeef")); err == nil {
		t.Fatal("expected error for missing block, got nil")
	}
}

func TestRestoreTxLookupEntries_NoReceipts(t *testing.T) {
	bc := initValSyncOpChain(t)
	tx := newTestTx(0)
	block := types.NewBlock(&types.Header{}, []*types.Transaction{tx}, nil, newTxTestHasher())
	rawdb.WriteBlock(bc.db, block)
	// No receipts written → GetReceiptsByHash returns nil.

	if err := bc.RestoreTxLookupEntries(block.Hash()); err == nil {
		t.Fatal("expected error for missing receipts, got nil")
	}
}

func TestRestoreTxLookupEntries_Success(t *testing.T) {
	bc := initValSyncOpChain(t)
	tx := newTestTx(0)
	block := types.NewBlock(&types.Header{}, []*types.Transaction{tx}, nil, newTxTestHasher())
	rawdb.WriteBlock(bc.db, block)

	receipts := types.Receipts{{
		Status:            types.ReceiptStatusSuccessful,
		CumulativeGasUsed: 21000,
	}}
	rawdb.WriteReceipts(bc.db, block.Hash(), receipts)
	// No TxLookup entries have been written → simulates a purged state.

	if err := bc.RestoreTxLookupEntries(block.Hash()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := bc.GetTxBlockHash(tx.Hash()); got != block.Hash() {
		t.Errorf("GetTxBlockHash after restore: got %v, want %v", got, block.Hash())
	}
}
