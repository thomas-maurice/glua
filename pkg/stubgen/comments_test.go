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
	"go/parser"
	"go/token"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// parseStructFields: parses a single struct type declaration and returns its
// *ast.StructType, for testing extractFieldDocs against exact source without
// depending on how gofmt happens to reformat multi-line comments elsewhere.
func parseStructFields(t *testing.T, src string) *ast.StructType {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "fixture.go", "package fixture\n\n"+src, parser.ParseComments)
	require.NoError(t, err)

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
			if st, ok := typeSpec.Type.(*ast.StructType); ok {
				return st
			}
		}
	}
	t.Fatal("no struct type found in source")
	return nil
}

// TestCommentResolver_FieldDoc: resolves doc/inline comments for docFixture's
// fields and returns "" for a field with neither.
func TestCommentResolver_FieldDoc(t *testing.T) {
	c := newCommentResolver()
	typ := reflect.TypeOf(docFixture{})

	assert.Equal(t, "Name is the fixture's display name, read from a doc comment above the field.", c.FieldDoc(typ, "Name"))
	assert.Equal(t, "Count is read from this inline comment.", c.FieldDoc(typ, "Count"))
	assert.Equal(t, "", c.FieldDoc(typ, "Silent"), "field with no doc/inline comment resolves to empty")
	assert.Equal(t, "", c.FieldDoc(typ, "DoesNotExist"), "unknown field name resolves to empty, not a panic")
}

// TestCommentResolver_FieldDoc_UnnamedType: an anonymous struct has an empty
// PkgPath; FieldDoc must degrade to "" rather than attempt to load "".
func TestCommentResolver_FieldDoc_UnnamedType(t *testing.T) {
	c := newCommentResolver()
	typ := reflect.TypeOf(struct{ X int }{})

	assert.Equal(t, "", c.FieldDoc(typ, "X"))
}

// TestCommentResolver_FieldDoc_CachesPerPackage: looking up two fields on
// the same type must only populate one cache entry for that package — this
// is what makes resolving many types from one package (e.g. k8s corev1)
// affordable instead of reloading the package per field.
func TestCommentResolver_FieldDoc_CachesPerPackage(t *testing.T) {
	c := newCommentResolver()
	typ := reflect.TypeOf(docFixture{})

	c.FieldDoc(typ, "Name")
	c.FieldDoc(typ, "Count")

	assert.Len(t, c.pkgs, 1, "only one package entry should exist for two fields on the same type")
}

// TestLoadPackageFieldDocs_UnresolvablePackage: stub generation must never
// fail because a package can't be resolved (e.g. a stripped module cache in
// a downstream build) — this must degrade to an empty map, not an error or
// panic.
func TestLoadPackageFieldDocs_UnresolvablePackage(t *testing.T) {
	result := loadPackageFieldDocs("this/package/definitely/does/not/exist/anywhere")
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

// TestExtractFieldDocs_MultiLineComment: a doc comment spread across several
// `//` lines is a single logical description — the newlines between them are
// formatting, not content, so they must collapse to spaces rather than leak
// literal "\n" characters into the ---@field annotation (which would break
// the single-line LuaLS comment syntax).
func TestExtractFieldDocs_MultiLineComment(t *testing.T) {
	st := parseStructFields(t, `type T struct {
	// Name is the display
	// name of the thing,
	// spread across lines.
	Name string
}`)

	docs := extractFieldDocs(st)
	assert.Equal(t, "Name is the display name of the thing, spread across lines.", docs["Name"])
	assert.NotContains(t, docs["Name"], "\n")
}

// TestExtractFieldDocs_ParagraphBreak: a doc comment with a blank comment
// line in the middle is two paragraphs to godoc, but still one logical
// ---@field description here — normalizeDoc must collapse the blank-line gap
// to a single space, not leave a double space or embedded newline.
func TestExtractFieldDocs_ParagraphBreak(t *testing.T) {
	st := parseStructFields(t, `type T struct {
	// First paragraph.
	//
	// Second paragraph.
	Name string
}`)

	docs := extractFieldDocs(st)
	assert.Equal(t, "First paragraph. Second paragraph.", docs["Name"])
	assert.NotContains(t, docs["Name"], "\n")
	assert.NotContains(t, docs["Name"], "  ", "no double space at the paragraph gap")
}

// TestExtractFieldDocs_BlockComment: a /* ... */ block comment spanning
// multiple physical lines must also normalize to one line, same as a run of
// "//" line comments.
func TestExtractFieldDocs_BlockComment(t *testing.T) {
	st := parseStructFields(t, `type T struct {
	/* Name is the display
	   name of the thing. */
	Name string
}`)

	docs := extractFieldDocs(st)
	assert.Equal(t, "Name is the display name of the thing.", docs["Name"])
	assert.NotContains(t, docs["Name"], "\n")
}

// TestExtractFieldDocs_EmbeddedField: an embedded field's own doc comment is
// captured under its declared (promoted) name; fields promoted from further
// inside the embedded type are out of scope (v1 limitation, not tested here
// since it requires chasing a second type declaration).
func TestExtractFieldDocs_EmbeddedField(t *testing.T) {
	docs := loadPackageFieldDocs("github.com/thomas-maurice/glua/pkg/stubgen")["embeddingFixture"]
	assert.Equal(t, "Base is embedded directly in embeddingFixture.", docs["docFixture"])

	// The AST key must match the field name reflect assigns to the embedded
	// field, since that is what FieldDoc will be queried with.
	assert.Equal(t, "docFixture", reflect.TypeOf(embeddingFixture{}).Field(0).Name)
}
