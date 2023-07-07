package main

import (
	"emulator/pkg"
	"emulator/pkg/types"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Node latest", func() {
	ginkgo.When("executing", func() {
		ginkgo.It("should fail execution because of syntax error", ginkgo.Label("go", "fail"), func() {
			emulator := pkg.NewEmulator(pkg.Options{
				GoLang: pkg.GoLang{
					Workers:    1,
					Containers: 1,
				},
				LogDirectory:       logDir,
				ExecutionDirectory: executionDir,
			})

			result := emulator.RunJob(string(types.GoLang.Name), `package main
import "fmt"
func main() {
	fmt.Printl("something")
}`)

			gomega.Expect(result.Success).Should(gomega.BeTrue())

			emulator.Close()
		})

		ginkgo.It("should run an endless loop and terminate on deadline", ginkgo.Label("go", "infinite_loop"), func() {
			emulator := pkg.NewEmulator(pkg.Options{
				GoLang: pkg.GoLang{
					Workers:    1,
					Containers: 1,
				},
				LogDirectory:       logDir,
				ExecutionDirectory: executionDir,
			})

			result := emulator.RunJob(string(types.GoLang.Name), `package main
func main() {
	for {}
}`)

			gomega.Expect(result.Success).Should(gomega.BeFalse())

			emulator.Close()
		})

		ginkgo.It("should run valid code with success", ginkgo.Label("go", "success"), func() {
			emulator := pkg.NewEmulator(pkg.Options{
				GoLang: pkg.GoLang{
					Workers:    1,
					Containers: 1,
				},
				LogDirectory:       logDir,
				ExecutionDirectory: executionDir,
			})

			result := emulator.RunJob(string(types.GoLang.Name), `package main
import "fmt"
func main() {
	fmt.Println("something")
}`)
			gomega.Expect(result.Success).Should(gomega.BeTrue())

			emulator.Close()
		})
	})
})
