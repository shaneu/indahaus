package iplookup

import (
	"log"

	"github.com/shaneu/indahaus/pkg/spamhaus"
)

type Store struct {
	log *log.Logger
}

func New(log *log.Logger) Store {
	return Store{
		log: log,
	}
}

// LookupIPAddresses is just a wrapper over spamhaus.QueryDNSBL that we can add our business logic and
// our cross cutting concerns like logging, which don't belong at the pkg level
func (s Store) LookupIPAddress(traceID string, ip string) ([]string, error) {
	s.log.Printf("%s : query dnsbl for : %s", traceID, ip)
	cs, err := spamhaus.QueryDNSBL(ip)
	if err != nil {
		s.log.Printf("%s : error %v", traceID, err)
		return nil, err
	}

	return cs, nil
}
