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

package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomas-maurice/glua/pkg/luareg"
)

// TestRegisterAll_AllModulesWired: A5 migrated all 17 original modules;
// S1/S2 added bit32 and strconv, bringing the count to 19. RegisterAll must
// wire every one of them. Locking in the exact count catches a module being
// silently dropped from the aggregator (which would silently disappear from
// the regenerated library/*.gen.lua too).
func TestRegisterAll_AllModulesWired(t *testing.T) {
	reg := luareg.NewRegistry()
	require.NotPanics(t, func() { RegisterAll(reg) })
	assert.Equal(t, 19, len(reg.Modules()),
		"RegisterAll must wire every module; missing module would silently drop from the stub regen")
}
