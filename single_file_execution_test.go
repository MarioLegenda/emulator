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
	"therebelsource/emulator/repository"
	"therebelsource/emulator/staticTypes"
	"therebelsource/emulator/var"
)

var _ = GinkgoDescribe("Single file execution tests", func() {
	GinkgoBeforeEach(func() {
		loadEnv()
		initRequiredDirectories(false)
	})

	GinkgoAfterEach(func() {
		gomega.Expect(os.RemoveAll(os.Getenv("EXECUTION_DIR"))).Should(gomega.BeNil())
	})

	GinkgoIt("Should execute a single file in a node LTS environment with imports", func() {
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.NodeEsm.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.NodeEsm.Name),
			EmulatorExtension: repository.NodeEsm.Extension,
			EmulatorTag:       string(repository.NodeEsm.Tag),
			EmulatorText:      "console.log('Hello World')",
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello World\n"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a node LTS environment if an infinite loop with a timeout with imports", func() {
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.NodeEsm.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.NodeEsm.Name),
			EmulatorExtension: repository.NodeEsm.Extension,
			EmulatorTag:       string(repository.NodeEsm.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.NodeEsm.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.NodeEsm.Name),
			EmulatorExtension: repository.NodeEsm.Extension,
			EmulatorTag:       string(repository.NodeEsm.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.NodeLts.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.NodeLts.Name),
			EmulatorExtension: repository.NodeLts.Extension,
			EmulatorTag:       string(repository.NodeLts.Tag),
			EmulatorText:      "console.log('Hello World')",
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello World\n"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a node LTS environment if an infinite loop with a timeout", func() {
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.NodeLts.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.NodeLts.Name),
			EmulatorExtension: repository.NodeLts.Extension,
			EmulatorTag:       string(repository.NodeLts.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.NodeLts.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.NodeLts.Name),
			EmulatorExtension: repository.NodeLts.Extension,
			EmulatorTag:       string(repository.NodeLts.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Ruby.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.Ruby.Name),
			EmulatorExtension: repository.Ruby.Extension,
			EmulatorTag:       string(repository.Ruby.Tag),
			EmulatorText:      `print "Hello World"`,
		})

		gomega.Expect(result.Result).Should(gomega.Equal("Hello World"))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Ruby environment that has a syntax error", func() {
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Ruby.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.Ruby.Name),
			EmulatorExtension: repository.Ruby.Extension,
			EmulatorTag:       string(repository.Ruby.Tag),
			EmulatorText:      `prit "Hello World"`,
		})

		gomega.Expect(result.Result).ShouldNot(gomega.BeEmpty())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Error).Should(gomega.BeNil())

		testExecutionDirEmpty()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute code in Ruby environment with an infinite loop", func() {
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Ruby.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.Ruby.Name),
			EmulatorExtension: repository.Ruby.Extension,
			EmulatorTag:       string(repository.Ruby.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Rust.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.Rust.Name),
			EmulatorExtension: repository.Rust.Extension,
			EmulatorTag:       string(repository.Rust.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Rust.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.Rust.Name),
			EmulatorExtension: repository.Rust.Extension,
			EmulatorTag:       string(repository.Rust.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Rust.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.Rust.Name),
			EmulatorExtension: repository.Rust.Extension,
			EmulatorTag:       string(repository.Rust.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.GoLang.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.GoLang.Name),
			EmulatorExtension: repository.GoLang.Extension,
			EmulatorTag:       string(repository.GoLang.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.GoLang.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.GoLang.Name),
			EmulatorExtension: repository.GoLang.Extension,
			EmulatorTag:       string(repository.GoLang.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.GoLang.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.GoLang.Name),
			EmulatorExtension: repository.GoLang.Extension,
			EmulatorTag:       string(repository.GoLang.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.CSharpMono.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.CSharpMono.Name),
			EmulatorExtension: repository.CSharpMono.Extension,
			EmulatorTag:       string(repository.CSharpMono.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.CSharpMono.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.CSharpMono.Name),
			EmulatorExtension: repository.CSharpMono.Extension,
			EmulatorTag:       string(repository.CSharpMono.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.CSharpMono.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.CSharpMono.Name),
			EmulatorExtension: repository.CSharpMono.Extension,
			EmulatorTag:       string(repository.CSharpMono.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.CLang.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.CLang.Name),
			EmulatorExtension: repository.CLang.Extension,
			EmulatorTag:       string(repository.CLang.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.CLang.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.CLang.Name),
			EmulatorExtension: repository.CLang.Extension,
			EmulatorTag:       string(repository.CLang.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.CLang.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.CLang.Name),
			EmulatorExtension: repository.CLang.Extension,
			EmulatorTag:       string(repository.CLang.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.CPlus.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.CPlus.Name),
			EmulatorExtension: repository.CPlus.Extension,
			EmulatorTag:       string(repository.CPlus.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.CPlus.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.CPlus.Name),
			EmulatorExtension: repository.CPlus.Extension,
			EmulatorTag:       string(repository.CPlus.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.CPlus.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.CPlus.Name),
			EmulatorExtension: repository.CPlus.Extension,
			EmulatorTag:       string(repository.CPlus.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Python2.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.Python2.Name),
			EmulatorExtension: repository.Python2.Extension,
			EmulatorTag:       string(repository.Python2.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Python3.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.Python3.Name),
			EmulatorExtension: repository.Python3.Extension,
			EmulatorTag:       string(repository.Python3.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Python2.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.Python2.Name),
			EmulatorExtension: repository.Python2.Extension,
			EmulatorTag:       string(repository.Python2.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Python3.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.Python3.Name),
			EmulatorExtension: repository.Python3.Extension,
			EmulatorTag:       string(repository.Python3.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Python2.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.Python2.Name),
			EmulatorExtension: repository.Python2.Extension,
			EmulatorTag:       string(repository.Python2.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Python3.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.Python3.Name),
			EmulatorExtension: repository.Python3.Extension,
			EmulatorTag:       string(repository.Python3.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Php74.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.Php74.Name),
			EmulatorExtension: repository.Php74.Extension,
			EmulatorTag:       string(repository.Php74.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Php74.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.Php74.Name),
			EmulatorExtension: repository.Php74.Extension,
			EmulatorTag:       string(repository.Php74.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Php74.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.Php74.Name),
			EmulatorExtension: repository.Php74.Extension,
			EmulatorTag:       string(repository.Php74.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Haskell.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.Haskell.Name),
			EmulatorExtension: repository.Haskell.Extension,
			EmulatorTag:       string(repository.Haskell.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Haskell.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.Haskell.Name),
			EmulatorExtension: repository.Haskell.Extension,
			EmulatorTag:       string(repository.Haskell.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Haskell.Tag),
			},
		})).Should(gomega.BeNil())

		result := execution.Service(_var.SINGLE_FILE_EXECUTION).RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(repository.Haskell.Name),
			EmulatorExtension: repository.Haskell.Extension,
			EmulatorTag:       string(repository.Haskell.Tag),
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
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.NodeLts.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `console.log("mile")`, repository.NodeLts, activeSession)
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

		var result repository.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a PHP environment", func() {
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Php74.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `
<?php

echo "mile";
`, repository.Php74, activeSession)
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

		var result repository.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("\nmile"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Ruby environment", func() {
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Ruby.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `puts "mile"`, repository.Ruby, activeSession)
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

		var result repository.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Go environment", func() {
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.GoLang.Tag),
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
`, repository.GoLang, activeSession)
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

		var result repository.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a C# (Mono) environment", func() {
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.CSharpMono.Tag),
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
`, repository.CSharpMono, activeSession)
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

		var result repository.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Python2 environment", func() {
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Python2.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `print("mile")`, repository.Python2, activeSession)
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

		var result repository.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Python3 environment", func() {
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Python3.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `print("mile")`, repository.Python3, activeSession)
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

		var result repository.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a Haskell environment", func() {
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Haskell.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `main = putStrLn "mile"`, repository.Haskell, activeSession)
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

		var result repository.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a C environment", func() {
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.CLang.Tag),
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
`, repository.CLang, activeSession)
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

		var result repository.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should execute a single file in a C++ environment", func() {
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.CPlus.Tag),
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
`, repository.CPlus, activeSession)
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

		var result repository.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("mile"))

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})

	GinkgoIt("Should gracefully fail multiple concurrent requests and stop containers", func() {
		ginkgo.Skip("")

		gomega.Expect(execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    10,
				ContainerNum: 1,
				Tag:          string(repository.NodeLts.Tag),
			},
		})).Should(gomega.BeNil())

		activeSession := testCreateAccount()

		wg := &sync.WaitGroup{}
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				pg := testCreateEmptyPage(activeSession)
				cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
				testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `while(true) {}`, repository.NodeLts, activeSession)
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

				var result repository.RunResult
				gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

				gomega.Expect(result.Success).Should(gomega.BeFalse())

				wg.Done()
			}(wg)
		}

		wg.Wait()

		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()
	})
})
