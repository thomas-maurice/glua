# glua Benchmarks

Performance benchmarks for the glua library, measuring conversion performance and Lua field access patterns.

## Latest Benchmark Results

```
goos: linux
goarch: amd64
pkg: github.com/thomas-maurice/glua/benchmarks
cpu: AMD Ryzen 7 7840U w/ Radeon  780M Graphics
BenchmarkGoToLuaSimple-16              372009       3119 ns/op     4327 B/op       45 allocs/op
BenchmarkGoToLuaComplex-16              77182      15232 ns/op    23979 B/op      221 allocs/op
BenchmarkGoToLuaPod-16                  30975      39203 ns/op    58908 B/op      468 allocs/op
BenchmarkLuaToGoSimple-16              601479       1904 ns/op     1000 B/op       23 allocs/op
BenchmarkLuaToGoComplex-16             126118       9552 ns/op     4901 B/op      118 allocs/op
BenchmarkRoundTripSimple-16            221824       5349 ns/op     5330 B/op       68 allocs/op
BenchmarkRoundTripPod-16                31137      37838 ns/op    42924 B/op      391 allocs/op
BenchmarkLuaFieldAccess-16              92990      13136 ns/op    33937 B/op      114 allocs/op
BenchmarkLuaNestedFieldAccess-16        49603      23374 ns/op    37745 B/op      271 allocs/op
BenchmarkLuaArrayIteration-16           45322      26555 ns/op    36633 B/op      335 allocs/op
BenchmarkLuaMapIteration-16             73569      16372 ns/op    34833 B/op      125 allocs/op
BenchmarkLuaFieldModification-16        80230      15426 ns/op    34705 B/op      154 allocs/op
BenchmarkLuaComplexOperation-16         19976      74346 ns/op   201064 B/op      458 allocs/op
PASS
