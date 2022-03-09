package tests

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifehelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("lifecycle lifecycle-deployment-scaling", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(lifeparameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		By("Enable intrusive tests")
		err = os.Setenv("TNF_NON_INTRUSIVE_ONLY", "false")
		Expect(err).ToNot(HaveOccurred())
	})

	// 47398
	specName := "One deployment, one pod, one container, scale in and out"
	It(specName, func() {
		tcNameForReport := globalhelper.ConvertSpecNameToFileName(specName)
		By("Define Deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(lifehelper.DefineDeployment(1, 1, "lifecycleput"),
			lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("start lifecycle lifecycle-deployment-scaling")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.DeploymentScalingDefaultName,
			tcNameForReport,
			lifeparameters.SkipAllButScalingRegex)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.DeploymentScalingDefaultName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})
})
