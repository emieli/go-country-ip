// package v6 uses a radix tree for even better memory usage?
package v6

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"net/netip"
	"os"
	"strings"
)

type CountryIPData struct {
	root          *Node
	countriesByCC map[[2]byte]string // An extra lookup but much more memory efficient
	length        int
}

type IPv4SubnetCountry struct {
	netAddr     uint32
	lastAddr    uint32
	countryCode [2]byte
}

type Node struct {
	children    []*Node
	value       byte
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
	return c.length
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

	// 224.0.0.0/4
	subnets = append(subnets, IPv4SubnetCountry{
		netAddr:     3758096384,
		lastAddr:    math.MaxUint32,
		countryCode: [2]byte{},
	})

	// Turn it into a tree structure
	var firstChild, secondChild, thirdChild, fourthChild *Node
	root := new(Node)
	for _, subnet := range subnets {
		var netAddrBytes, lastAddrBytes [4]byte
		binary.BigEndian.PutUint32(netAddrBytes[:], subnet.netAddr)
		binary.BigEndian.PutUint32(lastAddrBytes[:], subnet.lastAddr)

		// 1st byte
		var found bool
		for _, firstChild = range root.children {
			if netAddrBytes[0] == firstChild.value {
				found = true
				break
			}
		}
		if !found {
			firstChild = &Node{value: netAddrBytes[0]}
			root.children = append(root.children, firstChild)
		}

		// 2nd byte
		found = false
		for _, secondChild = range firstChild.children {
			if netAddrBytes[1] == secondChild.value {
				found = true
				break
			}
		}
		if !found {
			secondChild = &Node{value: netAddrBytes[1]}
			firstChild.children = append(firstChild.children, secondChild)
		}

		// 3rd byte
		found = false
		for _, thirdChild = range secondChild.children {
			if netAddrBytes[2] == thirdChild.value {
				found = true
				break
			}
		}
		if !found {
			thirdChild = &Node{value: netAddrBytes[2]}
			secondChild.children = append(secondChild.children, thirdChild)
		}

		// 4th byte
		found = false
		for _, fourthChild = range thirdChild.children {
			if netAddrBytes[3] == fourthChild.value {
				found = true
				break
			}
		}
		if !found {
			fourthChild = &Node{value: netAddrBytes[3], countryCode: subnet.countryCode}
			thirdChild.children = append(thirdChild.children, fourthChild)
		}
		// 		fmt.Println(firstChild.value, secondChild.value, thirdChild.value, fourthChild.value)
	}
	c.length = len(subnets)

	c.root = root
	return nil
}

func (c *CountryIPData) AddrCountry(ip string) (country string) {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return ""
	}

	var addrBytes = addr.As4()
	node := c.root
	var child *Node

	for _, addrByte := range addrBytes {
		for i := len(node.children) - 1; i >= 0; i-- {
			child = node.children[i]
			if addrByte < child.value {
				continue
			}
			//fmt.Println(addrByte, child.value)
			node = child
			break
		}
	}
	return c.countriesByCC[node.countryCode]
}
