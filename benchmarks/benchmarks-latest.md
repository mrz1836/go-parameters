running benchmarks...
goos: darwin
goarch: arm64
pkg: github.com/mrz1836/go-parameters
cpu: Apple M1 Max
BenchmarkUniqueUint64-10                12935953                83.69 ns/op           64 B/op          1 allocs/op
BenchmarkGetParams_ParseJSONBody-10     212284928                5.641 ns/op           0 B/op          0 allocs/op
BenchmarkGetParams-10                   213638424                5.615 ns/op           0 B/op          0 allocs/op
BenchmarkParams_GetStringOk-10          38213416                31.12 ns/op           16 B/op          1 allocs/op
BenchmarkParams_GetBoolOk-10            37702503                31.37 ns/op           16 B/op          1 allocs/op
BenchmarkParams_GetBytesOk-10           39913796                29.80 ns/op           16 B/op          1 allocs/op
BenchmarkParams_GetBool-10              38233860                31.38 ns/op           16 B/op          1 allocs/op
BenchmarkParams_GetFloatOk-10           22201544                53.77 ns/op           16 B/op          1 allocs/op
BenchmarkParams_GetIntOk-10             30041996                40.32 ns/op           16 B/op          1 allocs/op
BenchmarkParams_GetInt64Ok-10           29363125                40.25 ns/op           16 B/op          1 allocs/op
BenchmarkParams_GetIntSliceOk-10        39847030                30.12 ns/op           16 B/op          1 allocs/op
BenchmarkParams_GetUint64Ok-10          31657987                38.24 ns/op           16 B/op          1 allocs/op
