package upgradecluster

import (
	"fmt"

	"github.com/k3s-io/k3s/tests/acceptance/core/service/assert"
	"github.com/k3s-io/k3s/tests/acceptance/core/service/customflag"
	"github.com/k3s-io/k3s/tests/acceptance/core/testcase"
	"github.com/k3s-io/k3s/tests/acceptance/shared/util"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Test:", func() {

	Context("Build Cluster and validate", func() {

		It("Start Up with no issues", func() {
			testcase.TestBuildCluster(GinkgoT(), false)
		})

		It("Checks Node Status", func() {
			testcase.TestNodeStatus(
				assert.NodeAssertReadyStatus(),
				nil,
			)
		})

		It("Checks Pod Status", func() {
			testcase.TestPodStatus(
				assert.PodAssertRestarts(),
				assert.PodAssertReady(),
				assert.PodAssertStatus(),
			)
		})
	})

	It("Verifies ClusterIP Service", func() {
		testcase.TestServiceClusterIp(true)
	})

	It("Verifies NodePort Service", func() {
		testcase.TestServiceNodePort(true)
	})

	It("Verifies LoadBalancer Service", func() {
		testcase.TestServiceLoadBalancer(true)
	})

	It("Verifies Ingress", func() {
		testcase.TestIngress(true)
	})

	It("Verifies Daemonset", func() {
		testcase.TestDaemonset(true)
	})

	It("Verifies Local Path Provisioner storage", func() {
		testcase.TestLocalPathProvisionerStorage(true)
	})

	It("Verifies dns access", func() {
		testcase.TestDnsAccess(true)
	})

	It("Upgrade Manual", func() {
		_ = testcase.TestUpgradeClusterManually(customflag.InstallType.String())
	})

	It("Checks Node Status pos upgrade and validate version", func() {
		testcase.TestNodeStatus(
			assert.NodeAssertReadyStatus(),
			assert.NodeAssertVersionTypeUpgrade(&customflag.InstallType),
		)
	})

	It("Checks Pod Status pos upgrade", func() {
		testcase.TestPodStatus(
			assert.PodAssertRestarts(),
			assert.PodAssertReady(),
			assert.PodAssertStatus(),
		)
	})

	It("Verifies ClusterIP Service after upgrade", func() {
		testcase.TestServiceClusterIp(false)
	})

	It("Verifies NodePort Service after upgrade", func() {
		testcase.TestServiceNodePort(false)
	})

	It("Verifies Ingress after upgrade", func() {
		testcase.TestIngress(false)
	})

	It("Verifies Daemonset after upgrade", func() {
		testcase.TestDaemonset(false)
	})

	It("Verifies LoadBalancer Service after upgrade", func() {
		testcase.TestServiceLoadBalancer(false)
	})

	It("Verifies Local Path Provisioner storage after upgrade", func() {
		testcase.TestLocalPathProvisionerStorage(false)
	})

	It("Verifies dns access after upgrade", func() {
		testcase.TestDnsAccess(false)
	})
})

var _ = BeforeEach(func() {
	if *util.Destroy {
		Skip("Cluster is being Deleted")
	}
})

var _ = AfterEach(func() {
	if CurrentSpecReport().Failed() {
		fmt.Printf("\nFAILED! %s\n", CurrentSpecReport().FullText())
	} else {
		fmt.Printf("\nPASSED! %s\n", CurrentSpecReport().FullText())
	}
})
