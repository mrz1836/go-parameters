running benchmarks...
goos: darwin
goarch: arm64
pkg: github.com/mrz1836/go-parameters
cpu: Apple M1 Max
BenchmarkUniqueUint64-10                13989841                84.49 ns/op           64 B/op          1 allocs/op
BenchmarkGetParams_ParseJSONBody-10     209700817                5.721 ns/op           0 B/op          0 allocs/op
BenchmarkGetParams-10                   209668894                5.719 ns/op           0 B/op          0 allocs/op
BenchmarkParams_GetStringOk-10          37573434                31.96 ns/op           16 B/op          1 allocs/op
BenchmarkParams_GetBoolOk-10            36349316                33.15 ns/op           16 B/op          1 allocs/op
BenchmarkParams_GetBytesOk-10           38539616                31.52 ns/op           16 B/op          1 allocs/op
BenchmarkParams_GetBool-10              35637433                32.57 ns/op           16 B/op          1 allocs/op
BenchmarkParams_GetFloatOk-10           22064586                54.69 ns/op           16 B/op          1 allocs/op
BenchmarkParams_GetIntOk-10             28992304                41.78 ns/op           16 B/op          1 allocs/op
BenchmarkParams_GetInt64Ok-10           28752844                41.77 ns/op           16 B/op          1 allocs/op
BenchmarkParams_GetIntSliceOk-10        38099671                31.59 ns/op           16 B/op          1 allocs/op
BenchmarkParams_GetUint64Ok-10          29921580                39.68 ns/op           16 B/op          1 allocs/op