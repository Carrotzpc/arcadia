package impl

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.44

import (
	"context"

	"github.com/kubeagi/arcadia/apiserver/graph/generated"
	"github.com/kubeagi/arcadia/apiserver/pkg/modelservice"
)

// CreateModelService is the resolver for the createModelService field.
func (r *modelServiceMutationResolver) CreateModelService(ctx context.Context, obj *generated.ModelServiceMutation, input generated.CreateModelServiceInput) (*generated.ModelService, error) {
	c, err := getClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return modelservice.CreateModelService(ctx, c, input)
}

// UpdateModelService is the resolver for the updateModelService field.
func (r *modelServiceMutationResolver) UpdateModelService(ctx context.Context, obj *generated.ModelServiceMutation, input *generated.UpdateModelServiceInput) (*generated.ModelService, error) {
	c, err := getClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return modelservice.UpdateModelService(ctx, c, input)
}

// DeleteModelService is the resolver for the deleteModelService field.
func (r *modelServiceMutationResolver) DeleteModelService(ctx context.Context, obj *generated.ModelServiceMutation, input *generated.DeleteCommonInput) (*string, error) {
	c, err := getClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return modelservice.DeleteModelService(ctx, c, input)
}

// GetModelService is the resolver for the getModelService field.
func (r *modelServiceQueryResolver) GetModelService(ctx context.Context, obj *generated.ModelServiceQuery, name string, namespace string) (*generated.ModelService, error) {
	c, err := getClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return modelservice.ReadModelService(ctx, c, name, namespace)
}

// ListModelServices is the resolver for the listModelServices field.
func (r *modelServiceQueryResolver) ListModelServices(ctx context.Context, obj *generated.ModelServiceQuery, input *generated.ListModelServiceInput) (*generated.PaginatedResult, error) {
	c, err := getClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return modelservice.ListModelServices(ctx, c, input)
}

// CheckModelService is the resolver for the checkModelService field.
func (r *modelServiceQueryResolver) CheckModelService(ctx context.Context, obj *generated.ModelServiceQuery, input generated.CreateModelServiceInput) (*generated.ModelService, error) {
	c, err := getClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return modelservice.CheckModelService(ctx, c, input)
}

// ModelService is the resolver for the ModelService field.
func (r *mutationResolver) ModelService(ctx context.Context) (*generated.ModelServiceMutation, error) {
	return &generated.ModelServiceMutation{}, nil
}

// ModelService is the resolver for the ModelService field.
func (r *queryResolver) ModelService(ctx context.Context) (*generated.ModelServiceQuery, error) {
	return &generated.ModelServiceQuery{}, nil
}

// ModelServiceMutation returns generated.ModelServiceMutationResolver implementation.
func (r *Resolver) ModelServiceMutation() generated.ModelServiceMutationResolver {
	return &modelServiceMutationResolver{r}
}

// ModelServiceQuery returns generated.ModelServiceQueryResolver implementation.
func (r *Resolver) ModelServiceQuery() generated.ModelServiceQueryResolver {
	return &modelServiceQueryResolver{r}
}

type modelServiceMutationResolver struct{ *Resolver }
type modelServiceQueryResolver struct{ *Resolver }
