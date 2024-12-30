package fixValidatorsStates

import (
	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	"gitlab.waterfall.network/waterfall/protocol/gwat/log"
	"gitlab.waterfall.network/waterfall/protocol/gwat/params"
)

type blockChain interface {
	GetSlotInfo() *types.SlotInfo
	Config() *params.ChainConfig
	Genesis() *types.Block
}

type FixOp struct {
	OpType     types.ValidatorSyncOp
	Index      uint64
	Creator    common.Address
	InitTxHash common.Hash
	PubKey     common.BlsPubKey
}

func (f FixOp) CreateValidatorSync(procEpoch uint64) *types.ValidatorSync {
	return &types.ValidatorSync{
		InitTxHash: f.InitTxHash,
		OpType:     f.OpType,
		ProcEpoch:  procEpoch,
		Index:      f.Index,
		Creator:    f.Creator,
		Amount:     nil,
		TxHash:     nil,
		Balance:    nil,
	}
}

func GetFixValidatorSyncOps(bc blockChain, currEpoch uint64) []*types.ValidatorSync {
	si := bc.GetSlotInfo()
	bcConf := bc.Config()
	forkEpoch := si.SlotToEpoch(bcConf.ForkSlotValSyncProc)
	procEpoch := forkEpoch + 4

	if currEpoch < forkEpoch || currEpoch >= procEpoch {
		return []*types.ValidatorSync{}
	}

	var (
		ops []*FixOp
		res []*types.ValidatorSync
	)
	if bc.Genesis().Hash() == params.MainnetGenesisHash {
		ops = mainnetFixData
	}

	res = make([]*types.ValidatorSync, len(ops))
	for i, op := range ops {
		res[i] = op.CreateValidatorSync(procEpoch)
	}
	log.Info("Fix validator sync: add sync ops",
		"currEpoch", currEpoch,
		"forkEpoch", forkEpoch,
		"procEpoch", procEpoch,
		"isMainnet", bc.Genesis().Hash() == params.MainnetGenesisHash,
	)
	return res
}
