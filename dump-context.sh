#!/bin/bash

paths=(
	"cmd/wasm/main.go" 
	"cmd/wasm/static/main.js" 
	"cmd/wasm/universe.go"
)
allpaths=(
	"cmd/server/main.go"
	"cmd/wasm/main.go"
	"cmd/wasm/universe.go"
	"cmd/wasm/static/main.js"
	"cmd/wasm/static/stats.js"
	"cmd/wasm/color.go"
)

for path in "${paths[@]}"; do
    echo "--- File: $path ---"
    echo "\`\`\`"
    cat "$path"
    echo "\`\`\`"
    echo ""
done
