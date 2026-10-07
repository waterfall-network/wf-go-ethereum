// Copyright 2026 Digital Clever Solution LLC
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

package pluginimpl

import (
	wfcommon "github.com/LFDT-Iguazu/iguazu-types/common"
	gwatdownloader "gitlab.waterfall.network/waterfall/protocol/gwat/eth/downloader"
)

// downloaderWrapper wraps *gwatdownloader.Downloader to implement iface.Downloader.
// gwat/common.Hash and iguazu-types/common.Hash are both [32]byte; conversions are zero-cost.
type downloaderWrapper struct{ inner *gwatdownloader.Downloader }

func (w *downloaderWrapper) Synchronising() bool { return w.inner.Synchronising() }

func (w *downloaderWrapper) OptimisticSpineSync(spines wfcommon.HashArray) error {
	return w.inner.OptimisticSpineSync(gwatHashArray(spines))
}

func (w *downloaderWrapper) MainSync(baseSpine wfcommon.Hash, spines wfcommon.HashArray) error {
	return w.inner.MainSync(gwatHash(baseSpine), gwatHashArray(spines))
}

func (w *downloaderWrapper) DagSync(baseSpine wfcommon.Hash, spines wfcommon.HashArray) error {
	return w.inner.DagSync(gwatHash(baseSpine), gwatHashArray(spines))
}

func (w *downloaderWrapper) Terminate() { w.inner.Terminate() }
