package main

import (
	"emulator/pkg"
	"github.com/joho/godotenv"
	"log"
)

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}
}

func App() {
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
		LogDirectory:       "/home/mario/go/emulator/var/log",
		ExecutionDirectory: "/home/mario/go/go-emulator/var/execution",
	})

	emulator.Close()
}
