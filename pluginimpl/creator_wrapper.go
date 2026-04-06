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

package pluginimpl

import (
	gwatcommon "gitlab.waterfall.network/waterfall/protocol/gwat/common"
	gwatcoretypes "gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	dagcreator "gitlab.waterfall.network/waterfall/protocol/gwat/dag/creator"
	wftypes "gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/types"
	wfcommon "gitlab.waterfall.network/waterfall/protocol/wf-types/common"
)

// creatorWrapper wraps *dagcreator.Creator to implement iface.BlockCreator.
type creatorWrapper struct{ inner *dagcreator.Creator }

func (w *creatorWrapper) Start() { w.inner.Start() }
func (w *creatorWrapper) Stop()  { w.inner.Stop() }

func (w *creatorWrapper) IsRunning() bool { return w.inner.IsRunning() }

func (w *creatorWrapper) IsCreatorActive(coinbase wfcommon.Address) bool {
	return w.inner.IsCreatorActive(gwatcommon.Address(coinbase))
}

func (w *creatorWrapper) RunBlockCreation(
	slot uint64,
	slotCreators []wfcommon.Address,
	tips wftypes.Tips,
	checkpoint *wftypes.Checkpoint,
) error {
	gwatCreators := make([]gwatcommon.Address, len(slotCreators))
	for i, a := range slotCreators {
		gwatCreators[i] = gwatcommon.Address(a)
	}
	return w.inner.RunBlockCreation(slot, gwatCreators, gwatTips(tips), gwatCheckpoint(checkpoint))
}

// gwatTips converts wftypes.Tips to the gwat native Tips map.
func gwatTips(src wftypes.Tips) gwatcoretypes.Tips {
	if src == nil {
		return nil
	}
	dst := make(gwatcoretypes.Tips, len(src))
	for k, v := range src {
		dst[gwatHash(k)] = gwatBlockDAG(v)
	}
	return dst
}
