$ go test -bench=. -cpu=1
```
goos: darwin
goarch: arm64
pkg: country-ip
cpu: Apple M3
BenchmarkIPLookupV1 	     106	  11587029 ns/op
--- BENCH: BenchmarkIPLookupV1
    country_ip_test.go:102: countryIP: 223 MB
    country_ip_test.go:112: avg time per lookup: 11.587029ms
BenchmarkIPLookupV2 	   40951	     29161 ns/op
--- BENCH: BenchmarkIPLookupV2
    country_ip_test.go:125: countryIP: 227 MB
    country_ip_test.go:135: avg time per lookup: 29.16µs
248
BenchmarkIPLookupV3 	   41557	     28740 ns/op
--- BENCH: BenchmarkIPLookupV3
    country_ip_test.go:148: countryIP: 99 MB
    country_ip_test.go:158: avg time per lookup: 28.74µs
BenchmarkIPLookupV4 	 4317189	       278.8 ns/op
--- BENCH: BenchmarkIPLookupV4
    country_ip_test.go:171: countryIP: 15 MB
    country_ip_test.go:181: avg time per lookup: 278ns
PASS
ok  	country-ip	10.502s
```
