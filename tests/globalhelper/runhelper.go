package globalhelper

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/golang/glog"
	"github.com/sirupsen/logrus"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/container"
)

func LaunchTests(testSuites []string, skipRegEx string) error {
	containerEngine, err := container.SelectEngine()
	if err != nil {
		return err
	}

	err = os.Setenv("TNF_CONTAINER_CLIENT", containerEngine)

	if err != nil {
		return err
	}

	glog.V(5).Info(fmt.Sprintf("container engine set to %s", containerEngine))
	testArgs := []string{
		"-k", os.Getenv("KUBECONFIG"),
		"-t", Configuration.General.TnfConfigDir,
		"-o", Configuration.General.TnfReportDir,
		"-i", fmt.Sprintf("%s:%s", Configuration.General.TnfImage, Configuration.General.TnfImageTag),
	}

	if skipRegEx != "" {
		testArgs = append(testArgs, []string{"-s", skipRegEx}...)
		glog.V(5).Info(fmt.Sprintf("set skip regex to %s", skipRegEx))
	}

	if len(testSuites) > 0 {
		testArgs = append(testArgs, "-f")
		for _, testSuite := range testSuites {
			testArgs = append(testArgs, testSuite)
			glog.V(5).Info(fmt.Sprintf("add test suite %s", testSuite))
		}
	}

	argsString := fmt.Sprintf("-o %s ", Configuration.General.TnfReportDir) +
		"-f lifecycle " +
		"-s lifecycle-pod-high-availability lifecycle-pod-scheduling lifecycle-scaling lifecycle-pod-termination-grace-period lifecycle-pod-owner-type lifecycle-container-shutdown lifecycle-image-pull-policy"
	launchCmdArgs := strings.Split(argsString, " ")

	logrus.Infof("Args: %s", strings.Join(launchCmdArgs, " "))
	cmd := exec.Command("/home/sobarzan/go/src/github.com/test-network-function/test-network-function/run-cnf-suites.sh", launchCmdArgs...)

	// cmd := exec.Command(fmt.Sprintf("./%s", Configuration.General.TnfEntryPointScript))
	// fmt.Println("cmd = ", cmd)
	// cmd.Args = append(cmd.Args, testArgs...)
	cmd.Dir = Configuration.General.TnfRepoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
