package globalhelper

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/golang/glog"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/container"
)

func LaunchTests(testSuites []string, skipRegEx string) error {
	containerEngine, err := container.SelectEngine()
	fmt.Println("containerEngine = ", containerEngine)
	if err != nil {
		return err
	}

	err = os.Setenv("TNF_CONTAINER_CLIENT", containerEngine)

	if err != nil {
		return err
	}

	glog.V(5).Info(fmt.Sprintf("container engine set to %s", containerEngine))
	testArgs := []string{
		// "-s", skipRegmakeEx,
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

	cmd := exec.Command(fmt.Sprintf("./%s", Configuration.General.TnfEntryPointScript))
	fmt.Println("cmd = ", cmd)
	cmd.Args = append(cmd.Args, testArgs...)
	fmt.Println("cmd.Args = ", cmd.Args)
	cmd.Dir = Configuration.General.TnfRepoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
