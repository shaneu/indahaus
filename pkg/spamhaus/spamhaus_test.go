package spamhaus_test

import (
	"testing"

	"github.com/shaneu/indahaus/pkg/spamhaus"
)

// Success/Failure chars for nicer go test -v output
const (
	success = "\u2713"
	failure = "\u2717"
)

func TestQueryDNSBL(t *testing.T) {
	t.Log("Given the need to be able to perform DNSBL lookups")

	t.Run("notInDNSBL", queryNoResults)
	t.Run("knownSpam", queryKnownSpam)
	t.Run("invalidIPAddress", invalidIP)
	t.Run("ipv6", ipv6)
}

func queryNoResults(t *testing.T) {
	testID := 0
	t.Logf("\tTest %d:\tWhen looking up an address not in spamhaus.", testID)
	response, err := spamhaus.QueryDNSBL("127.0.0.1")
	if err == nil {
		t.Fatalf("\t%s\tTest %d:\tShould return error for no host : %v", failure, testID, response)
	}

	t.Logf("\t%s\tTest %d:\tShould return error for no host : %v", success, testID, err)
}

func queryKnownSpam(t *testing.T) {
	testID := 1
	t.Logf("\tTest %d:\tWhen looking up a known spammer address.", testID)
	response, err := spamhaus.QueryDNSBL("103.35.191.44")
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould return spamhaus response codes : %v", failure, testID, err)
	}

	t.Logf("\t%s\tTest %d:\tShould return spamhaus response codes : %v", success, testID, response)
}

func invalidIP(t *testing.T) {
	testID := 2
	t.Logf("\tTest %d:\tWhen the address is invalid.", testID)
	response, err := spamhaus.QueryDNSBL("123.4")
	if err == nil {
		t.Fatalf("\t%s\tTest %d:\tShould return error for invalid IP address : %v", failure, testID, response)
	}

	t.Logf("\t%s\tTest %d:\tShould return error for invalid IP address : %v", success, testID, err)
}

func ipv6(t *testing.T) {
	testID := 3
	t.Logf("\tTest %d:\tWhen the address is IPv6.", testID)
	response, err := spamhaus.QueryDNSBL("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	if err == nil {
		t.Fatalf("\t%s\tTest %d:\tShould return error for unsupported IPv6 lookup : %v", failure, testID, response)
	}

	t.Logf("\t%s\tTest %d:\tShould return error for unsupported IPv6 lookup : %v", success, testID, err)
}
