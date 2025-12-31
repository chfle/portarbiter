package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"portarbiter/internal/app"
)

func main() {
	var (
		doKill bool
		force  bool
		dryRun bool
	)

	flag.BoolVar(&doKill, "kill", false, "terminate the owner holding the port")
	flag.BoolVar(&force, "force", false, "force termination")
	flag.BoolVar(&dryRun, "dry-run", false, "show what would be done")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: portarbiter [--dry-run] [--kill] [--force] <port>")
		os.Exit(1)
	}

	if !doKill {
		dryRun = true
	}

	port, err := strconv.Atoi(flag.Arg(0))
	if err != nil || port <= 0 || port > 65535 {
		fmt.Println("Invalid port")
		os.Exit(2)
	}

	exitCode := app.Run(app.Options{
		Port:   port,
		DryRun: dryRun,
		Kill:   doKill,
		Force:  force,
	})

	os.Exit(exitCode)
}

