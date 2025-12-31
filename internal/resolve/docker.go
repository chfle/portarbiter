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

func (d *DockerOwner) Describe() string {
	if d.ComposeProj != "" {
		return fmt.Sprintf(
			"docker-compose project=%s service=%s container=%s pid=%d image=%s",
			d.ComposeProj, d.ComposeSvc, d.ContainerName, d.PID, d.Image,
		)
	}

	return fmt.Sprintf(
		"docker container=%s pid=%d image=%s",
		d.ContainerName, d.PID, d.Image,
	)
}

func (d *DockerOwner) Kill(force bool) error {
	if d.ComposeProj != "" {
		return fmt.Errorf(
			"container is part of docker-compose project %s (use: docker compose down)",
			d.ComposeProj,
		)
	}
	return fmt.Errorf(
		"docker-managed container %s (use: docker stop %s)",
		d.ContainerName,
	)
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

	// docker inspect
	cmd := exec.Command("docker", "inspect", containerID)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return nil, false, fmt.Errorf("docker inspect failed: %w", err)
	}

	output := out.String()

	name := extractJSONField(output, `"Name":`)
	image := extractJSONField(output, `"Image":`)

	composeProj := extractJSONLabel(output, "com.docker.compose.project")
	composeSvc := extractJSONLabel(output, "com.docker.compose.service")

	return &DockerOwner{
		ContainerID:   containerID[:12],
		ContainerName: strings.TrimPrefix(name, "/"),
		Image:         image,
		PID:           pid,
		ComposeProj:   composeProj,
		ComposeSvc:    composeSvc,
	}, true, nil
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

