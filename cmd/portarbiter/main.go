package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"portarbiter/internal/app"
	"portarbiter/internal/version"
)

func main() {
	var (
		doKill   bool
		force    bool
		dryRun   bool
		yes      bool
		showVer  bool
	)

	flag.BoolVar(&doKill, "kill", false, "terminate the owner holding the port")
	flag.BoolVar(&force, "force", false, "force termination")
	flag.BoolVar(&dryRun, "dry-run", false, "show what would be done")
	flag.BoolVar(&yes, "yes", false, "auto-confirm dangerous actions")
	flag.BoolVar(&yes, "y", false, "auto-confirm dangerous actions (shorthand)")
	flag.BoolVar(&showVer, "version", false, "show version information")
	flag.Parse()

	if showVer {
		fmt.Println(version.String())
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		fmt.Println("Usage: portarbiter [--dry-run] [--kill] [--force] [--yes] <port>")
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
		Yes:    yes,
	})

	os.Exit(exitCode)
}

