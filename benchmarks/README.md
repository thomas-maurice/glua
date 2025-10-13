# glua Benchmarks

Performance benchmarks for the glua library, measuring conversion performance and Lua field access patterns.

## Latest Benchmark Results

```
goos: linux
goarch: amd64
pkg: github.com/thomas-maurice/glua/benchmarks
cpu: AMD Ryzen 7 7840U w/ Radeon  780M Graphics
BenchmarkGoToLuaSimple-16              413799       3055 ns/op     4327 B/op       45 allocs/op
BenchmarkGoToLuaComplex-16              77994      15150 ns/op    23980 B/op      221 allocs/op
BenchmarkGoToLuaPod-16                  30216      39742 ns/op    58907 B/op      468 allocs/op
BenchmarkLuaToGoSimple-16              609015       1871 ns/op     1000 B/op       23 allocs/op
BenchmarkLuaToGoComplex-16             124642       9554 ns/op     4901 B/op      118 allocs/op
BenchmarkRoundTripSimple-16            223320       5406 ns/op     5330 B/op       68 allocs/op
BenchmarkRoundTripPod-16                31477      37718 ns/op    42924 B/op      391 allocs/op
BenchmarkLuaFieldAccess-16              92329      13088 ns/op    33937 B/op      114 allocs/op
BenchmarkLuaNestedFieldAccess-16        52149      23868 ns/op    37745 B/op      271 allocs/op
BenchmarkLuaArrayIteration-16           43675      27139 ns/op    36633 B/op      335 allocs/op
BenchmarkLuaMapIteration-16             78471      16252 ns/op    34833 B/op      125 allocs/op
BenchmarkLuaFieldModification-16        77530      15810 ns/op    34705 B/op      154 allocs/op
BenchmarkLuaComplexOperation-16         19971      74625 ns/op   201027 B/op      458 allocs/op
PASS
