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

package fs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// FileInfo: information about a file or directory returned by stat.
type FileInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir"`
	Mode    uint32 `json:"mode"`
	ModTime int64  `json:"mod_time"`
}

// readFile: reads the entire contents of a file; raises on error.
func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(content), nil
}

// writeFile: writes content to a file; raises on error.
func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// exists: reports whether a path exists.
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// mkdir: creates a directory; raises on error.
func mkdir(path string) error {
	if err := os.Mkdir(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}

// mkdirAll: creates a directory and all necessary parents; raises on error.
func mkdirAll(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}
	return nil
}

// remove: removes a file or empty directory; raises on error.
func remove(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to remove: %w", err)
	}
	return nil
}

// removeAll: removes a path and all its contents recursively; raises on error.
func removeAll(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to remove recursively: %w", err)
	}
	return nil
}

// list: lists all entries in a directory; raises on error.
func list(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list directory: %w", err)
	}
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name()
	}
	return names, nil
}

// stat: returns information about a file or directory; raises on error.
func stat(path string) (FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return FileInfo{}, fmt.Errorf("failed to stat: %w", err)
	}
	return FileInfo{
		Name:    filepath.Base(path),
		Size:    info.Size(),
		IsDir:   info.IsDir(),
		Mode:    uint32(info.Mode()), //nolint:gosec
		ModTime: info.ModTime().Unix(),
	}, nil
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("fs", "file system utilities")
	m.Fn("read_file", readFile, "reads the entire contents of a file, raises on error",
		luareg.Args("path"))
	m.Fn("write_file", writeFile, "writes content to a file, raises on error",
		luareg.Args("path", "content"))
	m.Fn("exists", exists, "reports whether a path exists",
		luareg.Args("path"))
	m.Fn("mkdir", mkdir, "creates a directory, raises on error",
		luareg.Args("path"))
	m.Fn("mkdir_all", mkdirAll, "creates a directory and all parents, raises on error",
		luareg.Args("path"))
	m.Fn("remove", remove, "removes a file or empty directory, raises on error",
		luareg.Args("path"))
	m.Fn("remove_all", removeAll, "removes a path and all its contents, raises on error",
		luareg.Args("path"))
	m.Fn("list", list, "lists all entries in a directory, raises on error",
		luareg.Args("path"))
	m.Fn("stat", stat, "returns file information, raises on error",
		luareg.Args("path"))
	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("fs", fs.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
