// package v6 uses a radix tree for even better memory usage?
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
	root          Node
	countriesByCC map[[2]byte]string // An extra lookup but much more memory efficient
}

type IPv4SubnetCountry struct {
	netAddr     uint32
	lastAddr    uint32
	countryCode [2]byte
}

type Node struct {
	children    []Node
	netAddr     uint32
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
	return len(c.subnets)
}

func (c *CountryIPData) parseIPInfoCSV() error {
	ipInfoCSV, err := os.Open("ipinfo_lite.csv")
	if err != nil {
		return fmt.Errorf("open ipinfo file: %v", err)
	}

	subnets := make([]IPv4SubnetCountry, 0, 1_000_000)
	subnets = append(subnets, IPv4SubnetCountry{
		netAddr:     0,        // 0.0.0.0
		lastAddr:    16777215, // 0.255.255.255
		countryCode: [2]byte{},
	})
	subnetsIdx := 0 // Tracks index of last added subnet in subnets
	countryCode := [2]byte{}
	scanner := bufio.NewScanner(ipInfoCSV)

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ",")

		prefix := fields[0]
		if !strings.Contains(prefix, "/") {
			prefix += "/32"
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

		copy(countryCode[:], fields[2])
		c.countriesByCC[countryCode] = fields[1] // country name

		// Convert subnet to IPv4SubnetCountry
		subnetBytes := subnet.Addr().As4()
		subnetCountry := IPv4SubnetCountry{countryCode: countryCode}
		subnetCountry.netAddr = binary.BigEndian.Uint32(subnetBytes[:])
		subnetCountry.lastAddr = subnetCountry.netAddr + 1<<(32-subnet.Bits()) - 1

		if subnets[subnetsIdx].countryCode != countryCode {
			// fmt.Println("add", string(countryCode[:]), subnet)
			subnets = append(subnets, subnetCountry)
			subnetsIdx++
			continue
		}

		if subnets[subnetsIdx].lastAddr+1 != subnetCountry.netAddr {
			// subnets not adjacent, place a "empty-country" entry inbetween

			// // Debug
			// var prevNetAddr, prevLastAddr [4]byte
			// binary.BigEndian.PutUint32(prevNetAddr[:], subnets[subnetsIdx].lastAddr+1)
			// binary.BigEndian.PutUint32(prevLastAddr[:], subnetCountry.netAddr-1)
			// prevNetIP := netip.AddrFrom4(prevNetAddr)
			// prevLastIP := netip.AddrFrom4(prevLastAddr)
			// fmt.Println("add --", prevNetIP, prevLastIP)
			// fmt.Println("add", string(countryCode[:]), subnet)

			// Create empty entry filling gap to previous entry
			// This typically adds RFC1918 ranges, etc
			subnets = append(subnets, IPv4SubnetCountry{
				netAddr:     subnets[subnetsIdx].lastAddr + 1,
				lastAddr:    subnetCountry.netAddr - 1,
				countryCode: [2]byte{},
			})
			subnetsIdx++

			// Create this entry
			subnets = append(subnets, subnetCountry)
			subnetsIdx++
			continue
		}

		// Debug
		// var cna, pla, tla [4]byte
		// binary.BigEndian.PutUint32(cna[:], subnets[subnetsIdx].netAddr)
		// binary.BigEndian.PutUint32(pla[:], subnets[subnetsIdx].lastAddr)
		// binary.BigEndian.PutUint32(tla[:], subnetCountry.lastAddr)
		// combNetIP := netip.AddrFrom4(cna)
		// thisLastIP := netip.AddrFrom4(tla)
		// fmt.Println("ext", string(countryCode[:]), combNetIP, thisLastIP)

		// Extend current entry as same country and subnets are adjacent
		subnets[subnetsIdx].lastAddr = subnetCountry.lastAddr
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

	if addr.IsPrivate() {
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
