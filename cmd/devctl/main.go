package main

import (
	"fmt"
	"os"

	"github.com/yeixio/yggdrasil-core/internal/config"
	"github.com/yeixio/yggdrasil-core/internal/version"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "version":
		fmt.Printf("Yggdrasil %s (commit %s, built %s)\n", version.Version, version.Commit, version.BuildDate)
	case "paths":
		cfg := config.DefaultConfig()
		fmt.Printf("data:     %s\n", cfg.DataDir)
		fmt.Printf("models:   %s\n", cfg.ModelsDir)
		fmt.Printf("runtimes: %s\n", cfg.RuntimesDir)
		fmt.Printf("logs:     %s\n", cfg.LogsDir)
		fmt.Printf("db:       %s\n", cfg.DBPath)
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: yggctl <version|paths>\n")
}
