// package v4 improves on v3 by using uint32 netaddr/lastAddr which reduces
// memory usage per entry from ~78 bytes to 12 bytes. This is a massive boost
// that lower memory usage from ~100MB to 15MB.
// We also use binary search to quickly search through the []IPv4SubnetCountry,
// performing lookups in less than 300 nanoseconds.
package v4

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net/netip"
	"os"
	"slices"
	"strings"
)

type CountryIPData struct {
	subnetCountries []IPv4SubnetCountry
	countriesByCC   map[[2]byte]string // An extra lookup but much more memory efficient
}

type IPv4SubnetCountry struct {
	netAddr     uint32
	lastAddr    uint32
	countryCode [2]byte
}

func NewCountryIPData() (*CountryIPData, error) {
	c := new(CountryIPData)
	c.countriesByCC = make(map[[2]byte]string)

	if err := c.parseIPInfoCSV(); err != nil {
		return nil, fmt.Errorf("parse ipinfo: %w", err)
	}
	return c, nil
}

func (c CountryIPData) Length() int {
	return len(c.subnetCountries)
}

func (c *CountryIPData) parseIPInfoCSV() error {

	ipInfoCSV, err := os.Open("ipinfo_lite.csv")
	if err != nil {
		return fmt.Errorf("open ipinfo file: %v", err)
	}

	var (
		subnetCountries = make([]IPv4SubnetCountry, 0, 1_500_000)
		scanner         = bufio.NewScanner(ipInfoCSV)
		countryCode     [2]byte
	)

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ",")

		prefix := fields[0]
		if !strings.Contains(prefix, "/") {
			prefix = prefix + "/32"
		}

		subnet, err := netip.ParsePrefix(prefix)
		if err != nil {
			if fields[0] == "network" {
				continue // first line contains CSV headers, first header is "network"
			}
			return fmt.Errorf("line %s: %v", scanner.Text(), err)
		} else if subnet.Addr().Is6() {
			break // We only care about ipv4 subnets
		}

		netAddrBytes := subnet.Addr().As4()
		netAddr := binary.BigEndian.Uint32(netAddrBytes[:])
		lastAddr := netAddr + 1<<(32-subnet.Bits()) - 1 // 1<<32-24 = 256

		copy(countryCode[:], fields[2])
		if _, exist := c.countriesByCC[countryCode]; !exist {
			c.countriesByCC[countryCode] = fields[1] // country name
		}

		subnetCountries = append(subnetCountries, IPv4SubnetCountry{
			netAddr:     netAddr,
			lastAddr:    lastAddr,
			countryCode: countryCode,
		})
	}

	// Create a new slice containing the exact number of entries
	subnets := make([]IPv4SubnetCountry, len(subnetCountries))
	copy(subnets, subnetCountries)

	// Sort entries in ascending order by the netAddr field
	slices.SortFunc(subnets, func(a, b IPv4SubnetCountry) int {
		return int(a.netAddr - b.netAddr)
	})
	c.subnetCountries = subnets
	return nil
}

// Search iteratively but make the search-window smaller as we go.
// Let's say we have a list of numbers from 0-100 and we want to find the value 60.
// Low starts at 0, high starts at 100. (100+0) / 2 = 50. That's our needle.
//
// The process:
// 1. 100+50 / 2 = 75. 75 is higher than 60. The new low is 50, the new high is 75.
// 2.  75+50 / 2 = 63. The low remains at 50, the new high is 63.
// 3.  63+50 / 2 = 57. ew low is 57, the new high is 63.
// 4.  63+57 / 2 = 60. The numbers match!
// We needed four iterations to find our value.
//
// Finding the correct IP-address takes <2 micro-seconds on my Macbook Pro.
// It takes about 20-21 iterations to find the correct IP.
func (c *CountryIPData) AddrCountry(ip string) (country string) {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return ""
	}

	var (
		addrBytes = addr.As4()
		addrInt   = binary.BigEndian.Uint32(addrBytes[:])
		low       = 0
		needle    = len(c.subnetCountries) >> 1 // start searching in the middle of the slice
		high      = len(c.subnetCountries)
	)

	for range 30 {
		subnet := c.subnetCountries[needle]
		// 	fmt.Println("addr", subnet.netAddr, "needle", needle, "low", low, "high", high)

		if high-low == 1 && low > 0 {
			return "" // No country matching, RFC1918
		} else if addrInt < subnet.netAddr {
			high = needle
			needle = (low + high) >> 1
			continue
		} else if addrInt > subnet.lastAddr {
			low = needle
			needle = (low + high) >> 1
			continue
		}
		return c.countriesByCC[subnet.countryCode]
	}
	return ""
}
