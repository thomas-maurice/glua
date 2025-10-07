package main

import (
	"fmt"
	"os"
	"time"

	"github.com/thomas-maurice/glua/example/sample"
	"github.com/thomas-maurice/glua/pkg/glua"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	lua "github.com/yuin/gopher-lua"
)

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
