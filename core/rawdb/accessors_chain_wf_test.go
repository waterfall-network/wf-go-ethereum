// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rawdb

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	"gitlab.waterfall.network/waterfall/protocol/gwat/tests/testutils"
)

// deprecated
// Tests that head headers and head blocks can be assigned, individually.
func TestHeadStorageWf(t *testing.T) {
	db := NewMemoryDatabase()

	blockHead := types.NewBlockWithHeader(&types.Header{Extra: []byte("test block header")})
	blockFull := types.NewBlockWithHeader(&types.Header{Extra: []byte("test block full")})
	blockFast := types.NewBlockWithHeader(&types.Header{Extra: []byte("test block fast")})

	// Check that no head entries are in a pristine database
	entry0 := ReadTipsHashes(db)
	testutils.AssertEqual(t, true, entry0.IsEqualTo(common.HashArray{}))
	entry1 := ReadLastCanonicalHash(db)
	testutils.AssertEqual(t, common.Hash{}, entry1)
	entry2 := ReadHeadFastBlockHash(db)
	testutils.AssertEqual(t, common.Hash{}, entry2)

	// Assign separate entries for the head header and block
	WriteTipsHashes(db, common.HashArray{blockHead.Hash()})
	WriteLastCanonicalHash(db, blockFull.Hash())
	WriteHeadFastBlockHash(db, blockFast.Hash())

	// Check that both heads are present, and different (i.e. two heads maintained)
	entry3 := ReadTipsHashes(db)
	testutils.AssertEqual(t, true, entry3.IsEqualTo(common.HashArray{blockHead.Hash()}))
	entry4 := ReadLastCanonicalHash(db)
	testutils.AssertEqual(t, blockFull.Hash(), entry4)
	entry5 := ReadHeadFastBlockHash(db)
	testutils.AssertEqual(t, blockFast.Hash(), entry5)

	WriteHeader(db, blockHead.Header())
	WriteBlock(db, blockHead)
}

// Tests that head headers and head blocks can be assigned, individually.
func TestLastFinalizedBlockWf(t *testing.T) {
	db := NewMemoryDatabase()

	// Check FinalizedHeightByHash
	entry0 := ReadFinalizedNumberByHash(db, common.Hash{})
	testutils.AssertNil(t, entry0)

	finHeight := uint64(111111111)
	finBlock := types.NewBlockWithHeader(&types.Header{Extra: []byte("test FinBlock")})
	_writeFinalizedNumberByHash(db, finBlock.Hash(), finHeight)
	entry1 := ReadFinalizedNumberByHash(db, finBlock.Hash())
	testutils.AssertEqual(t, finHeight, *entry1)

	// Check FinalizedHashByHeight
	entry2 := ReadFinalizedNumberByHash(db, common.Hash{})
	testutils.AssertNil(t, entry2)

	finHeight1 := uint64(252222222222)
	finBlock1 := types.NewBlockWithHeader(&types.Header{Extra: []byte("test FinBlock")})
	_writeFinalizedHashByNumber(db, finHeight1, finBlock1.Hash())
	entry3 := ReadFinalizedHashByNumber(db, finHeight1)
	testutils.AssertEqual(t, finBlock1.Hash(), entry3)

	// Check CpHash
	entry4 := ReadLastFinalizedHash(db)
	testutils.AssertEqual(t, common.Hash{}, entry4)

	lastFinBlock := types.NewBlockWithHeader(&types.Header{Extra: []byte("test lastFinBlock")})
	WriteLastFinalizedHash(db, lastFinBlock.Hash())
	entry5 := ReadLastFinalizedHash(db)
	testutils.AssertEqual(t, lastFinBlock.Hash(), entry5)

	// Check CpHeight WriteFinalizedHashNumber
	entry6 := ReadLastFinalizedNumber(db)
	testutils.AssertEqual(t, uint64(0), entry6)
	lastFinHeight1 := uint64(33333333)
	lastFinBlock1 := types.NewBlockWithHeader(&types.Header{Extra: []byte("test lastFinBlock1")})
	WriteFinalizedHashNumber(db, lastFinBlock1.Hash(), lastFinHeight1)
	WriteLastFinalizedHash(db, lastFinBlock1.Hash())
	entry7 := ReadLastFinalizedNumber(db)
	testutils.AssertEqual(t, lastFinHeight1, entry7)
}

// Tests that head headers and head blocks can be assigned, individually.
func TestBlockDAGWf(t *testing.T) {
	db := NewMemoryDatabase()

	// Check FinalizedHeightByHash
	entry0 := ReadBlockDag(db, common.Hash{})
	testutils.AssertNil(t, entry0)

	finBlock := types.NewBlockWithHeader(&types.Header{Extra: []byte("test FinBlock")})

	blockDag := &types.BlockDAG{
		Hash:                   finBlock.Hash(),
		Height:                 finBlock.Height(),
		CpHash:                 finBlock.Hash(),
		CpHeight:               1455646545646,
		OrderedAncestorsHashes: common.HashArray{common.Hash{}, finBlock.Hash(), common.Hash{}},
	}

	WriteBlockDag(db, blockDag)
	entry1 := ReadBlockDag(db, finBlock.Hash())
	testutils.AssertEqual(t, blockDag, entry1)

	DeleteBlockDag(db, blockDag.Hash)
	entry2 := ReadBlockDag(db, finBlock.Hash())
	testutils.AssertNil(t, entry2)
}

func TestValidatorSyncWf_Ok(t *testing.T) {
	db := NewMemoryDatabase()

	src1 := &types.ValidatorSync{
		OpType:          2,
		ProcEpoch:       45645,
		Index:           45645,
		Creator:         common.Address{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		Amount:          new(big.Int),
		TxHash:          &common.Hash{7, 8, 9},
		InitTxHash:      common.Hash{1, 2, 3},
		ActivationEpoch: 645,
		ExitEpoch:       math.MaxUint64,
	}
	src1.Amount.SetString("32789456000000", 10)

	WriteValidatorSync(db, src1)
	entry := ReadValidatorSync(db, src1.InitTxHash)
	testutils.AssertEqual(t, src1, entry)

	DeleteValidatorSync(db, src1.InitTxHash)
	entry1 := ReadValidatorSync(db, src1.InitTxHash)
	testutils.AssertNil(t, entry1)
}

func TestValidatorSyncWf_deprecatedJsonV0_Ok(t *testing.T) {
	db := NewMemoryDatabase()

	//deprecated json data
	src1 := []byte("{\"initTxHash\":\"0x0102030000000000000000000000000000000000000000000000000000000000\"," +
		"\"opType\":\"0x2\",\"procEpoch\":\"0xb24d\",\"index\":\"0xb24d\"," +
		"\"creator\":\"0xffffffffffffffffffffffffffffffffffffffff\",\"amount\":\"0x1dd263e09400\"," +
		"\"txHash\":null,\"balance\":\"0x3ab5108a71400\"}")

	exp := []byte("{\"initTxHash\":\"0x0102030000000000000000000000000000000000000000000000000000000000\",\"opType\":\"0x2\"," +
		"\"procEpoch\":\"0xb24d\",\"index\":\"0xb24d\",\"creator\":\"0xffffffffffffffffffffffffffffffffffffffff\"," +
		"\"amount\":\"0x1dd263e09400\",\"txHash\":null,\"balance\":\"0x3ab5108a71400\",\"activationEpoch\":\"0x0\",\"exitEpoch\":\"0x0\"}")

	initTx := common.HexToHash("0x0102030000000000000000000000000000000000000000000000000000000000")

	err := db.Put(validatorSyncKey(initTx), src1)
	testutils.AssertNoError(t, err)

	entry := ReadValidatorSync(db, initTx)
	assert.NotNil(t, entry)

	entryJson, _ := entry.MarshalJSON()
	testutils.AssertEqual(t, exp, entryJson)

	DeleteValidatorSync(db, initTx)
	entry = ReadValidatorSync(db, initTx)
	testutils.AssertNil(t, entry)
}

func TestValidatorSyncWf_Ok_noTxHash(t *testing.T) {
	db := NewMemoryDatabase()

	src1 := &types.ValidatorSync{
		OpType:    2,
		ProcEpoch: 45645,
		Index:     45645,
		Creator:   common.Address{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		Amount:    new(big.Int),
		//TxHash:     &common.Hash{7, 8, 9},
		InitTxHash: common.Hash{1, 2, 3},
	}
	src1.Amount.SetString("32789456000000", 10)

	WriteValidatorSync(db, src1)
	entry0 := ReadValidatorSync(db, src1.InitTxHash)
	testutils.AssertEqual(t, src1, entry0)

	DeleteValidatorSync(db, src1.InitTxHash)
	entry1 := ReadValidatorSync(db, src1.InitTxHash)
	testutils.AssertNil(t, entry1)
}

func TestValidatorSyncWf_Ok_noAmount(t *testing.T) {
	db := NewMemoryDatabase()

	src1 := &types.ValidatorSync{
		OpType:    1,
		ProcEpoch: 45645,
		Index:     45645,
		Creator:   common.Address{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		//Amount:     new(big.Int),
		TxHash:     &common.Hash{7, 8, 9},
		InitTxHash: common.Hash{1, 2, 3},
	}
	//src1.Amount.SetString("32789456000000", 10)

	WriteValidatorSync(db, src1)
	entry0 := ReadValidatorSync(db, src1.InitTxHash)
	testutils.AssertEqual(t, src1, entry0)

	DeleteValidatorSync(db, src1.InitTxHash)
	entry1 := ReadValidatorSync(db, src1.InitTxHash)
	testutils.AssertNil(t, entry1)
}

func TestNotProcessedValidatorSyncWf(t *testing.T) {
	db := NewMemoryDatabase()

	src1 := &types.ValidatorSync{
		OpType:     types.Activate,
		ProcEpoch:  45645,
		Index:      45645,
		Creator:    common.Address{0x11, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		Amount:     new(big.Int),
		TxHash:     nil,
		InitTxHash: common.Hash{1, 2, 3},
	}
	src2 := &types.ValidatorSync{
		OpType:     types.Deactivate,
		ProcEpoch:  45645,
		Index:      45645,
		Creator:    common.Address{0x22, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		Amount:     new(big.Int),
		TxHash:     nil,
		InitTxHash: common.Hash{1, 2, 3},
	}
	src3 := &types.ValidatorSync{
		OpType:     types.UpdateBalance,
		ProcEpoch:  45645,
		Index:      45645,
		Creator:    common.Address{0x33, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		Amount:     new(big.Int),
		TxHash:     nil,
		InitTxHash: common.Hash{1, 2, 3},
	}
	src3.Amount.SetString("32789456000000", 10)

	valSyncOps := []*types.ValidatorSync{src1, src2, src3}

	WriteNotProcessedValidatorSyncOps(db, valSyncOps)
	entry := ReadNotProcessedValidatorSyncOps(db)
	testutils.AssertEqual(t, len(valSyncOps), len(entry))
	for _, e := range entry {
		var vsop *types.ValidatorSync
		for _, vs := range valSyncOps {
			if vs.OpType == e.OpType {
				vsop = vs
			}
		}
		testutils.AssertEqual(t, vsop, e)
	}
}
