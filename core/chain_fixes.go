package core

import (
	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/rawdb"
	"gitlab.waterfall.network/waterfall/protocol/gwat/log"
	"gitlab.waterfall.network/waterfall/protocol/gwat/params"
	"gitlab.waterfall.network/waterfall/protocol/gwat/validator"
	"gitlab.waterfall.network/waterfall/protocol/gwat/validator/operation"
)

// resetBadChainSequences to call while NewBlockChain to check bad state of a block chain and correct if needed.
func resetBadChainSequences(bc *BlockChain) error {
	if err := fixMainnet0_resetBadChainSequences(bc); err != nil {
		return err
	}
	if err := fixMainnet1_resetBadChainSequences(bc); err != nil {
		return err
	}
	return nil
}

// FixValidatorSyncOpProcessing to call while NewBlockChain to check bad state of a block chain and correct if needed.
func (bc *BlockChain) FixValidatorSyncOpProcessing(processor *validator.Processor, opData operation.Operation, txHash common.Hash, from, to common.Address) (isApplied bool, ret []byte, err error) {
	isApplied, ret, err = fixMainnet0_FixValidatorSyncOpProcessing(bc, processor, opData, txHash, from, to)
	if isApplied {
		return true, ret, err
	}
	isApplied, ret, err = fixMainnet1_FixValidatorSyncOpProcessing(bc, processor, opData, txHash, from, to)
	if isApplied {
		return true, ret, err
	}
	return false, ret, err
}

func isMainnet(bc *BlockChain) bool {
	return bc.Genesis().Hash() == params.MainnetGenesisHash
}

// fixMainnet0_resetBadChainSequences fixes bad checkpoint state (mainntet block nr=3343671).
func fixMainnet0_resetBadChainSequences(bc *BlockChain) error {
	if !isMainnet(bc) {
		return nil
	}
	targetBlockHash := common.HexToHash("0xad0df4045483f44474f51511f7169779afc6cb437ebdc6729605892c4c0b9fb4")
	correctStateRoot := common.HexToHash("0xce0ae2bf519fa8a18be92f902d1fa201937e5868d2d1bf2d7412a9fec1bb07ee")
	targetHeader := bc.GetHeaderByHash(targetBlockHash)
	if targetHeader != nil && targetHeader.Root != correctStateRoot {
		// rollback to acceptable chain's state
		lfHash := rawdb.ReadLastFinalizedHash(bc.db)
		lfNr := rawdb.ReadFinalizedNumberByHash(bc.db, lfHash)
		if lfNr == nil {
			// sync stucking at block nr = 3343862
			*lfNr = uint64(3343862)
		}
		log.Warn("Hard fix mainnet 0: reset finalization of nr=3343671 hash=0xad0df4045483f44474f51511f7169779afc6cb437ebdc6729605892c4c0b9fb4")
		return bc.SetHead(targetHeader.CpHash)
	}
	return nil
}

// fixMainnet0_FixValidatorSyncOpProcessing fixes applying validator sync txs of mainntet block nr=3343671.
func fixMainnet1_FixValidatorSyncOpProcessing(bc *BlockChain, p *validator.Processor, opData operation.Operation, txHash common.Hash, from, to common.Address) (isApplied bool, ret []byte, err error) {
	if !isMainnet(bc) {
		return false, nil, nil
	}
	targetBlockHash := common.HexToHash("0xad0df4045483f44474f51511f7169779afc6cb437ebdc6729605892c4c0b9fb4")
	targetBlockNr := uint64(3343671)

	//Hard fix bad checkpoint state (mainntet block nr=3343671)
	blkCtx := p.GetBlockContext()
	if blkCtx.BlockNumber.Uint64() == targetBlockNr && blkCtx.BlockHash == targetBlockHash {
		log.Warn("Hard fix mainnet 0: process tx",
			"blkNr", blkCtx.BlockNumber,
			"blkHash", blkCtx.BlockHash.Hex(),
			"txHash", txHash.Hex(),
		)
		return true, nil, validator.ErrNoSavedValSyncOp
	}
	return false, nil, nil
}

// fixMainnet1_resetBadChainSequences fixes bad checkpoint state (mainntet block nr=3343937).
func fixMainnet1_resetBadChainSequences(bc *BlockChain) error {
	if !isMainnet(bc) {
		return nil
	}
	targetBlockHash := common.HexToHash("0xc8ab6c76d93ae2dcb8a54bb9cefdc8ecdd7af04996965532d1db0428603d9ea0")
	correctStateRoot := common.HexToHash("0x13357d3bea17190f0a3baf4afb7c61aac224e88d261c9932b9db42afc9ffd808")
	targetHeader := bc.GetHeaderByHash(targetBlockHash)
	if targetHeader != nil && targetHeader.Root != correctStateRoot {
		// rollback to acceptable chain's state
		lfHash := rawdb.ReadLastFinalizedHash(bc.db)
		lfNr := rawdb.ReadFinalizedNumberByHash(bc.db, lfHash)
		if lfNr == nil {
			// sync stucking at block nr = 3344081
			*lfNr = uint64(3344081)
		}
		log.Warn("Hard fix mainnet 1: reset finalization of nr=3343937 hash=0xc8ab6c76d93ae2dcb8a54bb9cefdc8ecdd7af04996965532d1db0428603d9ea0")
		return bc.SetHead(targetHeader.CpHash)
	}
	return nil
}

// fixMainnet0_FixValidatorSyncOpProcessing fixes applying validator sync txs of mainntet block nr=3343937.
func fixMainnet0_FixValidatorSyncOpProcessing(bc *BlockChain, p *validator.Processor, opData operation.Operation, txHash common.Hash, from, to common.Address) (isApplied bool, ret []byte, err error) {
	if !isMainnet(bc) {
		return false, nil, nil
	}
	targetBlockHash := common.HexToHash("0xc8ab6c76d93ae2dcb8a54bb9cefdc8ecdd7af04996965532d1db0428603d9ea0")
	targetBlockNr := uint64(3343937)

	//Hard fix bad checkpoint state (mainntet block nr=3343937)
	blkCtx := p.GetBlockContext()
	if blkCtx.BlockNumber.Uint64() == targetBlockNr && blkCtx.BlockHash == targetBlockHash {
		log.Warn("Hard fix mainnet 2: process tx",
			"blkNr", blkCtx.BlockNumber,
			"blkHash", blkCtx.BlockHash.Hex(),
			"txHash", txHash.Hex(),
		)
		return true, nil, validator.ErrNoSavedValSyncOp
	}
	return false, nil, nil
}
