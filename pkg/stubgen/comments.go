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

package stubgen

import (
	"go/ast"
	"go/token"
	"reflect"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"
)

// commentResolver resolves struct field descriptions from Go source doc
// comments, implementing glua.FieldDocFunc. It loads each package lazily via
// golang.org/x/tools/go/packages and caches the result per package import
// path, since many stubbed types (e.g. k8s corev1) share a package — a
// generator run should load each package at most once.
//
// x/tools/go/packages is a dev-time-only dependency: it belongs here, in
// pkg/stubgen, not in pkg/glua's runtime import graph. pkg/glua only knows
// about the FieldDocFunc signature; this type supplies the implementation.
//
// Package loading or type/field lookup failures degrade silently to "no
// description" — stub generation must never fail because Go source isn't
// resolvable (e.g. a stripped module cache in a downstream build).
type commentResolver struct {
	mu   sync.Mutex
	pkgs map[string]map[string]map[string]string // pkgPath -> type name -> field name -> doc
}

// newCommentResolver: creates an empty resolver. Safe for concurrent use.
func newCommentResolver() *commentResolver {
	return &commentResolver{
		pkgs: make(map[string]map[string]map[string]string),
	}
}

// FieldDoc: implements glua.FieldDocFunc. Returns "" if pkgPath is empty
// (unnamed/builtin types), the package failed to load, or the field has no
// doc/inline comment.
func (c *commentResolver) FieldDoc(t reflect.Type, fieldName string) string {
	pkgPath := t.PkgPath()
	if pkgPath == "" {
		return ""
	}

	c.mu.Lock()
	types, ok := c.pkgs[pkgPath]
	if !ok {
		types = loadPackageFieldDocs(pkgPath)
		c.pkgs[pkgPath] = types
	}
	c.mu.Unlock()

	return types[t.Name()][fieldName]
}

// loadPackageFieldDocs: loads pkgPath's syntax tree and extracts struct
// field doc comments for every struct type declared directly in the
// package. Returns an empty, non-nil map on any failure (load error, no
// syntax available, package not found in a stripped module cache, etc.) so
// the failure is cached and not retried for every field in the package.
//
// Tests: true is set so that types declared in in-package _test.go files
// (package foo, not foo_test) are also resolved — glua's own test fixtures
// rely on this, and it costs nothing extra for downstream packages that
// have no test files. packages.Load then returns multiple *packages.Package
// values sharing pkgPath (the plain build, the "[pkgPath.test]" variant with
// test files merged in, etc.); results from all of them are merged.
func loadPackageFieldDocs(pkgPath string) map[string]map[string]string {
	result := make(map[string]map[string]string)

	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Tests: true,
	}
	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil {
		return result
	}

	for _, pkg := range pkgs {
		if pkg.PkgPath != pkgPath || len(pkg.Errors) > 0 || pkg.Syntax == nil {
			continue
		}
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}
				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}
					if fields := extractFieldDocs(structType); len(fields) > 0 {
						result[typeSpec.Name.Name] = fields
					}
				}
			}
		}
	}

	return result
}

// extractFieldDocs: reads each field's doc comment, falling back to its
// inline comment, and normalizes whitespace to a single line. Embedded
// fields are keyed by their own declared name (e.g. "Bar" for an embedded
// `Bar` or `*pkg.Bar`) — fields promoted from an embedded type are not
// chased into that type's own declaration (v1 limitation).
func extractFieldDocs(st *ast.StructType) map[string]string {
	docs := make(map[string]string)
	if st.Fields == nil {
		return docs
	}

	for _, field := range st.Fields.List {
		text := ""
		if field.Doc != nil {
			text = field.Doc.Text()
		} else if field.Comment != nil {
			text = field.Comment.Text()
		}
		text = normalizeDoc(text)
		if text == "" {
			continue
		}

		names := field.Names
		if len(names) == 0 {
			if name := embeddedFieldName(field.Type); name != "" {
				docs[name] = text
			}
			continue
		}
		for _, n := range names {
			docs[n.Name] = text
		}
	}
	return docs
}

// embeddedFieldName: returns the Go field name reflect assigns to an
// embedded field, e.g. "Bar" for `Bar`, `*Bar`, or `pkg.Bar`.
func embeddedFieldName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return e.Sel.Name
	case *ast.StarExpr:
		return embeddedFieldName(e.X)
	default:
		return ""
	}
}

// normalizeDoc: collapses a (possibly multi-line) doc/line comment into a
// single whitespace-normalized line, suitable for a ---@field annotation.
func normalizeDoc(text string) string {
	return strings.Join(strings.Fields(text), " ")
}
