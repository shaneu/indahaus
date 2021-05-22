package ipresult

import (
	"context"
	"database/sql"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrInvalidIP = errors.New("IP is not in its proper form")
)

type Store struct {
	log *log.Logger
	db  *sqlx.DB
}

// New returns a configured Store
func New(log *log.Logger, db *sqlx.DB) Store {
	return Store{
		log: log,
		db:  db,
	}
}

// Create inserts a new row into the db
func (s Store) Create(ctx context.Context, traceId string, newIP NewIPResult, now time.Time) (IPResult, error) {
	ipRes := IPResult{
		CreatedAt:     now.UTC(),
		ID:            uuid.New().String(),
		IPAddress:     newIP.IPAddress,
		ResponseCodes: newIP.ResponseCodes,
		UpdatedAt:     now.UTC(),
	}

	const q = `INSERT INTO ip_results
		(id, created_at, updated_at, ip_address, response_codes)	
		VALUES ($1, $2, $3, $4, $5)`

	s.log.Printf("%s : query : %s ipresult.Create", traceId, newIP.IPAddress)

	if _, err := s.db.ExecContext(ctx, q, ipRes.ID, ipRes.CreatedAt, ipRes.UpdatedAt, ipRes.IPAddress, ipRes.ResponseCodes, ","); err != nil {
		return IPResult{}, errors.Wrap(err, "inserting ipresult")
	}

	return ipRes, nil
}

// Update an existing row
func (s Store) Update(ctx context.Context, traceId string, ip string, uIP UpdateIPResult, now time.Time) (IPResult, error) {
	ipRes, err := s.QueryByIP(ctx, traceId, ip)
	if err != nil {
		return IPResult{}, err
	}

	ipRes.UpdatedAt = now.UTC()
	ipRes.ResponseCodes = uIP.ResponseCodes

	const q = `UPDATE ip_results SET "updated_at" = $2,	"response_codes" = $3 WHERE ip_address = $1`

	s.log.Printf("%s : query : %s ipresult.Update", traceId, ip)

	if _, err := s.db.ExecContext(ctx, q, ip, ipRes.UpdatedAt, ipRes.ResponseCodes, ","); err != nil {
		return IPResult{}, errors.Wrap(err, "updating ipresult")
	}

	return ipRes, nil
}

// AddOrUpdate adds or, you guessed it, updates a row
func (s Store) AddOrUpdate(ctx context.Context, traceId string, ip string, uIP UpdateIPResult, now time.Time) (IPResult, error) {
	ipRes, err := s.QueryByIP(ctx, traceId, ip)
	if err != nil {
		if errors.Cause(err) != ErrNotFound {
			return IPResult{}, err
		}

		nIP := NewIPResult{
			IPAddress:     ip,
			ResponseCodes: uIP.ResponseCodes,
		}

		created, err := s.Create(ctx, traceId, nIP, now)
		if err != nil {
			return IPResult{}, err
		}

		return created, nil
	}

	ipRes.UpdatedAt = now.UTC()
	ipRes.ResponseCodes = uIP.ResponseCodes

	const q = `UPDATE ip_results SET "updated_at" = $1,	"response_codes" = $2 WHERE ip_address = $3`

	s.log.Printf("%s : query : %s ipresult.Update", traceId, ip)

	if _, err := s.db.ExecContext(ctx, q, ipRes.UpdatedAt, ipRes.ResponseCodes, ip, ","); err != nil {
		return IPResult{}, errors.Wrap(err, "updating ipresult")
	}

	return ipRes, nil
}

// QueryByIp finds a row by the ip address
func (s Store) QueryByIP(ctx context.Context, traceId string, ip string) (IPResult, error) {
	// we're leveraging net.ParseIP to do our IP validation
	addr := net.ParseIP(ip)
	if addr == nil {
		return IPResult{}, ErrInvalidIP
	}

	const q = `SELECT * FROM ip_results WHERE ip_address = $1`

	s.log.Printf("%s : query : %s ipresult.QueryByIP", traceId, ip)

	var ipRes IPResult
	if err := s.db.GetContext(ctx, &ipRes, q, addr.String()); err != nil {
		if err == sql.ErrNoRows {
			return IPResult{}, ErrNotFound
		}

		return IPResult{}, errors.Wrapf(err, "selecting ip address %q", addr.String())
	}

	return ipRes, nil
}