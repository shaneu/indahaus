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
)

func (r *iPDetailsResolver) UUID(ctx context.Context, obj *model.IPDetails) (string, error) {
	return obj.ID, nil
}

func (r *mutationResolver) Enqueue(ctx context.Context, ip []string) ([]*model.IPDetails, error) {
	v := ctx.Value(mid.RequestValueKey).(*mid.RequestValues)

	var results []*model.IPDetails

	// we use a buffered channel of len(ip) so in the event one of the
	// goroutines errors we don't want the other goroutines blocking trying to send to the
	// resultChan which won't have any receiver. With a buffered chan those writes won't block so
	// we avoid leaking goroutines. It also gives us the added benifit of supporting partial
	// failure - one goroutine might error but that doesn't mean the others won't succeed in the
	// background. There is a tradeoff that the user will get an error from the graphql response
	// and won't know that the writes were successful in the background, so they may try to enqueue
	// the same ip addresses again, but no data will be corrupted by doing that and hopefully on the
	// next attempt the user won't get an error
	resultChan := make(chan *model.IPDetails, len(ip))
	e := make(chan error, 1)

	for _, a := range ip {
		// kick off a goroutine to process each ip concurrently
		go func(ipAddr string) {
			codes, err := r.IPLookup.LookupIpAddress(v.TraceID, ipAddr)
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

			// The ipresult package returns a type of ipresult.IPResult while the resolver expects a pointer
			// to a type of model.IPDetails. It is my firm belief that the internal packages should have no knowledge
			// of the packages requiring them. In other words it would be incorrect for IPResult.QueryByIP to return a *model.IPDetails
			// as it should have no knowledge of anything to do with layers above it, such as the graphql layer. We could be lazy here
			// and simply cast model.IPDetails(result) since the types _just so happen_ to have the same fields
			// but its better to be explicit here: if the result returned from IPResult.QueryByIP ever has more fields we
			// shouldn't have to come back and refactor this code
			result := model.IPDetails{
				CreatedAt:     res.CreatedAt,
				ID:            res.ID,
				IPAddress:     res.IPAddress,
				ResponseCodes: res.ResponseCodes,
				UpdatedAt:     res.UpdatedAt,
			}

			resultChan <- &result
		}(a)
	}

	work := len(ip)

	for work > 0 {
		// while there are still IP addresses to be processed we
		// use a blocking select statement to get either a result or an error
		select {
		case res := <-resultChan:
			// append the result and keep going
			results = append(results, res)
		case err := <-e:
			// return the error to the caller, don't continue
			return nil, err
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

	// see comment in Enqueue for explanation
	response := model.IPDetails{
		CreatedAt:     result.CreatedAt,
		ID:            result.ID,
		IPAddress:     result.IPAddress,
		ResponseCodes: result.ResponseCodes,
		UpdatedAt:     result.UpdatedAt,
	}

	return &response, nil
}

// IPDetails returns generated.IPDetailsResolver implementation.
func (r *Resolver) IPDetails() generated.IPDetailsResolver { return &iPDetailsResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type iPDetailsResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
