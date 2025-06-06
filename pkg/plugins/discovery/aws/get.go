/*
Copyright 2020 The kconnect Authors.

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

package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/arn"

	"github.com/fidelity/kconnect/pkg/provider/discovery"
)

const (
	expectedNameParts = 2
)

// Get will get the details of a EKS cluster. The clusterID maps to a ARN
func (p *eksClusterProvider) GetCluster(ctx context.Context, input *discovery.GetClusterInput) (*discovery.GetClusterOutput, error) {
	if err := p.setup(input.ConfigSet, input.Identity); err != nil {
		return nil, fmt.Errorf("setting up eks provider: %w", err)
	}

	p.logger.Infow("getting EKS cluster", "id", input.ClusterID)
	clusterName, err := p.getClusterName(input.ClusterID)
	if err != nil {
		return nil, fmt.Errorf("getting cluster name for cluster id %s: %w", input.ClusterID, err)
	}

	cluster, err := p.getClusterConfig(clusterName)
	if err != nil {
		return nil, fmt.Errorf("getting cluster config for %s: %w", input.ClusterID, err)
	}

	return &discovery.GetClusterOutput{
		Cluster: cluster,
	}, nil
}

func (p *eksClusterProvider) getClusterName(clusterID string) (string, error) {
	clusterARN, err := arn.Parse(clusterID)
	if err != nil {
		return "", fmt.Errorf("parsing cluster id as ARN: %w", err)
	}

	parts := strings.Split(clusterARN.Resource, "/")
	if len(parts) != expectedNameParts {
		return "", ErrUnexpectedClusterFormat
	}

	return parts[1], nil
}
