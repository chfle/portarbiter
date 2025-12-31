package policy

import (
	"strings"

	"portarbiter/pkg/model"
)

func Evaluate(ownerType model.OwnerType, description string) Decision {
	switch ownerType {

	case model.OwnerDocker:
		return Decision{
			Level:  ConfirmRequired,
			Reason: "docker container or compose project",
		}

	case model.OwnerSystemd:
		// SSH must never be killed blindly
		if isSSH(description) {
			return Decision{
				Level:  Denied,
				Reason: "interactive SSH session",
			}
		}

		return Decision{
			Level:  ConfirmRequired,
			Reason: "systemd-managed service",
		}

	case model.OwnerProcess:
		return Decision{
			Level:  Allowed,
			Reason: "raw user process",
		}
	}

	return Decision{
		Level:  Denied,
		Reason: "unknown owner type",
	}
}

func isSSH(desc string) bool {
	d := strings.ToLower(desc)
	return strings.Contains(d, "ssh")
}

