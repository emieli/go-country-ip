// package v1 is the unoptimized and naive implementation.
// We put all subnets in a giant map where the key is the subnet and
// the value is the country name. This works, but is relatively slow.
// But hey, we have to start somewhere.
package v1

import (
	"bufio"
	"fmt"
	"net/netip"
	"os"
	"strings"
)

type CountryIPData struct {
	subnetCountries map[netip.Prefix]string
}

func NewCountryIPData() (*CountryIPData, error) {
	c := new(CountryIPData)
	if err := c.parseIPInfoCSV(); err != nil {
		return nil, fmt.Errorf("parse ipinfo: %w", err)
	}
	return c, nil
}

func (c CountryIPData) Length() int {
	return len(c.subnetCountries)
}

// slice of uint32, use "binary" search to find entry quickly?
func (c *CountryIPData) parseIPInfoCSV() error {

	ipInfoCSV, err := os.Open("ipinfo_lite.csv")
	if err != nil {
		return fmt.Errorf("open ipinfo file: %v", err)
	}

	subnetCountries := make(map[netip.Prefix]string, 1_500_000)
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

		country := fields[1]
		subnetCountries[subnet] = country
	}

	c.subnetCountries = subnetCountries
	return nil
}

func (c *CountryIPData) AddrCountry(ip string) (country string) {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return ""
	}
	for subnet, country := range c.subnetCountries {
		if subnet.Contains(addr) {
			return country
		}
	}
	return ""
}
