package tests

// import (
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// 	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
// 	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
// 	lifecycleparameters "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
// 	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
// 	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
// )

// var _ = Describe("lifecycle lifecycle-scaling", func() {

// 	BeforeEach(func() {
// 		By("Clean namespace before each test")
// 		err := namespaces.Clean(lifecycleparameters.LifecycleNamespace, globalhelper.APIClient)
// 		Expect(err).ToNot(HaveOccurred())
// 	})

// 	It("One deployment, one pod, one container, scale in & out", func() {

// 		deploymentStruct := deployment.DefineDeployment(lifecycleparameters.LifecycleNamespace,
// 			globalhelper.Configuration.General.TnfImage,
// 			lifecycleparameters.TestDeploymentLabels)

// 		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentStruct, lifecycleparameters.WaitingTime)
// 		Expect(err).ToNot(HaveOccurred())

// 		err = globalhelper.LaunchTests(
// 			[]string{lifecycleparameters.LifecycleTestSuiteName},
// 			lifecycleparameters.SkipAllButScalingRegex)

// 		By("Verify test case status in Junit and Claim reports")
// 		err = globalhelper.ValidateIfReportsAreValid(
// 			lifecycleparameters.ScalingDefaultName,
// 			globalparameters.TestCasePassed)
// 		Expect(err).ToNot(HaveOccurred())

// 	})

// })
