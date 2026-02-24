package fixValidatorsStates

import (
	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
)

var mainnetFixData = []*FixOp{
	{
		OpType:     types.Deactivate,
		Index:      2125,
		Creator:    common.HexToAddress("0xd374819f2f66d2828b06ef7777e8ab97e6a2ebca"),
		InitTxHash: common.HexToHash("0x4f2c8e7b7b9eb70fa1941236519714e5a670f4618efbcfa7a1a2327fc4285bed"),
		PubKey:     common.HexToBlsPubKey("0x8ecf1693ec3ef1a85401513c26cebf0079d4af0cbdcd4ac468641df2a5fed720bcda31903abb24c9cd9453bb10e2b82a"),
		// eth.getTransaction("0x4f2c8e7b7b9eb70fa1941236519714e5a670f4618efbcfa7a1a2327fc4285bed").blockHash
		TxBlock: common.HexToHash("0x5bb0f5335b15d5bd68bc7d56cbd6598ac3a2a49ad051616d16f16017af13bf9a"),
	},
	{
		OpType:     types.Activate,
		Index:      4015,
		Creator:    common.HexToAddress("0xf0908e8722dce042c8af29fc1ed15c659a406a6b"),
		InitTxHash: common.HexToHash("0x00d860db1ef65b539af8761193f7410dcd81db0ff2f13be058c989dfd532c861"),
		PubKey:     common.HexToBlsPubKey("0x85f960b8c73fc4d4232bc1b93e62d3dc7aa6f678e3d5b8a64d37a3931193ad0f744838d1e5fba2bc5d36cc09ea9e49ef"),
		// eth.getTransaction("0x00d860db1ef65b539af8761193f7410dcd81db0ff2f13be058c989dfd532c861").blockHash
		TxBlock: common.HexToHash("0x7d8186ef8f5609f2b1cb7f40ff0afec009af631dc2dc215a838b5f5a938dfc49"),
	},
	{
		OpType:     types.Activate,
		Index:      4016,
		Creator:    common.HexToAddress("0x1d48f9a83c7328cf7fb5e3c550e1220feb5c7b13"),
		InitTxHash: common.HexToHash("0xdb82ab18dd9472d51f9a38d8d79401333a21c0aa17360cb93c1efbf73cf03fa4"),
		PubKey:     common.HexToBlsPubKey("0xa2d4d31aec089d9dec6fe19f0938608380e59bedcb97f936e12e21882c80071765228ba8e82e4ae1be59eb5a53133891"),
		// eth.getTransaction("0xdb82ab18dd9472d51f9a38d8d79401333a21c0aa17360cb93c1efbf73cf03fa4").blockHash
		TxBlock: common.HexToHash("0xed7c8860db8db5f3a57eca617d145afacc3c29116b19b4194e7a1d51467d5a86"),
	},
	{
		OpType:     types.Activate,
		Index:      4017,
		Creator:    common.HexToAddress("0x5040ea3b763acc52fa45cf3b9ff6bd8fbaa50b19"),
		InitTxHash: common.HexToHash("0x924babb237a77f46541703cef25c9713724a53ca924f5d7b5182a7e316ae1ffa"),
		PubKey:     common.HexToBlsPubKey("0xb2ee0f6f78418e8cc305732ed7b7ca4129be7424a5bc24c8fefc70d4bc3a5ddf73e00651612817fb59383e396c1a3410"),
		// eth.getTransaction("0x924babb237a77f46541703cef25c9713724a53ca924f5d7b5182a7e316ae1ffa").blockHash
		TxBlock: common.HexToHash("0xa0082f5bbb79ba4f0eb7aff48f304fd633ecd1e797138c9252a189c46251ee40"),
	},
	{
		OpType:     types.Activate,
		Index:      4018,
		Creator:    common.HexToAddress("0xb7634fbf58438d5f231c6e343ce6a09015b4b122"),
		InitTxHash: common.HexToHash("0xfa590cb6c760b34b72caa6d942f91308edee2fc723ae66215e52d6a5695e1f66"),
		PubKey:     common.HexToBlsPubKey("0xb65a0c4ac947bfbb02526e2310f4f04a31adc497d86c9bbfa454289d6a9d63db4655323ff3197ea420280debb46153cf"),
		// eth.getTransaction("0xfa590cb6c760b34b72caa6d942f91308edee2fc723ae66215e52d6a5695e1f66").blockHash
		TxBlock: common.HexToHash("0x3e47d5635f8be9795824b6bac463d1cf0a311e612a26ff5a0e962dfabcc538c4"),
	},
	{
		OpType:     types.Activate,
		Index:      4019,
		Creator:    common.HexToAddress("0x6c30179993f98d02d819520c0ab7d93efc1a464f"),
		InitTxHash: common.HexToHash("0x22104335f100ccd555d3a234474608c854bc25d3c1b9886a80a44e38d7f36d3b"),
		PubKey:     common.HexToBlsPubKey("0x8d00c98aed533c3a298ecff56c02dac4575c28bd454282fe5dc3c12b181a7e4dc7ab9049cf7c3527bfd710e3b19ec8a0"),
		// eth.getTransaction("0x22104335f100ccd555d3a234474608c854bc25d3c1b9886a80a44e38d7f36d3b").blockHash
		TxBlock: common.HexToHash("0x59be597430ad7ecfcc9c80a083831dec680e459c92c02ed5839c076c517bed20"),
	},
	{
		OpType:     types.Activate,
		Index:      4020,
		Creator:    common.HexToAddress("0xf051499e8955feaf218dedeb46086af3f04341e4"),
		InitTxHash: common.HexToHash("0x53aff3b86fbcc548eb52e2ba4537602399e1f7e6bc8b214830040b816a542be7"),
		PubKey:     common.HexToBlsPubKey("0x80155f39d77f1e4aefa2a5f903f0d90e426c9fce535c915bab89c31fb39a4fcf1f7d1f7638477b1ce8fab186ed14718e"),
		// eth.getTransaction("0x53aff3b86fbcc548eb52e2ba4537602399e1f7e6bc8b214830040b816a542be7").blockHash
		TxBlock: common.HexToHash("0x4b13403942398f2e663ece1ca3ed0bb88a0a5a1ca670b683f73f95d62d180f26"),
	},
	{
		OpType:     types.Activate,
		Index:      4021,
		Creator:    common.HexToAddress("0xb0e42065991c5689ed3f570c64210f0174f5bf80"),
		InitTxHash: common.HexToHash("0x07233067bbb076bfc1c3e06a604e570e6bfa858e1d0fc2b1d74a5ccecd9b2b87"),
		PubKey:     common.HexToBlsPubKey("0x81468f2e5e448752f996b8607f22313d568ae6e44a13c7ac58553ef0f6085bd603b099f70028650e527f696776e1381a"),
		// eth.getTransaction("0x07233067bbb076bfc1c3e06a604e570e6bfa858e1d0fc2b1d74a5ccecd9b2b87").blockHash
		TxBlock: common.HexToHash("0x1eef44666340adf403e962d2cac09e08c514ff961eb7a179a4d7708321201bb9"),
	},
	{
		OpType:     types.Activate,
		Index:      4022,
		Creator:    common.HexToAddress("0xca94c692d4e5d6075c6eb2d1920d75169d8481ad"),
		InitTxHash: common.HexToHash("0x8b55f917100e95b6af49aeeba970a17e02de60d42e630364f3d8aa644acd8da2"),
		PubKey:     common.HexToBlsPubKey("0xb49e824b8d42e4d508bd165f22d1b27c61a3de1ae65f90be6d21a08b5f3fc4ddc6b761f7f00b869a8b667bb19c6cfd52"),
		// eth.getTransaction("0x8b55f917100e95b6af49aeeba970a17e02de60d42e630364f3d8aa644acd8da2").blockHash
		TxBlock: common.HexToHash("0xe8838d772e3332dd8edd2a5a05466fe1e064ffb4a429f127ca12bebde688502c"),
	},
	{
		OpType:     types.Activate,
		Index:      4023,
		Creator:    common.HexToAddress("0x6473b4c0552e46a80217a8ea619aa55fa8b5ad59"),
		InitTxHash: common.HexToHash("0x8b94ce4d80db11999d7cdaa4becd5697144b72d7d6b09a92e955327ab67d1092"),
		PubKey:     common.HexToBlsPubKey("0x825feabeb4fbac349b9760efd671c5d7101792b191886842bf9a178936d51fa34a4aac21abcdbe4dd387ef3b1fe890ad"),
		// eth.getTransaction("0x8b94ce4d80db11999d7cdaa4becd5697144b72d7d6b09a92e955327ab67d1092").blockHash
		TxBlock: common.HexToHash("0x2f622dcaafcf9c04b65d829d28df7058f37d54416f5b187e3534f59f8aa45b91"),
	},
	{
		OpType:     types.Activate,
		Index:      4024,
		Creator:    common.HexToAddress("0xcc5d2a3c5434eec89b8e87f71e693b35a0a603bb"),
		InitTxHash: common.HexToHash("0x83eb2d837463f405ae9f145a8f598363406c6968e685397e4e20d127ecb56138"),
		PubKey:     common.HexToBlsPubKey("0x86cdd7f8b7bfe63bad0deb27387eebb5d66e679e9cbb1ed86767c577d773070f8f0121c692b2cbaa57045a84f21db1e5"),
		// eth.getTransaction("0x83eb2d837463f405ae9f145a8f598363406c6968e685397e4e20d127ecb56138").blockHash
		TxBlock: common.HexToHash("0xa23d159354b9c4db71fcb7025cfc87acfe83abfa2bb1229f4c1e7fbb30085c69"),
	},
	{
		OpType:     types.Activate,
		Index:      4025,
		Creator:    common.HexToAddress("0x0927bbbab3cd3a03bd3996d60558399bfd69c5a2"),
		InitTxHash: common.HexToHash("0xf2d0d9db1648d3a05152c9c7fb4a6608fb3a17d2454152341c6caa376da6ab15"),
		PubKey:     common.HexToBlsPubKey("0x88996525b293f4ebd250633771d9f46f3cb84141dce93ef9ee4d419982f991ceee0704d783c810634ae534b144523f56"),
		// eth.getTransaction("0xf2d0d9db1648d3a05152c9c7fb4a6608fb3a17d2454152341c6caa376da6ab15").blockHash
		TxBlock: common.HexToHash("0xd578745bd2fe13e36127648341c03f3dfdcbea1e77fa8ae4961b27341f092ca0"),
	},
	{
		OpType:     types.Activate,
		Index:      4026,
		Creator:    common.HexToAddress("0x32e26f2c4ff439dbfd19e59d60cc4705cb128d0a"),
		InitTxHash: common.HexToHash("0x15b32c3ed82bb234be4af2c1a1450d41e9632d3066c5124a8848b3a00aaf2bea"),
		PubKey:     common.HexToBlsPubKey("0xb78841353bdc13711421550c09b13097330e7575c6367f19a99b59ba088be5335ceebd0fcb9de52aab07b3ce9cbef000"),
		// eth.getTransaction("0x15b32c3ed82bb234be4af2c1a1450d41e9632d3066c5124a8848b3a00aaf2bea").blockHash
		TxBlock: common.HexToHash("0x96025847e786d31564b904325d9333c257bd5aa55b3d1af14363d32610e5deb8"),
	},
	{
		OpType:     types.Deactivate,
		Index:      7447,
		Creator:    common.HexToAddress("0xed93439e8ae0ada76e7ebcd51557f2298971be89"),
		InitTxHash: common.HexToHash("0xbcbb29267fe44e7139734dc73ac7fb725adfa2e72df81c65ca48e8deba3be748"),
		PubKey:     common.HexToBlsPubKey("0x877719d2ffb738467bbc73b79c3ad685eb493f45e0528c1938fd85f2039dc059e99cbfc8cbfc54e5956f6454c1b18184"),
		// eth.getTransaction("0xbcbb29267fe44e7139734dc73ac7fb725adfa2e72df81c65ca48e8deba3be748").blockHash
		TxBlock: common.HexToHash("0x52d793dadc3d1c83c705ec1b10f8ae957d7261ab5e738b1a31f9025dc45a0aba"),
	},
	{
		OpType:     types.Deactivate,
		Index:      7518,
		Creator:    common.HexToAddress("0xe0dcb41add071a018b35826fcd0b73603fdefa16"),
		InitTxHash: common.HexToHash("0xff0111ad5b5b79dfcf64ca03179cdd2ccbcf9c7ba0c31494b41d534c8fd523db"),
		PubKey:     common.HexToBlsPubKey("0xa69ef383aba81510c1472fdb665394a54ee467813e5788c006cef4b935b7f3d3461a207e4e29b0461156b6f5653eb33f"),
		// eth.getTransaction("0xff0111ad5b5b79dfcf64ca03179cdd2ccbcf9c7ba0c31494b41d534c8fd523db").blockHash
		TxBlock: common.HexToHash("0xa9679e7416ab0be97cf64dff42a9c0b5f8164a6811188f8ba0615e6221162063"),
	},
	{
		OpType:     types.Deactivate,
		Index:      8125,
		Creator:    common.HexToAddress("0x0d08c4cc0c0bb3a96db751b1ca45580117edff88"),
		InitTxHash: common.HexToHash("0x393d56dc28a2f1fa5914f5e8db25c1e85cad85b7e012b78671624274abc84c82"),
		PubKey:     common.HexToBlsPubKey("0xa1a125b5095f8249f68837e0bc0f8564e31e475c9e2d3246d32cb3fe74b62873e80b726c8446b8202dcc6e6843644449"),
		// eth.getTransaction("0x393d56dc28a2f1fa5914f5e8db25c1e85cad85b7e012b78671624274abc84c82").blockHash
		TxBlock: common.HexToHash("0x6ccd21952206005f3a9e48ca5402f9943393628ffa6064ca9c89d7cd9a1a1852"),
	},
	{
		OpType:     types.Deactivate,
		Index:      8126,
		Creator:    common.HexToAddress("0xf0fadb906c2821e434425d0e0d7bbf42324d16c6"),
		InitTxHash: common.HexToHash("0xee430a29ff411df1ac28c6b9df71712ffa3e0bd73005f33d9415c3209871fef8"),
		PubKey:     common.HexToBlsPubKey("0x895cf27f0c5356a78c2dcdcde0e1af7b047b5a706e3b3fa27a9f6d7239b2a6f1b59f9977086cd4b34ad500c368b0bfb8"),
		// eth.getTransaction("0xee430a29ff411df1ac28c6b9df71712ffa3e0bd73005f33d9415c3209871fef8").blockHash
		TxBlock: common.HexToHash("0x6b32d83111ca6319384019523e9059baae5d4495c84ac5c41256383b85dcea9d"),
	},
	{
		OpType:     types.Deactivate,
		Index:      8574,
		Creator:    common.HexToAddress("0x316f9f26bcb0179c2efef93bd15639eb844f87cc"),
		InitTxHash: common.HexToHash("0x5b473fcface61c42209c9ecc5544b98eaf1c151ffcb46715a99673b3f7ab1d1b"),
		PubKey:     common.HexToBlsPubKey("0xb32837c02f841d9ccff2941591144d32dbf710b5a2d24d873425701a9f88380bb99d499285f578c7fbbfbab899a70f9d"),
		// eth.getTransaction("0x5b473fcface61c42209c9ecc5544b98eaf1c151ffcb46715a99673b3f7ab1d1b").blockHash
		TxBlock: common.HexToHash("0x8dbb6be538519f6afa5db40f9cebc9d9715d37d9ceaa2541c082622fd8a9c31b"),
	},
	{
		OpType:     types.Deactivate,
		Index:      8575,
		Creator:    common.HexToAddress("0xeba33e04faa49c222ad4f2e72a4b05605d478b02"),
		InitTxHash: common.HexToHash("0x27288f98b92e5c45643c852b34ae510a20497562d5e6b8298adfa8a855833320"),
		PubKey:     common.HexToBlsPubKey("0x882136dc219a7ba1d58e3cf4922b44d411c489caf5ca7839241141a43583d8750eb7897cd980f51f35c2ab76b244ecc1"),
		// eth.getTransaction("0x27288f98b92e5c45643c852b34ae510a20497562d5e6b8298adfa8a855833320").blockHash
		TxBlock: common.HexToHash("0xc439c901f845a03343ad8c01e8ef4e3868a8d681c65e12089f9c1d3ef6ac4c00"),
	},
	{
		OpType:     types.Deactivate,
		Index:      8576,
		Creator:    common.HexToAddress("0x0fb582447cdaeb99ec9b777d51fd0b343fe5f92c"),
		InitTxHash: common.HexToHash("0xbf30e0dec01431c9f63bbaa721816b650cccb79330f8d2e26d23eed0496497e9"),
		PubKey:     common.HexToBlsPubKey("0x8b355c84f93ed675371789faea402ba91299a7765babd54487f645d7fc6d321e23327d6a348da43a2984b343a3959dff"),
		// eth.getTransaction("0xbf30e0dec01431c9f63bbaa721816b650cccb79330f8d2e26d23eed0496497e9").blockHash
		TxBlock: common.HexToHash("0x0ec3a051f6114cc0c83794bb3f0c499e9b58a66b81bc2c70ef04a438946e449e"),
	},
	{
		OpType:     types.Deactivate,
		Index:      8577,
		Creator:    common.HexToAddress("0xda1d014b9ff590b4d0388f4f2cccbc7f974b13ad"),
		InitTxHash: common.HexToHash("0x4a0d3dda839c3b6691fa383431689fabd320f6c55dfb8851cbae8f30ed57cbbf"),
		PubKey:     common.HexToBlsPubKey("0xa637bb687a3c53dec88ba49e3d0e9b5e046857906f73fcaca2d052f8f202ecdb6919422c609acfe829635da3634b22cb"),
		// eth.getTransaction("0x4a0d3dda839c3b6691fa383431689fabd320f6c55dfb8851cbae8f30ed57cbbf").blockHash
		TxBlock: common.HexToHash("0x6e103c6ec7d9980e1d07772023cd337367c974139bd83e4945bc8c7ac3e9cda7"),
	},
	{
		OpType:     types.Activate,
		Index:      46793,
		Creator:    common.HexToAddress("0x134828c57593d6f913dcea89e5fb2dc233d65e63"),
		InitTxHash: common.HexToHash("0x1c893460442fdc5a78ce2de79d11a774b1c3345bfabcff7f427c2106aa8f0289"),
		PubKey:     common.HexToBlsPubKey("0x9363a750d2f53abb05d836b8039589ea11953f1dc0bc451a9ba74d0a1beec3ec16ad4547d6ce597de26504d40405dcfe"),
		// eth.getTransactionReceipt("0x1c893460442fdc5a78ce2de79d11a774b1c3345bfabcff7f427c2106aa8f0289").blockHash
		TxBlock: common.HexToHash("0xae0f21c95aaae244bf609d929927a0e51f1457e4a9b37fc04a147f5257df0e23"),
	},
	{
		OpType:     types.Activate,
		Index:      46794,
		Creator:    common.HexToAddress("0xe283a7858b8476ca9cacac957c7866ed8433518f"),
		InitTxHash: common.HexToHash("0xbe836831d31d723f7dcdf70f8bb37a278e47bfc060148a29ae0648a6099edefb"),
		PubKey:     common.HexToBlsPubKey("0xa385d2014493f69ce817fe8d8762549bb6b370b3defcf7960cba50884a3ebfbb78369572667f29d0198752a3b887276f"),
		// eth.getTransactionReceipt("0xbe836831d31d723f7dcdf70f8bb37a278e47bfc060148a29ae0648a6099edefb").blockHash
		TxBlock: common.HexToHash("0xdc0717140d7a9d10e35d1fcde0c73bf0ee3b116a7dfda589535064635d18c5ad"),
	},
	{
		OpType:     types.Activate,
		Index:      46984,
		Creator:    common.HexToAddress("0x998b226353208364b895a8be16a4a0e81d95988b"),
		InitTxHash: common.HexToHash("0x27c384c8ce5c53bb05c8476ba347562535a779b6ac11c3823031e9708683985b"),
		PubKey:     common.HexToBlsPubKey("0x85b54e200c1383ad40f6c87e08bdb05fd76609c0cc12df9a2fa24e4e6f44cd526cfd16e674084109db5ffe7189777e57"),
		// eth.getTransactionReceipt("0x27c384c8ce5c53bb05c8476ba347562535a779b6ac11c3823031e9708683985b").blockHash
		TxBlock: common.HexToHash("0x0fec6f39b48d3fa208914c485748072ec7e2ca8343ff6e7757f360f25c9a8938"),
	},
	{
		OpType:     types.Activate,
		Index:      46985,
		Creator:    common.HexToAddress("0x4c17d39a9db644329246cbeca08b11a4ac8ac50a"),
		InitTxHash: common.HexToHash("0x1a75a8f1b66742c5ddac96a5fd4a903061c416df07828403a64672a0976b3d4f"),
		PubKey:     common.HexToBlsPubKey("0xa6f8cf2ba3a25f7405a09e0834f9d15090f41cb27578653dab813bd69f31f7790563523a81f9d0038c97b3cd4f61bf42"),
		// eth.getTransactionReceipt("0x1a75a8f1b66742c5ddac96a5fd4a903061c416df07828403a64672a0976b3d4f").blockHash
		TxBlock: common.HexToHash("0xc2988ad4e21074cbcbee6b91450425508715e43a20bf06897271fda2b0afd71f"),
	},
	{
		OpType:     types.Activate,
		Index:      47541,
		Creator:    common.HexToAddress("0x2bdce0e0390bf1af69ad0d42dd14e19dc54098ba"),
		InitTxHash: common.HexToHash("0x6def891e6edf0d1bf173e1307eeaba303d1fae75a8f1394bad933e13c59e18aa"),
		PubKey:     common.HexToBlsPubKey("0x9685f73266362d665760c4197c642fe21119bfcd375c616618c6da9318b9bbc47077a48d5e6aeab96292e5d857223a03"),
		// eth.getTransactionReceipt("0x6def891e6edf0d1bf173e1307eeaba303d1fae75a8f1394bad933e13c59e18aa").blockHash
		TxBlock: common.HexToHash("0x55719e09fd96e4fbc1d5083165c511cda67b2fa54d98ae5d412d80b9ba7fb0cc"),
	},
	{
		OpType:     types.Activate,
		Index:      47542,
		Creator:    common.HexToAddress("0xbc5530a92f2bfa215df2fca66b89e61e57295315"),
		InitTxHash: common.HexToHash("0xd12b65aeb749d89f5b54565c8747dcf91964860d309d90e55ccbf7c36b39259a"),
		PubKey:     common.HexToBlsPubKey("0xb156a683c0ea4d1d29c04eb7087a113b3956e709948aa2b62e764ca158567792777cf436808fb24bf64f34e627232798"),
		// eth.getTransactionReceipt("0xd12b65aeb749d89f5b54565c8747dcf91964860d309d90e55ccbf7c36b39259a").blockHash
		TxBlock: common.HexToHash("0x1a2606641ad5aa89d452c67413f5a2c79920f72f11b1fe0fcc1c210101152a6d"),
	},
	{
		OpType:     types.Activate,
		Index:      47543,
		Creator:    common.HexToAddress("0x1ccfb4ecb2c37c9da6c98a7c6e793af97a86d8d0"),
		InitTxHash: common.HexToHash("0x5f0ad008f04f050a6bd89bcd1ca9b44d58b658b9f4aaf72f0dae4b5c930e9418"),
		PubKey:     common.HexToBlsPubKey("0x97b03d8cce800e68d092b5d560e544f94faf051f89cdee15a919c83d0f8cd814dd5ce9c5b1f15887acf4630ac4b8aeac"),
		// eth.getTransactionReceipt("0x5f0ad008f04f050a6bd89bcd1ca9b44d58b658b9f4aaf72f0dae4b5c930e9418").blockHash
		TxBlock: common.HexToHash("0xb714d6a1548598c3df2543d3a8fbba2a2d2f5d48ba54dd677109405a054b5b9e"),
	},
	{
		OpType:     types.Activate,
		Index:      47755,
		Creator:    common.HexToAddress("0xc9f5562a91ffeba99b8947e5aa500d556c96da1e"),
		InitTxHash: common.HexToHash("0xa7f9fc7c9d2f1a661c1ade6864d3da046457c036a1aa827ef6a26beba827810e"),
		PubKey:     common.HexToBlsPubKey("0xb2cad6fdcb96df249e85fa3acf5e2dd1ac5cd28b497186327d708c79ef2cea0eab6d8e9165b7b71b55b34b0b6330c320"),
		// eth.getTransactionReceipt("0xa7f9fc7c9d2f1a661c1ade6864d3da046457c036a1aa827ef6a26beba827810e").blockHash
		TxBlock: common.HexToHash("0x7928a27ec34f6c91ddb29941d375a70c5a400e04edfcf7bf31d4a617dd0e186d"),
	},
	{
		OpType:     types.Activate,
		Index:      47756,
		Creator:    common.HexToAddress("0x54f3a40ac2f077d82d955a996f6e87b809b17c43"),
		InitTxHash: common.HexToHash("0xeb5a1eb16c7a0a4e7438ccb07d74eb5514b84a21b626ad4be95fc26ea14f9d70"),
		PubKey:     common.HexToBlsPubKey("0xa819c26e239a0e925c82ea44b7c2fe7ea8c15c96d35f6e4185edebb040b4c16ca7a2fdc1d4637901087ed5dd985edca6"),
		// eth.getTransactionReceipt("0xeb5a1eb16c7a0a4e7438ccb07d74eb5514b84a21b626ad4be95fc26ea14f9d70").blockHash
		TxBlock: common.HexToHash("0xc6bdf387c561d8707d10f5c0478f5fd54451c29572d8331f05b17a5d199d5ef9"),
	},
	{
		OpType:     types.Activate,
		Index:      47757,
		Creator:    common.HexToAddress("0xb72ed19a6db5c117fd61b741a53d29422456dc62"),
		InitTxHash: common.HexToHash("0xf1d96d2c5455d1490a26513f24fe787cd388729316e164dfeda0fdcc170ecfd1"),
		PubKey:     common.HexToBlsPubKey("0xaf10a927a92ac50838ec2b32eb85b070d9f8c890dd79f30c431d560e277b5a5254fd2d65aac6b046097edc3bc72ae2d2"),
		// eth.getTransactionReceipt("0xf1d96d2c5455d1490a26513f24fe787cd388729316e164dfeda0fdcc170ecfd1").blockHash
		TxBlock: common.HexToHash("0xf939306c70ea5c2b47c58ebedbfa982d56fee43ce59756135b2dbc7932af4590"),
	},
	{
		OpType:     types.Activate,
		Index:      47758,
		Creator:    common.HexToAddress("0xc760f11cd00f78327a941d8daf4bec02cb3675f1"),
		InitTxHash: common.HexToHash("0x20acaed0f7ad03539cf6a9c595899334f4c0b622910b57f20923dc65327d11ed"),
		PubKey:     common.HexToBlsPubKey("0xa0851da391b59a3dbc4ea94842ec1ef9db21a5285be1ed8e81fe7f980c4ae6b4b378a014ab50f061e7658497b19b9909"),
		// eth.getTransactionReceipt("0x20acaed0f7ad03539cf6a9c595899334f4c0b622910b57f20923dc65327d11ed").blockHash
		TxBlock: common.HexToHash("0x15bbb824305fae6e265911e94a04b5d7d02f545ba2c3d3d1d80b2390f8585cdd"),
	},
	{
		OpType:     types.Activate,
		Index:      48160,
		Creator:    common.HexToAddress("0x73f3c03d7f3f6e2711f3aa486bdc60a03dc6463c"),
		InitTxHash: common.HexToHash("0x2b9410de2f37e9766bca6b14b7699a07be02f42bcb6c3429dd115a48c720f1fa"),
		PubKey:     common.HexToBlsPubKey("0xb697a275232e77eefa35a1d306e9c5957f6a28ca4ae3c674218551285ae89ebee80ce169214261ab0dd4cb166186a6df"),
		// eth.getTransactionReceipt("0x2b9410de2f37e9766bca6b14b7699a07be02f42bcb6c3429dd115a48c720f1fa").blockHash
		TxBlock: common.HexToHash("0x41fe9ef32ed24e22029743d62ede722d41f6eb074bda614a0f354da978029886"),
	},
	{
		OpType:     types.Activate,
		Index:      48161,
		Creator:    common.HexToAddress("0xd35f63e71a1d4f6d06c5dd91c05b27d304d059ee"),
		InitTxHash: common.HexToHash("0x110b99a81a05e6e2ec08371d065e4afe0c7b8a975bf61479e0d2656d863e1e69"),
		PubKey:     common.HexToBlsPubKey("0xb6cd5bf383f302f107d2ec4cc2557c818bcf2eaa72e1e4dd1fb193f3081b67b82bf6a1b5e70d044adc4cb2d19e031f70"),
		// eth.getTransactionReceipt("0x110b99a81a05e6e2ec08371d065e4afe0c7b8a975bf61479e0d2656d863e1e69").blockHash
		TxBlock: common.HexToHash("0xd50b5b54ca6e42688ce98f8948f200ea48142fc458b705e0804fcc68f70df58a"),
	},
	{
		OpType:     types.Activate,
		Index:      48162,
		Creator:    common.HexToAddress("0x145a2532913f7490e7885ac5572ad946d754fb10"),
		InitTxHash: common.HexToHash("0x599afb07d2fd7e3c079854a06e59a034d92f0e12cd91a91e3654c4258e573d5b"),
		PubKey:     common.HexToBlsPubKey("0x804e219aaeb7f73c3df9d8915408435b60cb450a47be76661c4b42fe3779884414afab5967b6dbfd47bcf9e2b6b94795"),
		// eth.getTransactionReceipt("0x599afb07d2fd7e3c079854a06e59a034d92f0e12cd91a91e3654c4258e573d5b").blockHash
		TxBlock: common.HexToHash("0xba3ac22d252f21c1a79142f29a521be905eb10fc6ea461d26a49297da6df788e"),
	},
	{
		OpType:     types.Activate,
		Index:      48373,
		Creator:    common.HexToAddress("0xfd1ed6e6f5b13aa47b9023e2edfdd7f3bfded43a"),
		InitTxHash: common.HexToHash("0xbcf0beffed43a40139fb1fb6dd7e9584aa2709eff83ac39a9f84439205a71f39"),
		PubKey:     common.HexToBlsPubKey("0xa23a0f5ccb2866e7a94adf6c803f98ce73d0352fb607e0231b843063054ddf02fef3c5127c902ba03a614b9b313ebf5e"),
		// eth.getTransactionReceipt("0xbcf0beffed43a40139fb1fb6dd7e9584aa2709eff83ac39a9f84439205a71f39").blockHash
		TxBlock: common.HexToHash("0xfd6875a28c4472f198cec1d69bfd4035226d579ecb84152274dad8085601e6ca"),
	},
	{
		OpType:     types.Activate,
		Index:      48374,
		Creator:    common.HexToAddress("0xd67d91840744ed2998e14ab54ab3e675e6ed09fc"),
		InitTxHash: common.HexToHash("0xbdbf54628dbaf2816a98af968214040ac466f63eb41e56b0f0fbf5c0f29ed106"),
		PubKey:     common.HexToBlsPubKey("0x927fd9419a32323535387be788bfb7fd2638f535315c1b167ed71cdbc0e17f9c8f386a092d0c250c44e02bae1ed54da7"),
		// eth.getTransactionReceipt("0xbdbf54628dbaf2816a98af968214040ac466f63eb41e56b0f0fbf5c0f29ed106").blockHash
		TxBlock: common.HexToHash("0xf3b90ed5118a2a07182531a2b9105ed6675b2138e029cc95122c38245387f0e2"),
	},
}
