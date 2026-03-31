// package v7 uses a four-level nested map where each level represents a byte in each IP-address.
// I thought this might reduce memory usage as identical bytes are only added once.
// But the memory is quite high, ~61 bytes per entry.
// The lookup is also quite slow, ~12 microseconds on avg. I guess generating the hash isn't free.
package v7

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net/netip"
	"os"
	"strings"
)

type CountryIPData struct {
	subnets        map[byte]map[byte]map[byte]map[byte][2]byte
	subnetsCounter int
	countriesByCC  map[[2]byte]string // An extra lookup but much more memory efficient
}

func NewCountryIPData() (*CountryIPData, error) {
	c := new(CountryIPData)
	c.subnets = make(map[byte]map[byte]map[byte]map[byte][2]byte)
	c.countriesByCC = make(map[[2]byte]string)

	if err := c.parseIPInfoCSV(); err != nil {
		return nil, fmt.Errorf("parse ipinfo: %w", err)
	}
	return c, nil
}

func (c CountryIPData) Length() int {
	return c.subnetsCounter
}

func (c *CountryIPData) parseIPInfoCSV() error {
	ipInfoCSV, err := os.Open("ipinfo_lite.csv")
	if err != nil {
		return fmt.Errorf("open ipinfo file: %v", err)
	}

	var countryCode, prevCountryCode [2]byte
	subnets := make(map[byte]map[byte]map[byte]map[byte][2]byte)
	subnets[0] = make(map[byte]map[byte]map[byte][2]byte)
	subnets[0][0] = make(map[byte]map[byte][2]byte)
	subnets[0][0][0] = make(map[byte][2]byte)
	subnets[0][0][0][0] = countryCode
	c.subnetsCounter++

	var prevLastAddr uint32 = 16777215 // 0.255.255.255

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

		subnetBytes := prefix.Addr().As4()
		subnetNetAddr := binary.BigEndian.Uint32(subnetBytes[:])
		subnetLastAddr := subnetNetAddr + 1<<(32-prefix.Bits()) - 1

		if prevCountryCode == countryCode && prevLastAddr+1 == subnetNetAddr {
			// By not creating this entry we extend the previous one
			prevCountryCode = countryCode
			prevLastAddr = subnetLastAddr
			continue
		}

		// Add filler entry if subnets are not adjacent
		if prevLastAddr+1 != subnetNetAddr {
			var fillerNetAddr [4]byte
			binary.BigEndian.PutUint32(fillerNetAddr[:], prevLastAddr+1)
			fillerBytes := netip.AddrFrom4(fillerNetAddr).As4()

			if _, exist := c.subnets[fillerBytes[0]]; !exist {
				c.subnets[fillerBytes[0]] = make(map[byte]map[byte]map[byte][2]byte)
			}
			if _, exist := c.subnets[fillerBytes[0]][fillerBytes[1]]; !exist {
				c.subnets[fillerBytes[0]][fillerBytes[1]] = make(map[byte]map[byte][2]byte)
			}
			if _, exist := c.subnets[fillerBytes[0]][fillerBytes[1]][fillerBytes[2]]; !exist {
				c.subnets[fillerBytes[0]][fillerBytes[1]][fillerBytes[2]] = make(map[byte][2]byte)
			}
			c.subnets[fillerBytes[0]][fillerBytes[1]][fillerBytes[2]][fillerBytes[3]] = [2]byte{}
			c.subnetsCounter++
			// fmt.Println("filler", fillerBytes)
		}

		// Add entry
		if _, exist := c.subnets[subnetBytes[0]]; !exist {
			c.subnets[subnetBytes[0]] = make(map[byte]map[byte]map[byte][2]byte)
		}
		if _, exist := c.subnets[subnetBytes[0]][subnetBytes[1]]; !exist {
			c.subnets[subnetBytes[0]][subnetBytes[1]] = make(map[byte]map[byte][2]byte)
		}
		if _, exist := c.subnets[subnetBytes[0]][subnetBytes[1]][subnetBytes[2]]; !exist {
			c.subnets[subnetBytes[0]][subnetBytes[1]][subnetBytes[2]] = make(map[byte][2]byte)
		}
		c.subnets[subnetBytes[0]][subnetBytes[1]][subnetBytes[2]][subnetBytes[3]] = countryCode
		c.subnetsCounter++
		// 		fmt.Println("subnet", subnetBytes, string(countryCode[:]))

		prevLastAddr = subnetLastAddr
		prevCountryCode = countryCode
	}

	return nil
}

func (c *CountryIPData) AddrCountry(ip string) (country string) {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return ""
	}
	addrBytes := addr.As4()
	// 	fmt.Println(addrBytes)

	firstByte := addrBytes[0]
	for firstByte > 0 {
		if _, exist := c.subnets[firstByte]; exist {
			break
		}
		firstByte--
	}
	// 	fmt.Println(firstByte)

	secondByte := addrBytes[1]
	for secondByte > 0 {
		if _, exist := c.subnets[firstByte][secondByte]; exist {
			break
		}
		secondByte--
	}
	// 	fmt.Println(firstByte, secondByte)

	thirdByte := addrBytes[2]
	for thirdByte > 0 {
		if _, exist := c.subnets[firstByte][secondByte][thirdByte]; exist {
			break
		}
		thirdByte--
	}
	// 	fmt.Println(firstByte, secondByte, thirdByte)

	fourthByte := addrBytes[3]
	for fourthByte > 0 {
		if _, exist := c.subnets[firstByte][secondByte][thirdByte][fourthByte]; exist {
			break
		}
		fourthByte--
	}
	countryCode := c.subnets[firstByte][secondByte][thirdByte][fourthByte]
	// 	fmt.Println(firstByte, secondByte, thirdByte, fourthByte, string(countryCode))
	return c.countriesByCC[countryCode]
}
