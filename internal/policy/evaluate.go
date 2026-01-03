package policy

import (
	"strings"

	"portarbiter/pkg/model"
)

/*
OwnerClass is a semantic classification independent of concrete model types.
*/
type OwnerClass int

const (
	OwnerUnknown OwnerClass = iota
	OwnerDocker
	OwnerSystemd
	OwnerProcess
)

func (c OwnerClass) String() string {
	switch c {
	case OwnerDocker:
		return "docker"
	case OwnerSystemd:
		return "systemd"
	case OwnerProcess:
		return "process"
	default:
		return "unknown"
	}
}

/*
Evaluate determines the policy decision for a resolved owner.
*/
func Evaluate(ownerType model.OwnerType, description string) Decision {
	class := classifyOwner(ownerType)

	switch class {

	case OwnerDocker:
		return Decision{
			Level:  ConfirmRequired,
			Reason: "docker container or compose project",
		}

	case OwnerSystemd:
		if isInteractiveSSH(description) {
			return Decision{
				Level:  Denied,
				Reason: "interactive SSH session",
			}
		}

		return Decision{
			Level:  ConfirmRequired,
			Reason: "systemd-managed service",
		}

	case OwnerProcess:
		return Decision{
			Level:  Allowed,
			Reason: "raw user process",
		}

	default:
		return Decision{
			Level:  Denied,
			Reason: "unknown owner type",
		}
	}
}

/*
classifyOwner maps concrete owner types to semantic classes.
*/
func classifyOwner(ownerType model.OwnerType) OwnerClass {
	t := strings.ToLower(ownerType.String())

	switch {
	case strings.Contains(t, "docker"):
		return OwnerDocker
	case strings.Contains(t, "systemd"):
		return OwnerSystemd
	case strings.Contains(t, "process"):
		return OwnerProcess
	default:
		return OwnerUnknown
	}
}

func isInteractiveSSH(description string) bool {
	d := strings.ToLower(description)
	return strings.Contains(d, "ssh")
}

