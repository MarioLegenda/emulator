package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"sync"
	"therebelsource/emulator/appErrors"
	errorHandler "therebelsource/emulator/appErrors"
	"therebelsource/emulator/execution"
	"therebelsource/emulator/logger"
	"therebelsource/emulator/projectExecution"
	"therebelsource/emulator/rateLimiter"
	"therebelsource/emulator/repository"
	"therebelsource/emulator/singleFileExecution"
	"therebelsource/emulator/slack"
	_var "therebelsource/emulator/var"
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
	err := execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
		{
			WorkerNum:    getEnvironmentWorkers("NODE_LTS"),
			ContainerNum: getEnvironmentContainers("NODE_LTS"),
			Tag:          string(repository.NodeLts.Tag),
		},
		{
			WorkerNum:    getEnvironmentWorkers("JULIA"),
			ContainerNum: getEnvironmentContainers("JULIA"),
			Tag:          string(repository.Julia.Tag),
		},
		{
			WorkerNum:    getEnvironmentWorkers("NODE_ESM"),
			ContainerNum: getEnvironmentContainers("NODE_ESM"),
			Tag:          string(repository.NodeEsm.Tag),
		},
		{
			WorkerNum:    getEnvironmentWorkers("RUBY"),
			ContainerNum: getEnvironmentContainers("RUBY"),
			Tag:          string(repository.Ruby.Tag),
		},
		{
			WorkerNum:    getEnvironmentWorkers("RUST"),
			ContainerNum: getEnvironmentContainers("RUST"),
			Tag:          string(repository.Rust.Tag),
		},
		{
			WorkerNum:    getEnvironmentWorkers("CPLUS"),
			ContainerNum: getEnvironmentContainers("CPLUS"),
			Tag:          string(repository.CPlus.Tag),
		},
		{
			WorkerNum:    getEnvironmentWorkers("HASKELL"),
			ContainerNum: getEnvironmentContainers("HASKELL"),
			Tag:          string(repository.Haskell.Tag),
		},
		{
			WorkerNum:    getEnvironmentWorkers("C"),
			ContainerNum: getEnvironmentContainers("C"),
			Tag:          string(repository.CLang.Tag),
		},
		{
			WorkerNum:    getEnvironmentWorkers("PERL"),
			ContainerNum: getEnvironmentContainers("PERL"),
			Tag:          string(repository.PerlLts.Tag),
		},
		{
			WorkerNum:    getEnvironmentWorkers("C_SHARP"),
			ContainerNum: getEnvironmentContainers("C_SHARP"),
			Tag:          string(repository.CSharpMono.Tag),
		},
		{
			WorkerNum:    getEnvironmentWorkers("PYTHON3"),
			ContainerNum: getEnvironmentContainers("PYTHON3"),
			Tag:          string(repository.Python3.Tag),
		},
		{
			WorkerNum:    getEnvironmentWorkers("LUA"),
			ContainerNum: getEnvironmentContainers("LUA"),
			Tag:          string(repository.Lua.Tag),
		},
		{
			WorkerNum:    getEnvironmentWorkers("PYTHON2"),
			ContainerNum: getEnvironmentContainers("PYTHON2"),
			Tag:          string(repository.Python2.Tag),
		},
		{
			WorkerNum:    getEnvironmentWorkers("PHP74"),
			ContainerNum: getEnvironmentContainers("PHP74"),
			Tag:          string(repository.Php74.Tag),
		},
		{
			WorkerNum:    getEnvironmentWorkers("GO"),
			ContainerNum: getEnvironmentContainers("GO"),
			Tag:          string(repository.GoLang.Tag),
		},
	})

	if err != nil {
		slack.SendErrorLog(err, "deploy_log")
		logger.Error(fmt.Sprintf("Cannot boot project execution: %s", err.Error()))

		if !execution.Service(_var.PROJECT_EXECUTION).Closed() {
			execution.Service(_var.PROJECT_EXECUTION).Close()
		}

		time.Sleep(5 * time.Second)

		if os.Getenv("APP_ENV") == "prod" {
			execution.FinalCleanup(true)
		}

		appErrors.TerminateWithMessage("Cannot boot executioner. Server cannot start!")
	}
}

func closeExecutioners() {
	wg := sync.WaitGroup{}
	for _, e := range []string{_var.PROJECT_EXECUTION} {
		wg.Add(1)

		go func(name string, wg *sync.WaitGroup) {
			if !execution.Service(name).Closed() {
				execution.Service(name).Close()
			}
			wg.Done()
		}(e, &wg)
	}
	wg.Wait()

	time.Sleep(5 * time.Second)
	execution.FinalCleanup(true)
}

func App() {
	loadEnv()
	logger.BuildLoggers()
	initRequiredDirectories(true)

	rateLimiter.InitRateLimiter()

	singleFileExecution.InitService()
	projectExecution.InitService()

	initExecutioners()

	WatchServerShutdown(InitServer(RegisterRoutes()))
}
