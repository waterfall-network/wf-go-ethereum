package core

import (
	"strconv"
	"testing"

	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/rawdb"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/vm"
	"gitlab.waterfall.network/waterfall/protocol/gwat/node"
	"gitlab.waterfall.network/waterfall/protocol/gwat/params"
	"gitlab.waterfall.network/waterfall/protocol/gwat/tests/testutils"
	"gitlab.waterfall.network/waterfall/protocol/gwat/validator/era"
)

func TestEpochToEra(t *testing.T) {
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

	genesisBlock := genesis.ToBlock(db)

	bc, err := NewBlockChain(db, nil, params.TestChainConfig, vm.Config{}, nil, &node.VerifiersKeystoreConfig{})
	testutils.AssertNoError(t, err)

	startEra := &era.Era{
		Number:    0,
		From:      0,
		To:        params.TestChainConfig.EpochsPerEra - 1,
		Root:      genesisBlock.Root(),
		BlockHash: genesisBlock.Hash(),
	}

	bc.SetNewEraInfo(startEra)

	err = bc.SetSlotInfo(&types.SlotInfo{
		GenesisTime:    0,
		SecondsPerSlot: 6,
		SlotsPerEpoch:  32,
	})

	testutils.AssertNoError(t, err)

	testCases := []struct {
		name       string
		epoch      uint64
		findingEra *era.Era
	}{
		{
			name:  "Epoch from current era",
			epoch: params.TestChainConfig.EpochsPerEra - 2,
			findingEra: &era.Era{
				Number:    0,
				From:      0,
				To:        params.TestChainConfig.EpochsPerEra - 1,
				Root:      genesisBlock.Root(),
				BlockHash: genesisBlock.Hash(),
			},
		},
		{
			name:  "Epoch from next era",
			epoch: params.TestChainConfig.EpochsPerEra + 1,
			findingEra: &era.Era{
				Number:    1,
				From:      params.TestChainConfig.EpochsPerEra,
				To:        params.TestChainConfig.EpochsPerEra*2 - 1,
				Root:      common.HexToHash("0xcefc25562123831f9d1fbcef7d30f03ba277e1b2c97ef259c69115c46f0fdebe"),
				BlockHash: common.HexToHash("0xcf3214ba22ec4c54637ce5b9bf3a16723a18d7a803f47a675377fbcd785db9ae"),
			},
		},
		{
			name:  "Epoch from previous era",
			epoch: params.TestChainConfig.EpochsPerEra / 2,
			findingEra: &era.Era{
				Number:    0,
				From:      0,
				To:        params.TestChainConfig.EpochsPerEra - 1,
				Root:      genesisBlock.Root(),
				BlockHash: genesisBlock.Hash(),
			},
		},
		{
			name:  "Epoch at boundary of current era",
			epoch: params.TestChainConfig.EpochsPerEra,
			findingEra: &era.Era{
				Number:    1,
				From:      8,
				To:        15,
				Root:      common.HexToHash("0xcefc25562123831f9d1fbcef7d30f03ba277e1b2c97ef259c69115c46f0fdebe"),
				BlockHash: common.HexToHash("0xcf3214ba22ec4c54637ce5b9bf3a16723a18d7a803f47a675377fbcd785db9ae"),
			},
		},
		{
			name:  "Epoch from future era",
			epoch: params.TestChainConfig.EpochsPerEra*2 + 1,
			findingEra: &era.Era{
				Number:    2,
				From:      params.TestChainConfig.EpochsPerEra * 2,
				To:        params.TestChainConfig.EpochsPerEra*3 - 1,
				Root:      common.HexToHash("0xcefc25562123831f9d1fbcef7d30f03ba277e1b2c97ef259c69115c46f0fdebe"),
				BlockHash: common.HexToHash("0xcf3214ba22ec4c54637ce5b9bf3a16723a18d7a803f47a675377fbcd785db9ae"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			epochEra := bc.EpochToEra(tc.epoch)
			testutils.AssertEqual(t, tc.findingEra, epochEra)
			bc.SetNewEraInfo(epochEra)
		})
	}
}

func TestStartTransitionPeriod(t *testing.T) {
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

	genesisBlock := genesis.ToBlock(db)

	bc, err := NewBlockChain(db, nil, params.TestChainConfig, vm.Config{}, nil, &node.VerifiersKeystoreConfig{})
	testutils.AssertNoError(t, err)

	curEra := &era.Era{
		Number:    1,
		From:      0,
		To:        8,
		Root:      genesisBlock.Root(),
		BlockHash: genesisBlock.Hash(),
	}
	rawdb.WriteEra(bc.db, curEra.Number, curEra)
	bc.SetNewEraInfo(curEra)

	err = bc.SetSlotInfo(&types.SlotInfo{
		GenesisTime:    0,
		SecondsPerSlot: 6,
		SlotsPerEpoch:  32,
	})
	testutils.AssertNoError(t, err)

	testCases := []struct {
		name       string
		checkpoint *types.Checkpoint
		eraNumber  uint64
		spineRoot  common.Hash
		spineHash  common.Hash
	}{
		{
			name: "Transition period processed",
			checkpoint: &types.Checkpoint{
				Epoch:    10,
				FinEpoch: 20,
			},
			eraNumber: 2,
			spineRoot: common.HexToHash("0xcefc25562123831f9d1fbcef7d30f03ba277e1b2c97ef259c69115c46f0fdebe"),
			spineHash: common.HexToHash("0xcf3214ba22ec4c54637ce5b9bf3a16723a18d7a803f47a675377fbcd785db9ae"),
		},
		{
			name: "Next era already exists",
			checkpoint: &types.Checkpoint{
				Epoch:    15,
				FinEpoch: 25,
			},
			eraNumber: 2,
			spineRoot: common.HexToHash("0xcefc25562123831f9d1fbcef7d30f03ba277e1b2c97ef259c69115c46f0fdebe"),
			spineHash: common.HexToHash("0xcf3214ba22ec4c54637ce5b9bf3a16723a18d7a803f47a675377fbcd785db9ae"),
		},
		{
			name: "Next era updated",
			checkpoint: &types.Checkpoint{
				Epoch:    15,
				FinEpoch: 25,
			},
			eraNumber: 2,
			spineRoot: common.HexToHash("0xcefc25562123831f9d1fbcef7d30f03ba277e1b2c97ef259c69115c46f0fdebe"),
			spineHash: common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = bc.StartTransitionPeriod(tc.checkpoint, tc.spineRoot, tc.spineHash)
			testutils.AssertNoError(t, err)
			nextEra := rawdb.ReadEra(bc.db, tc.eraNumber)
			testutils.AssertNotNil(t, nextEra)
			testutils.AssertEqual(t, tc.eraNumber, nextEra.Number)
			testutils.AssertEqual(t, tc.spineRoot, nextEra.Root)
			testutils.AssertEqual(t, tc.spineHash, nextEra.BlockHash)
		})
	}
}

func TestHandleEra(t *testing.T) {
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

	genesisBlock := genesis.ToBlock(db)

	bc, err := NewBlockChain(db, nil, params.TestChainConfig, vm.Config{}, nil, &node.VerifiersKeystoreConfig{})
	testutils.AssertNoError(t, err)

	curEra := &era.Era{
		Number:    1,
		From:      0,
		To:        10,
		Root:      genesisBlock.Root(),
		BlockHash: genesisBlock.Hash(),
	}
	rawdb.WriteEra(bc.db, curEra.Number, curEra)
	bc.SetNewEraInfo(curEra)

	err = bc.SetSlotInfo(&types.SlotInfo{
		GenesisTime:    0,
		SecondsPerSlot: 6,
		SlotsPerEpoch:  32,
	})
	testutils.AssertNoError(t, err)

	type TestCases struct {
		name          string
		checkpoint    *types.Checkpoint
		verifyFunc    func(t *testing.T, ts TestCases)
		expectedError error
	}

	testCases := []TestCases{
		{
			name: "Checkpoint without transition period",
			checkpoint: &types.Checkpoint{
				Epoch:    10,
				FinEpoch: 20,
				Spine:    genesisBlock.Hash(),
			},
			verifyFunc: func(t *testing.T, tc TestCases) {
				latestEra := rawdb.ReadEra(db, curEra.Number+1)
				testutils.AssertNotNil(t, latestEra)
			},
		},
		{
			name: "Checkpoint with transition period",
			checkpoint: &types.Checkpoint{
				Epoch:    0,
				FinEpoch: 9,
				Spine:    genesisBlock.Hash(),
			},
			verifyFunc: func(t *testing.T, tc TestCases) {
				latestEra := rawdb.ReadEra(db, curEra.Number+1)
				testutils.AssertNotNil(t, latestEra)
			},
		},
		{
			name: "Invalid checkpoint spine",
			checkpoint: &types.Checkpoint{
				Epoch:    10,
				FinEpoch: 20,
				Spine:    common.HexToHash("0xdeadbeef"),
			},
			verifyFunc: func(t *testing.T, tc TestCases) {
				// Nothing to verify, error expected
			},
			expectedError: errCheckpointInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = bc.HandleEra(tc.checkpoint)
			if tc.expectedError != nil {
				testutils.AssertError(t, err, tc.expectedError)
			} else {
				testutils.AssertNoError(t, err)
			}
			if tc.verifyFunc != nil {
				tc.verifyFunc(t, tc)
			}
		})
	}
}

func TestHandleEra_EraOverride(t *testing.T) {
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

	genesisBlock := genesis.ToBlock(db)

	bc, err := NewBlockChain(db, nil, params.TestChainConfig, vm.Config{}, nil, &node.VerifiersKeystoreConfig{})
	testutils.AssertNoError(t, err)

	curEra := &era.Era{
		Number:    1,
		From:      0,
		To:        10,
		Root:      genesisBlock.Root(),
		BlockHash: genesisBlock.Hash(),
	}
	rawdb.WriteEra(bc.db, curEra.Number, curEra)
	bc.SetNewEraInfo(curEra)

	err = bc.SetSlotInfo(&types.SlotInfo{
		GenesisTime:    0,
		SecondsPerSlot: 6,
		SlotsPerEpoch:  32,
	})
	testutils.AssertNoError(t, err)

	cpBlock := types.NewBlockWithHeader(&types.Header{
		TxHash:      common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111"),
		ReceiptHash: common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111"),
		Root:        genesisBlock.Root(),
	})
	rawdb.WriteBlock(bc.db, cpBlock)
	era3 := &era.Era{
		Number:    2,
		From:      50,
		To:        70,
		Root:      genesisBlock.Root(),
		BlockHash: genesisBlock.Hash(),
	}
	rawdb.WriteEra(bc.db, era3.Number, era3)

	type TestCases struct {
		name          string
		checkpoint    *types.Checkpoint
		verifyFunc    func(t *testing.T, ts TestCases)
		expectedError error
	}

	testCases := []TestCases{
		{
			name: "UPDATED Checkpoint without transition period",
			checkpoint: &types.Checkpoint{
				Epoch:    20,
				FinEpoch: 30,
				Spine:    cpBlock.Hash(),
				Root:     cpBlock.Root(),
			},
			verifyFunc: func(t *testing.T, tc TestCases) {
				latestEra := rawdb.ReadEra(bc.db, 2)
				testutils.AssertNotNil(t, latestEra)
				testutils.AssertEqual(t, tc.checkpoint.Root, latestEra.Root)
				testutils.AssertEqual(t, tc.checkpoint.Spine, latestEra.BlockHash)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = bc.HandleEra(tc.checkpoint)
			if tc.expectedError != nil {
				testutils.AssertError(t, err, tc.expectedError)
			} else {
				testutils.AssertNoError(t, err)
			}
			if tc.verifyFunc != nil {
				tc.verifyFunc(t, tc)
			}
		})
	}
}
