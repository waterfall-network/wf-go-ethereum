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
	GetTxBlockHash(txHash common.Hash) common.Hash
	RestoreTxLookupEntries(blHash common.Hash) error
}

type FixOp struct {
	OpType     types.ValidatorSyncOp
	Index      uint64
	Creator    common.Address
	InitTxHash common.Hash
	PubKey     common.BlsPubKey
	TxBlock    common.Hash //InitTx block hash
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
	// procEpoch is the epoch at which the fix ops will be executed.
	// During the window [forkEpoch, procEpoch), the ops are repeatedly injected into the
	// not-processed pool on every finalization so that all nodes have time to receive them
	// before execution. Actual processing happens when currEpoch reaches procEpoch.
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

	res = make([]*types.ValidatorSync, 0, len(ops))
	for _, op := range ops {
		//check init tx lookup entry
		if op.TxBlock != (common.Hash{}) {
			if txBlock := bc.GetTxBlockHash(op.InitTxHash); txBlock != op.TxBlock {
				if err := bc.RestoreTxLookupEntries(op.TxBlock); err != nil {
					log.Error("GetFixValidatorSyncOps: restore tx lookup failed",
						"err", err, "txBlock", op.TxBlock, "initTxHash", op.InitTxHash)
					continue
				}
			}
		}
		res = append(res, op.CreateValidatorSync(procEpoch))
	}
	log.Info("Fix validator sync: add sync ops",
		"currEpoch", currEpoch,
		"forkEpoch", forkEpoch,
		"procEpoch", procEpoch,
		"isMainnet", bc.Genesis().Hash() == params.MainnetGenesisHash,
	)
	return res
}
