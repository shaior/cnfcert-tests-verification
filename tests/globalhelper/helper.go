package globalhelper

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/golang/glog"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"gopkg.in/yaml.v2"

	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func defineTnfNamespaces(config *globalparameters.TnfConfig, namespaces []string) error {
	if len(namespaces) < 1 {
		return fmt.Errorf("target namespaces cannot be empty list")
	}

	if config == nil {
		return fmt.Errorf("config struct cannot be nil")
	}

	for _, namespace := range namespaces {
		config.TargetNameSpaces = append(config.TargetNameSpaces, globalparameters.TargetNameSpace{
			Name: namespace,
		})
	}

	return nil
}

func isDeploymentReady(operatorNamespace string, deploymentName string) (bool, error) {
	testDeployment, err := APIClient.Deployments(operatorNamespace).Get(
		context.Background(),
		deploymentName,
		metav1.GetOptions{},
	)
	if err != nil {
		return false, err
	}

	if testDeployment.Status.ReadyReplicas > 0 {
		if testDeployment.Status.Replicas == testDeployment.Status.ReadyReplicas {
			return true, nil
		}
	}

	return false, nil
}

// ValidateIfReportsAreValid test if report is valid for given test case.
func ValidateIfReportsAreValid(tcName string, tcExpectedStatus string) error {
	glog.V(5).Info("Verify test case status in Junit report")

	junitTestReport, err := OpenJunitTestReport()

	if err != nil {
		return err
	}

	claimReport, err := OpenClaimReport()

	if err != nil {
		return err
	}

	err = IsExpectedStatusParamValid(tcExpectedStatus)

	if err != nil {
		return err
	}

	isTestCaseInValidStatusInJunitReport := IsTestCasePassedInJunitReport
	isTestCaseInValidStatusInClaimReport := IsTestCasePassedInClaimReport

	if tcExpectedStatus == globalparameters.TestCaseFailed {
		isTestCaseInValidStatusInJunitReport = IsTestCaseFailedInJunitReport
		isTestCaseInValidStatusInClaimReport = IsTestCaseFailedInClaimReport
	}

	if tcExpectedStatus == globalparameters.TestCaseSkipped {
		isTestCaseInValidStatusInJunitReport = IsTestCaseSkippedInJunitReport
		isTestCaseInValidStatusInClaimReport = IsTestCaseSkippedInClaimReport
	}

	if !isTestCaseInValidStatusInJunitReport(junitTestReport, tcName) {
		return fmt.Errorf("test case %s is not in expected %s state in junit report", tcName, tcExpectedStatus)
	}

	glog.V(5).Info("Verify test case status in claim report file")

	testPassed, err := isTestCaseInValidStatusInClaimReport(tcName, *claimReport)

	if err != nil {
		return err
	}

	if !testPassed {
		return fmt.Errorf("test case %s is not in expected %s state in claim report", tcName, tcExpectedStatus)
	}

	return nil
}

// CreateAndWaitUntilDeploymentIsReady creates deployment and wait until all deployment replicas are up and running.
func CreateAndWaitUntilDeploymentIsReady(deployment *v1.Deployment, timeout time.Duration) error {
	runningDeployment, err := APIClient.Deployments(deployment.Namespace).Create(
		context.Background(),
		deployment,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}

	Eventually(func() bool {
		status, err := isDeploymentReady(runningDeployment.Namespace, runningDeployment.Name)
		if err != nil {
			glog.V(5).Info(fmt.Sprintf(
				"deployment %s is not ready, retry in 5 seconds", runningDeployment.Name))

			return false
		}

		return status
	}, timeout, 5*time.Second).Should(Equal(true), "Deployment is not ready")

	return nil
}

func defineTargetPodLabels(config *globalparameters.TnfConfig, targetPodLabels []string) error {
	if len(targetPodLabels) < 1 {
		return fmt.Errorf("target pod labels cannot be empty list")
	}

	for _, targetPodLabel := range targetPodLabels {
		prefixNameValue := strings.Split(targetPodLabel, "/")
		if len(prefixNameValue) != 2 {
			return fmt.Errorf(fmt.Sprintf("target pod label %s is invalid", targetPodLabel))
		}

		prefix := strings.TrimSpace(prefixNameValue[0])
		nameValue := strings.Split(prefixNameValue[1], ":")

		if len(nameValue) != 2 {
			return fmt.Errorf(fmt.Sprintf("target pod label %s is invalid", targetPodLabel))
		}

		name := strings.TrimSpace(nameValue[0])
		value := strings.TrimSpace(nameValue[1])

		config.TargetPodLabels = append(config.TargetPodLabels, globalparameters.PodLabel{
			Prefix: prefix,
			Name:   name,
			Value:  value,
		})
	}

	return nil
}

func defineCertifiedContainersInfo(config *globalparameters.TnfConfig, certifiedContainerInfo []string) error {
	if len(certifiedContainerInfo) < 1 {
		// do not add certifiedcontainerinfo to tnf_config at all in this case
		return nil
	}

	for _, certifiedContainerFields := range certifiedContainerInfo {
		nameRepository := strings.Split(certifiedContainerFields, "/")

		if len(nameRepository) == 1 {
			// certifiedContainerInfo item does not contain separation character
			// use this to add only the Certifiedcontainerinfo field with no sub fields
			var emptyInfo globalparameters.CertifiedContainerRepoInfo
			config.Certifiedcontainerinfo = append(config.Certifiedcontainerinfo, emptyInfo)

			return nil
		}

		if len(nameRepository) != 2 {
			return fmt.Errorf(fmt.Sprintf("certified container info %s is invalid", certifiedContainerFields))
		}

		name := strings.TrimSpace(nameRepository[0])
		repo := strings.TrimSpace(nameRepository[1])

		glog.V(5).Info(fmt.Sprintf("Adding container name:%s repository:%s to configuration", name, repo))

		config.Certifiedcontainerinfo = append(config.Certifiedcontainerinfo, globalparameters.CertifiedContainerRepoInfo{
			Name:       name,
			Repository: repo,
		})
	}

	return nil
}

// DefineTnfConfig creates tnf_config.yml file under tnf config directory.
func DefineTnfConfig(namespaces []string, targetPodLabels []string, certifiedContainerInfo []string) error {
	configFile, err := os.OpenFile(
		path.Join(
			Configuration.General.TnfConfigDir,
			globalparameters.DefaultTnfConfigFileName),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("error opening/creating file: %w", err)
	}
	defer configFile.Close()
	configFileEncoder := yaml.NewEncoder(configFile)
	tnfConfig := globalparameters.TnfConfig{}

	err = defineTnfNamespaces(&tnfConfig, namespaces)
	if err != nil {
		return err
	}

	err = defineTargetPodLabels(&tnfConfig, targetPodLabels)
	if err != nil {
		return err
	}

	err = defineCertifiedContainersInfo(&tnfConfig, certifiedContainerInfo)
	if err != nil {
		return err
	}

	err = configFileEncoder.Encode(tnfConfig)

	glog.V(5).Info(fmt.Sprintf("%s deployed under %s directory",
		globalparameters.DefaultTnfConfigFileName, Configuration.General.TnfConfigDir))

	return err
}

// IsExpectedStatusParamValid validates if requested test status is valid.
func IsExpectedStatusParamValid(status string) error {
	return validateIfParamInAllowedListOfParams(
		status,
		[]string{globalparameters.TestCaseFailed, globalparameters.TestCasePassed, globalparameters.TestCaseSkipped})
}

func validateIfParamInAllowedListOfParams(parameter string, listOfParameters []string) error {
	for _, allowedParameter := range listOfParameters {
		if allowedParameter == parameter {
			return nil
		}
	}

	return fmt.Errorf("parameter %s is not allowed. List of allowed parameters %s", parameter, listOfParameters)
}
