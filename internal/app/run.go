package app

import (
	"fmt"

	"portarbiter/internal/detect"
	"portarbiter/internal/policy"
	"portarbiter/internal/resolve"
	"portarbiter/pkg/model"
)

// Exit codes:
//   0  - success / allowed / dry-run
//   10 - termination failed
//   20 - denied or confirmation required
//   30 - detection failed

func Run(opts Options) int {
	if dockerOwner, ok, err := resolve.ResolveDockerByPort(opts.Port); err == nil && ok {
		return handleOwner(opts, dockerOwner)
	}

	pids, err := detect.FindPIDsByPort(opts.Port)
	if err != nil || len(pids) == 0 {
		if opts.ReasonOnly {
			fmt.Println("DETECTION_FAILED")
		}
		return 30
	}

	for _, pid := range pids {

		// Docker container process
		if dockerOwner, ok, err := resolve.ResolveDocker(pid); err == nil && ok {
			return handleOwner(opts, dockerOwner)
		}

		// systemd service
		if sysOwner, ok, err := resolve.ResolveSystemd(pid); err == nil && ok {
			return handleOwner(opts, sysOwner)
		}

		// raw process
		if procOwner, err := resolve.ResolveProcess(pid); err == nil {
			return handleOwner(opts, procOwner)
		}
	}

	if opts.ReasonOnly {
		fmt.Println("DETECTION_FAILED")
	}
	return 30
}

func handleOwner(opts Options, owner model.Owner) int {
	decision := policy.Evaluate(owner.Type(), owner.Describe())

	if opts.CheckOnly {
		if decision.Level == policy.Allowed {
			return 0
		}
		return 20
	}

	if opts.ReasonOnly {
		fmt.Printf("%s: %s\n", decision.Level.String(), decision.Reason)

		if decision.Level == policy.Allowed {
			return 0
		}
		return 20
	}

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

