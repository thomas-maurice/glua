# glua Benchmarks

Performance benchmarks for the glua library, measuring conversion performance and Lua field access patterns.

## Latest Benchmark Results

```
goos: darwin
goarch: arm64
pkg: github.com/thomas-maurice/glua/benchmarks
cpu: Apple M5
BenchmarkGoToLuaSimple-10              942836       1240 ns/op     4276 B/op       44 allocs/op
BenchmarkGoToLuaComplex-10             164640       6759 ns/op    23969 B/op      221 allocs/op
BenchmarkGoToLuaPod-10                  65240      17756 ns/op    58872 B/op      468 allocs/op
BenchmarkLuaToGoSimple-10             1350904        896.5 ns/op     1000 B/op       23 allocs/op
BenchmarkLuaToGoComplex-10             249796       4835 ns/op     4899 B/op      118 allocs/op
BenchmarkRoundTripSimple-10            508216       2390 ns/op     5279 B/op       67 allocs/op
BenchmarkRoundTripPod-10                70228      17374 ns/op    42888 B/op      391 allocs/op
BenchmarkLuaFieldAccess-10             185173       6622 ns/op    33904 B/op      112 allocs/op
BenchmarkLuaNestedFieldAccess-10        99068      11289 ns/op    37712 B/op      269 allocs/op
BenchmarkLuaArrayIteration-10           96232      12569 ns/op    36576 B/op      332 allocs/op
BenchmarkLuaMapIteration-10            154359       7850 ns/op    34776 B/op      122 allocs/op
BenchmarkLuaFieldModification-10       143973       7562 ns/op    34672 B/op      152 allocs/op
BenchmarkLuaComplexOperation-10         38978      52104 ns/op   334735 B/op      453 allocs/op
PASS
