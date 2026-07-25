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

package luareg

// Registry: aggregates module metadata for stub generation.
// Modules call Register to deposit their metadata here; the stub generator
// (A3) will iterate Modules() to produce output files.
type Registry struct {
	modules []*Module
}

// NewRegistry: creates a new empty Registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// add: appends a module to the registry (called by Module.Register).
func (r *Registry) add(m *Module) {
	r.modules = append(r.modules, m)
}

// Modules: returns all registered modules in registration order.
func (r *Registry) Modules() []*Module {
	return r.modules
}
