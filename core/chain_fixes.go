package core

import (
	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	"gitlab.waterfall.network/waterfall/protocol/gwat/log"
	"gitlab.waterfall.network/waterfall/protocol/gwat/params"
)

// applyFixesOnStart to call while NewBlockChain to .
func applyFixesOnStart(bc *BlockChain) error {
	if err := fixMainnet0_setFailedValSyncOps(bc); err != nil {
		return err
	}
	return nil
}

func isMainnet(bc *BlockChain) bool {
	return bc.Genesis().Hash() == params.MainnetGenesisHash
}

func fixMainnet0_setFailedValSyncOps(bc *BlockChain) error {
	if !isMainnet(bc) {
		return nil
	}

	/*
	   valSyncOp="{InitTxHash: 0x00d860db1ef65b539af8761193f7410dcd81db0ff2f13be058c989dfd532c861, OpType: 0, ProcEpoch: 12994, Index: 4015, Creator: 0xf0908e8722dce042c8af29fc1ed15c659a406a6b, Amount: <nil>, Balance: <nil>, TxHash: <nil>}"
	   valSyncOp="{InitTxHash: 0xdb82ab18dd9472d51f9a38d8d79401333a21c0aa17360cb93c1efbf73cf03fa4, OpType: 0, ProcEpoch: 12994, Index: 4016, Creator: 0x1d48f9a83c7328cf7fb5e3c550e1220feb5c7b13, Amount: <nil>, Balance: <nil>, TxHash: <nil>}"
	   valSyncOp="{InitTxHash: 0x924babb237a77f46541703cef25c9713724a53ca924f5d7b5182a7e316ae1ffa, OpType: 0, ProcEpoch: 12994, Index: 4017, Creator: 0x5040ea3b763acc52fa45cf3b9ff6bd8fbaa50b19, Amount: <nil>, Balance: <nil>, TxHash: <nil>}"
	   valSyncOp="{InitTxHash: 0xfa590cb6c760b34b72caa6d942f91308edee2fc723ae66215e52d6a5695e1f66, OpType: 0, ProcEpoch: 12994, Index: 4018, Creator: 0xb7634fbf58438d5f231c6e343ce6a09015b4b122, Amount: <nil>, Balance: <nil>, TxHash: <nil>}"
	   valSyncOp="{InitTxHash: 0x22104335f100ccd555d3a234474608c854bc25d3c1b9886a80a44e38d7f36d3b, OpType: 0, ProcEpoch: 12995, Index: 4019, Creator: 0x6c30179993f98d02d819520c0ab7d93efc1a464f, Amount: <nil>, Balance: <nil>, TxHash: <nil>}"
	   valSyncOp="{InitTxHash: 0x53aff3b86fbcc548eb52e2ba4537602399e1f7e6bc8b214830040b816a542be7, OpType: 0, ProcEpoch: 12995, Index: 4020, Creator: 0xf051499e8955feaf218dedeb46086af3f04341e4, Amount: <nil>, Balance: <nil>, TxHash: <nil>}"
	   valSyncOp="{InitTxHash: 0x07233067bbb076bfc1c3e06a604e570e6bfa858e1d0fc2b1d74a5ccecd9b2b87, OpType: 0, ProcEpoch: 12995, Index: 4021, Creator: 0xb0e42065991c5689ed3f570c64210f0174f5bf80, Amount: <nil>, Balance: <nil>, TxHash: <nil>}"
	   valSyncOp="{InitTxHash: 0x8b55f917100e95b6af49aeeba970a17e02de60d42e630364f3d8aa644acd8da2, OpType: 0, ProcEpoch: 12995, Index: 4022, Creator: 0xca94c692d4e5d6075c6eb2d1920d75169d8481ad, Amount: <nil>, Balance: <nil>, TxHash: <nil>}"
	   valSyncOp="{InitTxHash: 0x8b94ce4d80db11999d7cdaa4becd5697144b72d7d6b09a92e955327ab67d1092, OpType: 0, ProcEpoch: 12996, Index: 4023, Creator: 0x6473b4c0552e46a80217a8ea619aa55fa8b5ad59, Amount: <nil>, Balance: <nil>, TxHash: <nil>}"
	   valSyncOp="{InitTxHash: 0x83eb2d837463f405ae9f145a8f598363406c6968e685397e4e20d127ecb56138, OpType: 0, ProcEpoch: 12996, Index: 4024, Creator: 0xcc5d2a3c5434eec89b8e87f71e693b35a0a603bb, Amount: <nil>, Balance: <nil>, TxHash: <nil>}"
	   valSyncOp="{InitTxHash: 0xf2d0d9db1648d3a05152c9c7fb4a6608fb3a17d2454152341c6caa376da6ab15, OpType: 0, ProcEpoch: 12996, Index: 4025, Creator: 0x0927bbbab3cd3a03bd3996d60558399bfd69c5a2, Amount: <nil>, Balance: <nil>, TxHash: <nil>}"
	   valSyncOp="{InitTxHash: 0x15b32c3ed82bb234be4af2c1a1450d41e9632d3066c5124a8848b3a00aaf2bea, OpType: 0, ProcEpoch: 12996, Index: 4026, Creator: 0x32e26f2c4ff439dbfd19e59d60cc4705cb128d0a, Amount: <nil>, Balance: <nil>, TxHash: <nil>}"
	*/

	failedOpsByInitTx := []*types.ValidatorSync{
		{
			InitTxHash: common.HexToHash("0x00d860db1ef65b539af8761193f7410dcd81db0ff2f13be058c989dfd532c861"),
			OpType:     0,
			Index:      4015,
			Creator:    common.HexToAddress("0xf0908e8722dce042c8af29fc1ed15c659a406a6b"),
		},
		{
			InitTxHash: common.HexToHash("0xdb82ab18dd9472d51f9a38d8d79401333a21c0aa17360cb93c1efbf73cf03fa4"),
			OpType:     0,
			Index:      4016,
			Creator:    common.HexToAddress("0x1d48f9a83c7328cf7fb5e3c550e1220feb5c7b13"),
		},
		{
			InitTxHash: common.HexToHash("0x924babb237a77f46541703cef25c9713724a53ca924f5d7b5182a7e316ae1ffa"),
			OpType:     0,
			Index:      4017,
			Creator:    common.HexToAddress("0x5040ea3b763acc52fa45cf3b9ff6bd8fbaa50b19"),
		},
		{
			InitTxHash: common.HexToHash("0xfa590cb6c760b34b72caa6d942f91308edee2fc723ae66215e52d6a5695e1f66"),
			OpType:     0,
			Index:      4018,
			Creator:    common.HexToAddress("0xb7634fbf58438d5f231c6e343ce6a09015b4b122"),
		},
		{
			InitTxHash: common.HexToHash("0x22104335f100ccd555d3a234474608c854bc25d3c1b9886a80a44e38d7f36d3b"),
			OpType:     0,
			Index:      4019,
			Creator:    common.HexToAddress("0x6c30179993f98d02d819520c0ab7d93efc1a464f"),
		},
		{
			InitTxHash: common.HexToHash("0x53aff3b86fbcc548eb52e2ba4537602399e1f7e6bc8b214830040b816a542be7"),
			OpType:     0,
			Index:      4020,
			Creator:    common.HexToAddress("0xf051499e8955feaf218dedeb46086af3f04341e4"),
		},
		{
			InitTxHash: common.HexToHash("0x07233067bbb076bfc1c3e06a604e570e6bfa858e1d0fc2b1d74a5ccecd9b2b87"),
			OpType:     0,
			Index:      4021,
			Creator:    common.HexToAddress("0xb0e42065991c5689ed3f570c64210f0174f5bf80"),
		},
		{
			InitTxHash: common.HexToHash("0x8b55f917100e95b6af49aeeba970a17e02de60d42e630364f3d8aa644acd8da2"),
			OpType:     0,
			Index:      4022,
			Creator:    common.HexToAddress("0xca94c692d4e5d6075c6eb2d1920d75169d8481ad"),
		},
		{
			InitTxHash: common.HexToHash("0x8b94ce4d80db11999d7cdaa4becd5697144b72d7d6b09a92e955327ab67d1092"),
			OpType:     0,
			Index:      4023,
			Creator:    common.HexToAddress("0x6473b4c0552e46a80217a8ea619aa55fa8b5ad59"),
		},
		{
			InitTxHash: common.HexToHash("0x83eb2d837463f405ae9f145a8f598363406c6968e685397e4e20d127ecb56138"),
			OpType:     0,
			Index:      4024,
			Creator:    common.HexToAddress("0xcc5d2a3c5434eec89b8e87f71e693b35a0a603bb"),
		},
		{
			InitTxHash: common.HexToHash("0xf2d0d9db1648d3a05152c9c7fb4a6608fb3a17d2454152341c6caa376da6ab15"),
			OpType:     0,
			Index:      4025,
			Creator:    common.HexToAddress("0x0927bbbab3cd3a03bd3996d60558399bfd69c5a2"),
		},
		{
			InitTxHash: common.HexToHash("0x15b32c3ed82bb234be4af2c1a1450d41e9632d3066c5124a8848b3a00aaf2bea"),
			OpType:     0,
			Index:      4026,
			Creator:    common.HexToAddress("0x32e26f2c4ff439dbfd19e59d60cc4705cb128d0a"),
		},
	}

	for _, valSyncOp := range failedOpsByInitTx {
		valSyncOp.ProcEpoch = 0
		valSyncOp.TxHash = &common.Hash{}
		bc.SetValidatorSyncData(valSyncOp)
		log.Info("fixMainnet0_setFailedValSyncOps: applied", "op", valSyncOp.Print())
	}

	return nil
}
