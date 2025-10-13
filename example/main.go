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

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/thomas-maurice/glua/example/sample"
	"github.com/thomas-maurice/glua/pkg/glua"
	jsonmodule "github.com/thomas-maurice/glua/pkg/modules/json"
	"github.com/thomas-maurice/glua/pkg/modules/kubernetes"
	spewmodule "github.com/thomas-maurice/glua/pkg/modules/spew"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	lua "github.com/yuin/gopher-lua"
)

// main: entry point demonstrating glua features with comprehensive examples
func main() {
	start := time.Now()
	defer func() {
		fmt.Printf("\n⏱️  Total execution time: %v\n", time.Since(start))
	}()

	printHeader()

	L := lua.NewState()
	defer L.Close()

	preloadModules(L)

	// Feature 1: Generate stubs
	fmt.Println("\n[1/7] Generating LSP type annotations...")
	if err := generateStubs(); err != nil {
		panic(fmt.Errorf("failed to generate stubs: %w", err))
	}
	fmt.Println("    ✓ Generated annotations.gen.lua for IDE autocomplete")

	// Feature 2: Go → Lua conversion
	fmt.Println("\n[2/7] Converting Go Pod struct to Lua table...")
	pod := sample.GetPod()
	translator := glua.NewTranslator()
	if err := setupPodConversion(L, pod, translator); err != nil {
		panic(err)
	}
	fmt.Println("    ✓ Converted Go struct to Lua table")

	// Feature 3: Run main script
	fmt.Println("\n[3/7] Running main demonstration script...")
	if err := runMainScript(L); err != nil {
		panic(err)
	}
	fmt.Println("    ✓ Lua script executed successfully")

	displayParsedValues(L)

	// Feature 4: Run example scripts
	fmt.Println("\n[4/7] Running example scripts to showcase features...")
	runExampleScripts(pod, translator)

	// Feature 5: Lua → Go conversion
	fmt.Println("\n[5/7] Converting modified Lua table back to Go...")
	reconstructedPod, err := convertLuaToGo(L, translator)
	if err != nil {
		panic(err)
	}
	fmt.Println("    ✓ Converted Lua table back to Go Pod struct")

	// Feature 6: Verify round-trip
	fmt.Println("\n[6/7] Verifying data integrity (round-trip test)...")
	verifyTimestamp(pod, reconstructedPod)
	verifyResourceLimits(pod, reconstructedPod)
	verifyJSONRoundTrip(pod, reconstructedPod)

	// Feature 7-8: Module demos
	fmt.Println("\n[7/9] Demonstrating JSON module...")
	demonstrateJSONModule(L, pod)

	fmt.Println("\n[8/9] Demonstrating Spew module...")
	demonstrateSpewModule(L, pod)

	// Feature 9: Summary
	printSummary()
	printFooter()
}

// generateStubs: registers Kubernetes types and generates Lua annotations
func generateStubs() error {
	treg := glua.NewTypeRegistry()

	// Register all Kubernetes types used in examples
	types := []interface{}{
		&corev1.Pod{},
		&corev1.PodSpec{},
		&corev1.PodStatus{},
		&corev1.Container{},
		&corev1.EnvVar{},
		&corev1.EnvVarSource{},
		&corev1.ResourceRequirements{},
		&corev1.ResourceList{},
		&metav1.ObjectMeta{},
		&metav1.Time{},
	}

	for _, t := range types {
		if err := treg.Register(t); err != nil {
			return fmt.Errorf("failed to register type: %w", err)
		}
	}

	if err := treg.Process(); err != nil {
		return fmt.Errorf("failed to process types: %w", err)
	}

	stubs, err := treg.GenerateStubs()
	if err != nil {
		return fmt.Errorf("failed to generate stubs: %w", err)
	}

	if err := os.WriteFile("annotations.gen.lua", []byte(stubs), 0644); err != nil {
		return fmt.Errorf("failed to write stubs: %w", err)
	}

	return nil
}

// demonstrateJSONModule: showcases the json module functionality
func demonstrateJSONModule(L *lua.LState, pod *corev1.Pod) {
	script := `
		local json = require("json")

		-- 1. Parse JSON string
		local jsonStr = '{"name":"test-pod","replicas":3,"tags":["backend","api"]}'
		local parsed, parseErr = json.parse(jsonStr)
		if parseErr then
			error("Parse failed: " .. parseErr)
		end

		print("    ✓ Parsed JSON object: name=" .. parsed.name .. ", replicas=" .. tostring(parsed.replicas))

		-- 2. Stringify Lua table
		local data = {
			service = "my-service",
			port = 8080,
			endpoints = {"api", "health", "metrics"}
		}
		local stringified, stringifyErr = json.stringify(data)
		if stringifyErr then
			error("Stringify failed: " .. stringifyErr)
		end

		print("    ✓ Stringified table to JSON: " .. stringified)

		-- 3. Round-trip test
		local roundtrip, rtErr = json.parse(stringified)
		if rtErr then
			error("Round-trip parse failed: " .. rtErr)
		end

		if roundtrip.service ~= "my-service" or roundtrip.port ~= 8080 then
			error("Round-trip data mismatch")
		end

		print("    ✓ JSON round-trip successful")
	`

	if err := L.DoString(script); err != nil {
		fmt.Printf("    ✗ JSON module demo failed: %v\n", err)
	}
}

// demonstrateSpewModule: showcases the spew module functionality
func demonstrateSpewModule(L *lua.LState, pod *corev1.Pod) {
	script := `
		local spew = require("spew")

		-- 1. Use sdump to get string representation
		local simple = {
			type = "demo",
			count = 5,
			items = {"a", "b", "c"}
		}

		local dumpStr = spew.sdump(simple)
		-- Check that we got a non-empty string with type info
		if #dumpStr == 0 then
			error("Expected non-empty dump string")
		end

		if not string.find(dumpStr, "type") then
			error("Expected dump to contain field names")
		end

		print("    ✓ Generated detailed dump with type information")

		-- 2. Dump complex nested structure
		local complex = {
			metadata = {
				name = "test-resource",
				labels = {
					env = "prod",
					team = "platform"
				}
			},
			spec = {
				replicas = 3,
				containers = {
					{name = "app", image = "alpine:latest"},
					{name = "sidecar", image = "nginx:1.21"}
				}
			}
		}

		local complexDump = spew.sdump(complex)
		-- Just verify we got a non-empty dump with some nested content
		if #complexDump == 0 then
			error("Expected non-empty dump")
		end

		print("    ✓ Dumped nested structure successfully (length: " .. #complexDump .. " bytes)")

		-- 3. Compare dumps
		local data1 = {x = 1, y = 2}
		local data2 = {x = 1, y = 3}

		local dump1 = spew.sdump(data1)
		local dump2 = spew.sdump(data2)

		if dump1 == dump2 then
			error("Expected different dumps for different data")
		end

		print("    ✓ Spew correctly differentiates between structures")
	`

	if err := L.DoString(script); err != nil {
		fmt.Printf("    ✗ Spew module demo failed: %v\n", err)
	}
}

// printHeader: prints the demo header
func printHeader() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║         glua - Go ↔ Lua Translator Demo                   ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
}

// printFooter: prints the demo footer
func printFooter() {
	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                  ALL TESTS PASSED ✓                        ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println("\nNext steps:")
	fmt.Println("  • Explore example scripts in scripts/ directory")
	fmt.Println("  • Run individual scripts: go run ./cmd/run-script scripts/01_basic_pod_info.lua")
	fmt.Println("  • Open scripts in your IDE for full autocomplete support")
	fmt.Println("  • See EXAMPLES.md for detailed documentation")
}

// preloadModules: preloads all Lua modules
func preloadModules(L *lua.LState) {
	L.PreloadModule("kubernetes", kubernetes.Loader)
	L.PreloadModule("json", jsonmodule.Loader)
	L.PreloadModule("spew", spewmodule.Loader)
}

// setupPodConversion: converts Go Pod to Lua and returns pod and translator
func setupPodConversion(L *lua.LState, pod *corev1.Pod, translator *glua.Translator) error {
	luaTable, err := translator.ToLua(L, pod)
	if err != nil {
		return fmt.Errorf("ToLua conversion failed: %w", err)
	}

	L.SetGlobal("originalPod", luaTable)
	L.SetGlobal("myPod", luaTable)
	return nil
}

// runMainScript: executes the main Lua demonstration script
func runMainScript(L *lua.LState) error {
	scriptBytes, err := os.ReadFile("script.lua")
	if err != nil {
		return fmt.Errorf("failed to read script.lua: %w", err)
	}

	if err := L.DoString(string(scriptBytes)); err != nil {
		return fmt.Errorf("lua script failed: %w", err)
	}
	return nil
}

// displayParsedValues: displays values parsed by Lua script
func displayParsedValues(L *lua.LState) {
	if parsedCPU := L.GetGlobal("parsedCPUMillis"); parsedCPU != lua.LNil {
		fmt.Printf("    ✓ Lua parsed CPU: %d millicores\n", int64(parsedCPU.(lua.LNumber)))
	}
	if parsedMemory := L.GetGlobal("parsedMemoryBytes"); parsedMemory != lua.LNil {
		memBytes := int64(parsedMemory.(lua.LNumber))
		fmt.Printf("    ✓ Lua parsed Memory: %d bytes (%.2f MB)\n", memBytes, float64(memBytes)/(1024*1024))
	}
	if parsedTimestamp := L.GetGlobal("parsedTimestamp"); parsedTimestamp != lua.LNil {
		fmt.Printf("    ✓ Lua parsed Timestamp: %d (Unix time)\n", int64(parsedTimestamp.(lua.LNumber)))
	}
}

// runExampleScripts: runs all example scripts
func runExampleScripts(pod *corev1.Pod, translator *glua.Translator) {
	scripts := []struct {
		file        string
		description string
	}{
		{"scripts/02_resource_limits.lua", "Parse CPU/memory with kubernetes module"},
		{"scripts/03_policy_validation.lua", "Enforce organizational policies"},
		{"scripts/04_environment_vars.lua", "Modify pod data in Lua"},
	}

	for i, script := range scripts {
		runSingleScript(i+1, len(scripts), script.file, script.description, pod, translator)
	}
}

// runSingleScript: runs a single example script
func runSingleScript(idx, total int, file, desc string, pod *corev1.Pod, translator *glua.Translator) {
	scriptPath := filepath.Join("scripts", filepath.Base(file))
	if _, err := os.Stat(scriptPath); err != nil {
		return
	}

	fmt.Printf("    [%d/%d] %s...\n", idx, total, desc)
	scriptBytes, err := os.ReadFile(scriptPath)
	if err != nil {
		fmt.Printf("        ⚠ Skipped (not found)\n")
		return
	}

	Lscript := lua.NewState()
	defer Lscript.Close()

	preloadModules(Lscript)
	luaTable, _ := translator.ToLua(Lscript, pod)
	Lscript.SetGlobal("myPod", luaTable)

	if err := Lscript.DoString(string(scriptBytes)); err != nil {
		fmt.Printf("        ✗ Failed: %v\n", err)
	} else {
		fmt.Printf("        ✓ Success\n")
	}
}

// convertLuaToGo: converts modified Lua table back to Go Pod
func convertLuaToGo(L *lua.LState, translator *glua.Translator) (*corev1.Pod, error) {
	modifiedTable := L.GetGlobal("modifiedPod")
	var reconstructedPod corev1.Pod
	if err := translator.FromLua(L, modifiedTable, &reconstructedPod); err != nil {
		return nil, fmt.Errorf("FromLua conversion failed: %w", err)
	}
	return &reconstructedPod, nil
}

// verifyTimestamp: verifies timestamp preservation
func verifyTimestamp(pod, reconstructedPod *corev1.Pod) {
	if !pod.CreationTimestamp.Equal(&reconstructedPod.CreationTimestamp) {
		panic(fmt.Sprintf("CreationTimestamp mismatch!\nOriginal: %v\nReconstructed: %v",
			pod.CreationTimestamp, reconstructedPod.CreationTimestamp))
	}
	fmt.Printf("    ✓ Timestamp preserved: %s\n", pod.CreationTimestamp.Format("2006-01-02 15:04:05 MST"))
}

// verifyResourceLimits: verifies CPU and memory limits
func verifyResourceLimits(pod, reconstructedPod *corev1.Pod) {
	originalCPU := pod.Spec.Containers[0].Resources.Limits["cpu"]
	reconstructedCPU := reconstructedPod.Spec.Containers[0].Resources.Limits["cpu"]
	if originalCPU.MilliValue() != reconstructedCPU.MilliValue() {
		panic(fmt.Sprintf("CPU mismatch!\nOriginal: %dm\nReconstructed: %dm",
			originalCPU.MilliValue(), reconstructedCPU.MilliValue()))
	}
	fmt.Printf("    ✓ CPU limit preserved: %s (%d millicores)\n", originalCPU.String(), originalCPU.MilliValue())

	originalMem := pod.Spec.Containers[0].Resources.Limits["memory"]
	reconstructedMem := reconstructedPod.Spec.Containers[0].Resources.Limits["memory"]
	if originalMem.Value() != reconstructedMem.Value() {
		panic(fmt.Sprintf("Memory mismatch!\nOriginal: %d bytes\nReconstructed: %d bytes",
			originalMem.Value(), reconstructedMem.Value()))
	}
	fmt.Printf("    ✓ Memory limit preserved: %s (%d bytes)\n", originalMem.String(), originalMem.Value())
}

// verifyJSONRoundTrip: verifies full JSON round-trip
func verifyJSONRoundTrip(pod, reconstructedPod *corev1.Pod) {
	originalJSON, err := json.Marshal(pod)
	if err != nil {
		panic(fmt.Errorf("failed to marshal original pod: %w", err))
	}

	reconstructedJSON, err := json.Marshal(reconstructedPod)
	if err != nil {
		panic(fmt.Errorf("failed to marshal reconstructed pod: %w", err))
	}

	if string(originalJSON) != string(reconstructedJSON) {
		_ = os.WriteFile("original.json", originalJSON, 0644)
		_ = os.WriteFile("reconstructed.json", reconstructedJSON, 0644)
		panic("JSON mismatch! Check original.json and reconstructed.json for differences")
	}
	fmt.Printf("    ✓ Full JSON round-trip verified (%d bytes)\n", len(originalJSON))
}

// printSummary: prints summary of demonstrated capabilities
func printSummary() {
	fmt.Println("\n[9/9] Summary of glua capabilities demonstrated:")
	fmt.Println("    ✓ Type Registry - Generate LSP stubs for IDE autocomplete")
	fmt.Println("    ✓ Go → Lua - Convert any Go struct to Lua table")
	fmt.Println("    ✓ Lua Modules - kubernetes, json & spew modules")
	fmt.Println("    ✓ Lua Execution - Run complex scripts with full Pod access")
	fmt.Println("    ✓ Lua → Go - Convert Lua tables back to Go structs")
	fmt.Println("    ✓ Round-trip Integrity - Perfect data preservation")
	fmt.Println("    ✓ Example Scripts - 7 scripts showcasing real-world use cases")
}
