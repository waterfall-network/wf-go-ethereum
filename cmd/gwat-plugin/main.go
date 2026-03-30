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

// gwat-plugin builds gwat as a Go plugin (.so) that wf-engine can load via
// plugin.Open. The exported symbol GwatPlugin satisfies iface.GwatPlugin.
//
// Build with:
//
//	CGO_CFLAGS="-O2 -D__BLST_PORTABLE__" \
//	    go build -buildmode=plugin -o gwat.so ./cmd/gwat-plugin/
package main

import (
	"gitlab.waterfall.network/waterfall/protocol/gwat/pluginimpl"
	"gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/iface"
)

// GwatPlugin is the symbol looked up by wf-engine's loader package.
// It must be a package-level variable of type iface.GwatPlugin (or *iface.GwatPlugin).
var GwatPlugin iface.GwatPlugin = pluginimpl.New()

// main is required by the Go toolchain for package main, even in plugin builds.
func main() {}
