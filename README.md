# Benchmarks
```
$ go test -bench=. -cpu=1
goos: darwin
goarch: arm64
pkg: country-ip
cpu: Apple M3
BenchmarkIPLookupV1 	     112	  12940787 ns/op
--- BENCH: BenchmarkIPLookupV1
    country_ip_test.go:138: countryIP: 223 MB (174 bytes per entry)
    country_ip_test.go:143: entries: 1341613
    country_ip_test.go:153: avg time per lookup: 12.940787ms
BenchmarkIPLookupV2 	   25771	     46337 ns/op
--- BENCH: BenchmarkIPLookupV2
    country_ip_test.go:167: countryIP: 227 MB (178 bytes per entry)
    country_ip_test.go:172: entries: 1341613
    country_ip_test.go:182: avg time per lookup: 46.336µs
BenchmarkIPLookupV3 	   26342	     46414 ns/op
--- BENCH: BenchmarkIPLookupV3
    country_ip_test.go:196: countryIP: 99 MB (78 bytes per entry)
    country_ip_test.go:201: entries: 1341613
    country_ip_test.go:211: avg time per lookup: 46.413µs
BenchmarkIPLookupV4 	 4098352	       292.0 ns/op
--- BENCH: BenchmarkIPLookupV4
    country_ip_test.go:225: countryIP: 15 MB (12 bytes per entry)
    country_ip_test.go:230: entries: 1341613
    country_ip_test.go:240: avg time per lookup: 291ns
BenchmarkIPLookupV5 	 4127480	       293.6 ns/op
--- BENCH: BenchmarkIPLookupV5
    country_ip_test.go:254: countryIP: 7 MB (12 bytes per entry)
    country_ip_test.go:259: entries: 607604
    country_ip_test.go:269: avg time per lookup: 293ns
BenchmarkIPLookupV6 	 4209511	       288.2 ns/op
--- BENCH: BenchmarkIPLookupV6
    country_ip_test.go:283: countryIP: 4 MB (8 bytes per entry)
    country_ip_test.go:288: entries: 607617
    country_ip_test.go:298: avg time per lookup: 288ns
BenchmarkIPLookupV7 	   95426	     12395 ns/op
--- BENCH: BenchmarkIPLookupV7
    country_ip_test.go:327: countryIP: 35 MB (61 bytes per entry)
    country_ip_test.go:332: entries: 607617
    country_ip_test.go:342: avg time per lookup: 12.395µs
PASS
ok  	country-ip	12.723s
```

