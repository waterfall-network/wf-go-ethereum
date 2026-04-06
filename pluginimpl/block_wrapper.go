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
	gwatcoretypes "gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	"gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/iface"
	wfcommon "gitlab.waterfall.network/waterfall/protocol/wf-types/common"
)

// blockWrapper wraps *gwatcoretypes.Block to implement iface.Block.
// gwat/common.Hash and wf-types/common.Hash are both [32]byte; the only
// difference is the named type, so conversions are zero-cost.
type blockWrapper struct{ inner *gwatcoretypes.Block }

func wrapBlock(b *gwatcoretypes.Block) iface.Block {
	if b == nil {
		return nil
	}
	return &blockWrapper{inner: b}
}

func wrapBlockMap(src gwatcoretypes.BlockMap) iface.BlockMap {
	if src == nil {
		return nil
	}
	dst := make(iface.BlockMap, len(src))
	for k, v := range src {
		dst[wfHash(k)] = wrapBlock(v)
	}
	return dst
}

func (w *blockWrapper) Nr() uint64 { return w.inner.Nr() }

func (w *blockWrapper) Number() *uint64 { return w.inner.Number() }

func (w *blockWrapper) Slot() uint64 { return w.inner.Slot() }

func (w *blockWrapper) Height() uint64 { return w.inner.Height() }

func (w *blockWrapper) Hash() wfcommon.Hash { return wfHash(w.inner.Hash()) }

func (w *blockWrapper) ParentHashes() wfcommon.HashArray {
	return wfHashArray(w.inner.ParentHashes())
}
