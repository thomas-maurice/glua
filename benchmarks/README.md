# glua Benchmarks

Performance benchmarks for the glua library, measuring conversion performance and Lua field access patterns.

## Latest Benchmark Results

```
goos: linux
goarch: amd64
pkg: github.com/thomas-maurice/glua/benchmarks
cpu: AMD Ryzen 7 7840U w/ Radeon  780M Graphics
BenchmarkGoToLuaSimple-16              360781       3234 ns/op     4327 B/op       45 allocs/op
BenchmarkGoToLuaComplex-16              73927      16044 ns/op    23979 B/op      221 allocs/op
BenchmarkGoToLuaPod-16                  30454      40704 ns/op    58903 B/op      468 allocs/op
BenchmarkLuaToGoSimple-16              593233       1961 ns/op     1000 B/op       23 allocs/op
BenchmarkLuaToGoComplex-16             124947      10073 ns/op     4901 B/op      118 allocs/op
BenchmarkRoundTripSimple-16            205348       5714 ns/op     5330 B/op       68 allocs/op
BenchmarkRoundTripPod-16                30904      40151 ns/op    42921 B/op      391 allocs/op
BenchmarkLuaFieldAccess-16              91231      13363 ns/op    33937 B/op      114 allocs/op
BenchmarkLuaNestedFieldAccess-16        48250      24499 ns/op    37745 B/op      271 allocs/op
BenchmarkLuaArrayIteration-16           42354      27735 ns/op    36633 B/op      335 allocs/op
BenchmarkLuaMapIteration-16             74606      16353 ns/op    34833 B/op      125 allocs/op
BenchmarkLuaFieldModification-16        77610      15695 ns/op    34705 B/op      154 allocs/op
BenchmarkLuaComplexOperation-16         20061      74446 ns/op   201693 B/op      458 allocs/op
PASS
