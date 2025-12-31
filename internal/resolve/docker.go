package resolve

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"portarbiter/pkg/model"
)

type DockerOwner struct {
	ContainerID   string
	ContainerName string
	Image         string
	PID           int
	ComposeProj   string
	ComposeSvc    string
}

func (d *DockerOwner) Type() model.OwnerType {
	if d.ComposeProj != "" {
		return model.OwnerCompose
	}
	return model.OwnerDocker
}

func (d *DockerOwner) ID() string {
	return d.ContainerID
}

func (d *DockerOwner) shortID() string {
	if len(d.ContainerID) >= 12 {
		return d.ContainerID[:12]
	}
	return d.ContainerID
}

func (d *DockerOwner) Describe() string {
	if d.ComposeProj != "" {
		return fmt.Sprintf(
			"docker-compose project=%s service=%s container=%s id=%s pid=%d image=%s",
			d.ComposeProj, d.ComposeSvc, d.ContainerName, d.shortID(), d.PID, d.Image,
		)
	}

	return fmt.Sprintf(
		"docker container=%s id=%s pid=%d image=%s",
		d.ContainerName, d.shortID(), d.PID, d.Image,
	)
}

func (d *DockerOwner) Kill(force bool) error {
	// Compose project: stop the project, not a single container
	if d.ComposeProj != "" {
		// Use docker compose with project name so we don't need the compose directory.
		// This will target the project by name.
		args := []string{"compose", "-p", d.ComposeProj, "down"}
		if force {
			// "down" doesn't have a true "force" semantic; we speed up termination.
			args = append(args, "--timeout", "1")
		}

		if err := runDocker(args...); err != nil {
			return fmt.Errorf("failed to bring down compose project %s: %w", d.ComposeProj, err)
		}
		return nil
	}

	// Single container
	if force {
		if err := runDocker("kill", d.ContainerID); err != nil {
			return fmt.Errorf("docker kill %s failed: %w", d.shortID(), err)
		}
		return nil
	}

	if err := runDocker("stop", d.ContainerID); err != nil {
		return fmt.Errorf("docker stop %s failed: %w", d.shortID(), err)
	}
	return nil
}

// ResolveDocker checks if PID belongs to a Docker container via cgroups
func ResolveDocker(pid int) (*DockerOwner, bool, error) {
	cgroupPath := filepath.Join("/proc", fmt.Sprintf("%d", pid), "cgroup")
	data, err := os.ReadFile(cgroupPath)
	if err != nil {
		return nil, false, err
	}

	var containerID string
	for _, line := range strings.Split(string(data), "\n") {
		// typical:
		// 0::/docker/<container-id>
		if strings.Contains(line, "/docker/") {
			parts := strings.Split(line, "/docker/")
			if len(parts) == 2 {
				containerID = strings.TrimSpace(parts[1])
				break
			}
		}
	}

	if containerID == "" {
		return nil, false, nil
	}

	txt, err := dockerInspect(containerID)
	if err != nil {
		return nil, false, err
	}

	name := extractJSONField(txt, `"Name":`)
	image := extractJSONField(txt, `"Image":`)
	composeProj := extractJSONLabel(txt, "com.docker.compose.project")
	composeSvc := extractJSONLabel(txt, "com.docker.compose.service")

	return &DockerOwner{
		ContainerID:   containerID,
		ContainerName: strings.TrimPrefix(name, "/"),
		Image:         image,
		PID:           pid,
		ComposeProj:   composeProj,
		ComposeSvc:    composeSvc,
	}, true, nil
}

func dockerInspect(containerID string) (string, error) {
	cmd := exec.Command("docker", "inspect", containerID)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("docker inspect failed: %w (%s)", err, out.String())
	}
	return out.String(), nil
}

func runDocker(args ...string) error {
	cmd := exec.Command("docker", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w (%s)", err, out.String())
	}
	return nil
}

// minimal JSON scraping (we avoid full JSON parsing intentionally)
func extractJSONField(data, key string) string {
	idx := strings.Index(data, key)
	if idx == -1 {
		return ""
	}
	sub := data[idx+len(key):]
	sub = strings.TrimSpace(sub)
	sub = strings.TrimPrefix(sub, `"`)
	end := strings.Index(sub, `"`)
	if end == -1 {
		return ""
	}
	return sub[:end]
}

func extractJSONLabel(data, label string) string {
	key := fmt.Sprintf(`"%s":`, label)
	idx := strings.Index(data, key)
	if idx == -1 {
		return ""
	}
	sub := data[idx+len(key):]
	sub = strings.TrimSpace(sub)
	sub = strings.TrimPrefix(sub, `"`)
	end := strings.Index(sub, `"`)
	if end == -1 {
		return ""
	}
	return sub[:end]
}

