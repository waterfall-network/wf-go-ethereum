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

package ethapi

import (
	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
)

// DagServicer is the coordinator-facing consensus interface used by the RPC layer.
// By default it is satisfied by gwat's internal *dag.Dag; when running as a plugin,
// it is replaced by a wfDagAdapter that delegates to wf-engine's dag.
type DagServicer interface {
	HandleFinalize(data *types.FinalizationParams) *types.FinalizationResult
	HandleCoordinatedState() *types.FinalizationResult
	HandleGetCandidates(slot uint64) *types.CandidatesResult
	HandleGetOptimisticSpines(fromSpine common.Hash) *types.OptimisticSpinesResult
	HandleValidateSpines(spines common.HashArray) (bool, error)
	HandleValidateFinalization(spines common.HashArray) (bool, error)
	HandleSyncSpines(spines common.HashArray) (bool, error)
	HandleSyncSlotInfo(slotInfo types.SlotInfo) (bool, error)
}
