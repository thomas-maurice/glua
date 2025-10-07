# glua Benchmarks

Performance benchmarks for the glua library, measuring conversion performance and Lua field access patterns.

## Latest Benchmark Results

```
goos: linux
goarch: amd64
pkg: github.com/thomas-maurice/glua/benchmarks
cpu: AMD Ryzen 7 7840U w/ Radeon  780M Graphics

BenchmarkGoToLuaSimple-16              421356       3020 ns/op     4327 B/op       45 allocs/op
BenchmarkGoToLuaComplex-16              76162      15331 ns/op    23979 B/op      221 allocs/op
BenchmarkGoToLuaPod-16                  31567      39143 ns/op    58901 B/op      468 allocs/op
BenchmarkLuaToGoSimple-16              618954       1895 ns/op     1000 B/op       23 allocs/op
BenchmarkLuaToGoComplex-16             124588       9478 ns/op     4901 B/op      118 allocs/op
BenchmarkRoundTripSimple-16            202479       5375 ns/op     5330 B/op       68 allocs/op
BenchmarkRoundTripPod-16                31788      38159 ns/op    42926 B/op      391 allocs/op
BenchmarkLuaFieldAccess-16              88074      13417 ns/op    33937 B/op      114 allocs/op
BenchmarkLuaNestedFieldAccess-16        49584      23874 ns/op    37745 B/op      271 allocs/op
BenchmarkLuaArrayIteration-16           40807      29289 ns/op    36633 B/op      335 allocs/op
BenchmarkLuaMapIteration-16             74904      16628 ns/op    34833 B/op      125 allocs/op
BenchmarkLuaFieldModification-16        76243      15934 ns/op    34705 B/op      154 allocs/op
BenchmarkLuaComplexOperation-16         19380      75691 ns/op   196936 B/op      458 allocs/op
```

### Key Takeaways

- **Simple conversions**: ~3µs Go→Lua, ~1.9µs Lua→Go
- **Kubernetes Pod**: ~39µs for full conversion (both directions)
- **Round-trip overhead**: Minimal (~5µs for simple structs)
- **Lua field access**: Fast for simple operations (~13µs)
- **Complex operations**: ~76µs for multi-step Pod manipulation

## Running Benchmarks

Run all benchmarks:

```bash
go test -bench=. -benchmem ./benchmarks/
```

Run specific benchmark:

```bash
go test -bench=BenchmarkGoToLuaPod -benchmem ./benchmarks/
```

Run with CPU profiling:

```bash
go test -bench=. -benchmem -cpuprofile=cpu.prof ./benchmarks/
go tool pprof cpu.prof
```

## Benchmark Categories

### Conversion Benchmarks (`conversion_bench_test.go`)

**Go → Lua Conversion:**

- `BenchmarkGoToLuaSimple` - Simple struct (3 fields)
- `BenchmarkGoToLuaComplex` - Complex nested struct with maps, slices
- `BenchmarkGoToLuaPod` - Full Kubernetes Pod object

**Lua → Go Conversion:**

- `BenchmarkLuaToGoSimple` - Simple struct from Lua table
- `BenchmarkLuaToGoComplex` - Complex nested struct from Lua

**Round-trip Conversion:**

- `BenchmarkRoundTripSimple` - Simple struct Go→Lua→Go
- `BenchmarkRoundTripPod` - Kubernetes Pod Go→Lua→Go

### Lua Access Benchmarks (`lua_access_bench_test.go`)

**Field Access:**

- `BenchmarkLuaFieldAccess` - Direct field access
- `BenchmarkLuaNestedFieldAccess` - Deep nested field access (Pod metadata)

**Iteration:**

- `BenchmarkLuaArrayIteration` - Iterate over 100-item array
- `BenchmarkLuaMapIteration` - Iterate over 10-item map

**Modification:**

- `BenchmarkLuaFieldModification` - Modify fields in Lua
- `BenchmarkLuaComplexOperation` - Complex operations (count, filter, modify)

## Interpreting Results

Example output:

```
BenchmarkGoToLuaSimple-8        50000    25000 ns/op    15000 B/op    200 allocs/op
```

- `50000` - Number of iterations
- `25000 ns/op` - Nanoseconds per operation (lower is better)
- `15000 B/op` - Bytes allocated per operation (lower is better)
- `200 allocs/op` - Number of allocations per operation (lower is better)

## Performance Tips

1. **Reuse Lua States**: Creating new `lua.NewState()` is expensive
2. **Batch Conversions**: Convert once, use many times in Lua
3. **Avoid Repeated Conversions**: Cache converted objects when possible
4. **Profile First**: Use `-cpuprofile` to find bottlenecks before optimizing

## Updating Benchmarks

When adding new features to glua:

1. Add corresponding benchmarks for the new functionality
2. Run benchmarks before and after changes to measure impact
3. Update this README if adding new benchmark categories
4. Include benchmark results in PRs for performance-critical changes

## CI Integration

Benchmarks are run as part of the test suite but don't fail the build. They serve as:

- Performance regression detection
- Optimization verification
- Documentation of expected performance characteristics
