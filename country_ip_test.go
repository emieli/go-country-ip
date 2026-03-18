package main

import (
	v1 "country-ip/v1"
	v2 "country-ip/v2"
	v3 "country-ip/v3"
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
	{"1.0.0.0", "Australia"},
	{"1.0.0.1", "Australia"},
	{"1.0.0.2", "Australia"},
	{"1.0.0.255", "Australia"},
	{"1.0.1.0", "China"},
	{"1.255.240.200", "South Korea"},
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
	b.Log("countryIP:", m.Alloc/1024/1024, "MB")

	octet := 0
	for b.Loop() {
		countryIP.AddrCountry(fmt.Sprintf("%d.%d.%d.%d", octet, octet, octet, octet))
		octet += 10
		if octet > 230 {
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
	b.Log("countryIP:", m.Alloc/1024/1024, "MB")

	octet := 0
	for b.Loop() {
		countryIP.AddrCountry(fmt.Sprintf("%d.%d.%d.%d", octet, octet, octet, octet))
		octet += 10
		if octet > 230 {
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
	b.Log("countryIP:", m.Alloc/1024/1024, "MB")

	octet := 0
	for b.Loop() {
		countryIP.AddrCountry(fmt.Sprintf("%d.%d.%d.%d", octet, octet, octet, octet))
		octet += 10
		if octet > 230 {
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
	b.Log("countryIP:", m.Alloc/1024/1024, "MB")

	octet := 0
	for b.Loop() {
		countryIP.AddrCountry(fmt.Sprintf("%d.%d.%d.%d", octet, octet, octet, octet))
		octet += 10
		if octet > 230 {
			octet = 0
		}
	}
	b.Log("avg time per lookup:", b.Elapsed()/time.Duration(b.N))
}
