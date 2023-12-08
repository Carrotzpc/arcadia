/*
Copyright 2023 KubeAGI.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package worker

import (
	"context"
	"sort"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"github.com/kubeagi/arcadia/api/base/v1alpha1"
	"github.com/kubeagi/arcadia/graphql-server/go-server/graph/generated"
	gqlmodel "github.com/kubeagi/arcadia/graphql-server/go-server/pkg/model"
	"github.com/kubeagi/arcadia/pkg/utils"
)

const (
	NvidiaGPU = "nvidia.com/gpu"
)

var (
	scheme = schema.GroupVersionResource{Group: v1alpha1.GroupVersion.Group, Version: v1alpha1.GroupVersion.Version, Resource: "workers"}
)

func worker2model(ctx context.Context, c dynamic.Interface, obj *unstructured.Unstructured) *generated.Worker {
	worker := &v1alpha1.Worker{}
	if err := utils.UnstructuredToStructured(obj, worker); err != nil {
		return &generated.Worker{}
	}

	id := string(worker.GetUID())

	labels := make(map[string]interface{})
	for k, v := range obj.GetLabels() {
		labels[k] = v
	}
	annotations := make(map[string]interface{})
	for k, v := range obj.GetAnnotations() {
		annotations[k] = v
	}

	creationtimestamp := worker.GetCreationTimestamp().Time

	// conditioned status
	condition := worker.Status.GetCondition(v1alpha1.TypeReady)
	updateTime := condition.LastTransitionTime.Time

	// Unknown,Pending ,WorkerRunning ,Error
	status := string(condition.Reason)

	// resources
	cpu := worker.Spec.Resources.Limits[v1.ResourceCPU]
	cpuStr := cpu.String()
	memory := worker.Spec.Resources.Limits[v1.ResourceMemory]
	memoryStr := memory.String()
	nvidiaGPU := worker.Spec.Resources.Limits[NvidiaGPU]
	nvidiaGPUStr := nvidiaGPU.String()
	resources := generated.Resources{
		CPU:       &cpuStr,
		Memory:    &memoryStr,
		NvidiaGpu: &nvidiaGPUStr,
	}

	// wrap Worker
	w := generated.Worker{
		ID:                &id,
		Name:              worker.Name,
		Namespace:         worker.Namespace,
		Labels:            labels,
		Annotations:       annotations,
		DisplayName:       &worker.Spec.DisplayName,
		Description:       &worker.Spec.Description,
		Status:            &status,
		CreationTimestamp: &creationtimestamp,
		UpdateTimestamp:   &updateTime,
		Resources:         resources,
		Model:             worker.Spec.Model.Name,
		ModelTypes:        "unknown",
	}

	// read worker's models
	model, err := gqlmodel.ReadModel(ctx, c, worker.Spec.Model.Name, worker.Namespace)
	if err == nil {
		w.ModelTypes = model.Types
	}

	return &w
}

func CreateWorker(ctx context.Context, c dynamic.Interface, input generated.CreateWorkerInput) (*generated.Worker, error) {
	displayName, description := "", ""
	if input.DisplayName != nil {
		displayName = *input.DisplayName
	}
	if input.Description != nil {
		description = *input.Description
	}

	worker := v1alpha1.Worker{
		ObjectMeta: metav1.ObjectMeta{
			Name:      input.Name,
			Namespace: input.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Worker",
			APIVersion: v1alpha1.GroupVersion.String(),
		},
		Spec: v1alpha1.WorkerSpec{
			CommonSpec: v1alpha1.CommonSpec{
				DisplayName: displayName,
				Description: description,
			},
			Model: &v1alpha1.TypedObjectReference{
				Name: input.Model,
				Kind: "Model",
			},
		},
	}

	// cpu & memory
	resources := v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse(input.Resources.CPU),
			v1.ResourceMemory: resource.MustParse(input.Resources.Memory),
		},
	}
	// gpu (only nvidia gpu supported now)
	if input.Resources.NvidiaGpu != nil {
		resources.Limits[NvidiaGPU] = resource.MustParse(*input.Resources.NvidiaGpu)
	}
	worker.Spec.Resources = resources

	unstructuredWorker, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&worker)
	if err != nil {
		return nil, err
	}
	obj, err := c.Resource(schema.GroupVersionResource{Group: v1alpha1.GroupVersion.Group, Version: v1alpha1.GroupVersion.Version, Resource: "workers"}).
		Namespace(input.Namespace).Create(ctx, &unstructured.Unstructured{Object: unstructuredWorker}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return worker2model(ctx, c, obj), nil
}

func UpdateWorker(ctx context.Context, c dynamic.Interface, input *generated.UpdateWorkerInput) (*generated.Worker, error) {
	obj, err := c.Resource(scheme).Namespace(input.Namespace).Get(ctx, input.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	worker := &v1alpha1.Worker{}
	if err := utils.UnstructuredToStructured(obj, worker); err != nil {
		return nil, err
	}

	l := make(map[string]string)
	for k, v := range input.Labels {
		l[k] = v.(string)
	}
	worker.SetLabels(l)

	a := make(map[string]string)
	for k, v := range input.Annotations {
		a[k] = v.(string)
	}
	worker.SetAnnotations(a)

	if input.DisplayName != nil {
		worker.Spec.DisplayName = *input.DisplayName
	}
	if input.Description != nil {
		worker.Spec.Description = *input.Description
	}

	// resources
	if input.Resources != nil {
		// cpu & memory
		resources := v1.ResourceRequirements{
			Limits: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse(input.Resources.CPU),
				v1.ResourceMemory: resource.MustParse(input.Resources.Memory),
			},
		}
		// gpu (only nvidia gpu supported now)
		if input.Resources.NvidiaGpu != nil {
			resources.Limits["nvidia.com/gpu"] = resource.MustParse(*input.Resources.NvidiaGpu)
		}

		worker.Spec.Resources = resources
	}

	unstructuredWorker, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&worker)
	if err != nil {
		return nil, err
	}

	updatedObject, err := c.Resource(scheme).Namespace(input.Namespace).Update(ctx, &unstructured.Unstructured{Object: unstructuredWorker}, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return worker2model(ctx, c, updatedObject), nil
}

func DeleteWorkers(ctx context.Context, c dynamic.Interface, input *generated.DeleteCommonInput) (*string, error) {
	name := ""
	labelSelector, fieldSelector := "", ""
	if input.Name != nil {
		name = *input.Name
	}
	if input.FieldSelector != nil {
		fieldSelector = *input.FieldSelector
	}
	if input.LabelSelector != nil {
		labelSelector = *input.LabelSelector
	}
	resource := c.Resource(schema.GroupVersionResource{Group: v1alpha1.GroupVersion.Group, Version: v1alpha1.GroupVersion.Version, Resource: "workers"})
	if name != "" {
		err := resource.Namespace(input.Namespace).Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil {
			return nil, err
		}
	}
	err := resource.Namespace(input.Namespace).DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: labelSelector,
		FieldSelector: fieldSelector,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func ListWorkers(ctx context.Context, c dynamic.Interface, input generated.ListWorkerInput) (*generated.PaginatedResult, error) {
	keyword, modelTypes, labelSelector, fieldSelector := "", "", "", ""
	page, pageSize := 1, 10
	if input.Keyword != nil {
		keyword = *input.Keyword
	}
	if input.ModelTypes != nil {
		modelTypes = *input.ModelTypes
	}
	if input.FieldSelector != nil {
		fieldSelector = *input.FieldSelector
	}
	if input.LabelSelector != nil {
		labelSelector = *input.LabelSelector
	}
	if input.Page != nil && *input.Page > 0 {
		page = *input.Page
	}
	if input.PageSize != nil && *input.PageSize > 0 {
		pageSize = *input.PageSize
	}

	workerSchema := schema.GroupVersionResource{Group: v1alpha1.GroupVersion.Group, Version: v1alpha1.GroupVersion.Version, Resource: "workers"}
	listOptions := metav1.ListOptions{
		LabelSelector: labelSelector,
		FieldSelector: fieldSelector,
	}
	us, err := c.Resource(workerSchema).Namespace(input.Namespace).List(ctx, listOptions)
	if err != nil {
		return nil, err
	}
	// sort by creation time
	sort.Slice(us.Items, func(i, j int) bool {
		return us.Items[i].GetCreationTimestamp().After(us.Items[j].GetCreationTimestamp().Time)
	})

	totalCount := len(us.Items)

	result := make([]generated.PageNode, 0, pageSize)
	for _, u := range us.Items {
		m := worker2model(ctx, c, &u)
		// filter based on `keyword`
		if keyword != "" {
			if !strings.Contains(m.Name, keyword) && !strings.Contains(*m.DisplayName, keyword) {
				continue
			}
		}
		if modelTypes != "" {
			if !strings.Contains(m.ModelTypes, modelTypes) {
				continue
			}
		}

		result = append(result, m)

		// break if page size matches
		if len(result) == pageSize {
			break
		}
	}

	end := page * pageSize
	if end > totalCount {
		end = totalCount
	}

	return &generated.PaginatedResult{
		TotalCount:  totalCount,
		HasNextPage: end < totalCount,
		Nodes:       result,
	}, nil
}

func ReadWorker(ctx context.Context, c dynamic.Interface, name, namespace string) (*generated.Worker, error) {
	resource := c.Resource(schema.GroupVersionResource{Group: v1alpha1.GroupVersion.Group, Version: v1alpha1.GroupVersion.Version, Resource: "workers"})
	u, err := resource.Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return worker2model(ctx, c, u), nil
}