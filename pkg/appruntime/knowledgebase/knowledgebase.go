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

package knowledgebase

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kubeagi/arcadia/api/base/v1alpha1"
	"github.com/kubeagi/arcadia/pkg/appruntime/base"
)

type Knowledgebase struct {
	base.BaseNode
	Instance *v1alpha1.KnowledgeBase
}

func NewKnowledgebase(baseNode base.BaseNode) *Knowledgebase {
	return &Knowledgebase{
		BaseNode: baseNode,
	}
}

func (k *Knowledgebase) Init(ctx context.Context, cli client.Client, _ map[string]any) error {
	ns := base.GetAppNamespace(ctx)
	instance := &v1alpha1.KnowledgeBase{}
	if err := cli.Get(ctx, types.NamespacedName{Namespace: k.Ref.GetNamespace(ns), Name: k.Ref.Name}, instance); err != nil {
		return fmt.Errorf("can't find the knowledgebase in cluster: %w", err)
	}
	k.Instance = instance
	return nil
}

func (k *Knowledgebase) Run(_ context.Context, _ client.Client, args map[string]any) (map[string]any, error) {
	return args, nil
}
