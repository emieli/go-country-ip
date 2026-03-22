
# Benchmarks
```
go test -bench=. -cpu=1
goos: darwin
goarch: arm64
pkg: country-ip
cpu: Apple M3
--- BENCH: BenchmarkIPLookupV1
    country_ip_test.go:103: countryIP: 223 MB
    country_ip_test.go:113: avg time per lookup: 9.753066ms
--- BENCH: BenchmarkIPLookupV2
    country_ip_test.go:126: countryIP: 227 MB
    country_ip_test.go:136: avg time per lookup: 46.29µs
--- BENCH: BenchmarkIPLookupV3
    country_ip_test.go:149: countryIP: 100 MB
    country_ip_test.go:159: avg time per lookup: 47.369µs
--- BENCH: BenchmarkIPLookupV4
    country_ip_test.go:172: countryIP: 15 MB
    country_ip_test.go:182: avg time per lookup: 289ns
PASS
ok  	country-ip	10.309s
```

