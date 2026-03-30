package main

import (
	v1 "country-ip/v1"
	v2 "country-ip/v2"
	v3 "country-ip/v3"
	v4 "country-ip/v4"
	v5 "country-ip/v5"
	v6 "country-ip/v6"
	"fmt"
	"runtime"
	"testing"
	"time"
)

var ipTests = []struct {
	IP      string
	country string
}{
	{"0.0.0.0", ""},
	{"0.255.255.255", ""},
	{"1.0.0.0", "Australia"},
	{"1.0.0.1", "Australia"},
	{"1.0.0.2", "Australia"},
	{"1.0.0.255", "Australia"},
	{"1.0.1.0", "China"},
	{"1.7.168.174", "Singapore"},
	{"1.255.240.200", "South Korea"},
	{"1.200.0.200", "Taiwan"},
	{"162.0.0.200", "Canada"},
	{"10.10.10.10", ""},
	{"100.100.100.100", ""},
	{"192.168.123.123", ""},
	{"172.25.4.5", ""},
	{"223.255.255.253", "Australia"},
	{"223.255.255.254", "Australia"},
	{"223.255.255.255", "Australia"},
	{"29.1.2.3", "United States"}, // Department of defense
	{"83.250.95.17", "Sweden"},
}

func TestCountryIPDataIPLookupV1(t *testing.T) {
	countryIP, err := v1.NewCountryIPData()
	if err != nil {
		t.Fatalf("new CountryIPData: %v", err)
	}
	for _, tt := range ipTests {
		t.Run(tt.IP, func(t *testing.T) {
			if got, want := countryIP.AddrCountry(tt.IP), tt.country; got != want {
				t.Errorf("%s got %s, want %s", tt.IP, got, want)
			}
		})
	}
}

func TestCountryIPDataIPLookupV2(t *testing.T) {
	countryIP, err := v2.NewCountryIPData()
	if err != nil {
		t.Fatalf("new CountryIPData: %v", err)
	}

	for _, tt := range ipTests {
		t.Run(tt.IP, func(t *testing.T) {
			if got, want := countryIP.AddrCountry(tt.IP), tt.country; got != want {
				t.Errorf("%s got %s, want %s", tt.IP, got, want)
			}
		})
	}
}

func TestCountryIPDataIPLookupV3(t *testing.T) {
	countryIP, err := v3.NewCountryIPData()
	if err != nil {
		t.Fatalf("new CountryIPData: %v", err)
	}
	for _, tt := range ipTests {
		t.Run(tt.IP, func(t *testing.T) {
			if got, want := countryIP.AddrCountry(tt.IP), tt.country; got != want {
				t.Errorf("%s got %s, want %s", tt.IP, got, want)
			}
		})
	}
}

func TestCountryIPDataIPLookupV4(t *testing.T) {
	countryIP, err := v4.NewCountryIPData()
	if err != nil {
		t.Fatalf("new CountryIPData: %v", err)
	}
	for _, tt := range ipTests {
		t.Run(tt.IP, func(t *testing.T) {
			if got, want := countryIP.AddrCountry(tt.IP), tt.country; got != want {
				t.Errorf("%s got %s, want %s", tt.IP, got, want)
			}
		})
	}
}

func TestCountryIPDataIPLookupV5(t *testing.T) {
	countryIP, err := v5.NewCountryIPData()
	if err != nil {
		t.Fatalf("new CountryIPData: %v", err)
	}
	for _, tt := range ipTests {
		t.Run(tt.IP, func(t *testing.T) {
			if got, want := countryIP.AddrCountry(tt.IP), tt.country; got != want {
				t.Errorf("%s got %s, want %s", tt.IP, got, want)
			}
		})
	}
}

func TestCountryIPDataIPLookupV6(t *testing.T) {
	countryIP, err := v6.NewCountryIPData()
	if err != nil {
		t.Fatalf("new CountryIPData: %v", err)
	}
	for _, tt := range ipTests {
		t.Run(tt.IP, func(t *testing.T) {
			if got, want := countryIP.AddrCountry(tt.IP), tt.country; got != want {
				t.Errorf("%s got %s, want %s", tt.IP, got, want)
			}
		})
	}
}

func BenchmarkIPLookupV1(b *testing.B) {
	countryIP, err := v1.NewCountryIPData()
	if err != nil {
		fmt.Printf("new CountryIPData: %v\n", err)
		return
	}

	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	alloc := m.Alloc
	b.Logf(
		"countryIP: %d MB (%d bytes per entry)\n",
		alloc/1024/1024,
		alloc/uint64(countryIP.Length()),
	)

	octet := 0
	for b.Loop() {
		countryIP.AddrCountry(fmt.Sprintf("%d.%d.%d.%d", octet, octet, octet, octet))
		octet++
		if octet > 239 {
			octet = 0
		}
	}
	b.Log("avg time per lookup:", b.Elapsed()/time.Duration(b.N))
}

func BenchmarkIPLookupV2(b *testing.B) {
	countryIP, err := v2.NewCountryIPData()
	if err != nil {
		fmt.Printf("new CountryIPData: %v\n", err)
		return
	}

	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	alloc := m.Alloc
	b.Logf(
		"countryIP: %d MB (%d bytes per entry)\n",
		alloc/1024/1024,
		alloc/uint64(countryIP.Length()),
	)

	octet := 0
	for b.Loop() {
		countryIP.AddrCountry(fmt.Sprintf("%d.%d.%d.%d", octet, octet, octet, octet))
		octet++
		if octet > 239 {
			octet = 0
		}
	}
	b.Log("avg time per lookup:", b.Elapsed()/time.Duration(b.N))
}

func BenchmarkIPLookupV3(b *testing.B) {
	countryIP, err := v3.NewCountryIPData()
	if err != nil {
		fmt.Printf("new CountryIPData: %v\n", err)
		return
	}

	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	alloc := m.Alloc
	b.Logf(
		"countryIP: %d MB (%d bytes per entry)\n",
		alloc/1024/1024,
		alloc/uint64(countryIP.Length()),
	)

	octet := 0
	for b.Loop() {
		countryIP.AddrCountry(fmt.Sprintf("%d.%d.%d.%d", octet, octet, octet, octet))
		octet++
		if octet > 239 {
			octet = 0
		}
	}
	b.Log("avg time per lookup:", b.Elapsed()/time.Duration(b.N))
}

func BenchmarkIPLookupV4(b *testing.B) {
	countryIP, err := v4.NewCountryIPData()
	if err != nil {
		fmt.Printf("new CountryIPData: %v\n", err)
		return
	}

	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	alloc := m.Alloc
	b.Logf(
		"countryIP: %d MB (%d bytes per entry)\n",
		alloc/1024/1024,
		alloc/uint64(countryIP.Length()),
	)

	octet := 0
	for b.Loop() {
		countryIP.AddrCountry(fmt.Sprintf("%d.%d.%d.%d", octet, octet, octet, octet))
		octet++
		if octet > 239 {
			octet = 0
		}
	}
	b.Log("avg time per lookup:", b.Elapsed()/time.Duration(b.N))
}

func BenchmarkIPLookupV5(b *testing.B) {
	countryIP, err := v5.NewCountryIPData()
	if err != nil {
		fmt.Printf("new CountryIPData: %v\n", err)
		return
	}

	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	alloc := m.Alloc
	b.Logf(
		"countryIP: %d MB (%d bytes per entry)\n",
		alloc/1024/1024,
		alloc/uint64(countryIP.Length()),
	)

	octet := 0
	for b.Loop() {
		countryIP.AddrCountry(fmt.Sprintf("%d.%d.%d.%d", octet, octet, octet, octet))
		octet++
		if octet > 239 {
			octet = 0
		}
	}
	b.Log("avg time per lookup:", b.Elapsed()/time.Duration(b.N))
}

func BenchmarkIPLookupV6(b *testing.B) {
	countryIP, err := v6.NewCountryIPData()
	if err != nil {
		fmt.Printf("new CountryIPData: %v\n", err)
		return
	}

	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	alloc := m.Alloc
	b.Logf(
		"countryIP: %d MB (%d bytes per entry)\n",
		alloc/1024/1024,
		alloc/uint64(countryIP.Length()),
	)

	octet := 0
	for b.Loop() {
		countryIP.AddrCountry(fmt.Sprintf("%d.%d.%d.%d", octet, octet, octet, octet))
		octet++
		if octet > 239 {
			octet = 0
		}
	}
	b.Log("avg time per lookup:", b.Elapsed()/time.Duration(b.N))
}
