package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shaneu/indahaus/graph/generated"
	"github.com/shaneu/indahaus/graph/model"
	"github.com/shaneu/indahaus/internal/data/ipresult"
	"github.com/shaneu/indahaus/internal/mid"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *mutationResolver) Enqueue(ctx context.Context, ip []string) ([]*model.IPDetails, error) {
	v := ctx.Value(mid.RequestValueKey).(*mid.RequestValues)

	var results []*model.IPDetails

	// We use a buffered channel of len(ip) so in the event one of the child goroutines errors and we return
	// from the parent goroutine we don't want the other children blocking trying to send to the resultChan
	// which won't have a receiver anymore. With a buffered chan those writes won't block so the other children
	// will be able to exit and we avoid leaking goroutines. It also gives us the added benefit of supporting
	// partial failure - one goroutine might error but that doesn't mean the others won't succeed in the background.
	// There is a tradeoff that the user will get an error from the graphql response and won't know that the
	// other writes were successful, so they may try to enqueue the same ip addresses again, but no data will
	// be corrupted by doing that and hopefully on the next attempt the user won't get an error.
	resultChan := make(chan *model.IPDetails, len(ip))
	e := make(chan error, len(ip))

	for _, a := range ip {
		// kick off a goroutine to process each ip concurrently
		go func(ipAddr string) {
			codes, err := r.IPLookup.LookupIPAddress(v.TraceID, ipAddr)
			if err != nil {
				e <- err
			}

			up := ipresult.UpdateIPResult{
				ResponseCodes: strings.Join(codes, ","),
			}

			res, err := r.IPResult.AddOrUpdate(ctx, v.TraceID, ipAddr, up, time.Now())
			if err != nil {
				e <- err
			}

			result := model.IPDetails{
				CreatedAt:     res.CreatedAt,
				UUID:          res.ID,
				IPAddress:     res.IPAddress,
				ResponseCodes: res.ResponseCodes,
				UpdatedAt:     res.UpdatedAt,
			}

			resultChan <- &result
		}(a)
	}

	work := len(ip)

	for work > 0 {
		// while there are still IP addresses to be processed we use a blocking select statement to get either
		// a result or an error
		select {
		case res := <-resultChan:
			// append the result and keep going
			results = append(results, res)
		case <-e:
			// return the error to the caller, don't continue
			return results, gqlerror.Errorf("error processing ip addresses")
		}

		// if we've gotten here we have successfully processed one ip and we decrement our work counter
		work--
	}

	return results, nil
}

func (r *queryResolver) GetIPDetails(ctx context.Context, ip string) (*model.IPDetails, error) {
	v := ctx.Value(mid.RequestValueKey).(*mid.RequestValues)
	result, err := r.IPResult.QueryByIP(ctx, v.TraceID, ip)
	if err != nil {
		if errors.Cause(err) == ipresult.ErrNotFound {
			return nil, nil
		}

		return nil, gqlerror.Errorf("error getting ip address")
	}

	response := model.IPDetails{
		CreatedAt:     result.CreatedAt,
		UUID:          result.ID,
		IPAddress:     result.IPAddress,
		ResponseCodes: result.ResponseCodes,
		UpdatedAt:     result.UpdatedAt,
	}

	return &response, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
