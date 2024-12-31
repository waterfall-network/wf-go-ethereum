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
	GetBlockByHash(hash common.Hash) *types.Block
	WriteTxLookupEntry(txIndex int, txHash, blockHash common.Hash, receiptStatus uint64) bool
	GetTxBlockHash(txHash common.Hash) common.Hash
	GetReceiptsByHash(blHash common.Hash) types.Receipts
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
		//check init tx lookup entry
		if op.TxBlock != (common.Hash{}) {
			if txBlock := bc.GetTxBlockHash(op.InitTxHash); txBlock != op.TxBlock {
				restoreLookupEntry(bc, op.InitTxHash)
			}
		}
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

func restoreLookupEntry(bc blockChain, blHash common.Hash) {
	block := bc.GetBlockByHash(blHash)
	receipts := bc.GetReceiptsByHash(blHash)
	for i, tx := range block.Transactions() {
		receipt := receipts[i]
		if receipt == nil {
			log.Error("Fix validator sync: update tx lookup entry: no receipt", "blNr", block.Nr(), "blHash", blHash, "txI", i, "txHash", tx.Hash())
			continue
		}
		log.Info("Fix validator sync: update tx lookup entry", "blNr", block.Nr(), "blHash", blHash, "txI", i, "txHash", tx.Hash())
		bc.WriteTxLookupEntry(i, tx.Hash(), block.Hash(), receipt.Status)
	}
}
