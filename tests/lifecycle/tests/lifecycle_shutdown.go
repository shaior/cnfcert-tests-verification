package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	lifecycleparameters "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/networking/nethelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("lifecycle lifecycle-container-shutdown", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(lifecycleparameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	//47311
	It("One deployment, one pod with one container that has preStop field configured", func() {

		By("Define deployment with preStop field configured")
		preStopCommand := []string{"/bin/sh", "-c", "killall -0 tail"}
		preStopDeploymentStruct := deployment.RedefineWithPreStopSpec(deployment.DefineDeployment(
			lifecycleparameters.LifecycleNamespace,
			globalhelper.Configuration.General.TnfImage,
			lifecycleparameters.TestDeploymentLabels), preStopCommand)
		err := nethelper.CreateAndWaitUntilDeploymentIsReady(preStopDeploymentStruct, lifecycleparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			[]string{lifecycleparameters.LifecycleTestSuiteName},
			lifecycleparameters.SkipAllButShutdownRegex,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			lifecycleparameters.ShutdownDefaultName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// 47315
	It("One deployment, one pod with one container that does not have preStop field configured [negative]", func() {
		By("Define deployment without prestop field configured")
		_ = deployment.DefineDeployment(
			lifecycleparameters.LifecycleNamespace,
			globalhelper.Configuration.General.TestImage,
			lifecycleparameters.TestDeploymentLabels)

		By("Start lifecycle lifecycle-container-shutdown test")
		err := globalhelper.LaunchTests(
			[]string{lifecycleparameters.LifecycleTestSuiteName},
			lifecycleparameters.SkipAllButShutdownRegex,
		)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = nethelper.ValidateIfReportsAreValid(
			lifecycleparameters.ShutdownDefaultName,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})

})
