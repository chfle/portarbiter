package app

import (
	"fmt"

	"portarbiter/internal/detect"
	"portarbiter/internal/policy"
	"portarbiter/internal/resolve"
	"portarbiter/pkg/model"
)

func Run(opts Options) int {
	pids, err := detect.FindPIDsByPort(opts.Port)
	if err != nil || len(pids) == 0 {
		fmt.Printf("Failed to resolve port %d\n", opts.Port)
		return 30
	}

	if dockerOwner, ok, err := resolve.ResolveDockerByPort(opts.Port); err == nil && ok {
		return handleOwner(opts, dockerOwner)
	}

	for _, pid := range pids {

		if dockerOwner, ok, err := resolve.ResolveDocker(pid); err == nil && ok {
			return handleOwner(opts, dockerOwner)
		}

		if sysOwner, ok, err := resolve.ResolveSystemd(pid); err == nil && ok {
			return handleOwner(opts, sysOwner)
		}

		if procOwner, err := resolve.ResolveProcess(pid); err == nil {
			return handleOwner(opts, procOwner)
		}
	}

	fmt.Printf("Failed to resolve owner for port %d\n", opts.Port)
	return 30
}

func handleOwner(opts Options, owner model.Owner) int {
	decision := policy.Evaluate(owner.Type(), owner.Describe())

	fmt.Printf("Port %d is used by:\n", opts.Port)
	fmt.Printf("  %s\n", owner.Describe())

	if opts.DryRun {
		fmt.Printf("  action: DRY-RUN\n")
		return 0
	}

	switch decision.Level {

	case policy.Denied:
		fmt.Printf("  action: DENIED (%s)\n", decision.Reason)
		return 20

	case policy.ConfirmRequired:
		if !opts.Yes {
			fmt.Printf("  action: CONFIRM REQUIRED (%s)\n", decision.Reason)
			return 20
		}

	case policy.Allowed:
		// continue
	}

	if err := owner.Kill(opts.Force); err != nil {
		fmt.Printf("  action: FAILED: %v\n", err)
		return 10
	}

	fmt.Printf("  action: TERMINATED\n")
	return 0
}

