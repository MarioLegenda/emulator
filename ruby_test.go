package main

import (
	"emulator/pkg"
	"emulator/pkg/types"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Node latest", func() {
	ginkgo.When("executing", func() {
		ginkgo.It("should fail execution because of syntax error", ginkgo.Label("ruby", "fail"), func() {
			emulator := pkg.NewEmulator(pkg.Options{
				Ruby: pkg.Ruby{
					Workers:    1,
					Containers: 1,
				},
				LogDirectory:       logDir,
				ExecutionDirectory: executionDir,
			})

			result := emulator.RunJob(string(types.Ruby.Name), `put "hello`)

			gomega.Expect(result.Success).Should(gomega.BeTrue())

			emulator.Close()
		})

		ginkgo.It("should run an endless loop and terminate on deadline", ginkgo.Label("ruby", "infinite_loop"), func() {
			emulator := pkg.NewEmulator(pkg.Options{
				Ruby: pkg.Ruby{
					Workers:    1,
					Containers: 1,
				},
				LogDirectory:       logDir,
				ExecutionDirectory: executionDir,
			})

			result := emulator.RunJob(string(types.Ruby.Name), `loop do
end`)

			gomega.Expect(result.Success).Should(gomega.BeFalse())

			emulator.Close()
		})

		ginkgo.It("should run valid code with success", ginkgo.Label("ruby", "success"), func() {
			emulator := pkg.NewEmulator(pkg.Options{
				Ruby: pkg.Ruby{
					Workers:    1,
					Containers: 1,
				},
				LogDirectory:       logDir,
				ExecutionDirectory: executionDir,
			})

			result := emulator.RunJob(string(types.Ruby.Name), `puts "Hello World"`)

			gomega.Expect(result.Success).Should(gomega.BeTrue())

			emulator.Close()
		})
	})
})
