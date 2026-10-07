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
	"fmt"

	"gitlab.waterfall.network/waterfall/protocol/gwat/core"
	gwatcoretypes "gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
)

// txPoolWrapper wraps *core.TxPool and implements iface.TxPool.
type txPoolWrapper struct{ inner *core.TxPool }

// PendingCount returns the number of pending (executable) transactions.
func (w *txPoolWrapper) PendingCount() int {
	pending, _, _ := w.inner.Stats()
	return pending
}

// AddLocal decodes the RLP-encoded transaction and enqueues it as a local tx.
func (w *txPoolWrapper) AddLocal(raw []byte) error {
	tx := new(gwatcoretypes.Transaction)
	if err := tx.UnmarshalBinary(raw); err != nil {
		return fmt.Errorf("pluginimpl: AddLocal: decode tx: %w", err)
	}
	return w.inner.AddLocal(tx)
}

// Pending returns all pending transactions encoded as RLP bytes, keyed by
// sender address hex string. enforceTips=false mirrors gwat dag behaviour.
func (w *txPoolWrapper) Pending() (map[string][][]byte, error) {
	pending := w.inner.Pending(false)
	if pending == nil {
		return nil, nil
	}
	out := make(map[string][][]byte, len(pending))
	for addr, txs := range pending {
		key := addr.Hex()
		encoded := make([][]byte, 0, len(txs))
		for _, tx := range txs {
			b, err := tx.MarshalBinary()
			if err != nil {
				return nil, fmt.Errorf("pluginimpl: Pending: encode tx %s: %w", tx.Hash(), err)
			}
			encoded = append(encoded, b)
		}
		out[key] = encoded
	}
	return out, nil
}
