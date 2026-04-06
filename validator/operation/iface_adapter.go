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

package operation

import (
	"math/big"

	"gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/iface"
	wftypes "gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/types"
	wfcommon "gitlab.waterfall.network/waterfall/protocol/wf-types/common"
)

// valSyncIfaceAdapter adapts a gwat ValidatorSync to satisfy iface.ValidatorSyncOpFix.
// Both types have the same method names and semantics; only the concrete hash/address
// types differ (gwat/common vs wf-types/common — both are [32]byte / [20]byte).
type valSyncIfaceAdapter struct {
	op ValidatorSync
}

func (a *valSyncIfaceAdapter) InitTxHash() wfcommon.Hash {
	return wfcommon.Hash(a.op.InitTxHash())
}

func (a *valSyncIfaceAdapter) OpType() wftypes.ValidatorSyncOp {
	return wftypes.ValidatorSyncOp(a.op.OpType())
}

func (a *valSyncIfaceAdapter) ProcEpoch() uint64 { return a.op.ProcEpoch() }
func (a *valSyncIfaceAdapter) Index() uint64     { return a.op.Index() }

func (a *valSyncIfaceAdapter) Creator() wfcommon.Address {
	return wfcommon.Address(a.op.Creator())
}

func (a *valSyncIfaceAdapter) Amount() *big.Int { return a.op.Amount() }

// WrapForIface wraps a gwat ValidatorSync operation as iface.ValidatorSyncOpFix
// so it can be passed across the gwat/wf-engine boundary.
func WrapForIface(op ValidatorSync) iface.ValidatorSyncOpFix {
	return &valSyncIfaceAdapter{op: op}
}
