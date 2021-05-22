package spamhaus

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

const dnsZone = "zen.spamhaus.org"

// QueryDNSBL queries the spamhaus dns blacklist and returns any codes found for a given ip. 
// Because it is possible for an ip address to not be listed with spamhaus QueryDNSBL 
// we do not treat an IsNotFound error as an error to be reported, we instead return nil 
// to indicate there were no codes found
func QueryDNSBL(ip string) ([]string, error) {
	// ParseIP returns nil in the case of an ivalid IP and a 16 byte slice in the case of a valid
	// The first 12 bytes are the v4InV6Prefix defined in ip.go, the last 4 bytes are the IPv4 octets
	bs := net.ParseIP(ip)

	if bs == nil {
		return nil, errors.Errorf("invalid ip %s", ip)
	}

	// To4() gives us the nice benifit of removing the v4InV6Prefix _and_ validating that we have an IPv4
	// address, not an IPv6 address
	ipv4 := bs.To4()
	if ipv4 == nil {
		return nil, errors.Errorf("non ipv4 address %v", ip)
	}

	// format the host string with the IP address in reverse order with the dnsbl zone appended to the end
	// ex. 127.0.0.1 -> 1.0.0.127.zen.spamhaus.org
	host := fmt.Sprintf("%d.%d.%d.%d.%s", ipv4[3], ipv4[2], ipv4[1], ipv4[0], dnsZone)

	names, err := net.LookupHost(host)
	if err != nil {
		if v, ok := err.(*net.DNSError); ok {
			if v.IsNotFound {
				return nil, nil
			}
		}

		return nil, err
	}

	return names, nil
}
