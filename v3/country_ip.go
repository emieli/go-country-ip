// package v3 aims to reduce memory usage caused by storing each full country name as a string with each subnet.
// We save on space by storing each country name in a map where the key is the two-letter country-code
// as a [2]byte. This ensures we only store each country name as a string once, but it means we have to
// perform an extra map lookup to get the country name after we find the matching subnet.
// This lowers memory usage from ~225MB to ~100MB.
package v3

import (
	"bufio"
	"fmt"
	"net/netip"
	"os"
	"strings"
)

type CountryIPData struct {
	subnetCountryCodesByFirstByte map[byte]map[netip.Prefix][2]byte
	countriesByCC                 map[[2]byte]string
}

func NewCountryIPData() (*CountryIPData, error) {
	c := new(CountryIPData)
	c.subnetCountryCodesByFirstByte = make(map[byte]map[netip.Prefix][2]byte)
	c.countriesByCC = make(map[[2]byte]string)

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

	var countryCode [2]byte
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

		// Populate country in separate map
		copy(countryCode[:], fields[2])
		country := fields[1]
		if _, exist := c.countriesByCC[countryCode]; !exist {
			c.countriesByCC[countryCode] = country
		}

		// Populate subnet in the usual map
		firstByte := subnet.Addr().As4()[0]
		if _, exist := c.subnetCountryCodesByFirstByte[firstByte]; !exist {
			c.subnetCountryCodesByFirstByte[firstByte] = make(map[netip.Prefix][2]byte)
		}
		c.subnetCountryCodesByFirstByte[firstByte][subnet] = countryCode
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
		for subnet, countryCode := range subnetCountries {
			if subnet.Contains(addr) {
				country := c.countriesByCC[countryCode]
				return country
			}
		}
		return ""
	}
	return ""
}
