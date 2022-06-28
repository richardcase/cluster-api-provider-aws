/*
Copyright 2022 The Kubernetes Authors.

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

package gc

import (
	"context"
	"fmt"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/annotations"
	"sigs.k8s.io/cluster-api/util/patch"
)

// Enable is used to enable external resource garbage collection for a cluster.
func (c *CmdProcessor) Enable(ctx context.Context) error {
	infraObj, err := c.getInfraCluster(ctx)
	if err != nil {
		return err
	}

	patchHelper, err := patch.NewHelper(infraObj, c.client)
	if err != nil {
		return fmt.Errorf("creating patch helper: %w", err)
	}

	annotations.SetExternalResourceGC(infraObj, false)

	if err := patchHelper.Patch(ctx, infraObj); err != nil {
		return fmt.Errorf("patching infra cluster with gc annotation: %w", err)
	}

	return nil
}
