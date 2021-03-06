/*
Copyright 2018 The Kubernetes Authors.

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

package elb

import (
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi/resourcegroupstaggingapiiface"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
)

// Scope is a scope for use with the ELB reconciling service
type Scope interface {
	cloud.ClusterScoper

	// Network returns the cluster network object.
	Network() *infrav1.Network

	// Subnets returns the cluster subnets.
	Subnets() infrav1.Subnets

	// SecurityGroups returns the cluster security groups as a map, it creates the map if empty.
	SecurityGroups() map[infrav1.SecurityGroupRole]infrav1.SecurityGroup

	// VPC returns the cluster VPC.
	VPC() *infrav1.VPCSpec

	// ControlPlaneLoadBalancer returns the AWSLoadBalancerSpec
	ControlPlaneLoadBalancer() *infrav1.AWSLoadBalancerSpec

	// ControlPlaneLoadBalancerScheme returns the Classic ELB scheme (public or internal facing)
	ControlPlaneLoadBalancerScheme() infrav1.ClassicELBScheme
}

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type Service struct {
	scope                 Scope
	ELBClient             elbiface.ELBAPI
	ResourceTaggingClient resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI
}

// NewService returns a new service given the api clients.
func NewService(elbScope Scope) *Service {
	return &Service{
		scope:                 elbScope,
		ELBClient:             scope.NewELBClient(elbScope, elbScope, elbScope.InfraCluster()),
		ResourceTaggingClient: scope.NewResourgeTaggingClient(elbScope, elbScope, elbScope.InfraCluster()),
	}
}
