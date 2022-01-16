package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("lifecycle lifecycle-container-shutdown", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(lifeparameters.LifecycleNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47311
	It("One deployment, one pod with one container that has preStop field configured", func() {

		By("Define deployment with preStop field configured")
		preStopCommand := []string{"/bin/sh", "-c", "killall -0 tail"}
		preStopDeploymentStruct := deployment.RedefineAllContainersWithPreStopSpec(deployment.DefineDeployment(
			lifeparameters.LifecycleNamespace,
			globalhelper.Configuration.General.TnfImage,
			lifeparameters.TestDeploymentLabels), preStopCommand)
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(preStopDeploymentStruct, lifeparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle lifecycle-container-shutdown test")
		err = globalhelper.LaunchTests(
			[]string{lifeparameters.LifecycleTestSuiteName},
			lifeparameters.SkipAllButShutdownRegex)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			lifeparameters.ShutdownDefaultName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())

	})

	// // 47315
	// It("One deployment, one pod with one container that does not have preStop field configured [negative]", func() {

	// 	By("Define deployment without prestop field configured")
	// 	_ = deployment.DefineDeployment(
	// 		lifeparameters.LifecycleNamespace,
	// 		globalhelper.Configuration.General.TestImage,
	// 		lifeparameters.TestDeploymentLabels)

	// 	By("Start lifecycle lifecycle-container-shutdown test")
	// 	err := globalhelper.LaunchTests(
	// 		[]string{lifeparameters.LifecycleTestSuiteName},
	// 		lifeparameters.SkipAllButShutdownRegex)
	// 	Expect(err).To(HaveOccurred())

	// 	By("Verify test case status in Junit and Claim reports")
	// 	err = globalhelper.ValidateIfReportsAreValid(
	// 		lifeparameters.ShutdownDefaultName,
	// 		globalparameters.TestCaseFailed)
	// 	Expect(err).ToNot(HaveOccurred())

	// })

	// // 47382
	// It("One deployment, several pods, several containers that has preStop field configured", func() {

	// 	By("Define deployment with preStop field configured")
	// 	preStopCommand := []string{"/bin/sh", "-c", "killall -0 tail"}
	// 	deploymentName := "lifecycle"
	// 	replicaDefinedDeployment := deployment.RedefineWithReplicaNumber(
	// 		deployment.DefineDeploymentWithTwoContainers(
	// 			deploymentName,
	// 			lifeparameters.LifecycleNamespace,
	// 			globalhelper.Configuration.General.TnfImage,
	// 			lifeparameters.TestDeploymentLabels), 3)

	// 	preStopDeploymentStruct := deployment.RedefineAllContainersWithPreStopSpec(
	// 		replicaDefinedDeployment, preStopCommand)

	// 	err := globalhelper.CreateAndWaitUntilDeploymentIsReady(
	// 		preStopDeploymentStruct, lifeparameters.WaitingTime)
	// 	Expect(err).ToNot(HaveOccurred())

	// 	By("Start lifecycle lifecycle-container-shutdown test")
	// 	err = globalhelper.LaunchTests(
	// 		[]string{lifeparameters.LifecycleTestSuiteName},
	// 		lifeparameters.SkipAllButShutdownRegex)
	// 	Expect(err).ToNot(HaveOccurred())

	// 	By("Verify test case status in Junit and Claim reports")
	// 	err = globalhelper.ValidateIfReportsAreValid(
	// 		lifeparameters.ShutdownDefaultName,
	// 		globalparameters.TestCasePassed)
	// 	Expect(err).ToNot(HaveOccurred())
	// })

	// // 47383
	// It("Several deployments, several pods, several containers that has preStop field configured", func() {

	// 	preStopCommand := []string{"/bin/sh", "-c", "killall -0 tail"}

	// 	By("Define first deployment with preStop field configured")
	// 	deploymentA := "lifecycleone"
	// 	replicaDefinedDeploymentA := deployment.RedefineWithReplicaNumber(
	// 		deployment.DefineDeploymentWithTwoContainers(
	// 			deploymentA,
	// 			lifeparameters.LifecycleNamespace,
	// 			globalhelper.Configuration.General.TnfImage,
	// 			lifeparameters.TestDeploymentLabels), 3)
	// 	preStopDeploymentStructA := deployment.RedefineAllContainersWithPreStopSpec(
	// 		replicaDefinedDeploymentA, preStopCommand)

	// 	err := globalhelper.CreateAndWaitUntilDeploymentIsReady(
	// 		preStopDeploymentStructA, lifeparameters.WaitingTime)
	// 	Expect(err).ToNot(HaveOccurred())

	// 	By("Define second deployment with preStop field configured", func() {
	// 		deploymentB := "lifecycletwo"
	// 		replicaDefinedDeploymentB := deployment.RedefineWithReplicaNumber(
	// 			deployment.DefineDeploymentWithTwoContainers(
	// 				deploymentB,
	// 				lifeparameters.LifecycleNamespace,
	// 				globalhelper.Configuration.General.TnfImage,
	// 				lifeparameters.TestDeploymentLabels), 3)
	// 		preStopDeploymentStructB := deployment.RedefineAllContainersWithPreStopSpec(
	// 			replicaDefinedDeploymentB, preStopCommand)

	// 		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(
	// 			preStopDeploymentStructB, lifeparameters.WaitingTime)
	// 		Expect(err).ToNot(HaveOccurred())

	// 		By("Start lifecycle lifecycle-container-shutdown test")
	// 		err = globalhelper.LaunchTests(
	// 			[]string{lifeparameters.LifecycleTestSuiteName},
	// 			lifeparameters.SkipAllButShutdownRegex)
	// 		Expect(err).ToNot(HaveOccurred())

	// 		By("Verify test case status in Junit and Claim reports")
	// 		err = globalhelper.ValidateIfReportsAreValid(
	// 			lifeparameters.ShutdownDefaultName,
	// 			globalparameters.TestCasePassed)
	// 		Expect(err).ToNot(HaveOccurred())

	// 	})
	// })

	// // 47384
	// It("One deployment, several pods, several containers, one without preStop field configured [negative]", func() {

	// 	By("Define deployment with preStop field configured")
	// 	preStopCommand := []string{"/bin/sh", "-c", "killall -0 tail"}
	// 	deploymentName := "lifecycle"
	// 	replicaDefinedDeployment := deployment.RedefineWithReplicaNumber(
	// 		deployment.DefineDeploymentWithTwoContainers(
	// 			deploymentName,
	// 			lifeparameters.LifecycleNamespace,
	// 			globalhelper.Configuration.General.TnfImage,
	// 			lifeparameters.TestDeploymentLabels), 3)

	// 	preStopDeploymentStruct := deployment.RedefineFirstContainerWithPreStopSpec(
	// 		replicaDefinedDeployment, preStopCommand)

	// 	err := globalhelper.CreateAndWaitUntilDeploymentIsReady(
	// 		preStopDeploymentStruct, lifeparameters.WaitingTime)
	// 	Expect(err).ToNot(HaveOccurred())

	// 	By("Start lifecycle lifecycle-container-shutdown test")
	// 	err = globalhelper.LaunchTests(
	// 		[]string{lifeparameters.LifecycleTestSuiteName},
	// 		lifeparameters.SkipAllButShutdownRegex)
	// 	Expect(err).To(HaveOccurred())

	// 	By("Verify test case status in Junit and Claim reports")
	// 	err = globalhelper.ValidateIfReportsAreValid(
	// 		lifeparameters.ShutdownDefaultName,
	// 		globalparameters.TestCaseFailed)
	// 	Expect(err).ToNot(HaveOccurred())
	// })

	// // 47385
	// It("Several deployments, several pods, several containers that does not have preStop field configured [negative]", func() {

	// 	By("Define first deployment")
	// 	deploymentA := "lifecycleone"
	// 	replicaDefinedDeploymentA := deployment.RedefineWithReplicaNumber(
	// 		deployment.DefineDeploymentWithTwoContainers(
	// 			deploymentA,
	// 			lifeparameters.LifecycleNamespace,
	// 			globalhelper.Configuration.General.TnfImage,
	// 			lifeparameters.TestDeploymentLabels), 2)

	// 	err := globalhelper.CreateAndWaitUntilDeploymentIsReady(
	// 		replicaDefinedDeploymentA, lifeparameters.WaitingTime)
	// 	Expect(err).ToNot(HaveOccurred())

	// 	By("Define second deployment")
	// 	deploymentB := "lifecycletwo"
	// 	replicaDefinedDeploymentB := deployment.RedefineWithReplicaNumber(
	// 		deployment.DefineDeploymentWithTwoContainers(
	// 			deploymentB,
	// 			lifeparameters.LifecycleNamespace,
	// 			globalhelper.Configuration.General.TnfImage,
	// 			lifeparameters.TestDeploymentLabels), 2)

	// 	err = globalhelper.CreateAndWaitUntilDeploymentIsReady(
	// 		replicaDefinedDeploymentB, lifeparameters.WaitingTime)
	// 	Expect(err).ToNot(HaveOccurred())

	// 	By("Start lifecycle lifecycle-container-shutdown test")
	// 	err = globalhelper.LaunchTests(
	// 		[]string{lifeparameters.LifecycleTestSuiteName},
	// 		lifeparameters.SkipAllButShutdownRegex)
	// 	Expect(err).To(HaveOccurred())

	// 	By("Verify test case status in Junit and Claim reports")
	// 	err = globalhelper.ValidateIfReportsAreValid(
	// 		lifeparameters.ShutdownDefaultName,
	// 		globalparameters.TestCaseFailed)
	// 	Expect(err).ToNot(HaveOccurred())
	// })
})
