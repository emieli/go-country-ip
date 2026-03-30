// package v6 removes IPv4SubnetCountry.lastAddr that was used in v5.
// This reduces memory usage per entry from 12 bytes to 8 bytes.
package v6

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net/netip"
	"os"
	"strings"
)

type CountryIPData struct {
	subnets       []IPv4SubnetCountry
	countriesByCC map[[2]byte]string // An extra lookup but much more memory efficient
}

type IPv4SubnetCountry struct {
	netAddr     uint32
	countryCode [2]byte
}

func (s IPv4SubnetCountry) NetworkIP() string {
	var netAddr [4]byte
	binary.BigEndian.PutUint32(netAddr[:], s.netAddr)
	prevNetIP := netip.AddrFrom4(netAddr)
	return prevNetIP.String()
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
	return len(c.subnets)
}

func (c *CountryIPData) parseIPInfoCSV() error {
	ipInfoCSV, err := os.Open("ipinfo_lite.csv")
	if err != nil {
		return fmt.Errorf("open ipinfo file: %v", err)
	}

	subnets := make([]IPv4SubnetCountry, 0, 1_000_000)
	subnets = append(subnets, IPv4SubnetCountry{
		netAddr:     0, // 0.0.0.0 - 0.255.255.255
		countryCode: [2]byte{},
	})
	var prevSubnetLastAddr uint32 = 16777215
	countryCode := [2]byte{}
	scanner := bufio.NewScanner(ipInfoCSV)

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ",")
		if fields[0] == "network" {
			continue // first line contains CSV headers, first header is "network"
		}

		if !strings.Contains(fields[0], "/") {
			fields[0] += "/32"
		}
		prefix, err := netip.ParsePrefix(fields[0])
		if err != nil {
			return fmt.Errorf("line %s: %v", scanner.Text(), err)
		} else if prefix.Addr().Is6() {
			break // We only care about ipv4 subnets
		}

		copy(countryCode[:], fields[2])
		c.countriesByCC[countryCode] = fields[1] // country name

		// Convert subnet to IPv4SubnetCountry
		subnetBytes := prefix.Addr().As4()
		subnet := IPv4SubnetCountry{
			netAddr:     binary.BigEndian.Uint32(subnetBytes[:]),
			countryCode: countryCode,
		}
		subnetLastAddr := subnet.netAddr + 1<<(32-prefix.Bits()) - 1

		// subnets not adjacent, place a "empty-country" entry inbetween
		if prevSubnetLastAddr+1 != subnet.netAddr {

			// Create fill entry, filling in RFC1918 ranges etc
			filler := IPv4SubnetCountry{
				netAddr:     prevSubnetLastAddr + 1,
				countryCode: [2]byte{},
			}
			// fmt.Println("fill", filler.NetworkIP())
			subnets = append(subnets, filler)

			// Create this entry
			// fmt.Println("add", subnet.NetworkIP(), string(subnet.countryCode[:]))
			subnets = append(subnets, subnet)
			prevSubnetLastAddr = subnetLastAddr
			continue
		}

		// Subnets are adjacent but different countries
		if subnets[len(subnets)-1].countryCode != subnet.countryCode {
			// fmt.Println("add", subnet.NetworkIP(), string(subnet.countryCode[:]))
			subnets = append(subnets, subnet)
			prevSubnetLastAddr = subnetLastAddr
			continue
		}

		// Extend current entry as same country and subnets are adjacent
		// fmt.Println("ext", subnetCountry.NetworkIP(), string(subnetCountry.countryCode[:]))
		prevSubnetLastAddr = subnetLastAddr
	}

	copySubnets := make([]IPv4SubnetCountry, len(subnets))
	copy(copySubnets, subnets)
	c.subnets = copySubnets
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
		needle    = len(c.subnets) >> 1 // start searching in the middle of the slice
		high      = len(c.subnets)
	)

	for range 30 {
		subnet := c.subnets[needle]
		// 	fmt.Println( "addr", addr.String(), "netAddr", subnet.netAddr, "needle", needle, "low", low, "high", high,)

		if addrInt < subnet.netAddr {
			high = needle
			needle = (low + high) >> 1
			continue
		}

		if needle == len(c.subnets)-1 {
			return c.countriesByCC[subnet.countryCode]
		}

		if addrInt > c.subnets[needle+1].netAddr {
			low = needle
			needle = (low + high) >> 1
			continue
		}

		return c.countriesByCC[subnet.countryCode]
	}
	return ""
}
