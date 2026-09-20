# glua Benchmarks

Performance benchmarks for the glua library, measuring conversion performance and Lua field access patterns.

## Latest Benchmark Results

```
goos: darwin
goarch: arm64
pkg: github.com/thomas-maurice/glua/benchmarks
cpu: Apple M5
BenchmarkGoToLuaSimple-10              862110       1400 ns/op     4215 B/op       43 allocs/op
BenchmarkGoToLuaComplex-10             150016       8010 ns/op    23979 B/op      233 allocs/op
BenchmarkGoToLuaPod-10                  65001      20126 ns/op    59066 B/op      478 allocs/op
BenchmarkLuaToGoSimple-10             1448500        837.2 ns/op      666 B/op       19 allocs/op
BenchmarkLuaToGoComplex-10             302770       3978 ns/op     3641 B/op       82 allocs/op
BenchmarkRoundTripSimple-10            499888       2423 ns/op     4889 B/op       62 allocs/op
BenchmarkRoundTripPod-10                74859      16444 ns/op    41245 B/op      339 allocs/op
BenchmarkLuaFieldAccess-10             206894       5874 ns/op    33904 B/op      112 allocs/op
BenchmarkLuaNestedFieldAccess-10       120674      10016 ns/op    37712 B/op      269 allocs/op
BenchmarkLuaArrayIteration-10          105726      11608 ns/op    36576 B/op      332 allocs/op
BenchmarkLuaMapIteration-10            172785       7112 ns/op    34776 B/op      122 allocs/op
BenchmarkLuaFieldModification-10       184108       6643 ns/op    34672 B/op      152 allocs/op
BenchmarkLuaComplexOperation-10         44010      53299 ns/op   370047 B/op      453 allocs/op
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
