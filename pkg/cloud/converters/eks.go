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

package converters

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"
)

// AddonSDKToAddonState is used to convert an AWS SDK Addon to a control plane AddonState
func AddonSDKToAddonState(eksAddon *eks.Addon) *ekscontrolplanev1.AddonState {
	addonState := &ekscontrolplanev1.AddonState{
		Name:                  aws.StringValue(eksAddon.AddonName),
		Version:               aws.StringValue(eksAddon.AddonVersion),
		ARN:                   aws.StringValue(eksAddon.AddonArn),
		CreatedAt:             metav1.NewTime(*eksAddon.CreatedAt),
		ModifiedAt:            metav1.NewTime(*eksAddon.ModifiedAt),
		Status:                eksAddon.Status,
		ServiceAccountRoleArn: eksAddon.ServiceAccountRoleArn,
		Issues:                []*ekscontrolplanev1.Issue{},
	}
	if eksAddon.Health != nil {
		for _, issue := range eksAddon.Health.Issues {
			addonState.Issues = append(addonState.Issues, &ekscontrolplanev1.Issue{
				Code:        issue.Code,
				Message:     issue.Message,
				ResourceIDs: issue.ResourceIds,
			})
		}
	}

	return addonState
}

// UpdateSDKToUpdate will convert an AWS SDK EKS Update to a control plane Update
func UpdateSDKToUpdate(sdkUpdate *eks.Update) *ekscontrolplanev1.Update {
	update := &ekscontrolplanev1.Update{
		ID:        *sdkUpdate.Id,
		CreatedAt: metav1.NewTime(*sdkUpdate.CreatedAt),
		Status:    aws.StringValue(sdkUpdate.Status),
		Type:      aws.StringValue(sdkUpdate.Type),
		Params:    map[string]string{},
	}
	issues := []*ekscontrolplanev1.Issue{}
	for _, updateErr := range sdkUpdate.Errors {
		issues = append(issues, ErrorDetailSDKToIssue(updateErr))
	}
	for _, param := range sdkUpdate.Params {
		update.Params[*param.Type] = *param.Value
	}

	return update
}

func ErrorDetailSDKToIssue(errorDetails *eks.ErrorDetail) *ekscontrolplanev1.Issue {
	return &ekscontrolplanev1.Issue{
		Code:        errorDetails.ErrorCode,
		Message:     errorDetails.ErrorMessage,
		ResourceIDs: errorDetails.ResourceIds,
	}
}

// func MapSDKUpdateStatus(status *string) *ekscontrolplanev1.UpdateStatus {
// 	switch *status {
// 	case eks.UpdateStatusInProgress:
// 		return &ekscontrolplanev1.UpdateStatusInProgress
// 	case eks.UpdateStatusSuccessful:
// 		return &ekscontrolplanev1.UpdateStatusSuccessful
// 	case eks.UpdateStatusFailed:
// 		return &ekscontrolplanev1.UpdateStatusFailed
// 	case eks.UpdateStatusCancelled:
// 		return &ekscontrolplanev1.UpdateStatusCancelled
// 	default:
// 		return nil
// 	}
// }

// func MapSDKUpdateType(updateType *string) *ekscontrolplanev1.UpdateType {
// 	switch *updateType {
// 	case eks.UpdateTypeAddonUpdate:
// 		return &ekscontrolplanev1.UpdateTypeAddon
// 	case eks.UpdateTypeConfigUpdate:
// 		return &ekscontrolplanev1.UpdateTypeConfig
// 	case eks.UpdateTypeLoggingUpdate:
// 		return &ekscontrolplanev1.UpdateTypeLogging
// 	case eks.UpdateTypeEndpointAccessUpdate:
// 		return &ekscontrolplanev1.UpdateTypeEndpointAccess
// 	case eks.UpdateTypeVersionUpdate:
// 		return &ekscontrolplanev1.UpdateTypeVersion
// 	case eks.UpdateTypeAssociateIdentityProviderConfig:
// 		return &ekscontrolplanev1.UpdateTypeOIDCAssociate
// 	case eks.UpdateTypeDisassociateIdentityProviderConfig:
// 		return &ekscontrolplanev1.UpdateTypeOIDCDisassociate
// 	default:
// 		return nil
// 	}
// }
