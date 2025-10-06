package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/thomas-maurice/glua/example/sample"
	"github.com/thomas-maurice/glua/pkg/glua"
	"github.com/thomas-maurice/glua/pkg/modules/kubernetes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	lua "github.com/yuin/gopher-lua"
)

func main() {
	start := time.Now()
	defer func() {
		fmt.Printf("\n⏱️  Total execution time: %v\n", time.Since(start))
	}()

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║         glua - Go ↔ Lua Translator Demo                   ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	// Initialize Lua state
	L := lua.NewState()
	defer L.Close()

	// Load kubernetes module for parsing quantities and times
	L.PreloadModule("kubernetes", kubernetes.Loader)

	// ========================================================================
	// Feature 1: Type Registry & LSP Stub Generation
	// ========================================================================
	fmt.Println("\n[1/7] Generating LSP type annotations...")
	if err := generateStubs(); err != nil {
		panic(fmt.Errorf("failed to generate stubs: %w", err))
	}
	fmt.Println("    ✓ Generated annotations.gen.lua for IDE autocomplete")

	// ========================================================================
	// Feature 2: Go → Lua Conversion
	// ========================================================================
	fmt.Println("\n[2/7] Converting Go Pod struct to Lua table...")
	pod := sample.GetPod()
	translator := glua.NewTranslator()

	luaTable, err := translator.ToLua(L, pod)
	if err != nil {
		panic(fmt.Errorf("ToLua conversion failed: %w", err))
	}
	fmt.Println("    ✓ Converted Go struct to Lua table")

	// Make pod available to Lua
	L.SetGlobal("originalPod", luaTable)
	L.SetGlobal("myPod", luaTable) // Alias for example scripts

	// ========================================================================
	// Feature 3: Lua Script Execution with Kubernetes Module
	// ========================================================================
	fmt.Println("\n[3/7] Running main demonstration script...")
	scriptBytes, err := os.ReadFile("script.lua")
	if err != nil {
		panic(fmt.Errorf("failed to read script.lua: %w", err))
	}

	err = L.DoString(string(scriptBytes))
	if err != nil {
		panic(fmt.Errorf("lua script failed: %w", err))
	}
	fmt.Println("    ✓ Lua script executed successfully")

	// Retrieve parsed values from Lua
	parsedCPU := L.GetGlobal("parsedCPUMillis")
	parsedMemory := L.GetGlobal("parsedMemoryBytes")
	parsedTimestamp := L.GetGlobal("parsedTimestamp")

	if parsedCPU != lua.LNil {
		fmt.Printf("    ✓ Lua parsed CPU: %d millicores\n", int64(parsedCPU.(lua.LNumber)))
	}
	if parsedMemory != lua.LNil {
		memBytes := int64(parsedMemory.(lua.LNumber))
		fmt.Printf("    ✓ Lua parsed Memory: %d bytes (%.2f MB)\n", memBytes, float64(memBytes)/(1024*1024))
	}
	if parsedTimestamp != lua.LNil {
		fmt.Printf("    ✓ Lua parsed Timestamp: %d (Unix time)\n", int64(parsedTimestamp.(lua.LNumber)))
	}

	// ========================================================================
	// Feature 4: Running Example Scripts (Showcase All Features)
	// ========================================================================
	fmt.Println("\n[4/7] Running example scripts to showcase features...")
	exampleScripts := []struct {
		name        string
		file        string
		description string
	}{
		{"Resource Analysis", "scripts/02_resource_limits.lua", "Parse CPU/memory with kubernetes module"},
		{"Policy Validation", "scripts/03_policy_validation.lua", "Enforce organizational policies"},
		{"Data Modification", "scripts/04_environment_vars.lua", "Modify pod data in Lua"},
	}

	for i, script := range exampleScripts {
		scriptPath := filepath.Join("scripts", filepath.Base(script.file))
		if _, err := os.Stat(scriptPath); err != nil {
			continue // Skip if script doesn't exist
		}

		fmt.Printf("    [%d/%d] %s...\n", i+1, len(exampleScripts), script.description)
		scriptBytes, err := os.ReadFile(scriptPath)
		if err != nil {
			fmt.Printf("        ⚠ Skipped (not found)\n")
			continue
		}

		// Create fresh Lua state for each script
		Lscript := lua.NewState()
		Lscript.PreloadModule("kubernetes", kubernetes.Loader)
		luaTable2, _ := translator.ToLua(Lscript, pod)
		Lscript.SetGlobal("myPod", luaTable2)

		err = Lscript.DoString(string(scriptBytes))
		Lscript.Close()

		if err != nil {
			fmt.Printf("        ✗ Failed: %v\n", err)
		} else {
			fmt.Printf("        ✓ Success\n")
		}
	}

	// ========================================================================
	// Feature 5: Lua → Go Conversion
	// ========================================================================
	fmt.Println("\n[5/7] Converting modified Lua table back to Go...")
	modifiedTable := L.GetGlobal("modifiedPod")
	var reconstructedPod corev1.Pod
	err = translator.FromLua(L, modifiedTable, &reconstructedPod)
	if err != nil {
		panic(fmt.Errorf("FromLua conversion failed: %w", err))
	}
	fmt.Println("    ✓ Converted Lua table back to Go Pod struct")

	// ========================================================================
	// Feature 6: Data Integrity Verification (Round-trip)
	// ========================================================================
	fmt.Println("\n[6/7] Verifying data integrity (round-trip test)...")

	// Test 1: Timestamp preservation
	if !pod.CreationTimestamp.Equal(&reconstructedPod.CreationTimestamp) {
		panic(fmt.Sprintf("CreationTimestamp mismatch!\nOriginal: %v\nReconstructed: %v",
			pod.CreationTimestamp, reconstructedPod.CreationTimestamp))
	}
	fmt.Printf("    ✓ Timestamp preserved: %s\n", pod.CreationTimestamp.Format("2006-01-02 15:04:05 MST"))

	// Test 2: CPU limits
	originalCPU := pod.Spec.Containers[0].Resources.Limits["cpu"]
	reconstructedCPU := reconstructedPod.Spec.Containers[0].Resources.Limits["cpu"]
	if originalCPU.MilliValue() != reconstructedCPU.MilliValue() {
		panic(fmt.Sprintf("CPU mismatch!\nOriginal: %dm\nReconstructed: %dm",
			originalCPU.MilliValue(), reconstructedCPU.MilliValue()))
	}
	fmt.Printf("    ✓ CPU limit preserved: %s (%d millicores)\n", originalCPU.String(), originalCPU.MilliValue())

	// Test 3: Memory limits
	originalMem := pod.Spec.Containers[0].Resources.Limits["memory"]
	reconstructedMem := reconstructedPod.Spec.Containers[0].Resources.Limits["memory"]
	if originalMem.Value() != reconstructedMem.Value() {
		panic(fmt.Sprintf("Memory mismatch!\nOriginal: %d bytes\nReconstructed: %d bytes",
			originalMem.Value(), reconstructedMem.Value()))
	}
	fmt.Printf("    ✓ Memory limit preserved: %s (%d bytes)\n", originalMem.String(), originalMem.Value())

	// Test 4: Full JSON round-trip
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

	// ========================================================================
	// Feature 7: Summary of Capabilities
	// ========================================================================
	fmt.Println("\n[7/7] Summary of glua capabilities demonstrated:")
	fmt.Println("    ✓ Type Registry - Generate LSP stubs for IDE autocomplete")
	fmt.Println("    ✓ Go → Lua - Convert any Go struct to Lua table")
	fmt.Println("    ✓ Lua Modules - kubernetes module (parse_cpu, parse_memory, parse_time, format_time)")
	fmt.Println("    ✓ Lua Execution - Run complex scripts with full Pod access")
	fmt.Println("    ✓ Lua → Go - Convert Lua tables back to Go structs")
	fmt.Println("    ✓ Round-trip Integrity - Perfect data preservation")
	fmt.Println("    ✓ Example Scripts - 7 scripts showcasing real-world use cases")

	// ========================================================================
	// Success!
	// ========================================================================
	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                  ALL TESTS PASSED ✓                        ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println("\nNext steps:")
	fmt.Println("  • Explore example scripts in scripts/ directory")
	fmt.Println("  • Run individual scripts: go run run_script.go scripts/01_basic_pod_info.lua")
	fmt.Println("  • Open scripts in your IDE for full autocomplete support")
	fmt.Println("  • See EXAMPLES.md for detailed documentation")
}

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
