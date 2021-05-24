package processips

import (
	"log"
	"strings"
	"time"

	"github.com/shaneu/indahaus/internal/data/ipresult"
	"github.com/shaneu/indahaus/pkg/spamhaus"
)

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

			up := ipresult.UpdateIPResult{
				ResponseCodes: strings.Join(codes, ","),
			}

			_, err = s.dataStore.AddOrUpdate(traceID, ipAddr, up, time.Now())
			if err != nil {
				s.log.Printf("%s : ERROR    : AddOrUpdate for %s %v", traceID, ipAddr, err)
				return
			}
		}(a)
	}
}
