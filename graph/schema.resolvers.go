package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"
	"github.com/shaneu/indahaus/graph/generated"
	"github.com/shaneu/indahaus/graph/model"
	"github.com/shaneu/indahaus/internal/data/ipresult"
	"github.com/shaneu/indahaus/internal/mid"
)

func (r *mutationResolver) Enqueue(ctx context.Context, ip []string) ([]*model.IPDetails, error) {
	v := ctx.Value(mid.RequestValueKey).(*mid.RequestValues)

	var results []*model.IPDetails

	work := len(ip)

	// we're using a buffered channel meaning that writes will not block, thus allowing us to better leverage
	// concurrency in making our requests and persisting the data. We reduce the latency of the send receive
	// as we don't have to wait for the guarantee the data was received but there will be latenency on the
	// send as multiple goroutines trying to send at the same time will have to wait for their writes to be synchronized
	resultChan := make(chan *model.IPDetails, work)
	failedIPErrs := make(chan error, work)

	// limit the amount of concurrent process to avoid overwhelming resources in the event of a large number of ips
	// to process. We make a channel of empty struct as the type of value is meaningless and struct{}{} doesn't allocate
	// and can't be misinterpreted as having meaning beyond signaling. Starting with 100, we can adjust based on the
	// performance/limits of the spamhaus api
	sem := make(chan struct{}, 100)

	for _, a := range ip {
		// kick off a goroutine to process each ip concurrently
		go func(ipAddr string) {
			// push a value into the semaphore channel, once the channel reaches capacity the other goroutines
			// will block on the send until completed goroutines remove a value from the channel
			sem <- struct{}{}
			defer func() { <-sem }()

			codes, err := r.IPLookup.LookupIPAddress(v.TraceID, ipAddr)
			if err != nil {
				failedIPErrs <- errors.Wrap(err, ipAddr)
				return
			}

			up := ipresult.UpdateIPResult{
				ResponseCodes: strings.Join(codes, ","),
			}

			res, err := r.IPResult.AddOrUpdate(ctx, v.TraceID, ipAddr, up, time.Now())
			if err != nil {
				failedIPErrs <- errors.Wrap(err, ipAddr)
				return
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

	for work > 0 {
		// using a blocking select to get a result or a failed ip error
		select {
		case res := <-resultChan:
			// append the result and keep going
			results = append(results, res)
		case err := <-failedIPErrs:
			// add the failed ip address error to the `errors` key of the gql response
			graphql.AddError(ctx, err)
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

		return nil, err
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
