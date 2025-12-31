package detect

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// FindPIDsByPort returns a list of PIDs listening on the given TCP port.
// Linux-only. Uses `ss -lptn`.
func FindPIDsByPort(port int) ([]int, error) {
	cmd := exec.Command("ss", "-lptn")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to execute ss: %w", err)
	}

	lines := strings.Split(out.String(), "\n")
	pids := make(map[int]struct{})

	portStr := ":" + strconv.Itoa(port)

	for _, line := range lines {
		if !strings.Contains(line, portStr) {
			continue
		}

		// Look for pid=XXXX
		idx := strings.Index(line, "pid=")
		if idx == -1 {
			continue
		}

		rest := line[idx+4:]
		end := strings.IndexAny(rest, ",)")
		if end == -1 {
			continue
		}

		pidStr := rest[:end]
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		pids[pid] = struct{}{}
	}

	if len(pids) == 0 {
		return nil, errors.New("no process listening on given port")
	}

	result := make([]int, 0, len(pids))
	for pid := range pids {
		result = append(result, pid)
	}

	return result, nil
}

