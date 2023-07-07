package main

import (
	"emulator/pkg"
	"emulator/pkg/types"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Node latest", func() {
	ginkgo.When("executing", func() {
		ginkgo.It("should fail execution because of syntax error", ginkgo.Label("node_lts", "fail"), func() {
			emulator := pkg.NewEmulator(pkg.Options{
				NodeLts: pkg.NodeLts{
					Workers:    1,
					Containers: 1,
				},
				LogDirectory:       logDir,
				ExecutionDirectory: executionDir,
			})

			result := emulator.RunJob(string(types.NodeLts.Name), `console.lg('something');`)

			gomega.Expect(result.Success).Should(gomega.BeTrue())

			emulator.Close()
		})

		ginkgo.It("should run an endless loop and terminate on deadline", ginkgo.Label("node_lts", "infinite_loop"), func() {
			emulator := pkg.NewEmulator(pkg.Options{
				NodeLts: pkg.NodeLts{
					Workers:    1,
					Containers: 1,
				},
				LogDirectory:       logDir,
				ExecutionDirectory: executionDir,
			})

			result := emulator.RunJob(string(types.NodeLts.Name), `while (true) {}`)

			gomega.Expect(result.Success).Should(gomega.BeFalse())

			emulator.Close()
		})

		ginkgo.It("should run valid code with success", ginkgo.Label("node_lts", "success"), func() {
			emulator := pkg.NewEmulator(pkg.Options{
				NodeLts: pkg.NodeLts{
					Workers:    1,
					Containers: 1,
				},
				LogDirectory:       logDir,
				ExecutionDirectory: executionDir,
			})

			result := emulator.RunJob(string(types.NodeLts.Name), `console.log('code');`)

			gomega.Expect(result.Success).Should(gomega.BeTrue())

			emulator.Close()
		})
	})
})
