package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
	"therebelsource/emulator/appErrors"
	errorHandler "therebelsource/emulator/appErrors"
	"therebelsource/emulator/execution"
	"therebelsource/emulator/projectExecution"
	"therebelsource/emulator/rateLimiter"
	"therebelsource/emulator/runner"
	"therebelsource/emulator/singleFileExecution"
	_var "therebelsource/emulator/var"
)

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}
}

func initRequiredDirectories(output bool) {
	projectsDir := os.Getenv("PROJECTS_DIR")
	directoriesExist := true
	if _, err := os.Stat(projectsDir); os.IsNotExist(err) {
		directoriesExist = false

		if output {
			fmt.Println("")
			fmt.Println("Creating required directories...")
		}
		fsErr := os.Mkdir(projectsDir, os.ModePerm)

		if fsErr != nil {
			errorHandler.TerminateWithMessage(fmt.Sprintf("Cannot create %s directory", projectsDir))
		}
	}

	if !directoriesExist {
		rest := []string{
			os.Getenv("CODE_PROJECT_STATE_DIR"),
			os.Getenv("SINGLE_FILE_STATE_DIR"),
			os.Getenv("PACKAGES_DIR"),
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
			fmt.Println("Required directories already created! Skipping...")
			fmt.Println("")
		}
	}

	if !directoriesExist {
		if output {
			fmt.Println("Required directories created!")
			fmt.Println("")
		}
	}
}

func initExecutioners() {
	err := execution.Init(_var.SINGLE_FILE_EXECUTION, []execution.ContainerBlueprint{
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.NodeLts.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.NodeEsm.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.Ruby.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.Rust.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.CPlus.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.Haskell.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.CLang.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.CSharpMono.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.Python3.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.Python2.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.Php74.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.GoLang.Tag),
		},
	})

	if err != nil {
		appErrors.TerminateWithMessage("Cannot boot executioner. Server cannot start!")
	}

	err = execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.NodeLts.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.NodeEsm.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.Ruby.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.Rust.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.CPlus.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.Haskell.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.CLang.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.CSharpMono.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.Python3.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.Python2.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.Php74.Tag),
		},
		{
			WorkerNum:    1,
			ContainerNum: 1,
			Tag:          string(runner.GoLang.Tag),
		},
	})

	if err != nil {
		execution.Service(_var.SINGLE_FILE_EXECUTION).Close()

		appErrors.TerminateWithMessage("Cannot boot executioner. Server cannot start!")
	}
}

func closeExecutioners() {
	wg := sync.WaitGroup{}
	for _, e := range []string{_var.SINGLE_FILE_EXECUTION, _var.PROJECT_EXECUTION} {
		wg.Add(1)

		go func(name string, wg *sync.WaitGroup) {
			execution.Service(name).Close()
			wg.Done()
		}(e, &wg)
	}
	wg.Wait()
}

func App() {
	loadEnv()
	initRequiredDirectories(true)

	rateLimiter.InitRateLimiter()

	singleFileExecution.InitService()
	projectExecution.InitService()

	initExecutioners()

	WatchServerShutdown(InitServer(RegisterRoutes()))
}
