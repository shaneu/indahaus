package ipresult_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	"github.com/shaneu/indahaus/internal/data/ipresult"
	"github.com/shaneu/indahaus/internal/data/schema"
	"github.com/shaneu/indahaus/pkg/database"
)

// Success/Failure chars for nicer go test -v output
const (
	success = "\u2713"
	failure = "\u2717"
)

func setup(t *testing.T) (*log.Logger, *sqlx.DB, func()) {
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("unable to create temp file %v", err)
	}

	cfg := database.Config{
		Uri: fmt.Sprintf("file:%s", tempFile.Name()),
	}

	// We're testing with an actual database here as opposed to mocking. I've seen more bugs than
	// I can name (some of them my own) when a mocking implementation makes some assumption that doesn't
	// bear out in the real world, or proves to be brittle when it comes to changing implementation details.
	// Testing with the real thing gives us piece of mind that if our code passes the tests there's
	// a higher likelyhood it will perform the same in production. Were this a database like postgres,
	// mongodb, or mysql I'd spin up a container per test either by hand or using a lib like https://github.com/ory/dockertest
	db, err := database.Open(cfg)
	if err != nil {
		t.Fatalf("opening database connection: %v", err)
	}

	if err := schema.Migrate(db); err != nil {
		t.Fatalf("unable to migrate schema: %v", err)
	}

	teardown := func() {
		t.Helper()
		db.Close()
		tempFile.Close()
		if err := os.Remove(tempFile.Name()); err != nil {
			t.Fatalf("unable to remove temp file %s : %v", tempFile.Name(), err)
		}
	}

	log := log.New(os.Stdout, "TEST: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	return log, db, teardown
}

func TestIPResult(t *testing.T) {
	log, db, teardown := setup(t)
	t.Cleanup(teardown)

	t.Log("Given the need to work with IP Result records.")
	// ============================================================================
	// Setup: create a ipresult store
	s := ipresult.New(log, db)

	testID := 0

	t.Logf("\tTest %d:\tWhen inserting an IP result.", testID)
	// ============================================================================
	// Create an ip result
	now := time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC)
	traceID := "00000000-0000-0000-0000-000000000000"

	newIP := ipresult.NewIPResult{
		IPAddress:    "199.83.128.60",
		ResponseCode: "127.0.0.2,127.0.0.4",
	}

	ipRes, err := s.Create(traceID, newIP, now)
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to create IP result : %s.", failure, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to create IP result.", success, testID)

	// ============================================================================
	// Query by IP address
	saved, err := s.QueryByIP(traceID, ipRes.IPAddress)
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve result by IP: %s.", failure, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to retrieve record by IP.", success, testID)

	if diff := cmp.Diff(ipRes, saved); diff != "" {
		t.Fatalf("\t%s\tTest %d:\tShould get back the same IP result. Diff:\n %s.", failure, testID, diff)
	}
	t.Logf("\t%s\tTest %d:\tShould get back the same IP result.", success, testID)

	upd := ipresult.UpdateIPResult{
		ResponseCode: "127.0.0.4",
	}

	// ============================================================================
	// Update IP result
	if _, err := s.Update(traceID, ipRes.IPAddress, upd, now); err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to update IP result : %s.", failure, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to update IP result.", success, testID)

	// ============================================================================
	// AddOrUpdate (Add)
	upd = ipresult.UpdateIPResult{
		ResponseCode: "127.0.0.6",
	}
	newIPAddr := "18.205.180.52"

	if _, err := s.AddOrUpdate(traceID, newIPAddr, upd, now); err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to add or update : %s.", failure, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to add or update.", success, testID)

	upd = ipresult.UpdateIPResult{
		ResponseCode: "127.0.1.0",
	}

	// ============================================================================
	// AddOrUpdate  (Update)
	if _, err := s.AddOrUpdate(traceID, ipRes.IPAddress, upd, now); err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to add or update : %s.", failure, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to add or update.", success, testID)

	saved, err = s.QueryByIP(traceID, ipRes.IPAddress)
	if err != nil {
		t.Fatalf("unable to retrieve store result %v", err)
	}

	if diff := cmp.Diff(saved.ResponseCode, upd.ResponseCode); diff != "" {
		t.Fatalf("\t%s\tTest %d:\tShould get back the updated response codes. Diff:\n %s.", failure, testID, diff)
	}
	t.Logf("\t%s\tTest %d:\tShould get back the updated response codes.", success, testID)
}
