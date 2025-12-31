package resolve

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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
	sig := "TERM"
	if force {
		sig = "KILL"
	}
	return fmt.Errorf("kill not implemented yet (would send SIG%s to pid %d)", sig, p.PID)
}

// ResolveProcess creates a ProcessOwner from a PID using /proc
func ResolveProcess(pid int) (*ProcessOwner, error) {
	proc := fmt.Sprintf("/proc/%d", pid)

	// --- Name ---
	comm, err := os.ReadFile(filepath.Join(proc, "comm"))
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(string(comm))

	// --- Cmdline ---
	cmdBytes, _ := os.ReadFile(filepath.Join(proc, "cmdline"))
	cmdline := strings.ReplaceAll(string(cmdBytes), "\x00", " ")
	cmdline = strings.TrimSpace(cmdline)

	// --- Status (UID, PPID) ---
	status, err := os.ReadFile(filepath.Join(proc, "status"))
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

