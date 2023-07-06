package main

import (
	"emulator/cmd/http"
	"emulator/pkg"
	"emulator/pkg/logger"
	"emulator/pkg/types"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	loadEnv()
	logger.BuildLoggers()
	emulator := pkg.NewEmulator(pkg.Options{
		GoLang: pkg.GoLang{
			Workers:    10,
			Containers: 10,
		},
		NodeLts: pkg.NodeLts{
			Workers:    10,
			Containers: 10,
		},
		Ruby: pkg.Ruby{
			Workers:    10,
			Containers: 10,
		},
		ExecutionDirectory: os.Getenv("EXECUTION_DIR"),
	})

	result := emulator.RunJob(
		string(types.Ruby.Name),
		fmt.Sprintf(`puts "Hello world"`))

	fmt.Println(result)

	http.CloseExecutioners()
}
