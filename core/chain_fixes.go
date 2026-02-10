package core

import (
	"fmt"

	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	"gitlab.waterfall.network/waterfall/protocol/gwat/log"
	"gitlab.waterfall.network/waterfall/protocol/gwat/params"
	"gitlab.waterfall.network/waterfall/protocol/gwat/validator"
	"gitlab.waterfall.network/waterfall/protocol/gwat/validator/operation"
)

// FixValidatorSyncOpProcessing to call while NewBlockChain to check bad state of a block chain and correct if needed.
func (bc *BlockChain) FixValidatorSyncOpProcessing(processor *validator.Processor, opData operation.Operation, txHash common.Hash, from, to common.Address) (isApplied bool, ret []byte, err error) {
	isApplied, ret, err = fixMainnet0_FixValidatorSyncOpProcessing(bc, processor, opData, txHash, from, to)
	if isApplied {
		return true, ret, err
	}
	return false, ret, err
}

func isMainnet(bc *BlockChain) bool {
	return bc.Genesis().Hash() == params.MainnetGenesisHash
}

// fixMainnet0_FixValidatorSyncOpProcessing fixes applying validator sync txs of mainnet block nr=3343671.
func fixMainnet0_FixValidatorSyncOpProcessing(bc *BlockChain, p *validator.Processor, opData operation.Operation, txHash common.Hash, from, to common.Address) (isApplied bool, ret []byte, err error) {
	if !isMainnet(bc) {
		return false, nil, nil
	}

	/* validator sync ops that must be failed
	wat.validator.getBlockReceipts(3343671).forEach((v)=>{ console.log(JSON.stringify({
			idx: v.transactionIndex,
			hash: v.transactionHash,
			status: v.status,
			type: v.type,
			from: v.from,
			to: v.to,
			parsedData: v.logs?.[0]?.parsedData
	}, null, 2))})

	opCode=2 initTx=0x924babb237a77f46541703cef25c9713724a53ca924f5d7b5182a7e316ae1ffa creator=0x5040Ea3b763acC52Fa45CF3B9Ff6bD8fBAa50b19 procEpoch=12994
	opCode=2 initTx=0x8b55f917100e95b6af49aeeba970a17e02de60d42e630364f3d8aa644acd8da2 creator=0xCa94C692D4E5D6075C6Eb2d1920d75169D8481AD procEpoch=12995
	opCode=2 initTx=0xf2d0d9db1648d3a05152c9c7fb4a6608fb3a17d2454152341c6caa376da6ab15 creator=0x0927bBBAb3CD3a03bD3996d60558399bFd69c5A2 procEpoch=12996
	opCode=2 initTx=0x00d860db1ef65b539af8761193f7410dcd81db0ff2f13be058c989dfd532c861 creator=0xF0908E8722DCE042c8af29Fc1Ed15C659A406a6B procEpoch=12994
	opCode=2 initTx=0x15b32c3ed82bb234be4af2c1a1450d41e9632d3066c5124a8848b3a00aaf2bea creator=0x32E26f2C4fF439dBFd19e59d60CC4705CB128D0A procEpoch=12996
	opCode=2 initTx=0xfa590cb6c760b34b72caa6d942f91308edee2fc723ae66215e52d6a5695e1f66 creator=0xb7634fbF58438d5F231c6E343ce6a09015b4b122 procEpoch=12994
	opCode=2 initTx=0x07233067bbb076bfc1c3e06a604e570e6bfa858e1d0fc2b1d74a5ccecd9b2b87 creator=0xb0e42065991c5689eD3f570C64210F0174F5Bf80 procEpoch=12995
	opCode=2 initTx=0xdb82ab18dd9472d51f9a38d8d79401333a21c0aa17360cb93c1efbf73cf03fa4 creator=0x1d48f9A83C7328CF7fB5E3c550E1220FeB5c7b13 procEpoch=12994
	opCode=2 initTx=0x83eb2d837463f405ae9f145a8f598363406c6968e685397e4e20d127ecb56138 creator=0xCC5D2a3c5434EeC89B8E87F71E693b35a0A603bb procEpoch=12996
	opCode=2 initTx=0x8b94ce4d80db11999d7cdaa4becd5697144b72d7d6b09a92e955327ab67d1092 creator=0x6473b4C0552E46a80217A8EA619aa55Fa8b5Ad59 procEpoch=12996
	opCode=2 initTx=0x53aff3b86fbcc548eb52e2ba4537602399e1f7e6bc8b214830040b816a542be7 creator=0xf051499E8955FEaF218DEDEB46086Af3F04341E4 procEpoch=12995
	opCode=2 initTx=0x22104335f100ccd555d3a234474608c854bc25d3c1b9886a80a44e38d7f36d3b creator=0x6C30179993F98D02d819520C0AB7D93efC1A464f procEpoch=12995
	*/

	//block: 3343671 0xad0df4045483f44474f51511f7169779afc6cb437ebdc6729605892c4c0b9fb4
	targetBlockNr := uint64(3343671)
	blkCtx := p.GetBlockContext()
	if blkCtx.BlockNumber.Uint64() != targetBlockNr {
		return false, nil, nil
	}

	failedActivationsByInitTx := map[common.Hash]struct{}{
		common.HexToHash("0x924babb237a77f46541703cef25c9713724a53ca924f5d7b5182a7e316ae1ffa"): {},
		common.HexToHash("0x8b55f917100e95b6af49aeeba970a17e02de60d42e630364f3d8aa644acd8da2"): {},
		common.HexToHash("0xf2d0d9db1648d3a05152c9c7fb4a6608fb3a17d2454152341c6caa376da6ab15"): {},
		common.HexToHash("0x00d860db1ef65b539af8761193f7410dcd81db0ff2f13be058c989dfd532c861"): {},
		common.HexToHash("0x15b32c3ed82bb234be4af2c1a1450d41e9632d3066c5124a8848b3a00aaf2bea"): {},
		common.HexToHash("0xfa590cb6c760b34b72caa6d942f91308edee2fc723ae66215e52d6a5695e1f66"): {},
		common.HexToHash("0x07233067bbb076bfc1c3e06a604e570e6bfa858e1d0fc2b1d74a5ccecd9b2b87"): {},
		common.HexToHash("0xdb82ab18dd9472d51f9a38d8d79401333a21c0aa17360cb93c1efbf73cf03fa4"): {},
		common.HexToHash("0x83eb2d837463f405ae9f145a8f598363406c6968e685397e4e20d127ecb56138"): {},
		common.HexToHash("0x8b94ce4d80db11999d7cdaa4becd5697144b72d7d6b09a92e955327ab67d1092"): {},
		common.HexToHash("0x53aff3b86fbcc548eb52e2ba4537602399e1f7e6bc8b214830040b816a542be7"): {},
		common.HexToHash("0x22104335f100ccd555d3a234474608c854bc25d3c1b9886a80a44e38d7f36d3b"): {},
	}

	switch v := opData.(type) {
	case operation.ValidatorSync:
		if _, ok := failedActivationsByInitTx[v.InitTxHash()]; !ok {
			return false, nil, nil
		}

		if txValSyncOp, ok := bc.notProcValSyncOps[v.InitTxHash()]; ok {
			//1. set ValSync as done (to clear from caches)
			txValSyncOp = &types.ValidatorSync{
				InitTxHash: v.InitTxHash(),
				OpType:     v.OpType(),
				ProcEpoch:  v.ProcEpoch(),
				Index:      v.Index(),
				Creator:    v.Creator(),
				Amount:     v.Amount(),
				TxHash:     &common.Hash{},
			}
			bc.SetValidatorSyncData(txValSyncOp)

			log.Info("fixMainnet0_FixValidatorSyncOpProcessing: applied",
				"OpType", txValSyncOp.OpType,
				"ProcEpoch", txValSyncOp.ProcEpoch,
				"Index", txValSyncOp.Index,
				"Creator", fmt.Sprintf("%#x", txValSyncOp.Creator),
				"amount", txValSyncOp.Amount,
				"TxHash", fmt.Sprintf("%#x", txValSyncOp.TxHash),
				"InitTxHash", txValSyncOp.InitTxHash.Hex(),
				"currentTx", fmt.Sprintf("%#x", txHash),
			)
		} else {
			log.Warn("fixMainnet0_FixValidatorSyncOpProcessing: skipping processing",
				"ValSyncOp", nil,
				"InitTxHash", v.InitTxHash().Hex(),
			)
		}
		// action to quickly complete a transaction (not necessary)
		return true, nil, validator.ErrNoSavedValSyncOp
	}
	return false, nil, nil
}
