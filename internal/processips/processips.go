package processips

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/shaneu/indahaus/internal/data/ipresult"
	"github.com/shaneu/indahaus/pkg/spamhaus"
)

// the 127.255.255.0/24 network representing errors, see see https://www.spamhaus.org/faq/section/DNSBL%20Usage#200
var errIP = []byte("127.255.255.0")
var errIPMask = []byte("255.255.255.0")
var spamhausErrIPNet = net.IPNet{IP: errIP, Mask: errIPMask}

type Store struct {
	log       *log.Logger
	dataStore ipresult.Store
}

func New(log *log.Logger, dataStore ipresult.Store) Store {
	return Store{
		log:       log,
		dataStore: dataStore,
	}
}

func (Store) IsValid(ip string) bool {
	r := net.ParseIP(ip)
	if r == nil {
		return false
	}

	return true
}

// ProcessIPs takes the list of IP address and for each queries the spamhouse API and stores the results
// If the address is new it creates a new row, otherwise it updates the existing row with the latest response codes
func (s Store) ProcessIPs(ips []string, traceID string) {
	// limit the amount of concurrent process executing at the same time to avoid overwhelming resources in the event of a large number of ips
	// to process. We make a channel of empty struct as the type of value is meaningless and struct{}{} doesn't allocate
	// and can't be misinterpreted as having meaning beyond signaling. Starting with 50, we can adjust based on the
	// performance/limits of the spamhaus api
	sem := make(chan struct{}, 50)

	for _, a := range ips {
		// kick off a goroutine to process each ip concurrently
		go func(ipAddr string) {
			// push a value into the semaphore channel, once the channel reaches capacity the other goroutines
			// will block on the send until completed goroutines remove a value from the channel
			sem <- struct{}{}
			defer func() { <-sem }()

			codes, err := spamhaus.QueryDNSBL(ipAddr)
			if err != nil {
				s.log.Printf("%s : ERROR    : spamhaus.QueryDNSBL for %s %v", traceID, ipAddr, err)
				return
			}

			for _, code := range codes {
				// we check if the code returned is in the `127.255.255.0/24` IP network indicating a request error
				if spamhausErrIPNet.Contains(net.ParseIP(code)) {
					s.log.Printf("%s : ERROR    : spamhaus for %s %s", traceID, ipAddr, code)
					return
				}
			}

			up := ipresult.UpdateIPResult{}

			if codes != nil {
				codes := strings.Join(codes, ",")
				up.ResponseCode = &codes
			}

			_, err = s.dataStore.AddOrUpdate(traceID, ipAddr, up, time.Now())
			if err != nil {
				s.log.Printf("%s : ERROR    : AddOrUpdate for %s %v", traceID, ipAddr, err)
				return
			}
		}(a)
	}
}
