package resolve

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"portarbiter/pkg/model"
)

type SystemdOwner struct {
	Service string
	PID     int
}

func (s *SystemdOwner) Type() model.OwnerType {
	return model.OwnerSystemd
}

func (s *SystemdOwner) ID() string {
	return s.Service
}

func (s *SystemdOwner) Describe() string {
	return fmt.Sprintf(
		"systemd service=%s pid=%d",
		s.Service,
		s.PID,
	)
}

func (s *SystemdOwner) Kill(force bool) error {
	return fmt.Errorf(
		"systemd-managed process (service=%s); use systemctl stop %s",
		s.Service,
		s.Service,
	)
}

// ResolveSystemd checks whether a PID belongs to a systemd service.
// Returns (owner, true, nil) if yes
// Returns (nil, false, nil) if not systemd-managed
func ResolveSystemd(pid int) (*SystemdOwner, bool, error) {
	cmd := exec.Command(
		"systemctl",
		"status",
		strconv.Itoa(pid),
		"--no-pager",
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		// systemctl returns non-zero if PID is unknown
		return nil, false, nil
	}

	for _, line := range strings.Split(out.String(), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "●") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				service := fields[1]
				return &SystemdOwner{
					Service: service,
					PID:     pid,
				}, true, nil
			}
		}
	}

	return nil, false, nil
}

