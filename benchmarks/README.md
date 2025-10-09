# glua Benchmarks

Performance benchmarks for the glua library, measuring conversion performance and Lua field access patterns.

## Latest Benchmark Results

```
goos: linux
goarch: amd64
pkg: github.com/thomas-maurice/glua/benchmarks
cpu: AMD Ryzen 7 7840U w/ Radeon  780M Graphics
BenchmarkGoToLuaSimple-16              409915       3074 ns/op     4327 B/op       45 allocs/op
BenchmarkGoToLuaComplex-16              75747      15458 ns/op    23979 B/op      221 allocs/op
BenchmarkGoToLuaPod-16                  29530      40777 ns/op    58909 B/op      468 allocs/op
BenchmarkLuaToGoSimple-16              609123       2044 ns/op     1000 B/op       23 allocs/op
BenchmarkLuaToGoComplex-16             124522       9876 ns/op     4901 B/op      118 allocs/op
BenchmarkRoundTripSimple-16            202911       5640 ns/op     5330 B/op       68 allocs/op
BenchmarkRoundTripPod-16                30727      39612 ns/op    42924 B/op      391 allocs/op
BenchmarkLuaFieldAccess-16              88826      13688 ns/op    33937 B/op      114 allocs/op
BenchmarkLuaNestedFieldAccess-16        49378      24432 ns/op    37745 B/op      271 allocs/op
BenchmarkLuaArrayIteration-16           42882      28862 ns/op    36633 B/op      335 allocs/op
BenchmarkLuaMapIteration-16             72069      17139 ns/op    34833 B/op      125 allocs/op
BenchmarkLuaFieldModification-16        74764      16009 ns/op    34705 B/op      154 allocs/op
BenchmarkLuaComplexOperation-16         19200      76427 ns/op   195660 B/op      458 allocs/op
PASS
