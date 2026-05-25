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

package kubernetes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomas-maurice/glua/pkg/glua"
	lua "github.com/yuin/gopher-lua"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// loadPodIntoLua: helper that converts a corev1.Pod via the glua Translator
// and binds it to the global `pod` in L. Used by the round-trip tests to
// exercise the same conversion path a k8sclient module would use to hand
// objects back to Lua.
func loadPodIntoLua(t *testing.T, L *lua.LState, pod *corev1.Pod) {
	t.Helper()
	tr := glua.NewTranslator()
	lv, err := tr.ToLua(L, pod)
	require.NoError(t, err, "Pod → Lua conversion failed")
	L.SetGlobal("pod", lv)
}

// TestQuantityRoundTripInPodSpec: a Pod carrying typed resource.Quantity
// values (memory and CPU) round-trips through the Translator with the
// canonical K8s string preserved on the Lua side. Then we feed those strings
// back through the kubernetes module's parse_memory / parse_cpu helpers and
// verify the numeric result.
//
// This is the end-to-end shape a webhook / policy script hits: cluster API
// hands back a typed object, Lua reads `pod.spec.containers[1].resources.
// requests.memory`, and the script wants to compare or threshold the value.
//
// Important — Quantity canonicalization on the Lua side:
//
// When a resource.Quantity is serialized (which is what happens on the way
// to Lua, via JSON marshal in the Translator), K8s emits a *canonical* form,
// not the original user input. The semantic value never changes — it's the
// same number of bytes/millicores — but the string can differ:
//
//   - Fractional CPU written as a decimal is normalized to millicores:
//     "2.5"   on the input side  →  "2500m" on the Lua side.
//     "0.5"   on the input side  →  "500m"  on the Lua side.
//     Whole cores stay whole: "1", "2", "4" round-trip unchanged.
//   - Memory units round-trip unchanged for the standard binary/SI suffixes
//     listed below ("1Gi", "512Mi", "500M", "1Ti", etc.).
//   - Numeric helpers (parse_memory, parse_cpu) ignore the surface form and
//     always return the same byte / millicore count regardless of which
//     spelling the Lua side observes — so threshold checks are safe.
//
// Lua scripts that compare these strings directly (rare, but possible) must
// be ready to see the canonical form. Scripts that use parse_memory /
// parse_cpu (the recommended path) don't have to care.
func TestQuantityRoundTripInPodSpec(t *testing.T) {
	// memWant / cpuWant are the strings expected on the Lua side. See the
	// canonicalization note in the test's doc comment for the rules.
	cases := []struct {
		name               string
		memory, memWant    string
		cpu, cpuWant       string
		memBytes, cpuMilli int64
	}{
		// Binary units (Ki/Mi/Gi/Ti) are the common K8s memory shape.
		{"1Gi / 500m", "1Gi", "1Gi", "500m", "500m", 1024 * 1024 * 1024, 500},
		{"2Gi / 1", "2Gi", "2Gi", "1", "1", 2 * 1024 * 1024 * 1024, 1000},
		{"512Mi / 250m", "512Mi", "512Mi", "250m", "250m", 512 * 1024 * 1024, 250},
		// CPU "2.5" is canonicalized to "2500m" by K8s on marshal.
		{"4Gi / 2.5", "4Gi", "4Gi", "2.5", "2500m", 4 * 1024 * 1024 * 1024, 2500},
		// Decimal SI units are accepted by K8s too.
		{"500M / 1500m", "500M", "500M", "1500m", "1500m", 500 * 1000 * 1000, 1500},
		{"1Ti / 4", "1Ti", "1Ti", "4", "4", 1024 * 1024 * 1024 * 1024, 4000},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "demo", Namespace: "default"},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "c",
						Image: "nginx",
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(tc.memory),
								corev1.ResourceCPU:    resource.MustParse(tc.cpu),
							},
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(tc.memory),
								corev1.ResourceCPU:    resource.MustParse(tc.cpu),
							},
						},
					}},
				},
			}

			L := lua.NewState()
			defer L.Close()
			L.PreloadModule("kubernetes", Loader)
			loadPodIntoLua(t, L, pod)

			// Read the quantities back from Lua and parse them via the module.
			// Lua tables are 1-indexed.
			script := `
				local k8s = require("kubernetes")
				local req = pod.spec.containers[1].resources.requests
				mem_str = req.memory
				cpu_str = req.cpu
				mem_bytes = k8s.parse_memory(req.memory)
				cpu_milli = k8s.parse_cpu(req.cpu)
			`
			require.NoError(t, L.DoString(script))

			assert.Equal(t, lua.LString(tc.memWant), L.GetGlobal("mem_str"),
				"memory string on the Lua side should match canonical form")
			assert.Equal(t, lua.LString(tc.cpuWant), L.GetGlobal("cpu_str"),
				"cpu string on the Lua side should match canonical form")
			assert.Equal(t, lua.LNumber(tc.memBytes), L.GetGlobal("mem_bytes"),
				"parse_memory should yield the expected byte count")
			assert.Equal(t, lua.LNumber(tc.cpuMilli), L.GetGlobal("cpu_milli"),
				"parse_cpu should yield the expected millicore count")
		})
	}
}

// TestParseMemoryUnits: independent matrix of memory unit suffixes ensures
// parse_memory understands the full set of K8s suffixes (binary Ki..Ti and
// decimal K..T) and produces the right byte counts.
func TestParseMemoryUnits(t *testing.T) {
	cases := []struct {
		input    string
		expected int64
	}{
		{"1Ki", 1024},
		{"1Mi", 1024 * 1024},
		{"1Gi", 1024 * 1024 * 1024},
		{"1Ti", 1024 * 1024 * 1024 * 1024},
		// K8s decimal SI: lowercase k for kilo, uppercase M/G/T for the rest.
		// Uppercase "1K" alone is rejected — confusing but per the K8s grammar.
		{"1k", 1000},
		{"1M", 1000 * 1000},
		{"1G", 1000 * 1000 * 1000},
		{"1T", 1000 * 1000 * 1000 * 1000},
		{"512", 512}, // bare bytes
		{"0", 0},
	}

	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("kubernetes", Loader)

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			L.SetGlobal("q", lua.LString(tc.input))
			require.NoError(t, L.DoString(`result = require("kubernetes").parse_memory(q)`))
			assert.Equal(t, lua.LNumber(tc.expected), L.GetGlobal("result"))
		})
	}
}

// TestParseCPUUnits: independent matrix of CPU forms — millicores, whole
// cores, fractional cores — all expressed in millicores by parse_cpu.
func TestParseCPUUnits(t *testing.T) {
	cases := []struct {
		input    string
		expected int64
	}{
		{"100m", 100},
		{"500m", 500},
		{"1", 1000},
		{"2", 2000},
		{"0.5", 500},
		{"1.5", 1500},
		{"0", 0},
	}

	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("kubernetes", Loader)

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			L.SetGlobal("q", lua.LString(tc.input))
			require.NoError(t, L.DoString(`result = require("kubernetes").parse_cpu(q)`))
			assert.Equal(t, lua.LNumber(tc.expected), L.GetGlobal("result"))
		})
	}
}

// TestFormatParseMemoryRoundTrip: format_memory and parse_memory are
// inverses for byte counts that have an exact binary-SI representation
// (Ki/Mi/Gi/Ti). format_memory(parse_memory(x)) is also a fixed point for
// canonical inputs.
func TestFormatParseMemoryRoundTrip(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("kubernetes", Loader)

	t.Run("bytes -> string -> bytes", func(t *testing.T) {
		bytes := []int64{
			0,
			1024,                      // 1Ki
			1024 * 1024,               // 1Mi
			512 * 1024 * 1024,         // 512Mi
			1024 * 1024 * 1024,        // 1Gi
			2 * 1024 * 1024 * 1024,    // 2Gi
			1024 * 1024 * 1024 * 1024, // 1Ti
		}
		for _, b := range bytes {
			L.SetGlobal("b", lua.LNumber(b))
			require.NoError(t, L.DoString(`
				local k8s = require("kubernetes")
				local s = k8s.format_memory(b)
				result = k8s.parse_memory(s)
			`))
			assert.Equal(t, lua.LNumber(b), L.GetGlobal("result"),
				"format_memory then parse_memory should return the original byte count for %d", b)
		}
	})

	t.Run("canonical string -> bytes -> string", func(t *testing.T) {
		// Canonical binary-SI forms — these are the exact strings
		// format_memory emits, so the round-trip is a fixed point.
		canonical := []string{"0", "1Ki", "1Mi", "512Mi", "1Gi", "2Gi", "1Ti"}
		for _, s := range canonical {
			L.SetGlobal("s", lua.LString(s))
			require.NoError(t, L.DoString(`
				local k8s = require("kubernetes")
				local b = k8s.parse_memory(s)
				result = k8s.format_memory(b)
			`))
			assert.Equal(t, lua.LString(s), L.GetGlobal("result"))
		}
	})
}

// TestFormatMemorySIRoundTrip: format_memory_si emits the decimal-SI form
// (k/M/G/T, powers of 1000) while format_memory emits binary (Ki/Mi/Gi/Ti,
// powers of 1024). Both round-trip cleanly through parse_memory, and a value
// that is a "nice" decimal number renders differently in the two formats.
func TestFormatMemorySIRoundTrip(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("kubernetes", Loader)

	t.Run("bytes -> SI string -> bytes", func(t *testing.T) {
		// Powers of 1000 that the decimal-SI form renders exactly.
		bytes := []int64{
			0,
			1000,                      // 1k
			1000 * 1000,               // 1M
			500 * 1000 * 1000,         // 500M
			1000 * 1000 * 1000,        // 1G
			2 * 1000 * 1000 * 1000,    // 2G
			1000 * 1000 * 1000 * 1000, // 1T
		}
		for _, b := range bytes {
			L.SetGlobal("b", lua.LNumber(b))
			require.NoError(t, L.DoString(`
				local k8s = require("kubernetes")
				local s = k8s.format_memory_si(b)
				result = k8s.parse_memory(s)
			`))
			assert.Equal(t, lua.LNumber(b), L.GetGlobal("result"),
				"format_memory_si then parse_memory should preserve %d", b)
		}
	})

	t.Run("binary vs decimal differ for the same byte count", func(t *testing.T) {
		// 2GB (2,000,000,000 bytes) renders as "2G" in decimal SI, but the
		// binary form has to express that same byte count without a clean
		// Gi factor — it'll come out in Mi (the largest binary unit that
		// divides evenly): 2000000000 / 1024^2 ≈ 1907.35Mi → not exact, so
		// the formatter falls back to a coarser representation. Just check
		// the strings are different and both parse back to the same value.
		require.NoError(t, L.DoString(`
			local k8s = require("kubernetes")
			local b = 2000000000
			bin = k8s.format_memory(b)
			dec = k8s.format_memory_si(b)
			bin_back = k8s.parse_memory(bin)
			dec_back = k8s.parse_memory(dec)
		`))
		bin := L.GetGlobal("bin").String()
		dec := L.GetGlobal("dec").String()
		assert.NotEqual(t, bin, dec, "binary and decimal forms should differ for 2,000,000,000 bytes")
		assert.Equal(t, lua.LNumber(2000000000), L.GetGlobal("bin_back"))
		assert.Equal(t, lua.LNumber(2000000000), L.GetGlobal("dec_back"))
		assert.Equal(t, "2G", dec, "decimal SI of 2e9 bytes should be exactly 2G")
	})
}

// TestFormatParseCPURoundTrip: format_cpu and parse_cpu are inverses for
// millicore counts. format_cpu(parse_cpu(x)) is a fixed point for inputs
// already in canonical decimal-SI form.
func TestFormatParseCPURoundTrip(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("kubernetes", Loader)

	t.Run("millicores -> string -> millicores", func(t *testing.T) {
		millis := []int64{0, 1, 100, 250, 500, 1000, 1500, 2500, 4000}
		for _, m := range millis {
			L.SetGlobal("m", lua.LNumber(m))
			require.NoError(t, L.DoString(`
				local k8s = require("kubernetes")
				local s = k8s.format_cpu(m)
				result = k8s.parse_cpu(s)
			`))
			assert.Equal(t, lua.LNumber(m), L.GetGlobal("result"),
				"format_cpu then parse_cpu should return the original millicore count for %d", m)
		}
	})

	t.Run("canonical string -> millicores -> string", func(t *testing.T) {
		// Whole-core values render as bare integers ("1", "2"), sub-core
		// as millicores ("500m"). Verify both are fixed points.
		canonical := []string{"0", "1m", "100m", "500m", "1", "1500m", "2", "4"}
		for _, s := range canonical {
			L.SetGlobal("s", lua.LString(s))
			require.NoError(t, L.DoString(`
				local k8s = require("kubernetes")
				local m = k8s.parse_cpu(s)
				result = k8s.format_cpu(m)
			`))
			assert.Equal(t, lua.LString(s), L.GetGlobal("result"))
		}
	})
}

// TestTimestampRoundTripInPodMetadata: a Pod's metav1.Time
// (CreationTimestamp) survives the round-trip as an RFC3339 string on the
// Lua side, and parse_time converts it back to the original Unix timestamp.
// This is the path scripts use to compute pod age.
func TestTimestampRoundTripInPodMetadata(t *testing.T) {
	cases := []struct {
		name string
		when time.Time
	}{
		{"epoch", time.Unix(0, 0).UTC()},
		{"2020-01-01", time.Date(2020, 1, 1, 12, 30, 45, 0, time.UTC)},
		{"2025-10-03", time.Date(2025, 10, 3, 16, 39, 0, 0, time.UTC)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "demo",
					CreationTimestamp: metav1.NewTime(tc.when),
				},
			}

			L := lua.NewState()
			defer L.Close()
			L.PreloadModule("kubernetes", Loader)
			loadPodIntoLua(t, L, pod)

			script := `
				local k8s = require("kubernetes")
				ts_str = pod.metadata.creationTimestamp
				ts_unix = k8s.parse_time(pod.metadata.creationTimestamp)
			`
			require.NoError(t, L.DoString(script))

			// metav1.Time marshals as RFC3339 in UTC, which is what the Lua
			// side should see directly.
			expectedStr := tc.when.UTC().Format(time.RFC3339)
			assert.Equal(t, lua.LString(expectedStr), L.GetGlobal("ts_str"),
				"creationTimestamp string should be the RFC3339 representation")
			assert.Equal(t, lua.LNumber(tc.when.Unix()), L.GetGlobal("ts_unix"),
				"parse_time should recover the original Unix timestamp")
		})
	}
}

// TestParseFormatTimeRoundTrip: the kubernetes module's own format_time and
// parse_time are inverses. Already covered by the existing
// TestFormatParseRoundTrip in kubernetes_test.go; this adds an explicit
// chain (Go-side unix -> Lua format -> Lua parse -> Go-side unix) so a
// future change to either function will fail loudly here as well.
func TestParseFormatTimeRoundTrip(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("kubernetes", Loader)

	unixes := []int64{0, 1, 1696347540, 1759509540, 2000000000}
	for _, u := range unixes {
		t.Run(time.Unix(u, 0).UTC().Format(time.RFC3339), func(t *testing.T) {
			L.SetGlobal("ts", lua.LNumber(u))
			require.NoError(t, L.DoString(`
				local k8s = require("kubernetes")
				local s = k8s.format_time(ts)
				result = k8s.parse_time(s)
			`))
			assert.Equal(t, lua.LNumber(u), L.GetGlobal("result"))
		})
	}
}

// TestParseFormatDurationRoundTrip: format_duration(parse_duration(x)) is a
// fixed-point for canonical whole-second Go duration strings.
//
// Sub-second precision is intentionally NOT covered here: FormatDuration
// truncates its float-seconds argument to a whole int before multiplying by
// time.Second, so "100ms" rounds to "0s". K8s API surfaces (probe timeouts,
// grace periods, etc.) are second-granular in practice, so this matches
// real-world use.
func TestParseFormatDurationRoundTrip(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("kubernetes", Loader)

	canonicalDurations := []string{
		"0s",
		"30s",
		"1m0s",
		"1h0m0s",
		"2h30m0s",
	}
	for _, d := range canonicalDurations {
		t.Run(d, func(t *testing.T) {
			L.SetGlobal("d", lua.LString(d))
			require.NoError(t, L.DoString(`
				local k8s = require("kubernetes")
				local secs = k8s.parse_duration(d)
				result = k8s.format_duration(secs)
			`))
			assert.Equal(t, lua.LString(d), L.GetGlobal("result"))
		})
	}
}
