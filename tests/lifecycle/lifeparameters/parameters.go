package lifecycleparameters

import (
	"fmt"
	"time"
)

const (
	WaitingTime   = 5 * time.Minute
	RetryInterval = 5
)

var (
	LifecycleNamespace     = "lifecycle-tests"
	testPodLabelPrefixName = "test-network-function.com/generic"
	testPodLabelValue      = "target"
	TestPodLabel           = fmt.Sprintf("%s: %s", testPodLabelPrefixName, testPodLabelValue)
	TestDeploymentLabels   = map[string]string{
		testPodLabelPrefixName: testPodLabelValue,
		"app":                  "test",
	}
	LifecycleTestSuiteName  = "lifecycle"
	SkipAllButShutdownRegex = "lifecycle-pod-high-availability|lifecycle-pod-scheduling|lifecycle-pod-termination-grace-period|lifecycle-pod-owner-type|lifecycle-pod-recreation|lifecycle-scaling|lifecycle-image-pull-policy"
	ShutdownDefaultName     = "lifecycle lifecycle-container-shutdown"
)
