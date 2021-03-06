// Copyright 2016-2017 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k8s

import (
	"time"
)

const (
	// BackOffLoopTimeout is the default duration when trying to reach the
	// kube-apiserver.
	BackOffLoopTimeout = 2 * time.Minute

	// maxUpdateRetries is the maximum number of update retries when
	// updating k8s resources
	maxUpdateRetries = 30

	// EnvNodeNameSpec is the environment label used by Kubernetes to
	// specify the node's name.
	EnvNodeNameSpec = "K8S_NODE_NAME"

	// compatibleK8sVersions is the range of k8s versions this cilium is able to
	// work with. It will change as we add new support or deprecate older k8s
	// versions.
	compatibleK8sVersions = "> 1.6"
)
