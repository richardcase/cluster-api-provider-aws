/*
Copyright 2021 The Kubernetes Authors.

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

package eks

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/converters"
)

// ReconcileClusterUpdate checks for updates in progress for the EKS cluster
func (s *Service) ReconcileClusterUpdate() error {
	s.scope.Info("Checking for EKS cluster updates", "cluster-name", s.scope.Cluster.Name, "cluster-namespace", s.scope.Cluster.Namespace)

	eksClusterName := s.scope.KubernetesClusterName()

	input := &eks.ListUpdatesInput{
		Name: aws.String(eksClusterName),
	}
	output, err := s.EKSClient.ListUpdates(input)
	if err != nil {
		if awserrors.IsNotFound(err) {
			s.scope.V(4).Info("eks cluster doesn't exist")
			return nil
		}

		return fmt.Errorf("listing cluster updates: %w", err)
	}

	updates := []*ekscontrolplanev1.Update{}
	for _, updateID := range output.UpdateIds {
		describeInput := &eks.DescribeUpdateInput{
			Name:     aws.String(eksClusterName),
			UpdateId: updateID,
		}

		describeOutput, err := s.EKSClient.DescribeUpdate(describeInput)
		if err != nil {
			return fmt.Errorf("getting update details for %s: %w", *updateID, err)
		}

		if describeOutput.Update != nil {
			update := converters.UpdateSDKToUpdate(describeOutput.Update)
			updates = append(updates, update)
		}
	}

	s.scope.ControlPlane.Status.Updates = updates
	if err := s.scope.PatchObject(); err != nil {
		return fmt.Errorf("failed to update control plane: %w", err)
	}
	s.scope.V(2).Info("Reconcile EKS updates completed successfully")

	return nil
}
