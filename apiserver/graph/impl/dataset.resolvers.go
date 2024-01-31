package impl

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.44

import (
	"context"
	"fmt"

	"github.com/kubeagi/arcadia/api/base/v1alpha1"
	"github.com/kubeagi/arcadia/apiserver/graph/generated"
	"github.com/kubeagi/arcadia/apiserver/pkg/dataset"
	"github.com/kubeagi/arcadia/apiserver/pkg/versioneddataset"
)

// Versions is the resolver for the versions field.
func (r *datasetResolver) Versions(ctx context.Context, obj *generated.Dataset, input generated.ListVersionedDatasetInput) (*generated.PaginatedResult, error) {
	c, err := getClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	input.Name = nil
	input.Namespace = &obj.Namespace
	labelSelector := fmt.Sprintf("%s=%s", v1alpha1.LabelVersionedDatasetVersionOwner, obj.Name)
	input.LabelSelector = &labelSelector
	return versioneddataset.ListVersionedDatasets(ctx, c, &input)
}

// CreateDataset is the resolver for the createDataset field.
func (r *datasetMutationResolver) CreateDataset(ctx context.Context, obj *generated.DatasetMutation, input *generated.CreateDatasetInput) (*generated.Dataset, error) {
	c, err := getClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return dataset.CreateDataset(ctx, c, input)
}

// UpdateDataset is the resolver for the updateDataset field.
func (r *datasetMutationResolver) UpdateDataset(ctx context.Context, obj *generated.DatasetMutation, input *generated.UpdateDatasetInput) (*generated.Dataset, error) {
	c, err := getClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return dataset.UpdateDataset(ctx, c, input)
}

// DeleteDatasets is the resolver for the deleteDatasets field.
func (r *datasetMutationResolver) DeleteDatasets(ctx context.Context, obj *generated.DatasetMutation, input *generated.DeleteCommonInput) (*string, error) {
	c, err := getClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return dataset.DeleteDatasets(ctx, c, input)
}

// GetDataset is the resolver for the getDataset field.
func (r *datasetQueryResolver) GetDataset(ctx context.Context, obj *generated.DatasetQuery, name string, namespace string) (*generated.Dataset, error) {
	c, err := getClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return dataset.GetDataset(ctx, c, name, namespace)
}

// ListDatasets is the resolver for the listDatasets field.
func (r *datasetQueryResolver) ListDatasets(ctx context.Context, obj *generated.DatasetQuery, input *generated.ListDatasetInput) (*generated.PaginatedResult, error) {
	c, err := getClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return dataset.ListDatasets(ctx, c, input)
}

// Dataset is the resolver for the Dataset field.
func (r *mutationResolver) Dataset(ctx context.Context) (*generated.DatasetMutation, error) {
	return &generated.DatasetMutation{}, nil
}

// Dataset is the resolver for the Dataset field.
func (r *queryResolver) Dataset(ctx context.Context) (*generated.DatasetQuery, error) {
	return &generated.DatasetQuery{}, nil
}

// Dataset returns generated.DatasetResolver implementation.
func (r *Resolver) Dataset() generated.DatasetResolver { return &datasetResolver{r} }

// DatasetMutation returns generated.DatasetMutationResolver implementation.
func (r *Resolver) DatasetMutation() generated.DatasetMutationResolver {
	return &datasetMutationResolver{r}
}

// DatasetQuery returns generated.DatasetQueryResolver implementation.
func (r *Resolver) DatasetQuery() generated.DatasetQueryResolver { return &datasetQueryResolver{r} }

type datasetResolver struct{ *Resolver }
type datasetMutationResolver struct{ *Resolver }
type datasetQueryResolver struct{ *Resolver }
