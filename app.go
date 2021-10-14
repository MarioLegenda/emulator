package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	errorHandler "therebelsource/api/appErrors"
	"therebelsource/emulator/runner"
	"therebelsource/emulator/singleFileExecution"
	"therebelsource/emulator/staticTypes"
)

func LoadEnv(env string) {
	if env == "" {
		env = staticTypes.APP_DEV_ENV
	}

	err := godotenv.Load(fmt.Sprintf(".env.%s", env))

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
			os.Getenv("DEV_STATE_DIR"),
			os.Getenv("PROD_STATE_DIR"),
			os.Getenv("SESSION_STATE_DIR"),
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
	LoadEnv(staticTypes.APP_DEV_ENV)
	InitRequiredDirectories(true)

	singleFileExecution.InitService()

	go runner.WatchContainers()

	WatchServerShutdown(InitServer(RegisterRoutes()))
}
