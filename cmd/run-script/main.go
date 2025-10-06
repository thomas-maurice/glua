package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/thomas-maurice/glua/example/sample"
	"github.com/thomas-maurice/glua/pkg/glua"
	"github.com/thomas-maurice/glua/pkg/modules/kubernetes"
	lua "github.com/yuin/gopher-lua"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	// Check if script file is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run run_script.go <script.lua>")
		fmt.Println("\nAvailable example scripts:")
		fmt.Println("  scripts/01_basic_pod_info.lua")
		fmt.Println("  scripts/02_resource_limits.lua")
		fmt.Println("  scripts/03_policy_validation.lua")
		fmt.Println("  scripts/04_environment_vars.lua")
		fmt.Println("  scripts/05_timestamp_operations.lua")
		fmt.Println("  scripts/06_multi_container_analysis.lua")
		fmt.Println("  scripts/07_json_export.lua")
		os.Exit(1)
	}

	scriptPath := os.Args[1]

	// Generate stubs first
	fmt.Println("=== Generating LSP Type Annotations ===")
	if err := generateStubs(); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating stubs: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Generated annotations.gen.lua")
	fmt.Println()

	// Initialize Lua state
	L := lua.NewState()
	defer L.Close()

	// Load kubernetes module
	L.PreloadModule("kubernetes", kubernetes.Loader)

	// Create translator
	translator := glua.NewTranslator()

	// Create sample pod
	pod := sample.GetPod()

	// Convert pod to Lua
	luaTable, err := translator.ToLua(L, pod)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to convert pod to Lua: %v\n", err)
		os.Exit(1)
	}

	// Make pod available to Lua scripts
	L.SetGlobal("myPod", luaTable)
	L.SetGlobal("originalPod", luaTable) // Alias for compatibility

	// Read and execute script
	scriptBytes, err := os.ReadFile(scriptPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read script %s: %v\n", scriptPath, err)
		os.Exit(1)
	}

	fmt.Printf("=== Running Script: %s ===\n", filepath.Base(scriptPath))
	fmt.Println()

	err = L.DoString(string(scriptBytes))
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n❌ Lua script error: %v\n", err)
		os.Exit(1)
	}

	// Check if script exported a modified pod
	modifiedTable := L.GetGlobal("modifiedPod")
	if modifiedTable != lua.LNil {
		fmt.Println()
		fmt.Println("=== Script Modified Pod ===")
		var reconstructedPod corev1.Pod
		err = translator.FromLua(L, modifiedTable, &reconstructedPod)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to convert modified pod back to Go: %v\n", err)
		} else {
			fmt.Println("✓ Successfully converted modified pod back to Go")
			fmt.Printf("  Pod: %s/%s\n", reconstructedPod.Namespace, reconstructedPod.Name)
			fmt.Printf("  Containers: %d\n", len(reconstructedPod.Spec.Containers))
		}
	}

	// Check if script exported a report
	exportedReport := L.GetGlobal("exportedReport")
	if exportedReport != lua.LNil {
		fmt.Println()
		fmt.Println("=== Exported Data ===")
		fmt.Println("Script exported 'exportedReport' - accessible from Go")
	}

	fmt.Println()
	fmt.Println("✓ Script execution completed successfully")
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
