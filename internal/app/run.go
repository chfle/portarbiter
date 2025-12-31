package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"portarbiter/internal/detect"
	"portarbiter/internal/resolve"
	"portarbiter/pkg/model"
)

type Options struct {
	Port   int
	DryRun bool
	Kill   bool
	Force  bool
	Yes    bool
}

func Run(opts Options) int {
	if dock, ok, err := resolve.ResolveDockerByPort(opts.Port); err == nil && ok {
		fmt.Printf("Port %d is used by:\n", opts.Port)
		exitCode := 0
		if !confirmIfNeeded(dock, opts) {
			return 20
		}
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
	var owners []model.Owner
	exitCode := 0

	for _, pid := range pids {

		if dock, ok, err := resolve.ResolveDocker(pid); err == nil && ok {
			key := dock.Type().String() + ":" + dock.ID()
			if !seen[key] {
				seen[key] = true
				owners = append(owners, dock)
			}
			continue
		}

		if sys, ok, err := resolve.ResolveSystemd(pid); err == nil && ok {
			key := sys.Type().String() + ":" + sys.ID()
			if !seen[key] {
				seen[key] = true
				owners = append(owners, sys)
			}
			continue
		}

		proc, err := resolve.ResolveProcess(pid)
		if err != nil {
			fmt.Printf("  pid=%d error=%v\n", pid, err)
			exitCode = 4
			continue
		}

		key := proc.Type().String() + ":" + proc.ID()
		if !seen[key] {
			seen[key] = true
			owners = append(owners, proc)
		}
	}

	if opts.Kill && len(owners) > 1 && !opts.Yes {
		fmt.Println("WARNING: multiple owners will be affected:")
		for _, o := range owners {
			fmt.Println(" ", o.Describe())
		}
		if !askForConfirmation() {
			return 20
		}
	}

	for _, owner := range owners {
		if opts.Kill && !confirmIfNeeded(owner, opts) {
			continue
		}
		handleOwner(owner, opts, &exitCode)
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

func confirmIfNeeded(owner model.Owner, opts Options) bool {
	if opts.DryRun || opts.Yes {
		return true
	}

	if owner.Type() == model.OwnerCompose {
		fmt.Printf(
			"CONFIRM: docker-compose project will be brought DOWN (%s)\n",
			owner.ID(),
		)
		return askForConfirmation()
	}

	return true
}

func askForConfirmation() bool {
	fmt.Print("Proceed? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(strings.ToLower(line))
	return line == "y" || line == "yes"
}

