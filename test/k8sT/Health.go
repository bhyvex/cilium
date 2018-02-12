// Copyright 2018 Authors of Cilium
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

package k8sTest

import (
	"fmt"
	"sync"

	. "github.com/cilium/cilium/test/ginkgo-ext"
	"github.com/cilium/cilium/test/helpers"

	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var testName = "K8sValidatedHealthTest"

var _ = Describe(testName, func() {

	var kubectl *helpers.Kubectl
	var logger *logrus.Entry
	var once sync.Once
	initialize := func() {
		logger = log.WithFields(logrus.Fields{"testName": testName})
		logger.Info("Starting")

		kubectl = helpers.CreateKubectl(helpers.K8s1VMName(), logger)
		err := kubectl.DeployCiliumDS(helpers.DefaultK8sTCiliumOpts())
		Expect(err).Should(BeNil())
	}

	BeforeEach(func() {
		once.Do(initialize)
	})

	AfterFailed(func() {
		kubectl.CiliumReport(helpers.KubeSystemNamespace, []string{
			"cilium service list",
			"cilium endpoint list",
			"cilium policy get"})
	})

	JustAfterEach(func() {
		kubectl.ValidateNoErrorsOnLogs(CurrentGinkgoTestDescription().Duration)
	})

	AfterEach(func() {
		err := kubectl.WaitCleanAllTerminatingPods()
		Expect(err).To(BeNil(), "Terminating containers are not deleted after timeout")
	})

	checkIP := func(pod, ip string) {
		jsonpath := fmt.Sprintf("{.cluster.nodes[*].primary-address.*}")
		ciliumCmd := fmt.Sprintf("cilium status -o jsonpath='%s'", jsonpath)
		status := kubectl.CiliumExec(pod, ciliumCmd)
		Expect(status.Output().String()).Should(ContainSubstring(ip))
		status.ExpectSuccess()
	}

	It("checks cilium-health status between nodes", func() {
		cilium1, cilium1IP, err := kubectl.GetCilium(helpers.K8s1)
		Expect(err).To(BeNil(), "Unable to get cilium pod from node %s", helpers.K8s1)
		cilium2, cilium2IP, err := kubectl.GetCilium(helpers.K8s2)
		Expect(err).To(BeNil(), "Unable to get cilium pod from node %s", helpers.K8s1)

		By(fmt.Sprintf("checking that cilium API exposes health instances"))
		checkIP(cilium1, cilium1IP)
		checkIP(cilium1, cilium2IP)
		checkIP(cilium2, cilium1IP)
		checkIP(cilium2, cilium2IP)

		By(fmt.Sprintf("checking that `cilium-health --probe` succeeds"))
		healthCmd := fmt.Sprintf("cilium-health status --probe -o json")
		status := kubectl.CiliumExec(cilium1, healthCmd)
		Expect(status.Output()).ShouldNot(ContainSubstring("error"))
		status.ExpectSuccess()

		apiPaths := []string{
			"endpoint.icmp",
			"endpoint.http",
			"host.\"primary-address\".icmp",
			"host.\"primary-address\".http",
		}
		for node := 0; node <= 1; node++ {
			for _, path := range apiPaths {
				jqArg := fmt.Sprintf(".nodes[%d].%s.status", node, path)
				By(fmt.Sprintf("checking API response for '%s'", jqArg))
				healthCmd := fmt.Sprintf("cilium-health status -o json | jq '%s'", jqArg)
				status := kubectl.CiliumExec(cilium1, healthCmd)
				Expect(status.Output().String()).Should(ContainSubstring("null"))
				status.ExpectSuccess()
			}
		}
	}, 30)
})
