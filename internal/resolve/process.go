package resolve

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"portarbiter/pkg/model"
)

type ProcessOwner struct {
	PID     int
	Name    string
	Cmdline string
	UID     int
	PPID    int
}

func (p *ProcessOwner) Type() model.OwnerType {
	return model.OwnerProcess
}

func (p *ProcessOwner) ID() string {
	return strconv.Itoa(p.PID)
}

func (p *ProcessOwner) Describe() string {
	return fmt.Sprintf(
		"process pid=%d name=%s uid=%d ppid=%d cmd=%q",
		p.PID, p.Name, p.UID, p.PPID, p.Cmdline,
	)
}

func (p *ProcessOwner) Kill(force bool) error {
	proc, err := os.FindProcess(p.PID)
	if err != nil {
		return fmt.Errorf("cannot find process %d: %w", p.PID, err)
	}

	sig := syscall.SIGTERM
	if force {
		sig = syscall.SIGKILL
	}

	if err := proc.Signal(sig); err != nil {
		return fmt.Errorf("failed to send %s to pid %d: %w", sig.String(), p.PID, err)
	}

	return nil
}

// ResolveProcess creates a ProcessOwner from a PID using /proc
func ResolveProcess(pid int) (*ProcessOwner, error) {
	procPath := fmt.Sprintf("/proc/%d", pid)

	comm, err := os.ReadFile(procPath + "/comm")
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(string(comm))

	cmdBytes, _ := os.ReadFile(procPath + "/cmdline")
	cmdline := strings.ReplaceAll(string(cmdBytes), "\x00", " ")
	cmdline = strings.TrimSpace(cmdline)

	status, err := os.ReadFile(procPath + "/status")
	if err != nil {
		return nil, err
	}

	var uid, ppid int
	for _, line := range strings.Split(string(status), "\n") {
		if strings.HasPrefix(line, "Uid:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				uid, _ = strconv.Atoi(fields[1])
			}
		}
		if strings.HasPrefix(line, "PPid:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				ppid, _ = strconv.Atoi(fields[1])
			}
		}
	}

	return &ProcessOwner{
		PID:     pid,
		Name:    name,
		Cmdline: cmdline,
		UID:     uid,
		PPID:    ppid,
	}, nil
}

