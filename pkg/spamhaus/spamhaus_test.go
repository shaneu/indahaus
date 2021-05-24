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

	tests := []struct {
		when      string
		should    string
		shouldErr bool
		ip        string
	}{
		{
			"\tTest %d:\tWhen looking up an address not in spamhaus.",
			"\t%s\tTest %d:\tShould not return error for no host found : %v",
			false,
			"127.0.0.1",
		},
		{
			"\tTest %d:\tWhen looking up a known spammer address.",
			"\t%s\tTest %d:\tShould return spamhaus response codes : %v",
			false,
			"103.35.191.44",
		},
		{
			"\tTest %d:\tWhen the address is invalid.",
			"\t%s\tTest %d:\tShould return error for invalid IP address : %v",
			true,
			"123.4",
		},
		{
			"\tTest %d:\tWhen the address is IPv6.",
			"\t%s\tTest %d:\tShould return error for unsupported IPv6 lookup : %v",
			true,
			"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		},
	}

	for i, tt := range tests {
		t.Logf(tt.when, i)
		response, err := spamhaus.QueryDNSBL(tt.ip)

		if tt.shouldErr && err != nil {
			t.Logf(tt.should, success, i, err)
			continue
		}

		if tt.shouldErr && err == nil {
			t.Fatalf(tt.should, failure, i, response)
		}

		if !tt.shouldErr && err != nil {
			t.Fatalf(tt.should, failure, i, err)
		}

		if !tt.shouldErr && err == nil {
			t.Logf(tt.should, success, i, response)
		}
	}
}
