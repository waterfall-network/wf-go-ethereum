// Copyright 2024   Blue Wave Inc.
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

/*
Package era implements functionality for managing eras in the Waterfall blockchain.
*/
package era

import (
	"errors"
	"math"
	"time"

	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	"gitlab.waterfall.network/waterfall/protocol/gwat/log"
	"gitlab.waterfall.network/waterfall/protocol/gwat/params"
)

// Errors related to Era processing.
var (
	ErrCheckpointInvalid = errors.New("invalid checkpoint") // Error when the checkpoint is invalid.
	ErrHandleEraFailed   = errors.New("handle era failed")  // Error when handling the era fails.
)

// Blockchain defines an interface for interacting with the blockchain.
type Blockchain interface {
	GetSlotInfo() *types.SlotInfo                                                       // Retrieves slot information.
	GetLastCoordinatedCheckpoint() *types.Checkpoint                                    // Retrieves the last coordinated checkpoint.
	GetEraInfo() *EraInfo                                                               // Retrieves current era information.
	Config() *params.ChainConfig                                                        // Retrieves the chain configuration.
	GetHeaderByHash(common.Hash) *types.Header                                          // Retrieves a header by its hash.
	EnterNextEra(fromEpoch uint64, root, hash common.Hash) (*Era, error)                // Enters the next era.
	StartTransitionPeriod(cp *types.Checkpoint, spineRoot, spineHash common.Hash) error // Starts the transition period.
}

// Era represents a period of time in the blockchain.
type Era struct {
	Number    uint64      `json:"number"`    // Era number.
	From      uint64      `json:"fromEpoch"` // Start epoch of the era.
	To        uint64      `json:"toEpoch"`   // End epoch of the era.
	Root      common.Hash `json:"root"`      // Root hash of the era.
	BlockHash common.Hash `json:"blockHash"` // Block hash of the era.
}

// NewEra creates a new instance of Era.
func NewEra(number, from, to uint64, root, blockHash common.Hash) *Era {
	return &Era{
		Number:    number,
		From:      from,
		To:        to,
		Root:      root,
		BlockHash: blockHash,
	}
}

// NextEra calculates and returns the next era based on the current blockchain state.
func NextEra(bc Blockchain, root, blockHash common.Hash, numValidators uint64) *Era {
	nextEraNumber := bc.GetEraInfo().Number() + 1
	nextEraLength := EstimateEraLength(bc.Config(), numValidators, nextEraNumber)
	nextEraBegin := bc.GetEraInfo().ToEpoch() + 1
	nextEraEnd := bc.GetEraInfo().ToEpoch() + nextEraLength

	return NewEra(nextEraNumber, nextEraBegin, nextEraEnd, root, blockHash)
}

// Length calculates the number of epochs in the era.
func (e *Era) Length() uint64 {
	return e.To - e.From + 1
}

// IsContainsEpoch checks if a given epoch is within the era.
func (e *Era) IsContainsEpoch(epoch uint64) bool {
	return epoch >= e.From && epoch <= e.To
}

// EraInfo contains metadata about the current era.
type EraInfo struct {
	currentEra *Era   // The current era.
	length     uint64 // Length of the era in epochs.
}

// NewEraInfo creates a new EraInfo instance.
func NewEraInfo(era *Era) *EraInfo {
	return &EraInfo{
		currentEra: era,
		length:     era.Length(),
	}
}

// Number retrieves the number of the current era.
func (ei *EraInfo) Number() uint64 {
	return ei.currentEra.Number
}

// GetEra retrieves the current era object.
func (ei *EraInfo) GetEra() *Era {
	return ei.currentEra
}

// ToEpoch retrieves the last epoch of the current era.
func (ei *EraInfo) ToEpoch() uint64 {
	return ei.GetEra().To
}

// FromEpoch retrieves the first epoch of the current era.
func (ei *EraInfo) FromEpoch() uint64 {
	return ei.GetEra().From
}

// EpochsPerEra calculates the total number of epochs in the era.
func (ei *EraInfo) EpochsPerEra() uint64 {
	return ei.GetEra().To - ei.GetEra().From + 1
}

// FirstEpoch retrieves the first epoch of the era.
func (ei *EraInfo) FirstEpoch() uint64 {
	return ei.FromEpoch()
}

// FirstSlot calculates the first slot of the era.
func (ei *EraInfo) FirstSlot(bc Blockchain) uint64 {
	slot, err := bc.GetSlotInfo().SlotOfEpochStart(ei.FirstEpoch())
	if err != nil {
		return 0
	}
	return slot
}

// LastEpoch retrieves the last epoch of the era.
func (ei *EraInfo) LastEpoch() uint64 {
	return ei.ToEpoch()
}

// LastSlot calculates the last slot of the era.
func (ei *EraInfo) LastSlot(bc Blockchain) uint64 {
	slot, err := bc.GetSlotInfo().SlotOfEpochEnd(ei.LastEpoch())
	if err != nil {
		return 0
	}
	return slot
}

// IsTransitionPeriodEpoch checks if the given epoch is within the transition period of the era.
func (ei *EraInfo) IsTransitionPeriodEpoch(bc Blockchain, epoch uint64) bool {
	return epoch >= ei.ToEpoch()-bc.Config().TransitionPeriod && epoch <= ei.ToEpoch()
}

// IsTransitionPeriodStartEpoch checks if the given epoch is the start of the transition period of the era.
func (ei *EraInfo) IsTransitionPeriodStartEpoch(bc Blockchain, epoch uint64) bool {
	return epoch == (ei.ToEpoch() - bc.Config().TransitionPeriod)
}

// IsTransitionPeriodStartSlot checks if the given slot is the start of the transition period for the next era.
func (ei *EraInfo) IsTransitionPeriodStartSlot(bc Blockchain, slot uint64) bool {
	transitionEpoch := (ei.ToEpoch() + 1 - bc.Config().TransitionPeriod)
	currentEpoch := bc.GetSlotInfo().SlotToEpoch(slot)

	if currentEpoch == transitionEpoch {
		if bc.GetSlotInfo().IsEpochStart(slot) {
			return true
		}
	}

	return false
}

// NextEraFirstEpoch returns the first epoch of the next era.
func (ei *EraInfo) NextEraFirstEpoch() uint64 {
	return ei.ToEpoch() + 1
}

// NextEraFirstSlot returns the first slot of the next era.
func (ei *EraInfo) NextEraFirstSlot(bc Blockchain) uint64 {
	slot, err := bc.GetSlotInfo().SlotOfEpochStart(ei.NextEraFirstEpoch())
	if err != nil {
		return 0
	}

	return slot
}

// LenEpochs returns the total number of epochs in the current era.
func (ei *EraInfo) LenEpochs() uint64 {
	return ei.length
}

// LenSlots returns the total number of slots in the current era.
// Assumes each epoch contains 32 slots.
func (ei *EraInfo) LenSlots() uint64 {
	return ei.length * 32
}

// IsContainsEpoch checks if the given epoch is within the current era.
func (ei *EraInfo) IsContainsEpoch(epoch uint64) bool {
	if epoch >= ei.FromEpoch() && epoch <= ei.ToEpoch() {
		return true
	}
	return false
}

// EstimateEraLength calculates the length of an era based on the chain configuration.
func EstimateEraLength(chainConfig *params.ChainConfig, numberOfValidators, eraNumber uint64) (eraLength uint64) {
	if eraNumber >= chainConfig.StartEpochsPerEra {
		return chainConfig.EpochsPerEra
	}

	var (
		epochsPerEra      = chainConfig.EpochsPerEra
		slotsPerEpoch     = float64(chainConfig.SlotsPerEpoch)
		validatorsPerSlot = float64(chainConfig.ValidatorsPerSlot)
	)

	eraLength = epochsPerEra * roundUp(float64(numberOfValidators)/(float64(epochsPerEra)*slotsPerEpoch*validatorsPerSlot))

	return
}

// roundUp rounds a floating-point number up to the nearest integer.
func roundUp(num float64) uint64 {
	return uint64(math.Ceil(num))
}

// HandleEra manages transitions between eras in the blockchain.
func HandleEra(bc Blockchain, cp *types.Checkpoint) error {
	defer func(start time.Time) {
		log.Info("^^^^^^^^^^^^ TIME",
			"elapsed", common.PrettyDuration(time.Since(start)),
			"func:", "HandleEra",
		)
	}(time.Now())

	log.Info("ERA started for new cp", "cp", cp.Epoch, "finEpoch", cp.FinEpoch, "spine", cp.Spine.Hex())

	var spineRoot, spineHash common.Hash
	// if cp != nil {
	header := bc.GetHeaderByHash(cp.Spine)
	if header != nil {
		spineRoot = header.Root
		spineHash = header.Hash()
	} else {
		log.Error("Checkpoint spine header not found", "err", ErrCheckpointInvalid)
		return ErrCheckpointInvalid
	}

	curToEpoch := bc.GetEraInfo().ToEpoch()
	// New era
	if bc.GetEraInfo().ToEpoch()+1 <= cp.FinEpoch {
		for curToEpoch+1 <= cp.FinEpoch {
			nextEra, err := bc.EnterNextEra(curToEpoch+1, spineRoot, spineHash)
			if err != nil {
				return err
			}
			if nextEra != nil {
				curToEpoch = nextEra.To
			} else {
				return ErrHandleEraFailed
			}
		}
		log.Info("Handle era", "cpEpoch", cp.Epoch,
			"cpFinEpoch", cp.FinEpoch,
			"curEpoch", bc.GetSlotInfo().SlotInEpoch(bc.GetSlotInfo().CurrentSlot()),
			"curSlot", bc.GetSlotInfo().CurrentSlot(),
			"bc.GetEraInfo().ToEpoch", bc.GetEraInfo().ToEpoch(),
			"bc.GetEraInfo().FromEpoch", bc.GetEraInfo().FromEpoch(),
			"bc.GetEraInfo().Number", bc.GetEraInfo().Number(),
		)
		return nil
	} else if (bc.GetEraInfo().ToEpoch()+1)-bc.Config().TransitionPeriod == cp.FinEpoch && cp.FinEpoch <= bc.GetEraInfo().ToEpoch()+1 {
		return bc.StartTransitionPeriod(cp, spineRoot, spineHash)
	}
	return nil
}
