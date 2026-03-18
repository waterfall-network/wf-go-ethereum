package core

import (
	"math"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/rawdb"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/vm"
	"gitlab.waterfall.network/waterfall/protocol/gwat/node"
	"gitlab.waterfall.network/waterfall/protocol/gwat/params"
	"gitlab.waterfall.network/waterfall/protocol/gwat/tests/testutils"
	validatorOp "gitlab.waterfall.network/waterfall/protocol/gwat/validator/operation"
)

func newChainWithForkSlotValSyncProc(t *testing.T, forkSlot uint64) *BlockChain {
	t.Helper()
	cfg := *params.TestChainConfig
	cfg.ForkSlotValSyncProc = forkSlot
	depositData := make(DepositData, 0)
	for i := 0; i < 64; i++ {
		depositData = append(depositData, &ValidatorData{
			Pubkey:            common.HexToBlsPubKey(strconv.Itoa(i)).String(),
			CreatorAddress:    common.HexToAddress(strconv.Itoa(i)).String(),
			WithdrawalAddress: common.HexToAddress(strconv.Itoa(i)).String(),
			Amount:            32000,
		})
	}
	db := rawdb.NewMemoryDatabase()
	genesis := &Genesis{Config: &cfg, Validators: depositData}
	genesis.MustCommit(db)
	bc, err := NewBlockChain(db, nil, &cfg, vm.Config{}, nil, &node.VerifiersKeystoreConfig{})
	testutils.AssertNoError(t, err)
	return bc
}

func TestCheckValidatorOpWithdrawalFromValState(t *testing.T) {
	txData, err := validatorOp.EncodeToBytes(validatorOp.NewWithdrawalFromValStateOperation())
	testutils.AssertNoError(t, err)

	slotInfo := &types.SlotInfo{
		GenesisTime:    0,
		SecondsPerSlot: 1,
		SlotsPerEpoch:  32,
	}

	testCases := []struct {
		name           string
		bc             *BlockChain
		setSlotInfo    bool
		wantErr        error
		wantErrContain string
	}{
		{
			name:        "no slot info",
			bc:          newChainWithForkSlotValSyncProc(t, 0),
			setSlotInfo: false,
			wantErr:     ErrBadSlotInfo,
			// slotInfo is reset to nil in the test loop below
		},
		{
			name:           "before fork",
			bc:             newChainWithForkSlotValSyncProc(t, math.MaxUint64),
			setSlotInfo:    true,
			wantErrContain: "current fork does not support withdrawal from validators state address",
		},
		{
			name:        "after fork",
			bc:          newChainWithForkSlotValSyncProc(t, 0),
			setSlotInfo: true,
			wantErr:     nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setSlotInfo {
				testutils.AssertNoError(t, tc.bc.SetSlotInfo(slotInfo))
			} else {
				tc.bc.slotInfo = nil
			}
			got := tc.bc.CheckValidatorOp(txData, common.Address{}, nil)
			switch {
			case tc.wantErr != nil:
				testutils.AssertError(t, got, tc.wantErr)
			case tc.wantErrContain != "":
				if got == nil || !strings.Contains(got.Error(), tc.wantErrContain) {
					t.Fatalf("\n\tExpect error containing:\t%q\n\tGot:\t%v", tc.wantErrContain, got)
				}
			default:
				testutils.AssertNoError(t, got)
			}
		})
	}
}

func initValSyncOpChain(t *testing.T) (bc *BlockChain) {
	db := rawdb.NewMemoryDatabase()
	depositData := make(DepositData, 0)
	for i := 0; i < 64; i++ {
		valData := &ValidatorData{
			Pubkey:            common.HexToBlsPubKey(strconv.Itoa(i)).String(),
			CreatorAddress:    common.HexToAddress(strconv.Itoa(i)).String(),
			WithdrawalAddress: common.HexToAddress(strconv.Itoa(i)).String(),
			Amount:            32000,
		}
		depositData = append(depositData, valData)
	}
	genesis := &Genesis{
		Config:     params.TestChainConfig,
		Validators: depositData,
	}

	genesis.MustCommit(db)

	bc, err := NewBlockChain(db, nil, params.TestChainConfig, vm.Config{}, nil, &node.VerifiersKeystoreConfig{})
	testutils.AssertNoError(t, err)
	return bc
}

func TestAppendNotProcessedValidatorSyncData(t *testing.T) {
	type TestCase struct {
		name             string
		bc               *BlockChain
		validatorSyncOps []*types.ValidatorSync
		beforeTest       func(t *testing.T, tc *TestCase)
		testResult       func(t *testing.T, tc *TestCase)
	}

	testCases := []TestCase{
		{
			name: "Appended valSyncOp no tx",
			bc:   initValSyncOpChain(t),
			validatorSyncOps: []*types.ValidatorSync{
				{
					InitTxHash: *testutils.RandomHash(),
					OpType:     types.Activate,
					ProcEpoch:  1,
					Index:      2,
					Creator:    common.HexToAddress(strconv.Itoa(2)),
					Amount:     nil,
					TxHash:     nil,
					Balance:    nil,
				},
				{
					InitTxHash: *testutils.RandomHash(),
					OpType:     types.Deactivate,
					ProcEpoch:  1,
					Index:      3,
					Creator:    common.HexToAddress(strconv.Itoa(3)),
					Amount:     nil,
					TxHash:     nil,
					Balance:    nil,
				},
				{
					InitTxHash: *testutils.RandomHash(),
					OpType:     types.UpdateBalance,
					ProcEpoch:  1,
					Index:      1,
					Creator:    common.HexToAddress(strconv.Itoa(1)),
					Amount:     big.NewInt(0).SetUint64(1_000 * 1e9),
					TxHash:     nil,
					Balance:    big.NewInt(0).SetUint64(32_000 * 1e9),
				},
			},
			beforeTest: func(t *testing.T, tc *TestCase) {
			},
			testResult: func(t *testing.T, tc *TestCase) {
				bc := tc.bc

				// not in `not processed operations` from db
				dbNotProcOps := rawdb.ReadNotProcessedValidatorSyncOps(bc.db)
				testutils.AssertEqual(t, len(tc.validatorSyncOps), len(dbNotProcOps))

				var index int
				var opKey common.Hash

				// item 0
				index = 0
				opKey = tc.validatorSyncOps[index].Key()

				// exists in  `not processed operations` cache
				notProcOps := bc.GetNotProcessedValidatorSyncData()
				testutils.AssertEqual(t, tc.validatorSyncOps[index], notProcOps[opKey])

				// not exists in bc.valSyncCache
				cachedOp, ok := bc.valSyncCache.Get(opKey)
				testutils.AssertEqual(t, false, ok)
				testutils.AssertNil(t, cachedOp)

				// exists in db
				dbOp := rawdb.ReadValidatorSync(bc.db, tc.validatorSyncOps[index].InitTxHash)
				testutils.AssertEqual(t, tc.validatorSyncOps[index], dbOp)

				// item 1
				index = 1
				opKey = tc.validatorSyncOps[index].Key()

				// exists in  `not processed operations` cache
				notProcOps = bc.GetNotProcessedValidatorSyncData()
				testutils.AssertEqual(t, tc.validatorSyncOps[index], notProcOps[opKey])

				// not exists in bc.valSyncCache
				cachedOp, ok = bc.valSyncCache.Get(opKey)
				testutils.AssertEqual(t, false, ok)
				testutils.AssertNil(t, cachedOp)

				// exists in db
				dbOp = rawdb.ReadValidatorSync(bc.db, tc.validatorSyncOps[index].InitTxHash)
				testutils.AssertEqual(t, tc.validatorSyncOps[index], dbOp)

				// item 2
				index = 2
				opKey = tc.validatorSyncOps[index].Key()

				// exists in  `not processed operations` cache
				notProcOps = bc.GetNotProcessedValidatorSyncData()
				testutils.AssertEqual(t, tc.validatorSyncOps[index], notProcOps[opKey])

				// not exists in bc.valSyncCache
				cachedOp, ok = bc.valSyncCache.Get(opKey)
				testutils.AssertEqual(t, false, ok)
				testutils.AssertNil(t, cachedOp)

				// exists in db
				dbOp = rawdb.ReadValidatorSync(bc.db, tc.validatorSyncOps[index].InitTxHash)
				testutils.AssertEqual(t, tc.validatorSyncOps[index], dbOp)
			},
		},

		{
			name: "Appended valSyncOp with tx",
			bc:   initValSyncOpChain(t),
			validatorSyncOps: []*types.ValidatorSync{
				{
					InitTxHash: *testutils.RandomHash(),
					OpType:     types.Activate,
					ProcEpoch:  1,
					Index:      2,
					Creator:    common.HexToAddress(strconv.Itoa(2)),
					Amount:     nil,
					TxHash:     testutils.RandomHash(),
					Balance:    nil,
				},
				{
					InitTxHash: *testutils.RandomHash(),
					OpType:     types.Deactivate,
					ProcEpoch:  1,
					Index:      3,
					Creator:    common.HexToAddress(strconv.Itoa(3)),
					Amount:     nil,
					TxHash:     testutils.RandomHash(),
					Balance:    nil,
				},
				{
					InitTxHash: *testutils.RandomHash(),
					OpType:     types.UpdateBalance,
					ProcEpoch:  1,
					Index:      1,
					Creator:    common.HexToAddress(strconv.Itoa(1)),
					Amount:     big.NewInt(0).SetUint64(1_000 * 1e9),
					TxHash:     testutils.RandomHash(),
					Balance:    big.NewInt(0).SetUint64(32_000 * 1e9),
				},
			},
			beforeTest: func(t *testing.T, tc *TestCase) {
			},
			testResult: func(t *testing.T, tc *TestCase) {
				bc := tc.bc

				// not in `not processed operations` from db
				dbNotProcOps := rawdb.ReadNotProcessedValidatorSyncOps(bc.db)
				testutils.AssertEqual(t, len(tc.validatorSyncOps), len(dbNotProcOps))

				var index int
				var opKey common.Hash

				// item 0
				index = 0
				opKey = tc.validatorSyncOps[index].Key()

				// exists in  `not processed operations` cache
				notProcOps := bc.GetNotProcessedValidatorSyncData()
				testutils.AssertEqual(t, tc.validatorSyncOps[index], notProcOps[opKey])

				// not exists in bc.valSyncCache
				cachedOp, ok := bc.valSyncCache.Get(opKey)
				testutils.AssertEqual(t, false, ok)
				testutils.AssertNil(t, cachedOp)

				// exists in db
				dbOp := rawdb.ReadValidatorSync(bc.db, tc.validatorSyncOps[index].InitTxHash)
				testutils.AssertEqual(t, tc.validatorSyncOps[index], dbOp)

				// item 1
				index = 1
				opKey = tc.validatorSyncOps[index].Key()

				// exists in  `not processed operations` cache
				notProcOps = bc.GetNotProcessedValidatorSyncData()
				testutils.AssertEqual(t, tc.validatorSyncOps[index], notProcOps[opKey])

				// not exists in bc.valSyncCache
				cachedOp, ok = bc.valSyncCache.Get(opKey)
				testutils.AssertEqual(t, false, ok)
				testutils.AssertNil(t, cachedOp)

				// exists in db
				dbOp = rawdb.ReadValidatorSync(bc.db, tc.validatorSyncOps[index].InitTxHash)
				testutils.AssertEqual(t, tc.validatorSyncOps[index], dbOp)

				// item 2
				index = 2
				opKey = tc.validatorSyncOps[index].Key()

				// exists in  `not processed operations` cache
				notProcOps = bc.GetNotProcessedValidatorSyncData()
				testutils.AssertEqual(t, tc.validatorSyncOps[index], notProcOps[opKey])

				// not exists in bc.valSyncCache
				cachedOp, ok = bc.valSyncCache.Get(opKey)
				testutils.AssertEqual(t, false, ok)
				testutils.AssertNil(t, cachedOp)

				// exists in db
				dbOp = rawdb.ReadValidatorSync(bc.db, tc.validatorSyncOps[index].InitTxHash)
				testutils.AssertEqual(t, tc.validatorSyncOps[index], dbOp)
			},
		},

		{
			name: "Appended valSyncOp with tx and SetValidatorSyncData",
			bc:   initValSyncOpChain(t),
			validatorSyncOps: []*types.ValidatorSync{
				{
					InitTxHash: *testutils.RandomHash(),
					OpType:     types.Activate,
					ProcEpoch:  1,
					Index:      2,
					Creator:    common.HexToAddress(strconv.Itoa(2)),
					Amount:     nil,
					TxHash:     testutils.RandomHash(),
					Balance:    nil,
				},
				{
					InitTxHash: *testutils.RandomHash(),
					OpType:     types.Deactivate,
					ProcEpoch:  1,
					Index:      3,
					Creator:    common.HexToAddress(strconv.Itoa(3)),
					Amount:     nil,
					TxHash:     testutils.RandomHash(),
					Balance:    nil,
				},
				{
					InitTxHash: *testutils.RandomHash(),
					OpType:     types.UpdateBalance,
					ProcEpoch:  1,
					Index:      1,
					Creator:    common.HexToAddress(strconv.Itoa(1)),
					Amount:     big.NewInt(0).SetUint64(1_000 * 1e9),
					TxHash:     testutils.RandomHash(),
					Balance:    big.NewInt(0).SetUint64(32_000 * 1e9),
				},
			},
			beforeTest: func(t *testing.T, tc *TestCase) {
				for _, op := range tc.validatorSyncOps {
					tc.bc.SetValidatorSyncData(op)
				}
			},
			testResult: func(t *testing.T, tc *TestCase) {
				bc := tc.bc

				// not in `not processed operations` from db
				dbNotProcOps := rawdb.ReadNotProcessedValidatorSyncOps(bc.db)
				testutils.AssertEqual(t, 0, len(dbNotProcOps))

				var index int
				var opKey common.Hash

				// item 0
				index = 0
				opKey = tc.validatorSyncOps[index].Key()

				// not exists in  `not processed operations` cache
				notProcOps := bc.GetNotProcessedValidatorSyncData()
				testutils.AssertNil(t, notProcOps[opKey])

				// exists in bc.valSyncCache
				cachedOp, ok := bc.valSyncCache.Get(opKey)
				testutils.AssertEqual(t, true, ok)
				testutils.AssertEqual(t, tc.validatorSyncOps[index], cachedOp)

				// exists in db
				dbOp := rawdb.ReadValidatorSync(bc.db, tc.validatorSyncOps[index].InitTxHash)
				testutils.AssertEqual(t, tc.validatorSyncOps[index], dbOp)

				// item 1
				index = 1
				opKey = tc.validatorSyncOps[index].Key()

				// not exists in  `not processed operations` cache
				notProcOps = bc.GetNotProcessedValidatorSyncData()
				testutils.AssertNil(t, notProcOps[opKey])

				// exists in bc.valSyncCache
				cachedOp, ok = bc.valSyncCache.Get(opKey)
				testutils.AssertEqual(t, true, ok)
				testutils.AssertEqual(t, tc.validatorSyncOps[index], cachedOp)

				// exists in db
				dbOp = rawdb.ReadValidatorSync(bc.db, tc.validatorSyncOps[index].InitTxHash)
				testutils.AssertEqual(t, tc.validatorSyncOps[index], dbOp)

				// item 2
				index = 2
				opKey = tc.validatorSyncOps[index].Key()

				// not exists in  `not processed operations` cache
				notProcOps = bc.GetNotProcessedValidatorSyncData()
				testutils.AssertNil(t, notProcOps[opKey])

				// exists in bc.valSyncCache
				cachedOp, ok = bc.valSyncCache.Get(opKey)
				testutils.AssertEqual(t, true, ok)
				testutils.AssertEqual(t, tc.validatorSyncOps[index], cachedOp)

				// exists in db
				dbOp = rawdb.ReadValidatorSync(bc.db, tc.validatorSyncOps[index].InitTxHash)
				testutils.AssertEqual(t, tc.validatorSyncOps[index], dbOp)
			},
		},

		{
			name: "Appended valSyncOp with tx and SetValidatorSyncData no tx",
			bc:   initValSyncOpChain(t),
			validatorSyncOps: []*types.ValidatorSync{
				{
					InitTxHash: *testutils.RandomHash(),
					OpType:     types.Activate,
					ProcEpoch:  1,
					Index:      2,
					Creator:    common.HexToAddress(strconv.Itoa(2)),
					Amount:     nil,
					TxHash:     testutils.RandomHash(),
					Balance:    nil,
				},
				{
					InitTxHash: *testutils.RandomHash(),
					OpType:     types.Deactivate,
					ProcEpoch:  1,
					Index:      3,
					Creator:    common.HexToAddress(strconv.Itoa(3)),
					Amount:     nil,
					TxHash:     testutils.RandomHash(),
					Balance:    nil,
				},
				{
					InitTxHash: *testutils.RandomHash(),
					OpType:     types.UpdateBalance,
					ProcEpoch:  1,
					Index:      1,
					Creator:    common.HexToAddress(strconv.Itoa(1)),
					Amount:     big.NewInt(0).SetUint64(1_000 * 1e9),
					TxHash:     testutils.RandomHash(),
					Balance:    big.NewInt(0).SetUint64(32_000 * 1e9),
				},
			},
			beforeTest: func(t *testing.T, tc *TestCase) {
				for _, op := range tc.validatorSyncOps {
					cpy := op.Copy()
					cpy.TxHash = nil
					tc.bc.SetValidatorSyncData(op)
				}
			},
			testResult: func(t *testing.T, tc *TestCase) {
				bc := tc.bc

				// not in `not processed operations` from db
				dbNotProcOps := rawdb.ReadNotProcessedValidatorSyncOps(bc.db)
				testutils.AssertEqual(t, 0, len(dbNotProcOps))

				var index int
				var opKey common.Hash

				// item 0
				index = 0
				opKey = tc.validatorSyncOps[index].Key()

				// not exists in  `not processed operations` cache
				notProcOps := bc.GetNotProcessedValidatorSyncData()
				testutils.AssertNil(t, notProcOps[opKey])

				// exists in bc.valSyncCache
				cachedOp, ok := bc.valSyncCache.Get(opKey)
				testutils.AssertEqual(t, true, ok)
				testutils.AssertEqual(t, tc.validatorSyncOps[index], cachedOp)

				// exists in db
				dbOp := rawdb.ReadValidatorSync(bc.db, tc.validatorSyncOps[index].InitTxHash)
				testutils.AssertEqual(t, tc.validatorSyncOps[index], dbOp)

				// item 1
				index = 1
				opKey = tc.validatorSyncOps[index].Key()

				// not exists in  `not processed operations` cache
				notProcOps = bc.GetNotProcessedValidatorSyncData()
				testutils.AssertNil(t, notProcOps[opKey])

				// exists in bc.valSyncCache
				cachedOp, ok = bc.valSyncCache.Get(opKey)
				testutils.AssertEqual(t, true, ok)
				testutils.AssertEqual(t, tc.validatorSyncOps[index], cachedOp)

				// exists in db
				dbOp = rawdb.ReadValidatorSync(bc.db, tc.validatorSyncOps[index].InitTxHash)
				testutils.AssertEqual(t, tc.validatorSyncOps[index], dbOp)

				// item 2
				index = 2
				opKey = tc.validatorSyncOps[index].Key()

				// not exists in  `not processed operations` cache
				notProcOps = bc.GetNotProcessedValidatorSyncData()
				testutils.AssertNil(t, notProcOps[opKey])

				// exists in bc.valSyncCache
				cachedOp, ok = bc.valSyncCache.Get(opKey)
				testutils.AssertEqual(t, true, ok)
				testutils.AssertEqual(t, tc.validatorSyncOps[index], cachedOp)

				// exists in db
				dbOp = rawdb.ReadValidatorSync(bc.db, tc.validatorSyncOps[index].InitTxHash)
				testutils.AssertEqual(t, tc.validatorSyncOps[index], dbOp)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.beforeTest(t, &tc)
			tc.bc.AppendNotProcessedValidatorSyncData(tc.validatorSyncOps)
			tc.testResult(t, &tc)
		})
	}
}

func TestSetValidatorSyncData(t *testing.T) {
	type TestCase struct {
		name          string
		bc            *BlockChain
		validatorSync *types.ValidatorSync
		beforeTest    func(t *testing.T, tc *TestCase)
		testResult    func(t *testing.T, tc *TestCase)
	}

	testCases := []TestCase{
		{
			name: "Set not appended valSyncOp no tx",
			bc:   initValSyncOpChain(t),
			validatorSync: &types.ValidatorSync{
				InitTxHash: *testutils.RandomHash(),
				OpType:     types.UpdateBalance,
				ProcEpoch:  1,
				Index:      1,
				Creator:    common.HexToAddress(strconv.Itoa(1)),
				Amount:     big.NewInt(0).SetUint64(1_000 * 1e9),
				TxHash:     nil,
				Balance:    big.NewInt(0).SetUint64(32_000 * 1e9),
			},
			beforeTest: func(t *testing.T, tc *TestCase) {
			},
			testResult: func(t *testing.T, tc *TestCase) {
				bc := tc.bc
				opKey := tc.validatorSync.Key()

				// not in cached `not processed operations`
				notProcOps := bc.GetNotProcessedValidatorSyncData()
				testutils.AssertNil(t, notProcOps[opKey])

				// not in `not processed operations` from db
				dbNotProcOps := rawdb.ReadNotProcessedValidatorSyncOps(bc.db)
				opCount := 0
				for _, op := range dbNotProcOps {
					if op.InitTxHash == tc.validatorSync.InitTxHash {
						opCount++
					}
				}
				testutils.AssertEqual(t, 0, opCount)

				// exists in bc.valSyncCache
				cachedOp, ok := bc.valSyncCache.Get(opKey)
				testutils.AssertEqual(t, true, ok)
				testutils.AssertEqual(t, tc.validatorSync, cachedOp)

				// exists in db
				dbOp := rawdb.ReadValidatorSync(bc.db, tc.validatorSync.InitTxHash)
				testutils.AssertEqual(t, tc.validatorSync, dbOp)
			},
		},
		{
			name: "Set not appended valSyncOp with tx",
			bc:   initValSyncOpChain(t),
			validatorSync: &types.ValidatorSync{
				InitTxHash: *testutils.RandomHash(),
				OpType:     types.UpdateBalance,
				ProcEpoch:  1,
				Index:      1,
				Creator:    common.HexToAddress(strconv.Itoa(1)),
				Amount:     big.NewInt(0).SetUint64(1_000 * 1e9),
				TxHash:     testutils.RandomHash(),
				Balance:    big.NewInt(0).SetUint64(32_000 * 1e9),
			},
			beforeTest: func(t *testing.T, tc *TestCase) {
			},
			testResult: func(t *testing.T, tc *TestCase) {
				bc := tc.bc
				opKey := tc.validatorSync.Key()

				// not in cached `not processed operations`
				notProcOps := bc.GetNotProcessedValidatorSyncData()
				testutils.AssertNil(t, notProcOps[opKey])

				// not in `not processed operations` from db
				dbNotProcOps := rawdb.ReadNotProcessedValidatorSyncOps(bc.db)
				opCount := 0
				for _, op := range dbNotProcOps {
					if op.InitTxHash == tc.validatorSync.InitTxHash {
						opCount++
					}
				}
				testutils.AssertEqual(t, 0, opCount)

				// exists in bc.valSyncCache
				cachedOp, ok := bc.valSyncCache.Get(opKey)
				testutils.AssertEqual(t, true, ok)
				testutils.AssertEqual(t, tc.validatorSync, cachedOp)

				// exists in db
				dbOp := rawdb.ReadValidatorSync(bc.db, tc.validatorSync.InitTxHash)
				testutils.AssertEqual(t, tc.validatorSync, dbOp)
			},
		},
		{
			name: "Set appended valSyncOp no tx",
			bc:   initValSyncOpChain(t),
			validatorSync: &types.ValidatorSync{
				InitTxHash: *testutils.RandomHash(),
				OpType:     types.UpdateBalance,
				ProcEpoch:  1,
				Index:      1,
				Creator:    common.HexToAddress(strconv.Itoa(1)),
				Amount:     big.NewInt(0).SetUint64(1_000 * 1e9),
				TxHash:     nil,
				Balance:    big.NewInt(0).SetUint64(32_000 * 1e9),
			},
			beforeTest: func(t *testing.T, tc *TestCase) {
				bc := tc.bc

				//append valSyncOp
				cpy := tc.validatorSync.Copy()
				cpy.TxHash = nil
				bc.AppendNotProcessedValidatorSyncData([]*types.ValidatorSync{cpy})
			},
			testResult: func(t *testing.T, tc *TestCase) {
				bc := tc.bc
				opKey := tc.validatorSync.Key()

				// not in cached `not processed operations`
				notProcOps := bc.GetNotProcessedValidatorSyncData()
				testutils.AssertEqual(t, tc.validatorSync, notProcOps[opKey])

				// not in `not processed operations` from db
				dbNotProcOps := rawdb.ReadNotProcessedValidatorSyncOps(bc.db)
				opCount := 0
				for _, op := range dbNotProcOps {
					if op.InitTxHash == tc.validatorSync.InitTxHash {
						opCount++
					}
				}
				testutils.AssertEqual(t, 1, opCount)

				// exists in bc.valSyncCache
				cachedOp, ok := bc.valSyncCache.Get(opKey)
				testutils.AssertEqual(t, true, ok)
				testutils.AssertEqual(t, tc.validatorSync, cachedOp)

				// exists in db
				dbOp := rawdb.ReadValidatorSync(bc.db, tc.validatorSync.InitTxHash)
				testutils.AssertEqual(t, tc.validatorSync, dbOp)
			},
		},
		{
			name: "Set appended valSyncOp with tx",
			bc:   initValSyncOpChain(t),
			validatorSync: &types.ValidatorSync{
				InitTxHash: common.BytesToHash(testutils.RandomData(common.HashLength)),
				OpType:     types.UpdateBalance,
				ProcEpoch:  1,
				Index:      1,
				Creator:    common.HexToAddress(strconv.Itoa(1)),
				Amount:     big.NewInt(0).SetUint64(1_000 * 1e9),
				//
				TxHash:  testutils.RandomHash(),
				Balance: big.NewInt(0).SetUint64(32_000 * 1e9),
			},
			beforeTest: func(t *testing.T, tc *TestCase) {
				bc := tc.bc

				//append valSyncOp
				cpy := tc.validatorSync.Copy()
				cpy.TxHash = nil
				bc.AppendNotProcessedValidatorSyncData([]*types.ValidatorSync{cpy})

				//	check current bc state
				opKey := tc.validatorSync.Key()

				// not in `not processed operations` cache
				notProcOps := bc.GetNotProcessedValidatorSyncData()
				testutils.AssertNotNil(t, notProcOps[opKey])

				// not in `not processed operations` from db
				dbNotProcOps := rawdb.ReadNotProcessedValidatorSyncOps(bc.db)
				opCount := 0
				for _, op := range dbNotProcOps {
					if op.InitTxHash == tc.validatorSync.InitTxHash {
						opCount++
					}
				}
				testutils.AssertEqual(t, 1, opCount)

				// exists in bc.valSyncCache
				cachedOp, ok := bc.valSyncCache.Get(opKey)
				testutils.AssertEqual(t, false, ok)
				testutils.AssertNil(t, cachedOp)

				// exists in db
				dbOp := rawdb.ReadValidatorSync(bc.db, tc.validatorSync.InitTxHash)
				testutils.AssertNotNil(t, dbOp)

			},
			testResult: func(t *testing.T, tc *TestCase) {
				bc := tc.bc
				opKey := tc.validatorSync.Key()

				// not in cached `not processed operations`
				notProcOps := bc.GetNotProcessedValidatorSyncData()
				testutils.AssertNil(t, notProcOps[opKey])

				// not in `not processed operations` from db
				dbNotProcOps := rawdb.ReadNotProcessedValidatorSyncOps(bc.db)
				opCount := 0
				for _, op := range dbNotProcOps {
					if op.InitTxHash == tc.validatorSync.InitTxHash {
						opCount++
					}
				}
				testutils.AssertEqual(t, 0, opCount)

				// exists in bc.valSyncCache
				cachedOp, ok := bc.valSyncCache.Get(opKey)
				testutils.AssertEqual(t, true, ok)
				testutils.AssertEqual(t, tc.validatorSync, cachedOp)

				// exists in db
				dbOp := rawdb.ReadValidatorSync(bc.db, tc.validatorSync.InitTxHash)
				testutils.AssertEqual(t, tc.validatorSync, dbOp)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.beforeTest(t, &tc)
			tc.bc.SetValidatorSyncData(tc.validatorSync)
			tc.testResult(t, &tc)
		})
	}
}
