package main

import (
	"bytes"
	"encoding/json"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"therebelsource/emulator/execution"
	"therebelsource/emulator/httpUtil"
	"therebelsource/emulator/runner"
	"therebelsource/emulator/staticTypes"
	"therebelsource/emulator/var"
)

var _ = GinkgoDescribe("Single file execution tests", func() {
	GinkgoBeforeEach(func() {
		LoadEnv()
		InitRequiredDirectories(false)
	})

	GinkgoAfterEach(func() {
		gomega.Expect(os.RemoveAll(os.Getenv("PROJECTS_DIR"))).Should(gomega.BeNil())
	})

	GinkgoIt("Should execute a single file in a node LTS environment with imports", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.NodeEsm.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.NodeEsm.Name),
			EmulatorExtension: runner.NodeEsm.Extension,
			EmulatorTag:       string(runner.NodeEsm.Tag),
			EmulatorText:      "console.log('Hello World')",
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello World\n"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a node LTS environment if an infinite loop with a timeout with imports", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.NodeEsm.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.NodeEsm.Name),
			EmulatorExtension: runner.NodeEsm.Extension,
			EmulatorTag:       string(runner.NodeEsm.Tag),
			EmulatorText: `
while(true) {
}
`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal(""))
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should gracefully fail to execute a single file in a node LTS environment because of a syntax error with imports", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.NodeEsm.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.NodeEsm.Name),
			EmulatorExtension: runner.NodeEsm.Extension,
			EmulatorTag:       string(runner.NodeEsm.Tag),
			EmulatorText: `
while(true {
}
`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a node LTS environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.NodeLts.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.NodeLts.Name),
			EmulatorExtension: runner.NodeLts.Extension,
			EmulatorTag:       string(runner.NodeLts.Tag),
			EmulatorText:      "console.log('Hello World')",
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello World\n"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a node LTS environment if an infinite loop with a timeout", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.NodeLts.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.NodeLts.Name),
			EmulatorExtension: runner.NodeLts.Extension,
			EmulatorTag:       string(runner.NodeLts.Tag),
			EmulatorText: `
while(true) {
}
`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal(""))
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should gracefully fail to execute a single file in a node LTS environment because of a syntax error", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.NodeLts.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.NodeLts.Name),
			EmulatorExtension: runner.NodeLts.Extension,
			EmulatorTag:       string(runner.NodeLts.Tag),
			EmulatorText: `
while(true {
}
`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Ruby environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Ruby.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.Ruby.Name),
			EmulatorExtension: runner.Ruby.Extension,
			EmulatorTag:       string(runner.Ruby.Tag),
			EmulatorText:      `print "Hello World"`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello World"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Ruby environment that has a syntax error", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Ruby.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.Ruby.Name),
			EmulatorExtension: runner.Ruby.Extension,
			EmulatorTag:       string(runner.Ruby.Tag),
			EmulatorText:      `prit "Hello World"`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Ruby environment with an infinite loop", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Ruby.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.Ruby.Name),
			EmulatorExtension: runner.Ruby.Extension,
			EmulatorTag:       string(runner.Ruby.Tag),
			EmulatorText: `
loop do
end
`,
		})

		gomega.Expect(result.Result).Should(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Rust environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Rust.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.Rust.Name),
			EmulatorExtension: runner.Rust.Extension,
			EmulatorTag:       string(runner.Rust.Tag),
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

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Rust environment with a syntax error", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Rust.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.Rust.Name),
			EmulatorExtension: runner.Rust.Extension,
			EmulatorTag:       string(runner.Rust.Tag),
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

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Rust environment with a syntax error", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Rust.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.Rust.Name),
			EmulatorExtension: runner.Rust.Extension,
			EmulatorTag:       string(runner.Rust.Tag),
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

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Golang environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.GoLang.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.GoLang.Name),
			EmulatorExtension: runner.GoLang.Extension,
			EmulatorTag:       string(runner.GoLang.Tag),
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

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Golang environment with syntax error", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.GoLang.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.GoLang.Name),
			EmulatorExtension: runner.GoLang.Extension,
			EmulatorTag:       string(runner.GoLang.Tag),
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

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Golang environment with timeout", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.GoLang.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.GoLang.Name),
			EmulatorExtension: runner.GoLang.Extension,
			EmulatorTag:       string(runner.GoLang.Tag),
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

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Mono environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.CSharpMono.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.CSharpMono.Name),
			EmulatorExtension: runner.CSharpMono.Extension,
			EmulatorTag:       string(runner.CSharpMono.Tag),
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

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Mono environment with syntax error", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.CSharpMono.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.CSharpMono.Name),
			EmulatorExtension: runner.CSharpMono.Extension,
			EmulatorTag:       string(runner.CSharpMono.Tag),
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

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Mono environment with infinite loop", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.CSharpMono.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.CSharpMono.Name),
			EmulatorExtension: runner.CSharpMono.Extension,
			EmulatorTag:       string(runner.CSharpMono.Tag),
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

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in C environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.CLang.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.CLang.Name),
			EmulatorExtension: runner.CLang.Extension,
			EmulatorTag:       string(runner.CLang.Tag),
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

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in C environment with syntax error", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.CLang.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.CLang.Name),
			EmulatorExtension: runner.CLang.Extension,
			EmulatorTag:       string(runner.CLang.Tag),
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

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in C environment with infinite loop", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.CLang.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.CLang.Name),
			EmulatorExtension: runner.CLang.Extension,
			EmulatorTag:       string(runner.CLang.Tag),
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

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in C++ environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.CPlus.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.CPlus.Name),
			EmulatorExtension: runner.CPlus.Extension,
			EmulatorTag:       string(runner.CPlus.Tag),
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

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in C++ environment with syntax error", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.CPlus.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.CPlus.Name),
			EmulatorExtension: runner.CPlus.Extension,
			EmulatorTag:       string(runner.CPlus.Tag),
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

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in C++ environment with infinite loop", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.CPlus.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.CPlus.Name),
			EmulatorExtension: runner.CPlus.Extension,
			EmulatorTag:       string(runner.CPlus.Tag),
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

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Python2 environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Python2.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.Python2.Name),
			EmulatorExtension: runner.Python2.Extension,
			EmulatorTag:       string(runner.Python2.Tag),
			EmulatorText: `
print("Hello world")
`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello world\n"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Python3 environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Python3.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.Python3.Name),
			EmulatorExtension: runner.Python3.Extension,
			EmulatorTag:       string(runner.Python3.Tag),
			EmulatorText: `
print("Hello world")
`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello world\n"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Python2 environment with syntax error", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Python2.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.Python2.Name),
			EmulatorExtension: runner.Python2.Extension,
			EmulatorTag:       string(runner.Python2.Tag),
			EmulatorText: `
prit("Hello world")
`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Python3 environment with syntax error", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Python3.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.Python3.Name),
			EmulatorExtension: runner.Python3.Extension,
			EmulatorTag:       string(runner.Python3.Tag),
			EmulatorText: `
prit("Hello world")
`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Python2 environment with na infinite loop", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Python2.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.Python2.Name),
			EmulatorExtension: runner.Python2.Extension,
			EmulatorTag:       string(runner.Python2.Tag),
			EmulatorText: `
while True:
    print("hello")
`,
		})

		gomega.Expect(result.Result).Should(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Python3 environment with na infinite loop", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Python3.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.Python3.Name),
			EmulatorExtension: runner.Python3.Extension,
			EmulatorTag:       string(runner.Python3.Tag),
			EmulatorText: `
while True:
    print("hello")
`,
		})

		gomega.Expect(result.Result).Should(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeFalse())
		gomega.Expect(result.Error).ShouldNot(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in PHP environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Php74.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.Php74.Name),
			EmulatorExtension: runner.Php74.Extension,
			EmulatorTag:       string(runner.Php74.Tag),
			EmulatorText: `
<?php

echo "Hello world";
`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal("\nHello world"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in PHP environment with syntax error", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Php74.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.Php74.Name),
			EmulatorExtension: runner.Php74.Extension,
			EmulatorTag:       string(runner.Php74.Tag),
			EmulatorText: `
<?php

ech "Hello world";
`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in PHP environment with infinite loop", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Php74.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.Php74.Name),
			EmulatorExtension: runner.Php74.Extension,
			EmulatorTag:       string(runner.Php74.Tag),
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

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Haskell environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Haskell.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.Haskell.Name),
			EmulatorExtension: runner.Haskell.Extension,
			EmulatorTag:       string(runner.Haskell.Tag),
			EmulatorText: `
main :: IO ()
main = putStrLn "Hello world"
`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello world\n"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Haskell environment with syntax error", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Haskell.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.Haskell.Name),
			EmulatorExtension: runner.Haskell.Extension,
			EmulatorTag:       string(runner.Haskell.Tag),
			EmulatorText: `
main :: IO ()
man = putStrLn "Hello world"
`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Haskell environment with infinite loop", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Haskell.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(runner.Haskell.Name),
			EmulatorExtension: runner.Haskell.Extension,
			EmulatorTag:       string(runner.Haskell.Tag),
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

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Node latest environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.NodeLts.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `console.log("mile")`, runner.NodeLts, activeSession)
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
		handler := http.HandlerFunc(executeSingleCodeBlockHandler)

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

		var result runner.SingleFileRunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a PHP environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Php74.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `
<?php

echo "mile";
`, runner.Php74, activeSession)
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
		handler := http.HandlerFunc(executeSingleCodeBlockHandler)

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

		var result runner.SingleFileRunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("\nmile"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Ruby environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Ruby.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `puts "mile"`, runner.Ruby, activeSession)
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
		handler := http.HandlerFunc(executeSingleCodeBlockHandler)

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

		var result runner.SingleFileRunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Go environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.GoLang.Tag),
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
`, runner.GoLang, activeSession)
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
		handler := http.HandlerFunc(executeSingleCodeBlockHandler)

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

		var result runner.SingleFileRunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a C# (Mono) environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.CSharpMono.Tag),
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
`, runner.CSharpMono, activeSession)
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
		handler := http.HandlerFunc(executeSingleCodeBlockHandler)

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

		var result runner.SingleFileRunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Python2 environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Python2.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `print("mile")`, runner.Python2, activeSession)
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
		handler := http.HandlerFunc(executeSingleCodeBlockHandler)

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

		var result runner.SingleFileRunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Python3 environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Python3.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `print("mile")`, runner.Python3, activeSession)
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
		handler := http.HandlerFunc(executeSingleCodeBlockHandler)

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

		var result runner.SingleFileRunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Haskell environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.Haskell.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `main = putStrLn "mile"`, runner.Haskell, activeSession)
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
		handler := http.HandlerFunc(executeSingleCodeBlockHandler)

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

		var result runner.SingleFileRunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a C environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.CLang.Tag),
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
`, runner.CLang, activeSession)
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
		handler := http.HandlerFunc(executeSingleCodeBlockHandler)

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

		var result runner.SingleFileRunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a C++ environment", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(runner.CPlus.Tag),
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
`, runner.CPlus, activeSession)
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
		handler := http.HandlerFunc(executeSingleCodeBlockHandler)

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

		var result runner.SingleFileRunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should gracefully fail multiple concurrent requests and stop containers", func() {
		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    10,
				ContainerNum: 1,
				Tag:          string(runner.NodeLts.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		wg := &sync.WaitGroup{}
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				pg := testCreateEmptyPage(activeSession)
				cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
				testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `while(true) {}`, runner.NodeLts, activeSession)
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
				handler := http.HandlerFunc(executeSingleCodeBlockHandler)

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

				var result runner.SingleFileRunResult
				gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

				gomega.Expect(result.Success).Should(gomega.BeFalse())

				wg.Done()
			}(wg)
		}

		wg.Wait()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})
})
