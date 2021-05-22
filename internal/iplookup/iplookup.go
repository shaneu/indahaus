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

// LookupIpAddresses is just a wrapper over spamhaus.QueryDNSBL that we do our business logic and
// our cross cutting concerns like logging, which don't belong at the pkg level
func (s Store) LookupIpAddress(traceId string, ip string) ([]string, error) {
	s.log.Printf("%s : query dnsbl for : %s", traceId, ip)
	cs, err := spamhaus.QueryDNSBL(ip)
	if err != nil {
		return nil, err
	}

	return cs, nil
}
