// Copyright 2024 Blue Wave Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package gwatapp provides shared CLI flag definitions and node-bootstrap logic
// used by both cmd/gwat (the standalone binary) and pluginimpl (the .so plugin).
// Keeping flags and MakeNode here avoids duplicating setup code in two places.
package gwatapp

import (
	"fmt"
	"strings"

	"gitlab.waterfall.network/waterfall/protocol/gwat/accounts"
	"gitlab.waterfall.network/waterfall/protocol/gwat/accounts/external"
	"gitlab.waterfall.network/waterfall/protocol/gwat/accounts/keystore"
	"gitlab.waterfall.network/waterfall/protocol/gwat/accounts/scwallet"
	"gitlab.waterfall.network/waterfall/protocol/gwat/accounts/usbwallet"
	"gitlab.waterfall.network/waterfall/protocol/gwat/cmd/utils"
	"gitlab.waterfall.network/waterfall/protocol/gwat/eth"
	"gitlab.waterfall.network/waterfall/protocol/gwat/eth/downloader"
	"gitlab.waterfall.network/waterfall/protocol/gwat/eth/ethconfig"
	"gitlab.waterfall.network/waterfall/protocol/gwat/internal/debug"
	"gitlab.waterfall.network/waterfall/protocol/gwat/log"
	"gitlab.waterfall.network/waterfall/protocol/gwat/node"
	"gopkg.in/urfave/cli.v1"
)

// ConfigFileFlag is the --config TOML flag shared between node mode and plugin mode.
var ConfigFileFlag = cli.StringFlag{
	Name:  "config",
	Usage: "TOML configuration file",
}

// NodeFlags lists the CLI flags used to configure the gwat node and P2P layer.
var NodeFlags = []cli.Flag{
	utils.IdentityFlag,
	utils.UnlockedAccountFlag,
	utils.PasswordFileFlag,
	utils.BootnodesFlag,
	utils.DataDirFlag,
	utils.AncientFlag,
	utils.MinFreeDiskSpaceFlag,
	utils.KeyStoreDirFlag,
	utils.UnlockVerifiersFlag,
	utils.VerifiersKeystoreFlag,
	utils.VerifiersPasswordsFlag,
	utils.ExternalSignerFlag,
	utils.NoUSBFlag,
	utils.USBFlag,
	utils.SmartCardDaemonPathFlag,
	utils.OverrideDelegatingStakeFlag,
	utils.OverridePrefixFinFlag,
	utils.TxPoolLocalsFlag,
	utils.TxPoolNoLocalsFlag,
	utils.TxPoolJournalFlag,
	utils.TxPoolRejournalFlag,
	utils.TxPoolPriceLimitFlag,
	utils.TxPoolPriceBumpFlag,
	utils.TxPoolAccountSlotsFlag,
	utils.TxPoolGlobalSlotsFlag,
	utils.TxPoolAccountQueueFlag,
	utils.TxPoolGlobalQueueFlag,
	utils.TxPoolLifetimeFlag,
	utils.TxPoolDroppedFlag,
	utils.SyncModeFlag,
	utils.ExitWhenSyncedFlag,
	utils.GCModeFlag,
	utils.SnapshotFlag,
	utils.TxLookupLimitFlag,
	utils.WhitelistFlag,
	utils.BloomFilterSizeFlag,
	utils.CacheFlag,
	utils.CacheDatabaseFlag,
	utils.CacheTrieFlag,
	utils.CacheTrieJournalFlag,
	utils.CacheTrieRejournalFlag,
	utils.CacheGCFlag,
	utils.CacheSnapshotFlag,
	utils.CacheNoPrefetchFlag,
	utils.CachePreimagesFlag,
	utils.ListenPortFlag,
	utils.MaxPeersFlag,
	utils.MaxPendingPeersFlag,
	utils.MinerEnabledFlag,
	utils.MinerThreadsFlag,
	utils.MinerNotifyFlag,
	utils.LegacyMinerGasTargetFlag,
	utils.MinerGasLimitFlag,
	utils.MinerGasLimitForceFlag,
	utils.MinerGasPriceFlag,
	utils.MinerExtraDataFlag,
	utils.MinerRecommitIntervalFlag,
	utils.MinerNoVerifyFlag,
	utils.NATFlag,
	utils.NoDiscoverFlag,
	utils.DiscoveryV5Flag,
	utils.NetrestrictFlag,
	utils.NodeKeyFileFlag,
	utils.NodeKeyHexFlag,
	utils.DNSDiscoveryFlag,
	utils.MainnetFlag,
	utils.TestNet8Flag,
	utils.DeveloperFlag,
	utils.Testnet5Flag,
	utils.Testnet9Flag,
	utils.VMEnableDebugFlag,
	utils.NetworkIdFlag,
	utils.EthStatsURLFlag,
	utils.FakePoWFlag,
	utils.NoCompactionFlag,
	utils.GpoBlocksFlag,
	utils.GpoPercentileFlag,
	utils.GpoMaxGasPriceFlag,
	utils.GpoIgnoreGasPriceFlag,
	utils.MinerNotifyFullFlag,
	ConfigFileFlag,
}

// RPCFlags lists the CLI flags used to configure the HTTP/WS/IPC RPC servers.
var RPCFlags = []cli.Flag{
	utils.HTTPEnabledFlag,
	utils.HTTPListenAddrFlag,
	utils.HTTPPortFlag,
	utils.HTTPCORSDomainFlag,
	utils.HTTPVirtualHostsFlag,
	utils.GraphQLEnabledFlag,
	utils.GraphQLCORSDomainFlag,
	utils.GraphQLVirtualHostsFlag,
	utils.HTTPApiFlag,
	utils.HTTPPathPrefixFlag,
	utils.WSEnabledFlag,
	utils.WSListenAddrFlag,
	utils.WSPortFlag,
	utils.WSApiFlag,
	utils.WSAllowedOriginsFlag,
	utils.WSPathPrefixFlag,
	utils.IPCDisabledFlag,
	utils.IPCPathFlag,
	utils.InsecureUnlockAllowedFlag,
	utils.RPCGlobalGasCapFlag,
	utils.RPCGlobalEVMTimeoutFlag,
	utils.RPCGlobalTxFeeCapFlag,
	utils.AllowUnprotectedTxs,
	utils.AuthListenFlag,
	utils.AuthPortFlag,
	utils.AuthVirtualHostsFlag,
	utils.JWTSecretFlag,
}

// MetricsFlags lists the CLI flags used to configure metrics/telemetry export.
var MetricsFlags = []cli.Flag{
	utils.MetricsEnabledFlag,
	utils.MetricsEnabledExpensiveFlag,
	utils.MetricsHTTPFlag,
	utils.MetricsPortFlag,
	utils.MetricsEnableInfluxDBFlag,
	utils.MetricsInfluxDBEndpointFlag,
	utils.MetricsInfluxDBDatabaseFlag,
	utils.MetricsInfluxDBUsernameFlag,
	utils.MetricsInfluxDBPasswordFlag,
	utils.MetricsInfluxDBTagsFlag,
	utils.MetricsEnableInfluxDBV2Flag,
	utils.MetricsInfluxDBTokenFlag,
	utils.MetricsInfluxDBBucketFlag,
	utils.MetricsInfluxDBOrganizationFlag,
}

// AppFlags returns the combined set of all gwat node-mode CLI flags.
// Console-only flags (--jspath, --exec, --preload) are excluded since they
// are not meaningful in plugin/daemon mode.
func AppFlags() []cli.Flag {
	var f []cli.Flag
	f = append(f, NodeFlags...)
	f = append(f, RPCFlags...)
	f = append(f, debug.Flags...)
	f = append(f, MetricsFlags...)
	return f
}

// MakeNode initializes a gwat node and Ethereum backend from raw CLI args.
// It does NOT start the node — call stack.Start() separately.
// The returned *cli.Context must be passed to UnlockAccounts after start.
// Subcommand tokens (e.g. "console") in args are silently ignored.
func MakeNode(args []string) (stack *node.Node, ethereum *eth.Ethereum, cliCtx *cli.Context, err error) {
	app := cli.NewApp()
	app.Name = "gwat"
	app.HideVersion = true
	app.HideHelp = true
	app.Flags = AppFlags()
	app.Before = func(ctx *cli.Context) error {
		return debug.Setup(ctx)
	}

	app.Action = func(ctx *cli.Context) error {
		cliCtx = ctx

		// Build node config from flags.
		nodeCfg := node.DefaultConfig
		nodeCfg.Name = "gwat"
		nodeCfg.IPCPath = "gwat.ipc"
		utils.SetNodeConfig(ctx, &nodeCfg)

		s, nodeErr := node.New(&nodeCfg)
		if nodeErr != nil {
			err = fmt.Errorf("gwatapp: create node: %w", nodeErr)
			return nil
		}

		if amErr := SetAccountManagerBackends(s); amErr != nil {
			s.Close()
			err = fmt.Errorf("gwatapp: account manager backends: %w", amErr)
			return nil
		}

		// Build eth config from flags.
		ethCfg := ethconfig.Defaults
		ethCfg.SyncMode = downloader.FullSync
		utils.SetEthConfig(ctx, s, &ethCfg)

		// Apply verifiers keystore overrides.
		if ctx.GlobalIsSet(utils.UnlockVerifiersFlag.Name) {
			nodeCfg.VerifiersKeystore.UnlockAllVerifiers = ctx.GlobalBool(utils.UnlockVerifiersFlag.Name)
		}
		if ctx.GlobalIsSet(utils.VerifiersKeystoreFlag.Name) {
			nodeCfg.VerifiersKeystore.KeyStoreDir = ctx.GlobalString(utils.VerifiersKeystoreFlag.Name)
		}
		if ctx.GlobalIsSet(utils.VerifiersPasswordsFlag.Name) {
			nodeCfg.VerifiersKeystore.PasswordFile = ctx.GlobalString(utils.VerifiersPasswordsFlag.Name)
		}
		nodeCfg.VerifiersKeystore.DataDir = s.Config().DataDir
		if keyDir, kdErr := s.Config().KeyDirConfig(); kdErr == nil {
			nodeCfg.VerifiersKeystore.OriginalKeyStoreDir = keyDir
		}
		if ctx.GlobalIsSet(utils.PasswordFileFlag.Name) {
			nodeCfg.VerifiersKeystore.OriginalPasswordFile = ctx.GlobalString(utils.PasswordFileFlag.Name)
		}

		// Apply hard-fork slot overrides.
		if ctx.GlobalIsSet(utils.OverrideDelegatingStakeFlag.Name) {
			val := ctx.GlobalUint64(utils.OverrideDelegatingStakeFlag.Name)
			ethCfg.OverrideDelegatingStake = &val
		}
		if ctx.GlobalIsSet(utils.OverridePrefixFinFlag.Name) {
			val := ctx.GlobalUint64(utils.OverridePrefixFinFlag.Name)
			ethCfg.OverridePrefixFin = &val
		}

		_, ethereum = utils.RegisterEthService(s, &ethCfg)
		stack = s
		return nil
	}

	// cli.v1 expects args[0] to be the program name.
	runArgs := append([]string{"gwat"}, args...)
	if runErr := app.Run(runArgs); runErr != nil && err == nil {
		err = fmt.Errorf("gwatapp: parse args: %w", runErr)
	}
	return stack, ethereum, cliCtx, err
}

// UnlockAccounts unlocks keystore accounts requested via --unlock / --password.
// Must be called after stack.Start() so the keystore backends are ready.
func UnlockAccounts(ctx *cli.Context, stack *node.Node) {
	var unlocks []string
	for _, input := range strings.Split(ctx.GlobalString(utils.UnlockedAccountFlag.Name), ",") {
		if trimmed := strings.TrimSpace(input); trimmed != "" {
			unlocks = append(unlocks, trimmed)
		}
	}
	if len(unlocks) == 0 {
		return
	}
	if !stack.Config().InsecureUnlockAllowed && stack.Config().ExtRPCEnabled() {
		utils.Fatalf("Account unlock with HTTP access is forbidden!")
	}
	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	passwords := utils.MakePasswordList(ctx)
	for i, addr := range unlocks {
		if stack.IsClosed() {
			return
		}
		unlockOne(ks, addr, i, passwords)
	}
}

// unlockOne unlocks a single keystore account, retrying up to 3 times on wrong password.
func unlockOne(ks *keystore.KeyStore, address string, i int, passwords []string) {
	account, err := utils.MakeAddress(ks, address)
	if err != nil {
		log.Warn("Could not find account to unlock", "address", address, "err", err)
		return
	}
	for trial := 0; trial < 3; trial++ {
		prompt := fmt.Sprintf("Unlocking account %s | Attempt %d/%d", address, trial+1, 3)
		password := utils.GetPassPhraseWithList(prompt, false, i, passwords)
		if unlockErr := ks.Unlock(account, password); unlockErr == nil {
			log.Info("Unlocked account", "address", account.Address.Hex())
			return
		} else if unlockErr != keystore.ErrDecrypt {
			log.Warn("Failed to unlock account", "address", address, "err", unlockErr)
			return
		}
	}
	log.Warn("Failed to unlock account after 3 attempts", "address", address)
}

// SetAccountManagerBackends adds keystore and hardware-wallet backends to the
// node's account manager. Mirrors cmd/gwat/config.go setAccountManagerBackends.
func SetAccountManagerBackends(stack *node.Node) error {
	conf := stack.Config()
	am := stack.AccountManager()
	keydir := stack.KeyStoreDir()

	scryptN, scryptP := keystore.StandardScryptN, keystore.StandardScryptP
	if conf.UseLightweightKDF {
		scryptN, scryptP = keystore.LightScryptN, keystore.LightScryptP
	}

	if len(conf.ExternalSigner) > 0 {
		log.Info("Using external signer", "url", conf.ExternalSigner)
		extapi, extErr := external.NewExternalBackend(conf.ExternalSigner)
		if extErr != nil {
			return fmt.Errorf("error connecting to external signer: %w", extErr)
		}
		am.AddBackend(extapi)
		return nil
	}

	am.AddBackend(keystore.NewKeyStore(keydir, scryptN, scryptP))
	if conf.USB {
		if ledger, ledErr := usbwallet.NewLedgerHub(); ledErr != nil {
			log.Warn("Failed to start Ledger hub, disabling", "err", ledErr)
		} else {
			am.AddBackend(ledger)
		}
		if trezorHID, tErr := usbwallet.NewTrezorHubWithHID(); tErr != nil {
			log.Warn("Failed to start HID Trezor hub, disabling", "err", tErr)
		} else {
			am.AddBackend(trezorHID)
		}
		if trezorUSB, tErr := usbwallet.NewTrezorHubWithWebUSB(); tErr != nil {
			log.Warn("Failed to start WebUSB Trezor hub, disabling", "err", tErr)
		} else {
			am.AddBackend(trezorUSB)
		}
	}

	// Set up smart card daemon if configured.
	if len(conf.SmartCardDaemonPath) > 0 {
		if scHub, scErr := scwallet.NewHub(conf.SmartCardDaemonPath, scwallet.Scheme, keydir); scErr != nil {
			log.Warn("Failed to start smart card hub, disabling", "err", scErr)
		} else {
			am.AddBackend(scHub)
		}
	}
	_ = accounts.WalletEvent{} // ensure accounts is used
	return nil
}
