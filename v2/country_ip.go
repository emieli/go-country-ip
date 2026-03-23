// package v2 improves on v1 by splitting the prefixes into smaller chunks based on the first byte.
// So instead of one map containing 1M+ prefix entries, we have ~240 maps that contain at most ~50k entries each.
// Iterating through 50k prefixes is way faster than 1M+.
package v2

import (
	"bufio"
	"fmt"
	"net/netip"
	"os"
	"strings"
)

type CountryIPData struct {
	subnetCountryCodesByFirstByte map[byte]map[netip.Prefix]string
}

func NewCountryIPData() (*CountryIPData, error) {
	c := new(CountryIPData)
	c.subnetCountryCodesByFirstByte = make(map[byte]map[netip.Prefix]string)
	if err := c.parseIPInfoCSV(); err != nil {
		return nil, fmt.Errorf("parse ipinfo: %w", err)
	}
	return c, nil
}

func (c CountryIPData) Length() int {
	var total int
	for _, nestedMap := range c.subnetCountryCodesByFirstByte {
		total += len(nestedMap)
	}
	return total
}

func (c *CountryIPData) parseIPInfoCSV() error {
	ipInfoCSV, err := os.Open("ipinfo_lite.csv")
	if err != nil {
		return fmt.Errorf("open ipinfo file: %v", err)
	}
	scanner := bufio.NewScanner(ipInfoCSV)

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

		firstByte := subnet.Addr().As4()[0]
		if _, exist := c.subnetCountryCodesByFirstByte[firstByte]; !exist {
			c.subnetCountryCodesByFirstByte[firstByte] = make(map[netip.Prefix]string)
		}

		country := fields[1]
		c.subnetCountryCodesByFirstByte[firstByte][subnet] = country
	}

	return nil
}

func (c *CountryIPData) AddrCountry(ip string) (country string) {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return ""
	}
	firstByte := uint8(addr.As4()[0])
	for firstByte > 0 {
		// 		fmt.Println(firstByte)
		subnetCountries, exist := c.subnetCountryCodesByFirstByte[firstByte]
		if !exist {
			firstByte -= 1
			continue
		}
		for subnet, country := range subnetCountries {
			if subnet.Contains(addr) {
				return country
			}
		}
		return ""
	}
	return ""
}
