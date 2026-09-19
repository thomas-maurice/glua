// Copyright (c) 2024-2025 Thomas Maurice
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package modules aggregates every Lua module shipped with glua. The
// RegisterAll entrypoint wires each module into a luareg.Registry — primarily
// used by cmd/glua-gen to regenerate the checked-in library/*.gen.lua stubs.
//
// Downstream projects that want a sandboxed subset of modules should not use
// RegisterAll. Instead, import the individual module packages
// (e.g. github.com/thomas-maurice/glua/pkg/modules/strings) and call their
// Register functions directly.
package modules

import (
	"github.com/thomas-maurice/glua/pkg/luareg"
	"github.com/thomas-maurice/glua/pkg/modules/base64"
	"github.com/thomas-maurice/glua/pkg/modules/bit32"
	"github.com/thomas-maurice/glua/pkg/modules/collections"
	"github.com/thomas-maurice/glua/pkg/modules/compress"
	"github.com/thomas-maurice/glua/pkg/modules/filepath"
	"github.com/thomas-maurice/glua/pkg/modules/fs"
	"github.com/thomas-maurice/glua/pkg/modules/hash"
	"github.com/thomas-maurice/glua/pkg/modules/hex"
	"github.com/thomas-maurice/glua/pkg/modules/http"
	"github.com/thomas-maurice/glua/pkg/modules/json"
	"github.com/thomas-maurice/glua/pkg/modules/jsonpath"
	k8sclient "github.com/thomas-maurice/glua/pkg/modules/k8sclient"
	"github.com/thomas-maurice/glua/pkg/modules/kubernetes"
	logmod "github.com/thomas-maurice/glua/pkg/modules/log"
	"github.com/thomas-maurice/glua/pkg/modules/netaddr"
	"github.com/thomas-maurice/glua/pkg/modules/osmod"
	"github.com/thomas-maurice/glua/pkg/modules/random"
	"github.com/thomas-maurice/glua/pkg/modules/regexp"
	"github.com/thomas-maurice/glua/pkg/modules/spew"
	"github.com/thomas-maurice/glua/pkg/modules/strconv"
	"github.com/thomas-maurice/glua/pkg/modules/strings"
	"github.com/thomas-maurice/glua/pkg/modules/template"
	"github.com/thomas-maurice/glua/pkg/modules/text"
	"github.com/thomas-maurice/glua/pkg/modules/time"
	neturl "github.com/thomas-maurice/glua/pkg/modules/url"
	"github.com/thomas-maurice/glua/pkg/modules/uuid"
	"github.com/thomas-maurice/glua/pkg/modules/yaml"
)

// RegisterAll: registers every module shipped with glua into reg. After the
// call, reg.Modules() contains one *luareg.Module per built-in module, ready
// for stub generation via pkg/stubgen.
//
// As modules are migrated to luareg in chunk A5, each is wired in here. Until
// A5 lands a given module, that module is not yet available through this
// aggregator.
func RegisterAll(reg *luareg.Registry) {
	base64.Register(reg)
	bit32.Register(reg)
	collections.Register(reg)
	compress.Register(reg)
	filepath.Register(reg)
	fs.Register(reg)
	hash.Register(reg)
	hex.Register(reg)
	http.Register(reg)
	json.Register(reg)
	jsonpath.Register(reg)
	k8sclient.Register(reg)
	kubernetes.Register(reg)
	logmod.Register(reg)
	netaddr.Register(reg)
	osmod.Register(reg)
	random.Register(reg)
	regexp.Register(reg)
	spew.Register(reg)
	strconv.Register(reg)
	strings.Register(reg)
	template.Register(reg)
	text.Register(reg)
	time.Register(reg)
	neturl.Register(reg)
	uuid.Register(reg)
	yaml.Register(reg)
}
