package core

import (
	"strconv"
	"testing"

	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/rawdb"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/vm"
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

	bc, err := NewBlockChain(db, nil, params.TestChainConfig, vm.Config{}, nil)
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			epochEra := bc.EpochToEra(tc.epoch)
			testutils.AssertEqual(t, tc.findingEra, epochEra)
			bc.SetNewEraInfo(epochEra)
		})
	}
}
