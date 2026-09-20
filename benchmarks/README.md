# glua Benchmarks

Performance benchmarks for the glua library, measuring conversion performance and Lua field access patterns.

## Latest Benchmark Results

```
goos: darwin
goarch: arm64
pkg: github.com/thomas-maurice/glua/benchmarks
cpu: Apple M5
BenchmarkGoToLuaSimple-10              175927       8157 ns/op     4215 B/op       43 allocs/op
BenchmarkGoToLuaComplex-10              32445      36076 ns/op    23988 B/op      233 allocs/op
BenchmarkGoToLuaPod-10                  15172      77784 ns/op    59088 B/op      478 allocs/op
BenchmarkLuaToGoSimple-10              322880       3516 ns/op      667 B/op       19 allocs/op
BenchmarkLuaToGoComplex-10              82521      16281 ns/op     3642 B/op       82 allocs/op
BenchmarkRoundTripSimple-10            100527      11636 ns/op     4890 B/op       62 allocs/op
BenchmarkRoundTripPod-10                16885      71506 ns/op    41262 B/op      339 allocs/op
BenchmarkLuaFieldAccess-10              25604      45352 ns/op    33904 B/op      112 allocs/op
BenchmarkLuaNestedFieldAccess-10        21942      53456 ns/op    37712 B/op      269 allocs/op
BenchmarkLuaArrayIteration-10           21526      52783 ns/op    36576 B/op      332 allocs/op
BenchmarkLuaMapIteration-10             26996      47664 ns/op    34776 B/op      122 allocs/op
BenchmarkLuaFieldModification-10        25705      43479 ns/op    34672 B/op      152 allocs/op
BenchmarkLuaComplexOperation-10         10000     184072 ns/op   129831 B/op      453 allocs/op
PASS
```

### Key Takeaways

- **The JSON round-trip dominates.** Both directions marshal to JSON and back
  (see the Translator section in the root README), so cost scales with the size
  of the object graph, not with how many fields Lua actually touches.
  `GoToLuaPod` moves a full Pod spec and costs ~13x `GoToLuaSimple`.
- **Go -> Lua is the expensive direction**, roughly 1.7x `LuaToGo` on the simple
  case and 2.1x on the complex one, allocating ~6.6x the memory across ~2.8x the
  allocations on that complex case. Converting a large Go struct into
  Lua is the thing to avoid in a hot path — hoist it out of the loop, or expose
  the value as a class (userdata) instead, which skips serialisation entirely.
- **Lua-side field access is not free.** `LuaFieldAccess` and friends allocate
  ~34 KB per iteration because each benchmark rebuilds its table; the per-access
  cost itself is small. Read repeatedly from one converted table rather than
  reconverting.
- **Nested access costs about 1.7x flat access**, and iteration (array or map)
  lands in the same range. None of these are pathological; the conversion at the
  boundary is what to watch.

Comparability note: these numbers were produced under Go 1.26+ with the current
dependency set. Comparing them against results recorded before a toolchain or
gopher-lua bump is not meaningful — re-run `make bench-update` on an idle
machine and compare like for like.
