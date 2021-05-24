package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/pkg/errors"

	"github.com/shaneu/indahaus/graph/generated"
	"github.com/shaneu/indahaus/graph/model"
	"github.com/shaneu/indahaus/internal/data/ipresult"
	"github.com/shaneu/indahaus/internal/mid"
)

func (r *mutationResolver) Enqueue(ctx context.Context, ip []string) ([]string, error) {
	v := ctx.Value(mid.RequestValueKey).(*mid.RequestValues)

	// Fire and forget ProcessIPs to let it run in the background
	go r.ProcessIPStore.ProcessIPs(ip, v.TraceID)

	return ip, nil
}

func (r *queryResolver) GetIPDetails(ctx context.Context, ip string) (*model.IPDetails, error) {
	v := ctx.Value(mid.RequestValueKey).(*mid.RequestValues)
	result, err := r.IPResultStore.QueryByIP(v.TraceID, ip)
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
