#!/bin/bash

# Check MIT license headers in Go and Lua files

EXIT_CODE=0

# Expected headers
GO_HEADER="// Copyright (c) 2024-2025 Thomas Maurice"
LUA_HEADER="-- Copyright (c) 2024-2025 Thomas Maurice"

# Check Go files
echo "Checking Go files for MIT license headers..."
for file in "$@"; do
    if [[ "$file" == *.go ]]; then
        if ! head -n 1 "$file" | grep -q "^$GO_HEADER"; then
            echo "ERROR: $file is missing MIT license header"
            EXIT_CODE=1
        fi
    fi
done

# Check Lua files (excluding generated files)
echo "Checking Lua files for MIT license headers..."
for file in "$@"; do
    if [[ "$file" == *.lua ]] && [[ "$file" != *.gen.lua ]]; then
        if ! head -n 1 "$file" | grep -q "^$LUA_HEADER"; then
            echo "ERROR: $file is missing MIT license header"
            EXIT_CODE=1
        fi
    fi
done

if [ $EXIT_CODE -eq 0 ]; then
    echo "✓ All files have proper MIT license headers"
fi

exit $EXIT_CODE
