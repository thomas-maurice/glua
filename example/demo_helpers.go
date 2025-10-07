package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/thomas-maurice/glua/pkg/glua"
	jsonmodule "github.com/thomas-maurice/glua/pkg/modules/json"
	"github.com/thomas-maurice/glua/pkg/modules/kubernetes"
	spewmodule "github.com/thomas-maurice/glua/pkg/modules/spew"
	corev1 "k8s.io/api/core/v1"

	lua "github.com/yuin/gopher-lua"
)

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
	fmt.Println("  • Run individual scripts: go run run_script.go scripts/01_basic_pod_info.lua")
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
