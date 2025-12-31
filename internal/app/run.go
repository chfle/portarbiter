package app

import (
	"fmt"

	"portarbiter/internal/detect"
	"portarbiter/internal/resolve"
	"portarbiter/pkg/model"
)

type Options struct {
	Port   int
	DryRun bool
	Kill   bool
	Force  bool
}

func Run(opts Options) int {
	if dock, ok, err := resolve.ResolveDockerByPort(opts.Port); err == nil && ok {
		fmt.Printf("Port %d is used by:\n", opts.Port)
		exitCode := 0
		handleOwner(dock, opts, &exitCode)
		return exitCode
	}

	pids, err := detect.FindPIDsByPort(opts.Port)
	if err != nil {
		fmt.Println("Error:", err)
		return 3
	}

	fmt.Printf("Port %d is used by:\n", opts.Port)

	seen := make(map[string]bool)
	exitCode := 0

	for _, pid := range pids {

		if dock, ok, err := resolve.ResolveDocker(pid); err == nil && ok {
			key := dock.Type().String() + ":" + dock.ID()
			if seen[key] {
				continue
			}
			seen[key] = true
			handleOwner(dock, opts, &exitCode)
			continue
		}

		if sys, ok, err := resolve.ResolveSystemd(pid); err == nil && ok {
			key := sys.Type().String() + ":" + sys.ID()
			if seen[key] {
				continue
			}
			seen[key] = true
			handleOwner(sys, opts, &exitCode)
			continue
		}

		proc, err := resolve.ResolveProcess(pid)
		if err != nil {
			fmt.Printf("  pid=%d error=%v\n", pid, err)
			exitCode = 4
			continue
		}

		key := proc.Type().String() + ":" + proc.ID()
		if seen[key] {
			continue
		}
		seen[key] = true
		handleOwner(proc, opts, &exitCode)
	}

	return exitCode
}

func handleOwner(owner model.Owner, opts Options, exitCode *int) {
	fmt.Println(" ", owner.Describe())

	if opts.DryRun {
		fmt.Println("   action: DRY-RUN")
		return
	}

	if opts.Kill {
		if err := owner.Kill(opts.Force); err != nil {
			fmt.Println("   action: FAILED:", err)
			*exitCode = 10
			return
		}
		fmt.Println("   action: TERMINATED")
	}
}

