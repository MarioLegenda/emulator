package main

import (
	"emulator/cmd/http"
	errorHandler "emulator/pkg/appErrors"
	execution2 "emulator/pkg/execution"
	"emulator/pkg/logger"
	"emulator/pkg/projectExecution"
	"emulator/pkg/repository"
	"emulator/pkg/singleFileExecution"
	_var "emulator/var"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}
}

func getEnvironmentWorkers(part string) int {
	workers, _ := strconv.Atoi(os.Getenv(fmt.Sprintf("%s_WORKERS", part)))

	return workers
}

func getEnvironmentContainers(part string) int {
	containers, _ := strconv.Atoi(os.Getenv(fmt.Sprintf("%s_CONTAINERS", part)))

	return containers
}

func createBlueprint(name, tag string) execution2.ContainerBlueprint {
	return execution2.ContainerBlueprint{
		WorkerNum:    getEnvironmentWorkers(name),
		ContainerNum: getEnvironmentContainers(name),
		Tag:          tag,
	}
}

func initRequiredDirectories(output bool) {
	projectsDir := os.Getenv("EXECUTION_DIR")
	directoriesExist := true
	if _, err := os.Stat(projectsDir); os.IsNotExist(err) {
		directoriesExist = false

		if output {
			fmt.Println("")
			logger.Info("Creating required directories...")
		}
		fsErr := os.Mkdir(projectsDir, os.ModePerm)

		if fsErr != nil {
			errorHandler.TerminateWithMessage(fmt.Sprintf("Cannot create %s directory", projectsDir))
		}
	}

	if !directoriesExist {
		rest := []string{
			os.Getenv("EXECUTION_DIR"),
		}

		for _, dir := range rest {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				fsErr := os.Mkdir(dir, os.ModePerm)

				if fsErr != nil {
					errorHandler.TerminateWithMessage(fmt.Sprintf("Cannot create %s directory", dir))
				}
			}
		}
	} else {
		if output {
			fmt.Println("")
			logger.Info("Required directories already created! Skipping...")
			fmt.Println("")
		}
	}

	if !directoriesExist {
		if output {
			logger.Info("Required directories created!")
			fmt.Println("")
		}
	}
}

func initExecutioners() {
	err := execution2.Init(_var.PROJECT_EXECUTION, []execution2.ContainerBlueprint{
		createBlueprint("NODE_LTS", string(repository.NodeLts.Tag)),
		createBlueprint("JULIA", string(repository.Julia.Tag)),
		createBlueprint("NODE_ESM", string(repository.NodeEsm.Tag)),
		createBlueprint("RUBY", string(repository.Ruby.Tag)),
		createBlueprint("RUST", string(repository.Rust.Tag)),
		createBlueprint("CPLUS", string(repository.CPlus.Tag)),
		createBlueprint("HASKELL", string(repository.Haskell.Tag)),
		createBlueprint("C", string(repository.CLang.Tag)),
		createBlueprint("PERL", string(repository.PerlLts.Tag)),
		createBlueprint("C_SHARP", string(repository.CSharpMono.Tag)),
		createBlueprint("PYTHON3", string(repository.Python3.Tag)),
		createBlueprint("LUA", string(repository.Lua.Tag)),
		createBlueprint("PYTHON2", string(repository.Python2.Tag)),
		createBlueprint("PHP74", string(repository.Php74.Tag)),
		createBlueprint("GO", string(repository.GoLang.Tag)),
	})

	if err != nil {
		logger.Error(fmt.Sprintf("Cannot boot project execution: %s", err.Error()))

		if !execution2.Service(_var.PROJECT_EXECUTION).Closed() {
			execution2.Service(_var.PROJECT_EXECUTION).Close()
		}

		time.Sleep(5 * time.Second)

		if os.Getenv("APP_ENV") == "prod" {
			execution2.FinalCleanup(true)
		}

		errorHandler.TerminateWithMessage("Cannot boot executioner. Server cannot start!")
	}
}

func main() {
	loadEnv()
	logger.BuildLoggers()
	initRequiredDirectories(true)

	singleFileExecution.InitService()
	projectExecution.InitService()

	initExecutioners()

	time.Sleep(2 * time.Second)

	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(current int) {
			res := execution2.Service(_var.PROJECT_EXECUTION).RunJob(execution2.Job{
				BuilderType:       "single_file",
				ExecutionType:     "single_file",
				EmulatorName:      string(repository.NodeLts.Name),
				EmulatorTag:       string(repository.NodeLts.Tag),
				EmulatorExtension: string(repository.NodeLts.Extension),
				EmulatorText:      fmt.Sprintf("console.log('Hello World -> %d')", current),
			})

			fmt.Println(res)

			wg.Done()
		}(i)
	}

	wg.Wait()

	http.CloseExecutioners()
}
