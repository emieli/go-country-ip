# Benchmarks
```
$ go test -bench=. -cpu=1
goos: darwin
goarch: arm64
pkg: country-ip
cpu: Apple M3
BenchmarkIPLookupV1 	     126	  10199400 ns/op
--- BENCH: BenchmarkIPLookupV1
    country_ip_test.go:120: countryIP: 223 MB
    country_ip_test.go:130: avg time per lookup: 10.1994ms
BenchmarkIPLookupV2 	   25707	     46609 ns/op
--- BENCH: BenchmarkIPLookupV2
    country_ip_test.go:143: countryIP: 227 MB
    country_ip_test.go:153: avg time per lookup: 46.608µs
BenchmarkIPLookupV3 	   26227	     46398 ns/op
--- BENCH: BenchmarkIPLookupV3
    country_ip_test.go:166: countryIP: 100 MB
    country_ip_test.go:176: avg time per lookup: 46.397µs
BenchmarkIPLookupV4 	 4127764	       288.1 ns/op
--- BENCH: BenchmarkIPLookupV4
    country_ip_test.go:189: countryIP: 15 MB
    country_ip_test.go:190: entries: 1341613
    country_ip_test.go:200: avg time per lookup: 288ns
BenchmarkIPLookupV5 	 4050750	       294.8 ns/op
--- BENCH: BenchmarkIPLookupV5
    country_ip_test.go:213: countryIP: 7 MB
    country_ip_test.go:214: entries: 607604
    country_ip_test.go:224: avg time per lookup: 294ns
PASS
ok  	country-ip	12.723s
```

