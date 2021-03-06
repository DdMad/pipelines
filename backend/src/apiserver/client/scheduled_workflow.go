// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"time"

	"github.com/cenkalti/backoff"
	"github.com/golang/glog"
	swfclient "github.com/kubeflow/pipelines/backend/src/crd/pkg/client/clientset/versioned"
	"github.com/kubeflow/pipelines/backend/src/crd/pkg/client/clientset/versioned/typed/scheduledworkflow/v1beta1"
	"k8s.io/client-go/rest"
)

// creates a new client for the Kubernetes ScheduledWorkflow CRD.
func CreateScheduledWorkflowClientOrFatal(namespace string, initConnectionTimeout time.Duration) v1beta1.ScheduledWorkflowInterface {
	var swfClient v1beta1.ScheduledWorkflowInterface
	var operation = func() error {
		restConfig, err := rest.InClusterConfig()
		if err != nil {
			return err
		}
		swfClientSet := swfclient.NewForConfigOrDie(restConfig)
		swfClient = swfClientSet.ScheduledworkflowV1beta1().ScheduledWorkflows(namespace)
		return nil
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = initConnectionTimeout
	if err := backoff.Retry(operation, b); err != nil {
		glog.Fatalf("Failed to create scheduled workflow client. Error: %v", err)
	}

	return swfClient
}
