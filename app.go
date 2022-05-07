package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"therebelsource/emulator/appErrors"
	errorHandler "therebelsource/emulator/appErrors"
	"therebelsource/emulator/execution"
	"therebelsource/emulator/projectExecution"
	"therebelsource/emulator/rateLimiter"
	"therebelsource/emulator/singleFileExecution"
)

func LoadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}
}

func InitRequiredDirectories(output bool) {
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

func App() {
	LoadEnv()
	InitRequiredDirectories(true)

	rateLimiter.InitRateLimiter()

	singleFileExecution.InitService()
	projectExecution.InitService()

	err := execution.Init(20, 1)
	if err != nil {
		appErrors.TerminateWithMessage("Cannot boot executioner. Server cannot start!")
	}

	WatchServerShutdown(InitServer(RegisterRoutes()))
}
