// Copyright 2026 Digital Clever Solution Inc.
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

// Package pluginimpl wires gwat's concrete types to the wf-types/iface interfaces
// so that gwat can be loaded as a plugin by wf-engine.
package pluginimpl

import (
	"fmt"

	"gitlab.waterfall.network/waterfall/protocol/gwat/eth"
	"gitlab.waterfall.network/waterfall/protocol/gwat/internal/gwatapp"
	"gitlab.waterfall.network/waterfall/protocol/gwat/node"
	"gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/iface"
	"gopkg.in/urfave/cli.v1"
)

// GwatPluginImpl is the concrete implementation of iface.GwatPlugin exported
// from the gwat shared library. The exported plugin symbol in cmd/gwat-plugin/main.go
// holds a pointer to this type.
type GwatPluginImpl struct {
	stack    *node.Node
	ethereum *eth.Ethereum
	cliCtx   *cli.Context // retained for UnlockAccounts in Start()
	devMode  bool
}

// New returns a zero-value GwatPluginImpl. The caller must call Init before
// using any other method.
func New() *GwatPluginImpl { return &GwatPluginImpl{} }

// Init configures and creates the gwat node and Ethereum backend from the
// provided NodeConfig.Args, which are forwarded verbatim to gwat's own CLI
// flag parser. Must be called exactly once before Start.
func (p *GwatPluginImpl) Init(cfg *iface.NodeConfig) error {
	if cfg == nil {
		return fmt.Errorf("pluginimpl: Init: nil config")
	}

	stack, ethereum, cliCtx, err := gwatapp.MakeNode(cfg.Args)
	if err != nil {
		return fmt.Errorf("pluginimpl: Init: %w", err)
	}

	p.stack = stack
	p.ethereum = ethereum
	p.cliCtx = cliCtx
	p.devMode = cfg.DevMode
	return nil
}

// Start starts gwat's background services (P2P networking, DAG work, etc.),
// unlocks any accounts requested via --unlock / --password, and starts the
// block creator so that wf-engine's dag can call Creator().RunBlockCreation.
func (p *GwatPluginImpl) Start() error {
	if p.stack == nil {
		return fmt.Errorf("pluginimpl: Start: plugin not initialized")
	}
	if err := p.stack.Start(); err != nil {
		return err
	}
	if p.cliCtx != nil {
		gwatapp.UnlockAccounts(p.cliCtx, p.stack)
	}
	if err := p.ethereum.StartMining(0); err != nil {
		return fmt.Errorf("pluginimpl: Start: %w", err)
	}
	return nil
}

// Stop shuts down the gwat node gracefully.
func (p *GwatPluginImpl) Stop() error {
	if p.stack == nil {
		return nil
	}
	return p.stack.Close()
}

// BlockChain returns an iface.BlockChain adapter over gwat's *core.BlockChain.
func (p *GwatPluginImpl) BlockChain() iface.BlockChain {
	return &blockChainWrapper{inner: p.ethereum.BlockChain()}
}

// ValidatorChain returns an iface.ValidatorChain adapter — the same underlying
// wrapper as BlockChain since gwat's *core.BlockChain satisfies both interfaces.
func (p *GwatPluginImpl) ValidatorChain() iface.ValidatorChain {
	return &blockChainWrapper{inner: p.ethereum.BlockChain()}
}

// TxPool returns an iface.TxPool adapter over gwat's *core.TxPool.
func (p *GwatPluginImpl) TxPool() iface.TxPool {
	return &txPoolWrapper{inner: p.ethereum.TxPool()}
}

// Downloader returns an iface.Downloader adapter over gwat's eth/downloader.
func (p *GwatPluginImpl) Downloader() iface.Downloader {
	return &downloaderWrapper{inner: p.ethereum.Downloader()}
}

// Creator returns an iface.BlockCreator adapter over gwat's dag/creator.
func (p *GwatPluginImpl) Creator() iface.BlockCreator {
	return &creatorWrapper{inner: p.ethereum.DagCreator()}
}

// IsDevMode reports whether the plugin was started in development mode.
func (p *GwatPluginImpl) IsDevMode() bool { return p.devMode }

// SetProcessorFactories injects external token and validator processor factories
// into gwat's BlockChain. After this call, every block execution will use the
// provided factories instead of gwat's default implementations.
func (p *GwatPluginImpl) SetProcessorFactories(tp iface.TokenProcessorFactory, vp iface.ValidatorProcessorFactory) {
	p.ethereum.BlockChain().SetProcessorFactories(tp, vp)
}

// SetDag replaces gwat's internal dag workloop with the provided wf-engine dag.
// Must be called after Start(). gwat stops its own dag and routes all coordinator
// IPC calls to d via a thin type-conversion adapter.
func (p *GwatPluginImpl) SetDag(d iface.Dag) {
	p.ethereum.SetDagServicer(&wfDagAdapter{inner: d})
}
