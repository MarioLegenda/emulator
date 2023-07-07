package main

import (
	"context"
	"emulator/pkg/appErrors"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"os"
	"os/exec"
	"testing"
)

var GomegaRegisterFailHandler = gomega.RegisterFailHandler
var GinkgoFail = ginkgo.Fail
var GinkgoRunSpecs = ginkgo.RunSpecs
var GinkgoBeforeEach = ginkgo.BeforeEach
var GinkgoAfterEach = ginkgo.AfterEach
var GinkgoAfterSuite = ginkgo.AfterSuite
var GinkgoBeforeSuite = ginkgo.BeforeSuite
var GinkgoDescribe = ginkgo.Describe
var GinkgoIt = ginkgo.It

var cancelFn context.CancelFunc

var executionDir = "/home/mario/go/emulator/var/execution"
var logDir = "/home/mario/go/emulator/var/log"

func TestApi(t *testing.T) {
	GomegaRegisterFailHandler(GinkgoFail)
	GinkgoRunSpecs(t, "API Suite")
}

var _ = GinkgoAfterSuite(func() {
	cmd := exec.Command("/usr/bin/docker", "rm", "-f", "$(docker ps -a -q)")

	err := cmd.Start()

	gomega.Expect(err).Should(gomega.BeNil())
	err = cmd.Wait()
	gomega.Expect(err).Should(gomega.BeNil())
})

func testExecutionDirEmpty(dir string) {
	containerDir, err := os.ReadDir(dir)

	gomega.Expect(err).Should(gomega.BeNil())
	gomega.Expect(len(containerDir)).Should(gomega.Equal(0))
}

func testCleanup() {
	cmd := exec.Command("bash", "-c", "/usr/bin/rm -rf /var/www/execution")
	_, err := cmd.Output()

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot do cleanup: %s", err.Error()))
		return
	}
}
