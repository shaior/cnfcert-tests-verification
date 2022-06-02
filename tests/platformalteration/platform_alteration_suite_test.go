package platformalteration

import (
	"flag"
	"fmt"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/tests"

	"github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
)

func TestPlatformAlteration(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.Configuration.General.VerificationLogLevel)
	junitPath := globalhelper.Configuration.GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	rr := append([]Reporter{}, reporters.NewJUnitReporter(junitPath))
	RunSpecsWithDefaultAndCustomReporters(t, "CNFCert platform-alteration tests", rr)
}

var _ = BeforeSuite(func() {

	By("Create namespace")
	err := namespaces.Create(parameters.PlatformAlterationNamespace, globalhelper.APIClient)
	Expect(err).ToNot(HaveOccurred())

	By("Define TNF config file")
	err = globalhelper.DefineTnfConfig(
		[]string{parameters.PlatformAlterationNamespace},
		[]string{parameters.TestPodLabel},
		[]string{})
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {

	By(fmt.Sprintf("Remove %s namespace", parameters.PlatformAlterationNamespace))
	err := namespaces.DeleteAndWait(
		globalhelper.APIClient,
		parameters.PlatformAlterationNamespace,
		parameters.WaitingTime,
	)
	Expect(err).ToNot(HaveOccurred())

	By("Remove reports from reports directory")
	err = globalhelper.RemoveContentsFromReportDir()
	Expect(err).ToNot(HaveOccurred())

})
