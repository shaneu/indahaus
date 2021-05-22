package ipresult_test

import (
	"context"
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
		Username: "sqlite",
		Password: "sqlite",
		Uri:      fmt.Sprintf("file:%s", tempFile.Name()),
	}

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
	ctx := context.Background()
	now := time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC)
	traceID := "00000000-0000-0000-0000-000000000000"

	newIP := ipresult.NewIPResult{
		IPAddress:     "199.83.128.60",
		ResponseCodes: "127.0.0.2,127.0.0.4",
	}

	ipRes, err := s.Create(ctx, traceID, newIP, now)
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to create IP result : %s.", failure, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to create IP result.", success, testID)

	// ============================================================================
	// Query by IP address
	saved, err := s.QueryByIP(ctx, traceID, ipRes.IPAddress)
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve result by IP: %s.", failure, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by IP.", success, testID)

	if diff := cmp.Diff(ipRes, saved); diff != "" {
		t.Fatalf("\t%s\tTest %d:\tShould get back the same IP result. Diff:\n %s.", failure, testID, diff)
	}
	t.Logf("\t%s\tTest %d:\tShould get back the same IP result.", success, testID)

	upd := ipresult.UpdateIPResult{
		ResponseCodes: "127.0.0.4",
	}

	// ============================================================================
	// Update IP result
	if _, err := s.Update(ctx, traceID, ipRes.IPAddress, upd, now); err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to update IP result : %s.", failure, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to update IP result.", success, testID)

	// ============================================================================
	// AddOrUpdate (Add)
	upd = ipresult.UpdateIPResult{
		ResponseCodes: "127.0.0.6",
	}
	newIPAddr := "18.205.180.52"

	if _, err := s.AddOrUpdate(ctx, traceID, newIPAddr, upd, now); err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to add or update : %s.", failure, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to add or update.", success, testID)

	upd = ipresult.UpdateIPResult{
		ResponseCodes: "127.0.1.0",
	}

	// ============================================================================
	// AddOrUpdate  (Update)
	if _, err := s.AddOrUpdate(ctx, traceID, ipRes.IPAddress, upd, now); err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to add or update : %s.", failure, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to add or update.", success, testID)

	saved, err = s.QueryByIP(ctx, traceID, ipRes.IPAddress)
	if err != nil {
		t.Fatalf("unable to retrieve store result %v", err)
	}

	if diff := cmp.Diff(saved.ResponseCodes, upd.ResponseCodes); diff != "" {
		t.Fatalf("\t%s\tTest %d:\tShould get back the updated response codes. Diff:\n %s.", failure, testID, diff)
	}
	t.Logf("\t%s\tTest %d:\tShould get back the updated response codes.", success, testID)
}
