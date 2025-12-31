package resolve

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

// ResolveSystemd determines whether a PID truly belongs to a systemd service.
// IMPORTANT:
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

	if err := cmd.Run(); err != nil {
		return nil, false, nil
	}

	var service string
	for _, line := range strings.Split(out.String(), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "●") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				service = fields[1]
				break
			}
		}
	}

	if service == "" {
		return nil, false, nil
	}

	// --- SSH SERVICE SPECIAL RULE ---
	if service == "ssh.service" {
		commPath := filepath.Join("/proc", strconv.Itoa(pid), "comm")
		comm, err := os.ReadFile(commPath)
		if err == nil {
			procName := strings.TrimSpace(string(comm))
			// Only sshd itself is service-owned
			if procName != "sshd" {
				return nil, false, nil
			}
		}
	}

	return &SystemdOwner{
		Service: service,
		PID:     pid,
	}, true, nil
}

