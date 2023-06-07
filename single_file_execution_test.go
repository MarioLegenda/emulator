package main

import (
	"bytes"
	http2 "emulator/cmd/http"
	"emulator/pkg/execution"
	"emulator/pkg/httpUtil"
	"emulator/pkg/logger"
	repository2 "emulator/pkg/repository"
	"emulator/pkg/staticTypes"
	"emulator/var"
	"encoding/json"
	"github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
)

var _ = GinkgoDescribe("Single file execution tests", func() {
	GinkgoBeforeEach(func() {
		loadEnv()
		logger.BuildLoggers()
		initRequiredDirectories(false)
	})

	GinkgoAfterEach(func() {
		gomega.Expect(os.RemoveAll(os.Getenv("EXECUTION_DIR"))).Should(gomega.BeNil())
	})

	GinkgoIt("Should execute a single file in a Perl environment", ginkgo.Label("single_file", "perl", "1"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.PerlLts.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.PerlLts.Name),
			EmulatorExtension: repository2.PerlLts.Extension,
			EmulatorTag:       string(repository2.PerlLts.Tag),
			EmulatorText: `#!/usr/bin/perl
  
use strict;
use warnings;
  
print("Hello World\n");
`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello World\n"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Perl environment with an infinite loop", ginkgo.Label("single_file", "perl", "1"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.PerlLts.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.PerlLts.Name),
			EmulatorExtension: repository2.PerlLts.Extension,
			EmulatorTag:       string(repository2.PerlLts.Tag),
			EmulatorText: `#!/usr/bin/perl
for( ; ; ) {
}
`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal(""))
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Lua environment", ginkgo.Label("single_file", "lua", "1"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Lua.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Lua.Name),
			EmulatorExtension: repository2.Lua.Extension,
			EmulatorTag:       string(repository2.Lua.Tag),
			EmulatorText:      `print("Hello World")`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello World\n"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Lua environment with infinite loop", ginkgo.Label("single_file", "lua", "1"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Lua.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Lua.Name),
			EmulatorExtension: repository2.Lua.Extension,
			EmulatorTag:       string(repository2.Lua.Tag),
			EmulatorText: `while true do
end`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal(""))
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a node ESM environment with imports", ginkgo.Label("single_file", "1"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.NodeEsm.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.NodeEsm.Name),
			EmulatorExtension: repository2.NodeEsm.Extension,
			EmulatorTag:       string(repository2.NodeEsm.Tag),
			EmulatorText:      "console.log('Hello World')",
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello World\n"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a node ESM environment if an infinite loop with a timeout with imports", ginkgo.Label("single_file", "3"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.NodeEsm.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.NodeEsm.Name),
			EmulatorExtension: repository2.NodeEsm.Extension,
			EmulatorTag:       string(repository2.NodeEsm.Tag),
			EmulatorText: `
while(true) {
}
`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal(""))
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should gracefully fail to execute a single file in a node ESM environment because of a syntax error with imports", ginkgo.Label("single_file", "4"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.NodeEsm.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.NodeEsm.Name),
			EmulatorExtension: repository2.NodeEsm.Extension,
			EmulatorTag:       string(repository2.NodeEsm.Tag),
			EmulatorText: `
while(true {
}
`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Julia environment", ginkgo.Label("single_file", "2"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Julia.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Julia.Name),
			EmulatorExtension: repository2.Julia.Extension,
			EmulatorTag:       string(repository2.Julia.Tag),
			EmulatorText:      "println(\"Hello World\")",
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello World\n"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a node LTS environment", ginkgo.Label("single_file", "5"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.NodeLts.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.NodeLts.Name),
			EmulatorExtension: repository2.NodeLts.Extension,
			EmulatorTag:       string(repository2.NodeLts.Tag),
			EmulatorText:      "console.log('Hello World')",
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello World\n"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a node LTS environment if an infinite loop with a timeout", ginkgo.Label("single_file", "6"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.NodeLts.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.NodeLts.Name),
			EmulatorExtension: repository2.NodeLts.Extension,
			EmulatorTag:       string(repository2.NodeLts.Tag),
			EmulatorText: `
while(true) {
}
`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal(""))
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should gracefully fail to execute a single file in a node LTS environment because of a syntax error", ginkgo.Label("single_file", "7"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.NodeLts.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.NodeLts.Name),
			EmulatorExtension: repository2.NodeLts.Extension,
			EmulatorTag:       string(repository2.NodeLts.Tag),
			EmulatorText: `
while(true {
}
`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Ruby environment", ginkgo.Label("single_file", "8"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Ruby.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Ruby.Name),
			EmulatorExtension: repository2.Ruby.Extension,
			EmulatorTag:       string(repository2.Ruby.Tag),
			EmulatorText:      `print "Hello World"`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello World"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Ruby environment that has a syntax error", ginkgo.Label("single_file", "9"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Ruby.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Ruby.Name),
			EmulatorExtension: repository2.Ruby.Extension,
			EmulatorTag:       string(repository2.Ruby.Tag),
			EmulatorText:      `prit "Hello World"`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Ruby environment with an infinite loop", ginkgo.Label("single_file", "10"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Ruby.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Ruby.Name),
			EmulatorExtension: repository2.Ruby.Extension,
			EmulatorTag:       string(repository2.Ruby.Tag),
			EmulatorText: `
loop do
end
`,
		})

		gomega.Expect(result.Result).Should(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Rust environment", ginkgo.Label("single_file", "11"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Rust.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Rust.Name),
			EmulatorExtension: repository2.Rust.Extension,
			EmulatorTag:       string(repository2.Rust.Tag),
			EmulatorText: `
fn main() {
    println!("Hello World!");
}
`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello World!\n"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Rust environment with a syntax error", ginkgo.Label("single_file", "12"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Rust.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Rust.Name),
			EmulatorExtension: repository2.Rust.Extension,
			EmulatorTag:       string(repository2.Rust.Tag),
			EmulatorText: `
fn main() {
    printn!("Hello World!");
}
`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Rust environment with a syntax error", ginkgo.Label("single_file", "13"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Rust.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Rust.Name),
			EmulatorExtension: repository2.Rust.Extension,
			EmulatorTag:       string(repository2.Rust.Tag),
			EmulatorText: `
fn main() {
	loop {
	}
}
`,
		})

		gomega.Expect(result.Result).Should(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Golang environment", ginkgo.Label("single_file", "14"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.GoLang.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.GoLang.Name),
			EmulatorExtension: repository2.GoLang.Extension,
			EmulatorTag:       string(repository2.GoLang.Tag),
			EmulatorText: `
package main

import "fmt"

func main() {
  fmt.Println("Hello world")
}
`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello world\n"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Golang environment with syntax error", ginkgo.Label("single_file", "15"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.GoLang.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.GoLang.Name),
			EmulatorExtension: repository2.GoLang.Extension,
			EmulatorTag:       string(repository2.GoLang.Tag),
			EmulatorText: `
package main

import "fmt

func main() {
  fmt.Println("Hello world")
}
`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Golang environment with timeout", ginkgo.Label("single_file", "16"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.GoLang.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.GoLang.Name),
			EmulatorExtension: repository2.GoLang.Extension,
			EmulatorTag:       string(repository2.GoLang.Tag),
			EmulatorText: `
package main

func main() {
	for {
	}
}
`,
		})

		gomega.Expect(result.Result).Should(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Mono environment", ginkgo.Label("single_file", "17"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.CSharpMono.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.CSharpMono.Name),
			EmulatorExtension: repository2.CSharpMono.Extension,
			EmulatorTag:       string(repository2.CSharpMono.Tag),
			EmulatorText: `
namespace HelloWorld
{
    class Hello {         
        static void Main(string[] args)
        {
            System.Console.WriteLine("Hello world");
        }
    }
}
`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello world\n"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Mono environment with syntax error", ginkgo.Label("single_file", "18"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.CSharpMono.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.CSharpMono.Name),
			EmulatorExtension: repository2.CSharpMono.Extension,
			EmulatorTag:       string(repository2.CSharpMono.Tag),
			EmulatorText: `
namespae HelloWorld
{
    class Hello {         
        static void Main(string[] args)
        {
            System.Console.WriteLine("Hello world");
        }
    }
}
`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Mono environment with infinite loop", ginkgo.Label("single_file", "19"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.CSharpMono.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.CSharpMono.Name),
			EmulatorExtension: repository2.CSharpMono.Extension,
			EmulatorTag:       string(repository2.CSharpMono.Tag),
			EmulatorText: `
namespace HelloWorld
{
    class Hello {         
        static void Main(string[] args)
        {
			while(true) {}
        }
    }
}
`,
		})

		gomega.Expect(result.Result).Should(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in C environment", ginkgo.Label("single_file", "20"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.CLang.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.CLang.Name),
			EmulatorExtension: repository2.CLang.Extension,
			EmulatorTag:       string(repository2.CLang.Tag),
			EmulatorText: `
#include <stdio.h>
int main() {
   // printf() displays the string inside quotation
   printf("Hello world");
   return 0;
}

`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello world"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in C environment with syntax error", ginkgo.Label("single_file", "21"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.CLang.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.CLang.Name),
			EmulatorExtension: repository2.CLang.Extension,
			EmulatorTag:       string(repository2.CLang.Tag),
			EmulatorText: `
#include <stdio.h>
int man() {
   // printf() displays the string inside quotation
   printf("Hello world");
   return 0;
}

`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in C environment with infinite loop", ginkgo.Label("single_file", "22"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.CLang.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.CLang.Name),
			EmulatorExtension: repository2.CLang.Extension,
			EmulatorTag:       string(repository2.CLang.Tag),
			EmulatorText: `
#include <stdio.h>
int main() {
	for (;;) {}
}

`,
		})

		gomega.Expect(result.Result).Should(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in C++ environment", ginkgo.Label("single_file", "23"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.CPlus.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.CPlus.Name),
			EmulatorExtension: repository2.CPlus.Extension,
			EmulatorTag:       string(repository2.CPlus.Tag),
			EmulatorText: `
#include <iostream>

int main() {
    std::cout << "Hello world";
    return 0;
}
`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello world"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in C++ environment with syntax error", ginkgo.Label("single_file", "24"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.CPlus.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.CPlus.Name),
			EmulatorExtension: repository2.CPlus.Extension,
			EmulatorTag:       string(repository2.CPlus.Tag),
			EmulatorText: `
#include <iostream>

int man() {
    std::cout << "Hello world";
    return 0;
}
`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in C++ environment with infinite loop", ginkgo.Label("single_file", "25"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.CPlus.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.CPlus.Name),
			EmulatorExtension: repository2.CPlus.Extension,
			EmulatorTag:       string(repository2.CPlus.Tag),
			EmulatorText: `
#include <iostream>

int main() {
	for (;;) {}
}
`,
		})

		gomega.Expect(result.Result).Should(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Python2 environment", ginkgo.Label("single_file", "26"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Python2.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Python2.Name),
			EmulatorExtension: repository2.Python2.Extension,
			EmulatorTag:       string(repository2.Python2.Tag),
			EmulatorText: `
print("Hello world")
`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello world\n"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Python3 environment", ginkgo.Label("single_file", "27"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Python3.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Python3.Name),
			EmulatorExtension: repository2.Python3.Extension,
			EmulatorTag:       string(repository2.Python3.Tag),
			EmulatorText: `
print("Hello world")
`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello world\n"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Python2 environment with syntax error", ginkgo.Label("single_file", "28"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Python2.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Python2.Name),
			EmulatorExtension: repository2.Python2.Extension,
			EmulatorTag:       string(repository2.Python2.Tag),
			EmulatorText: `
prit("Hello world")
`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Python3 environment with syntax error", ginkgo.Label("single_file", "29"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Python3.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Python3.Name),
			EmulatorExtension: repository2.Python3.Extension,
			EmulatorTag:       string(repository2.Python3.Tag),
			EmulatorText: `
prit("Hello world")
`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Python2 environment with na infinite loop", ginkgo.Label("single_file", "30"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Python2.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Python2.Name),
			EmulatorExtension: repository2.Python2.Extension,
			EmulatorTag:       string(repository2.Python2.Tag),
			EmulatorText: `
while True:
    print("hello")
`,
		})

		gomega.Expect(result.Result).Should(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Python3 environment with na infinite loop", ginkgo.Label("single_file", "31"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Python3.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Python3.Name),
			EmulatorExtension: repository2.Python3.Extension,
			EmulatorTag:       string(repository2.Python3.Tag),
			EmulatorText: `
while True:
    print("hello")
`,
		})

		gomega.Expect(result.Result).Should(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in PHP environment", ginkgo.Label("single_file", "32"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Php74.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Php74.Name),
			EmulatorExtension: repository2.Php74.Extension,
			EmulatorTag:       string(repository2.Php74.Tag),
			EmulatorText: `
<?php

echo "Hello world";
`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal("\nHello world"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in PHP environment with syntax error", ginkgo.Label("single_file", "33"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Php74.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Php74.Name),
			EmulatorExtension: repository2.Php74.Extension,
			EmulatorTag:       string(repository2.Php74.Tag),
			EmulatorText: `
<?php

ech "Hello world";
`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in PHP environment with infinite loop", ginkgo.Label("single_file", "34"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Php74.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Php74.Name),
			EmulatorExtension: repository2.Php74.Extension,
			EmulatorTag:       string(repository2.Php74.Tag),
			EmulatorText: `
<?php

    while (1){
    }
`,
		})

		gomega.Expect(result.Result).Should(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Haskell environment", ginkgo.Label("single_file", "35"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Haskell.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Haskell.Name),
			EmulatorExtension: repository2.Haskell.Extension,
			EmulatorTag:       string(repository2.Haskell.Tag),
			EmulatorText: `
main :: IO ()
main = putStrLn "Hello world"
`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello world\n"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Haskell environment with syntax error", ginkgo.Label("single_file", "36"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Haskell.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Haskell.Name),
			EmulatorExtension: repository2.Haskell.Extension,
			EmulatorTag:       string(repository2.Haskell.Tag),
			EmulatorText: `
main :: IO ()
man = putStrLn "Hello world"
`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Haskell environment with infinite loop", ginkgo.Label("single_file", "37"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Haskell.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository2.Haskell.Name),
			EmulatorExtension: repository2.Haskell.Extension,
			EmulatorTag:       string(repository2.Haskell.Tag),
			EmulatorText: `
infi = 
  do
   infi

main = 
  do
    infi

`,
		})

		gomega.Expect(result.Result).Should(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Node latest environment", ginkgo.Label("single_file", "38"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.NodeLts.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `console.log("mile")`, repository2.NodeLts, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(http2.executeSingleCodeBlockHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_RESOURCE))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("Emulator run result"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusOK))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		b, err = json.Marshal(apiResponse.Data)

		gomega.Expect(err).To(gomega.BeNil())

		var result repository2.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a PHP environment", ginkgo.Label("single_file", "39"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Php74.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `
<?php

echo "mile";
`, repository2.Php74, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(http2.executeSingleCodeBlockHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_RESOURCE))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("Emulator run result"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusOK))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		b, err = json.Marshal(apiResponse.Data)

		gomega.Expect(err).To(gomega.BeNil())

		var result repository2.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("\nmile"))

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Ruby environment", ginkgo.Label("single_file", "40"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Ruby.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `puts "mile"`, repository2.Ruby, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(http2.executeSingleCodeBlockHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_RESOURCE))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("Emulator run result"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusOK))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		b, err = json.Marshal(apiResponse.Data)

		gomega.Expect(err).To(gomega.BeNil())

		var result repository2.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Go environment", ginkgo.Label("single_file", "41"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.GoLang.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `
package main

import "fmt"

func main() {
	fmt.Println("mile")
}
`, repository2.GoLang, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(http2.executeSingleCodeBlockHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_RESOURCE))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("Emulator run result"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusOK))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		b, err = json.Marshal(apiResponse.Data)

		gomega.Expect(err).To(gomega.BeNil())

		var result repository2.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a C# (Mono) environment", ginkgo.Label("single_file", "42"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.CSharpMono.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `
using System;

public class HelloWorld
{
    public static void Main(string[] args)
    {
        Console.WriteLine ("mile");
    }
}
`, repository2.CSharpMono, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(http2.executeSingleCodeBlockHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_RESOURCE))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("Emulator run result"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusOK))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		b, err = json.Marshal(apiResponse.Data)

		gomega.Expect(err).To(gomega.BeNil())

		var result repository2.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Python2 environment", ginkgo.Label("single_file", "43"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Python2.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `print("mile")`, repository2.Python2, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(http2.executeSingleCodeBlockHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_RESOURCE))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("Emulator run result"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusOK))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		b, err = json.Marshal(apiResponse.Data)

		gomega.Expect(err).To(gomega.BeNil())

		var result repository2.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Python3 environment", ginkgo.Label("single_file", "44"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Python3.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `print("mile")`, repository2.Python3, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(http2.executeSingleCodeBlockHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_RESOURCE))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("Emulator run result"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusOK))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		b, err = json.Marshal(apiResponse.Data)

		gomega.Expect(err).To(gomega.BeNil())

		var result repository2.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Haskell environment", ginkgo.Label("single_file", "45"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.Haskell.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `main = putStrLn "mile"`, repository2.Haskell, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(http2.executeSingleCodeBlockHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_RESOURCE))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("Emulator run result"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusOK))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		b, err = json.Marshal(apiResponse.Data)

		gomega.Expect(err).To(gomega.BeNil())

		var result repository2.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a C environment", ginkgo.Label("single_file", "46"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.CLang.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `
#include <stdio.h>
int main() {
   printf("mile");
   return 0;
}
`, repository2.CLang, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(http2.executeSingleCodeBlockHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_RESOURCE))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("Emulator run result"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusOK))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		b, err = json.Marshal(apiResponse.Data)

		gomega.Expect(err).To(gomega.BeNil())

		var result repository2.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile"))

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a C++ environment", ginkgo.Label("single_file", "47"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository2.CPlus.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `
#include <iostream>

int main() {
    std::cout << "mile";
    return 0;
}
`, repository2.CPlus, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(http2.executeSingleCodeBlockHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_RESOURCE))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("Emulator run result"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusOK))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		b, err = json.Marshal(apiResponse.Data)

		gomega.Expect(err).To(gomega.BeNil())

		var result repository2.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile"))

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should gracefully fail multiple concurrent requests and stop containers", ginkgo.Label("single_file", "48"), func() {
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    10,
				ContainerNum: 1,
				Tag:          string(repository2.NodeLts.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		wg := &sync.WaitGroup{}
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				pg := testCreateEmptyPage(activeSession)
				cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
				testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `while(true) {}`, repository2.NodeLts, activeSession)
				sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

				bm := map[string]interface{}{
					"uuid": sessionUuid,
				}

				body, err := json.Marshal(bm)

				gomega.Expect(err).To(gomega.BeNil())

				req, err := http.NewRequest("POST", "/api/environment-emulator/execute/single-file", bytes.NewReader(body))

				if err != nil {
					ginkgo.Fail(err.Error())

					return
				}

				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(http2.executeSingleCodeBlockHandler)

				handler.ServeHTTP(rr, req)

				b := rr.Body.Bytes()

				var apiResponse httpUtil.ApiResponse
				err = json.Unmarshal(b, &apiResponse)

				gomega.Expect(err).To(gomega.BeNil())

				gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
				gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

				gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
				gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_RESOURCE))
				gomega.Expect(apiResponse.Message).To(gomega.Equal("Emulator run result"))
				gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(0))
				gomega.Expect(apiResponse.Code).To(gomega.Equal(0))
				gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusOK))
				gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

				b, err = json.Marshal(apiResponse.Data)

				gomega.Expect(err).To(gomega.BeNil())

				var result repository2.RunResult
				gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

				gomega.Expect(result.Success).Should(gomega.BeFalse())

				wg.Done()
			}(wg)
		}

		wg.Wait()

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})
})
